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


_ELF_MAGIC = b'\x7fELF'


# Known machines
_EM_386 = 3
_EM_ARM = 40
_EM_X86_64 = 62
_EM_AARCH64 = 183

_KNOWN_MACHINES = {_EM_386, _EM_ARM, _EM_X86_64, _EM_AARCH64}


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


ELF = collections.namedtuple(
  'ELF',
  ('dt_soname', 'dt_needed', 'imported', 'exported', 'header'))


def _get_os_name():
  """Get the host OS name."""
  if sys.platform == 'linux2':
    return 'linux'
  if sys.platform == 'darwin':
    return 'darwin'
  raise ValueError(sys.platform + ' is not supported')


def _get_build_top():
  """Find the build top of the source tree ($ANDROID_BUILD_TOP)."""
  prev_path = None
  curr_path = os.path.abspath(os.getcwd())
  while prev_path != curr_path:
    if os.path.exists(os.path.join(curr_path, '.repo')):
      return curr_path
    prev_path = curr_path
    curr_path = os.path.dirname(curr_path)
  return None


def _select_latest_llvm_version(versions):
  """Select the latest LLVM prebuilts version from a set of versions."""
  pattern = re.compile('clang-r([0-9]+)([a-z]?)')
  found_rev = 0
  found_ver = None
  for curr_ver in versions:
    match = pattern.match(curr_ver)
    if not match:
      continue
    curr_rev = int(match.group(1))
    if not found_ver or curr_rev > found_rev or (
        curr_rev == found_rev and curr_ver > found_ver):
      found_rev = curr_rev
      found_ver = curr_ver
  return found_ver


def _get_latest_llvm_version(llvm_dir):
  """Find the latest LLVM prebuilts version from `llvm_dir`."""
  return _select_latest_llvm_version(os.listdir(llvm_dir))


def _get_llvm_dir():
  """Find the path to LLVM prebuilts."""
  build_top = _get_build_top()

  llvm_prebuilts_base = os.environ.get('LLVM_PREBUILTS_BASE')
  if not llvm_prebuilts_base:
    llvm_prebuilts_base = os.path.join('prebuilts', 'clang', 'host')

  llvm_dir = os.path.join(
    build_top, llvm_prebuilts_base, _get_os_name() + '-x86')

  if not os.path.exists(llvm_dir):
    return None

  llvm_prebuilts_version = os.environ.get('LLVM_PREBUILTS_VERSION')
  if not llvm_prebuilts_version:
    llvm_prebuilts_version = _get_latest_llvm_version(llvm_dir)

  llvm_dir = os.path.join(llvm_dir, llvm_prebuilts_version)

  if not os.path.exists(llvm_dir):
    return None

  return llvm_dir


def _get_llvm_readobj():
  """Find the path to llvm-readobj executable."""
  llvm_dir = _get_llvm_dir()
  llvm_readobj = os.path.join(llvm_dir, 'bin', 'llvm-readobj')
  return llvm_readobj if os.path.exists(llvm_readobj) else 'llvm-readobj'


class ELFError(ValueError):
  """Generic ELF parse error."""
  pass


class ELFInvalidMagicError(ELFError):
  """Invalid ELF magic word error."""
  def __init__(self):
    super(ELFInvalidMagicError, self).__init__('bad ELF magic')


class ELFParser(object):
  """A parser that parses an ELF file."""

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
  def open(cls, elf_file_path, llvm_readobj):
    """Open and parse the ELF file."""
    # Parse the ELF header for simple sanity checks.
    header = cls._read_elf_header(elf_file_path)
    if not header or header.ei_magic != _ELF_MAGIC:
      raise ELFInvalidMagicError()

    # Run llvm-readobj and parse the output.
    return cls._read_llvm_readobj(elf_file_path, header, llvm_readobj)


  @classmethod
  def _find_prefix(cls, pattern, lines_it):
    """Iterate `lines_it` until finding a string that starts with `pattern`."""
    for line in lines_it:
      if line.startswith(pattern):
        return True
    return False


  @classmethod
  def _read_llvm_readobj(cls, elf_file_path, header, llvm_readobj):
    """Run llvm-readobj and parse the output."""
    proc = subprocess.Popen(
      [llvm_readobj, '-dynamic-table', '-dyn-symbols', elf_file_path],
      stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, _ = proc.communicate()
    lines = out.splitlines()
    return cls._parse_llvm_readobj(elf_file_path, header, lines)


  @classmethod
  def _parse_llvm_readobj(cls, elf_file_path, header, lines):
    """Parse the output of llvm-readobj."""
    lines_it = iter(lines)
    imported, exported = cls._parse_dynamic_symbols(lines_it)
    dt_soname, dt_needed = cls._parse_dynamic_table(elf_file_path, lines_it)
    return ELF(dt_soname, dt_needed, imported, exported, header)


  _DYNAMIC_SECTION_START_PATTERN = 'DynamicSection ['

  _DYNAMIC_SECTION_NEEDED_PATTERN = re.compile(
    '^  0x[0-9a-fA-F]+\\s+NEEDED\\s+Shared library: \\[(.*)\\]$')

  _DYNAMIC_SECTION_SONAME_PATTERN = re.compile(
    '^  0x[0-9a-fA-F]+\\s+SONAME\\s+Library soname: \\[(.*)\\]$')

  _DYNAMIC_SECTION_END_PATTERN = ']'


  @classmethod
  def _parse_dynamic_table(cls, elf_file_path, lines_it):
    """Parse the dynamic table section."""
    dt_soname = os.path.basename(elf_file_path)
    dt_needed = []

    dynamic = cls._find_prefix(cls._DYNAMIC_SECTION_START_PATTERN, lines_it)
    if not dynamic:
      return (dt_soname, dt_needed)

    for line in lines_it:
      if line == cls._DYNAMIC_SECTION_END_PATTERN:
        break

      match = cls._DYNAMIC_SECTION_NEEDED_PATTERN.match(line)
      if match:
        dt_needed.append(match.group(1))
        continue

      match = cls._DYNAMIC_SECTION_SONAME_PATTERN.match(line)
      if match:
        dt_soname = match.group(1)
        continue

    return (dt_soname, dt_needed)


  _DYNAMIC_SYMBOLS_START_PATTERN = 'DynamicSymbols ['
  _DYNAMIC_SYMBOLS_END_PATTERN = ']'

  _SYMBOL_ENTRY_START_PATTERN = '  Symbol {'
  _SYMBOL_ENTRY_PATTERN = re.compile('^    ([A-Za-z0-9_]+): (.*)$')
  _SYMBOL_ENTRY_PAREN_PATTERN = re.compile(
    '\\s+\\((?:(?:\\d+)|(?:0x[0-9a-fA-F]+))\\)$')
  _SYMBOL_ENTRY_END_PATTERN = '  }'


  @classmethod
  def _parse_symbol_name(cls, name_with_version):
    """Split `name_with_version` into name and version. This function may split
    at last occurrence of `@@` or `@`."""
    name, version = name_with_version.rsplit('@', 1)
    if name and name[-1] == '@':
      name = name[:-1]
    return (name, version)


  @classmethod
  def _parse_dynamic_symbols(cls, lines_it):
    """Parse dynamic symbol table and collect imported and exported symbols."""
    imported = collections.defaultdict(set)
    exported = collections.defaultdict(set)

    for symbol in cls._parse_dynamic_symbols_internal(lines_it):
      name, version = cls._parse_symbol_name(symbol['Name'])
      if name:
        if symbol['Section'] == 'Undefined':
          if symbol['Binding'] != 'Weak':
            imported[name].add(version)
        else:
          if symbol['Binding'] != 'Local':
            exported[name].add(version)

    return (imported, exported)


  @classmethod
  def _parse_dynamic_symbols_internal(cls, lines_it):
    """Parse symbols entries and yield each symbols."""

    if not cls._find_prefix(cls._DYNAMIC_SYMBOLS_START_PATTERN, lines_it):
      return

    for line in lines_it:
      if line == cls._DYNAMIC_SYMBOLS_END_PATTERN:
        return

      if line == cls._SYMBOL_ENTRY_START_PATTERN:
        symbol = {}
        continue

      if line == cls._SYMBOL_ENTRY_END_PATTERN:
        yield symbol
        symbol = None
        continue

      match = cls._SYMBOL_ENTRY_PATTERN.match(line)
      if match:
        key = match.group(1)
        value = cls._SYMBOL_ENTRY_PAREN_PATTERN.sub('', match.group(2))
        symbol[key] = value
        continue


def _check_dt_needed(filename, file_under_test, shared_libs):
  """Check whether all DT_NEEDED entries are specified in the build
  system."""

  missing_shared_libs = False

  # Collect the DT_SONAMEs from shared libs specified in the build system.
  specified_sonames = {lib.dt_soname for lib in shared_libs}

  # Chech whether all DT_NEEDED entries are specified.
  for lib in file_under_test.dt_needed:
    if lib not in specified_sonames:
      print(filename + ':',
            'error: DT_NEEDED "{}" is not specified in shared_libs.'
            .format(lib.decode('utf-8')), file=sys.stderr)
      missing_shared_libs = True

  if missing_shared_libs:
    dt_needed = sorted(set(file_under_test.dt_needed))
    print(filename + ':',
          'note: Fix suggestion: LOCAL_SHARED_LIBRARIES := ' + ' '.join(
            re.sub('\\.so$', '', lib) for lib in dt_needed),
          file=sys.stderr)
    sys.exit(2)


def _find_symbol(lib, name, version):
  """Check whether the symbol name and version matches a definition in lib."""
  try:
    lib_sym_vers = lib.exported[name]
  except KeyError:
    return False
  if version == '':  # Symbol version is not requested
    return True
  return version in lib_sym_vers


def _find_symbol_from_libs(libs, name, version):
  """Check whether the symbol name and version is defined in one of the shared
  libraries in libs."""
  for lib in libs:
    if _find_symbol(lib, name, version):
      return lib
  return None


def _check_symbols(filename, file_under_test, shared_libs):
  """Check whether all undefined symbols are resolved to a definition."""
  all_elf_files = [file_under_test] + shared_libs
  missing_symbols = []
  for sym, imported_vers in file_under_test.imported.iteritems():
    for imported_ver in imported_vers:
      lib = _find_symbol_from_libs(all_elf_files, sym, imported_ver)
      if not lib:
        missing_symbols.append((sym, imported_ver))

  if missing_symbols:
    for sym, ver in sorted(missing_symbols):
      print(filename + ':',
            'error: Unresolved symbol: {} @ {}'.format(
              sym.decode('utf-8'), ver.decode('utf-8')),
            file=sys.stderr)
    sys.exit(2)


def _parse_args():
  """Parse command line options."""
  parser = argparse.ArgumentParser()

  # Input file
  parser.add_argument('filename',
                      help='Path to the input file to be checked')
  parser.add_argument('--soname',
                      help='Shared object name of the input file')

  # Shared library dependencies
  parser.add_argument('--shared-lib', action='append', default=[],
                      help='Path to shared library dependencies')

  # Check options
  parser.add_argument('--skip-bad-elf-magic', action='store_true',
                      help='Ignore the input file without the ELF magic word')
  parser.add_argument('--skip-unknown-elf-machine', action='store_true',
                      help='Ignore the input file with unknown machine ID')
  parser.add_argument('--allow-undefined-symbols', action='store_true',
                      help='Ignore unresolved undefined symbols')

  # Other options
  parser.add_argument('--llvm-readobj',
                      help='Path to the llvm-readobj executable')

  return parser.parse_args()


def main():
  """Main function"""
  args = _parse_args()

  llvm_readobj = args.llvm_readobj
  if not llvm_readobj:
    llvm_readobj = _get_llvm_readobj()

  # Load the ELF file under test (either an executable or a shared lib).
  try:
    file_under_test = ELFParser.open(args.filename, llvm_readobj)
    if args.skip_unknown_elf_machine and \
        file_under_test.header.e_machine not in _KNOWN_MACHINES:
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
      shared_libs.append(ELFParser.open(name, llvm_readobj))
    except Exception:
      print(args.filename + ':',
            'error: Failed to read shared lib "{}"'.format(name),
            file=sys.stderr)
      sys.exit(2)

  # Run checks.
  _check_dt_needed(args.filename, file_under_test, shared_libs)

  if not args.allow_undefined_symbols:
    _check_symbols(args.filename, file_under_test, shared_libs)


if __name__ == '__main__':
  main()
