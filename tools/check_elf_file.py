#!/usr/bin/env python
#
# Copyright (C) 2018 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""ELF file checker.

This command ensures all undefined symbols in an ELF file can be resolved to
global (or weak) symbols defined in shared objects specified in DT_NEEDED
entries.
"""

from __future__ import print_function

import argparse
import collections
import os
import os.path
import re
import struct
import subprocess
import sys


ELF = collections.namedtuple(
  'ELF',
  ('dt_soname', 'dt_needed', 'imported', 'exported', 'header'))


_ELF_MAGIC = b'\x7fELF'


# Known architectures
_EM_386 = 3
_EM_ARM = 40
_EM_X86_64 = 62
_EM_AARCH64 = 183


# ELF header struct
_ELF_HEADER_STRUCT = (
  ('ei_magic', '4s'),
  ('ei_class', 'B'),
  ('ei_data', 'B'),
  ('ei_version', 'B'),
  ('ei_osabi', 'B'),
  ('ei_pad', '8s'),
  ('e_type', 'H'),
  ('e_machine', 'H'),
  ('e_version', 'I'),
)

_ELF_HEADER_STRUCT_FMT = ''.join(_fmt for _, _fmt in _ELF_HEADER_STRUCT)


ELFHeader = collections.namedtuple(
  'ELFHeader', [_name for _name, _ in _ELF_HEADER_STRUCT])


# Toolchain for different architectures
_TOOLCHAINS = {
  _EM_386: (
    'x86_64-linux-android',
    os.path.join('x86', 'x86_64-linux-android-4.9')
  ),
  _EM_ARM: (
    'arm-linux-androideabi',
    os.path.join('arm', 'arm-linux-androideabi-4.9'),
  ),
  _EM_X86_64: (
    'x86_64-linux-android',
    os.path.join('x86', 'x86_64-linux-android-4.9'),
  ),
  _EM_AARCH64: (
    'aarch64-linux-android',
    os.path.join('aarch64', 'aarch64-linux-android-4.9'),
  ),
}


def _get_os_name():
  """Get the host OS name."""
  if sys.platform == 'linux2':
    return 'linux'
  if sys.platform == 'darwin':
    return 'darwin'
  raise ValueError(sys.platform + ' is not supported')


def _get_build_top():
  """Find the build top of the source tree (equivalent to
  ${ANDROID_BUILD_TOP})."""
  prev_path = None
  curr_path = os.path.abspath(os.getcwd())
  while prev_path != curr_path:
    if os.path.exists(os.path.join(curr_path, '.repo')):
      return curr_path
    prev_path = curr_path
    curr_path = os.path.dirname(curr_path)
  return None


def get_objdump(elf_machine):
  """Get the best matching objdump command."""

  # Select a matching objdump according to the machine specified in the ELF
  # file.
  try:
    triple, toolchain_dir = _TOOLCHAINS[elf_machine]
  except KeyError:
    triple, toolchain_dir = _TOOLCHAINS[_EM_AARCH64]

  # Build the path to objdump executable relative to ${ANDROID_BUILD_TOP}.
  exe_path = os.path.join(
    'prebuilts', 'gcc', _get_os_name() + '-x86', toolchain_dir, 'bin',
    triple + '-objdump')
  if os.path.exists(exe_path):
    return exe_path

  # Find the ${ANDROID_BUILD_TOP} and look for the objdump executable again.
  exe_path = os.path.join(_get_build_top(), exe_path)
  if os.path.exists(exe_path):
    return exe_path

  # If we cannot find the prebuilt objdump executable, return plain 'objdump'.
  return 'objdump'


class ELFError(ValueError):
  pass


class ELFInvalidMagicError(ELFError):
  def __init__(self):
    super(ELFInvalidMagicError, self).__init__('bad ELF magic')


class ELFParser(object):
  @classmethod
  def _read_elf_header(cls, elf_file_path):
    """Read the ELF magic word from the beginning of the file."""
    with open(elf_file_path, 'rb') as elf_file:
      buf = elf_file.read(struct.calcsize(_ELF_HEADER_STRUCT_FMT))
      try:
        return ELFHeader(*struct.unpack(_ELF_HEADER_STRUCT_FMT, buf))
      except struct.error:
        return None


  @classmethod
  def open(cls, elf_file_path):
    """Open and parse the ELF file."""
    # Parse the ELF header.  The code below checks the magic word of the
    # ELF file and picks a matching objdump command.
    header = cls._read_elf_header(elf_file_path)
    if not header or header.ei_magic != _ELF_MAGIC:
      raise ELFInvalidMagicError()

    # Parse the output of objdump command.
    objdump = get_objdump(header.e_machine)
    return cls._parse_objdump(elf_file_path, header, objdump)


  # Objdump output patterns
  _DYNAMIC_PATTERN = re.compile(b'^Dynamic Section:$')

  _DYNAMIC_ENTRY_PATTERN = re.compile(b'^\\s{2,}([^\\s]+)\\s+(.*)$')

  _DYNSYM_PATTERN = re.compile(b'^DYNAMIC SYMBOL TABLE:$')

  _DYNSYM_ENTRY_PATTERN = re.compile(
    b'([0-9a-fA-F]+)'  # VMA
    b' '
    b'(?P<local>[!lgu ])'  # Local-and-global, local, global, or GNU unique
    b'(?P<weak>[w ])'  # Weak
    b'([C ])'  # Constructor
    b'([W ])'  # Warning
    b'([Ii ])'  # Indirect or GNU indirect function
    b'([dD ])'  # Debugging or dynamic
    b'([FfO ])'  # Function, file, or object
    b' '
    b'(?P<section>[^\\t]+)'  # Section name
    b'\t'
    b'([0-9a-fA-F]+)'  # Alignment or size
    b'\\s+'
    b'(?P<last>.*)'  # Version, visibility, and name (separated by space)
  )


  @classmethod
  def _find_pattern(cls, pattern, iterable):
    """Search for a regular expression pattern from lines iterator."""
    for line in iterable:
      match = pattern.match(line)
      if match:
        return match
    return None


  @classmethod
  def _parse_objdump(cls, elf_file_path, header, objdump):
    """Parse the output of the objdump."""

    # Default dt_soname
    dt_soname = os.path.basename(elf_file_path)

    # Run objdump and get the output.
    proc = subprocess.Popen([objdump, '-p', '-T', elf_file_path],
                            stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, _ = proc.communicate()

    lines = out.splitlines()
    lines_it = iter(lines)

    # Parse .dynamic section.
    dynamic = cls._find_pattern(cls._DYNAMIC_PATTERN, lines_it)
    if not dynamic:
      return ELF(dt_soname, [], {}, {}, header)

    dt_needed = []
    for line in lines_it:
      if not line:
        break

      match = cls._DYNAMIC_ENTRY_PATTERN.match(line)
      if not match:
        continue

      key = match.group(1)
      value = match.group(2)

      if key == 'SONAME':
        dt_soname = value
      elif key == 'NEEDED':
        dt_needed.append(value)

    # Parse .dynsym section.
    dynsym = cls._find_pattern(cls._DYNSYM_PATTERN, lines_it)
    if not dynsym:
      return ELF(dt_soname, dt_needed, {}, {}, header)

    exported = collections.defaultdict(set)
    imported = collections.defaultdict(set)
    for line in lines_it:
      if not line or line == b'no symbols':
        break

      match = cls._DYNSYM_ENTRY_PATTERN.match(line)
      if not match:
        continue

      section = match.group('section')
      is_undef = section == b'*UND*'
      is_local = match.group('local') in b'!l'
      is_weak = match.group('weak') == b'w'

      # Parse the version, visibility, and name tuple.
      version = None
      name = None

      last = re.split(b'\\s+', match.group('last'), 2)
      if len(last) == 3:
        version = last[0]
        name = last[2]
      elif len(last) == 2:
        if last[0] not in {b'.internal', b'.hidden', b'.protected'}:
          version = last[0]
        name = last[1]
      else:
        name = last[0]

      # Add the imported and exported symbols.
      if is_undef:
        if not is_weak:
          imported[name].add(version)
      elif not is_local:
        exported[name].add(version)

    return ELF(dt_soname, dt_needed, imported, exported, header)


def _check_dt_needed(filename, main_file, shared_libs):
  """Check whether all DT_NEEDED entries are specified in the build
  system."""

  missing_shared_libs = False

  # Collect the DT_SONAMEs from shared libs specified in the build system.
  specified_sonames = {lib.dt_soname for lib in shared_libs}

  # Chech whether all DT_NEEDED entries are specified.
  for lib in main_file.dt_needed:
    if lib not in specified_sonames:
      print(filename + ':',
            'error: DT_NEEDED "{}" is not specified in shared_libs.'
            .format(lib.decode('utf-8')), file=sys.stderr)
      missing_shared_libs = True

  if missing_shared_libs:
    dt_needed = sorted(set(main_file.dt_needed))
    print(filename + ':',
          'note: Fix suggestion: LOCAL_SHARED_LIBRARIES := ' + ' '.join(
            re.sub('\\.so$', '', lib) for lib in dt_needed),
          file=sys.stderr)
    sys.exit(2)


def _find_symbol(lib, name, version):
  """Check whether the symbol name and version matches a definition."""
  try:
    lib_sym_vers = lib.exported[name]
  except KeyError:
    return False
  if version is None:
    return True
  if version == 'Base' and None in lib_sym_vers:
    # User is versioned but the dependency is unversioned.
    return True
  return version in lib_sym_vers


def _check_symbols(filename, main_file, shared_libs):
  """Check whether all undefined symbols are resolved to a definition."""

  all_elf_files = [main_file] + shared_libs
  missing_symbols = []
  for sym, imported_vers in main_file.imported.iteritems():
    for imported_ver in imported_vers:
      found = False
      for lib in all_elf_files:
        if _find_symbol(lib, sym, imported_ver):
          found = True
          break
      if not found:
        missing_symbols.append((sym, imported_ver))

  if missing_symbols:
    for sym, ver in sorted(missing_symbols):
      print(filename + ':',
            'error: Unresolved symbol: {} @ {}'.format(
              sym.decode('utf-8'),
              '(unversioned)' if ver is None else ver.decode('utf-8')),
            file=sys.stderr)
    sys.exit(2)


def _parse_args():
  """Parse command line options."""
  parser = argparse.ArgumentParser()
  parser.add_argument('filename')
  parser.add_argument('--shared-lib', action='append', default=[])
  parser.add_argument('--soname')
  parser.add_argument('--skip-bad-elf-magic', action='store_true')
  parser.add_argument('--skip-unknown-elf-machine', action='store_true')
  parser.add_argument('--allow-undefined-symbols', action='store_true')
  return parser.parse_args()


def main():
  """Main function"""
  args = _parse_args()

  known_archs = {'arm', 'aarch64', 'i386:x86-64', 'i386'}

  # Load the main ELF file (either an executable or a shared lib).
  try:
    main_file = ELFParser.open(args.filename)
    if args.skip_unknown_elf_machine and \
        main_file.header.e_machine not in known_archs:
      return
  except ELFInvalidMagicError:
    if not args.skip_bad_elf_magic:
      print(args.filename + ':',
            'error: Failed to read "{}"'.format(args.filename),
            file=sys.stderr)
      sys.exit(2)
    return

  # Load shared libraries.
  shared_libs = []
  for name in args.shared_lib:
    try:
      shared_libs.append(ELFParser.open(name))
    except Exception:
      print(args.filename + ':',
            'error: Failed to read shared lib "{}"'.format(name),
            file=sys.stderr)
      sys.exit(2)

  # Run checks.
  _check_dt_needed(args.filename, main_file, shared_libs)

  if not args.allow_undefined_symbols:
    _check_symbols(args.filename, main_file, shared_libs)


if __name__ == '__main__':
  main()
