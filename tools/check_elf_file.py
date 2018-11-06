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
import os.path
import re
import struct
import sys
import traceback


ELF_MAGIC = b'\x7fELF'

EI_CLASS = 4
EI_DATA  = 5

ELFCLASSNONE = 0
ELFCLASS32   = 1
ELFCLASS64   = 2

ELFDATANONE = 0
ELFDATA2LSB = 1
ELFDATA2MSB = 2

PT_NULL         = 0
PT_LOAD         = 1
PT_DYNAMIC      = 2
PT_INTERP       = 3
PT_NOTE         = 4
PT_SHLIB        = 5
PT_PHDR         = 6

DT_NULL            = 0
DT_NEEDED          = 1
DT_PLTRELSZ        = 2
DT_PLTGOT          = 3
DT_HASH            = 4
DT_STRTAB          = 5
DT_SYMTAB          = 6
DT_RELA            = 7
DT_RELASZ          = 8
DT_RELAENT         = 9
DT_STRSZ           = 10
DT_SYMENT          = 11
DT_INIT            = 12
DT_FINI            = 13
DT_SONAME          = 14
DT_RPATH           = 15
DT_SYMBOLIC        = 16
DT_REL             = 17
DT_RELSZ           = 18
DT_RELENT          = 19
DT_PLTREL          = 20
DT_DEBUG           = 21
DT_TEXTREL         = 22
DT_JMPREL          = 23
DT_BIND_NOW        = 24
DT_INIT_ARRAY      = 25
DT_FINI_ARRAY      = 26
DT_INIT_ARRAYSZ    = 27
DT_FINI_ARRAYSZ    = 28
DT_RUNPATH         = 29
DT_FLAGS           = 30
DT_PREINIT_ARRAY   = 32
DT_PREINIT_ARRAYSZ = 33
DT_SYMTAB_SHNDX    = 34
DT_NUM             = 35
DT_GNU_HASH        = 0x6ffffef5
DT_VERSYM          = 0x6ffffff0
DT_RELCOUNT        = 0x6ffffff9
DT_VERDEF          = 0x6ffffffc
DT_VERDEFNUM       = 0x6ffffffd
DT_VERNEED         = 0x6ffffffe
DT_VERNEEDNUM      = 0x6fffffff

EM_NONE    = 0
EM_386     = 3
EM_MIPS    = 8
EM_ARM     = 40
EM_X86_64  = 62
EM_AARCH64 = 183

STB_LOCAL  = 0
STB_GLOBAL = 1
STB_WEAK   = 2

SHN_UNDEF  = 0
SHN_ABS    = 0xfff1
SHN_COMMON = 0xfff2


def struct_class(name, fields):
    """Create a namedtuple with unpack_from() function.

    >>> Point = struct_class('Point', [('x', 'H'), ('y', 'H')])
    >>> pt = Point.unpack_from(b'\\x00\\x00\\x01\\x00', 0, '<')
    >>> pt.x
    0
    >>> pt.y
    1
    """
    field_names = [name for name, ty in fields]
    cls = collections.namedtuple(name, field_names)
    cls.struct_fields = fields
    cls.struct_fmt = ''.join(ty for name, ty in fields)
    cls.struct_size = struct.calcsize(cls.struct_fmt)
    def unpack_from(cls, buf, offset=0, endianness='='):
        assert endianness in '=<>'
        unpacked = struct.unpack_from(endianness + cls.struct_fmt, buf, offset)
        return cls.__new__(cls, *unpacked)
    cls.unpack_from = classmethod(unpack_from)
    return cls


Elf_Ident = struct_class('Elf_Ident', (
    ('ei_magic', '4s'),
    ('ei_class', 'B'),
    ('ei_data', 'B'),
    ('ei_version', 'B'),
    ('ei_osabi', 'B'),
    ('ei_pad', '8s'),
))


Elf32_Hdr = struct_class('Elf32_Hdr', Elf_Ident.struct_fields + (
    ('e_type', 'H'),
    ('e_machine', 'H'),
    ('e_version', 'I'),
    ('e_entry', 'I'),  # addr
    ('e_phoff', 'I'),  # offset
    ('e_shoff', 'I'),  # offset
    ('e_flags', 'I'),
    ('e_ehsize', 'H'),
    ('e_phentsize', 'H'),
    ('e_phnum', 'H'),
    ('e_shentsize', 'H'),
    ('e_shnum', 'H'),
    ('e_shstrndx', 'H'),
))


Elf64_Hdr = struct_class('Elf64_Hdr', Elf_Ident.struct_fields + (
    ('e_type', 'H'),
    ('e_machine', 'H'),
    ('e_version', 'I'),
    ('e_entry', 'Q'),  # addr
    ('e_phoff', 'Q'),  # offset
    ('e_shoff', 'Q'),  # offset
    ('e_flags', 'I'),
    ('e_ehsize', 'H'),
    ('e_phentsize', 'H'),
    ('e_phnum', 'H'),
    ('e_shentsize', 'H'),
    ('e_shnum', 'H'),
    ('e_shstrndx', 'H'),
))


Elf32_Phdr = struct_class('Elf32_Phdr', (
    ('p_type', 'I'),
    ('p_offset', 'I'),  # offset
    ('p_vaddr', 'I'),  # addr
    ('p_paddr', 'I'),  # addr
    ('p_filesz', 'I'),
    ('p_memsz', 'I'),
    ('p_flags', 'I'),
    ('p_align', 'I'),
))


Elf64_Phdr = struct_class('Elf64_Phdr', (
    ('p_type', 'I'),
    ('p_flags', 'I'),
    ('p_offset', 'Q'),  # offset
    ('p_vaddr', 'Q'),  # addr
    ('p_paddr', 'Q'),  # addr
    ('p_filesz', 'Q'),
    ('p_memsz', 'Q'),
    ('p_align', 'Q'),
))


Elf32_Dyn = struct_class('Elf32_Dyn', (
    ('d_tag', 'i'),
    ('d_val', 'I'),
))


Elf64_Dyn = struct_class('Elf64_Dyn', (
    ('d_tag', 'q'),
    ('d_val', 'Q'),
))


def _decorate_elf_rel_info(cls, symbol_shift_right, type_mask):
    def r_sym(self):
        return self.r_info >> symbol_shift_right

    def r_type(self):
        return self.r_info & type_mask

    cls.r_sym = property(r_sym)
    cls.r_type = property(r_type)
    return cls


def _decoreate_elf32_rel_info(cls):
    return _decorate_elf_rel_info(cls, 8, 0xff)


def _decoreate_elf64_rel_info(cls):
    return _decorate_elf_rel_info(cls, 32, 0xffffffff)


Elf32_Rel = _decoreate_elf32_rel_info(struct_class('Elf32_Rel', (
    ('r_offset', 'I'),
    ('r_info', 'I'),
)))


Elf64_Rel = _decoreate_elf64_rel_info(struct_class('Elf64_Rel', (
    ('r_offset', 'Q'),
    ('r_info', 'Q'),
)))


Elf32_Rela = _decoreate_elf32_rel_info(struct_class('Elf32_Rela', (
    ('r_offset', 'I'),
    ('r_info', 'I'),
    ('r_addend', 'i'),
)))


Elf64_Rela = _decoreate_elf64_rel_info(struct_class('Elf64_Rela', (
    ('r_offset', 'Q'),
    ('r_info', 'Q'),
    ('r_addend', 'q'),
)))


def _decorate_elf_sym(cls):
    def st_bind(self):
        return self.st_info >> 4

    def is_local(self):
        return self.st_bind == STB_LOCAL

    def is_global(self):
        return self.st_bind == STB_GLOBAL

    def is_weak(self):
        return self.st_bind == STB_WEAK

    def is_undef(self):
        return self.st_shndx == SHN_UNDEF

    cls.st_bind = property(st_bind)
    cls.is_local = property(is_local)
    cls.is_global = property(is_global)
    cls.is_weak = property(is_weak)
    cls.is_undef = property(is_undef)
    return cls


Elf32_Sym = _decorate_elf_sym(struct_class('Elf32_Sym', (
    ('st_name', 'I'),
    ('st_value', 'I'),  # addr
    ('st_size', 'I'),
    ('st_info', 'B'),
    ('st_other', 'B'),
    ('st_shndx', 'H'),
)))


Elf64_Sym = _decorate_elf_sym(struct_class('Elf64_Sym', (
    ('st_name', 'I'),
    ('st_info', 'B'),
    ('st_other', 'B'),
    ('st_shndx', 'H'),
    ('st_value', 'Q'),  # addr
    ('st_size', 'Q'),
)))


Elf_Verdef = struct_class('Elf_Verdef', (
    ('vd_version', 'H'),
    ('vd_flags', 'H'),
    ('vd_ndx', 'H'),
    ('vd_cnt', 'H'),
    ('vd_hash', 'I'),
    ('vd_aux', 'I'),
    ('vd_next', 'I'),
))


Elf_Verdaux = struct_class('Elf_Verdaux', (
    ('vda_name', 'I'),
    ('vda_next', 'I'),
))


Elf_Verneed = struct_class('Elf_Verneed', (
    ('vn_version', 'H'),
    ('vn_cnt', 'H'),
    ('vn_file', 'I'),
    ('vn_aux', 'I'),
    ('vn_next', 'I'),
))


Elf_Vernaux = struct_class('Elf_Vernaux', (
    ('vna_hash', 'I'),
    ('vna_flags', 'H'),
    ('vna_other', 'H'),
    ('vna_name', 'I'),
    ('vna_next', 'I'),
))


class Elf32(object):
    WORD_SIZE_IN_BYTES = 4

    Hdr = Elf32_Hdr
    Phdr = Elf32_Phdr
    Dyn = Elf32_Dyn
    Rel = Elf32_Rel
    Rela = Elf32_Rela
    Sym = Elf32_Sym


class Elf64(object):
    WORD_SIZE_IN_BYTES = 8

    Hdr = Elf64_Hdr
    Phdr = Elf64_Phdr
    Dyn = Elf64_Dyn
    Rel = Elf64_Rel
    Rela = Elf64_Rela
    Sym = Elf64_Sym


ParsedELF = collections.namedtuple(
    'ParsedELF',
    ('dt_soname', 'dt_needed', 'imported', 'exported', 'header',
     'first_version_name'))


class ELFError(ValueError):
    pass


class ELFInvalidMagicError(ELFError):
    def __init__(self):
        super(ELFInvalidMagicError, self).__init__('bad ELF magic')


class ELF(object):
    @classmethod
    def open(cls, elf_file_path):
        with open(elf_file_path, 'rb') as elf_fp:
            return cls._open_file(elf_fp, os.path.basename(elf_file_path))


    @classmethod
    def _open_file(cls, elf_file, file_name):
        return cls.parse(elf_file.read(), file_name)


    @classmethod
    def parse(cls, buf, file_name):
        elf, endianness = cls._parse_ident(buf)

        header = elf.Hdr.unpack_from(buf, 0, endianness)

        phdrs = cls._parse_program_headers(
            buf, header.e_phoff, header.e_phentsize, header.e_phnum, elf,
            endianness)

        dyns = cls._parse_dynamics(buf, phdrs, elf, endianness)

        if not dyns:
            # Some ELF files (e.g. statically linked executables) do not have
            # the .dynamic section.
            return ParsedELF(file_name, [], {}, {}, header, '')

        # Parse .strtab
        dt_strtab = cls._offset_from_vma(
            phdrs, cls._get_dynamic(dyns, DT_STRTAB))

        # Parse .rel, .rela, .rel.plt, and .rela.plt
        sym_ent_size = cls._get_dynamic(dyns, DT_SYMENT)
        sym_idx_max = -1

        dt_rel = cls._parse_relocs(
            buf, phdrs, dyns, DT_REL, DT_RELENT, DT_RELSZ, elf.Rel)
        for rel in dt_rel:
            sym_idx_max = max(sym_idx_max, rel.r_sym)

        dt_rela = cls._parse_relocs(
            buf, phdrs, dyns, DT_RELA, DT_RELAENT, DT_RELASZ, elf.Rela)
        for rel in dt_rela:
            sym_idx_max = max(sym_idx_max, rel.r_sym)

        dt_pltrel = cls._get_dynamic(dyns, DT_PLTREL)
        if dt_pltrel is not None:
            if dt_pltrel == DT_REL:
                rel_ent_size_tag = DT_RELENT
                rel_struct = elf.Rel
            elif dt_pltrel == DT_RELA:
                rel_ent_size_tag = DT_RELAENT
                rel_struct = elf.Rela
            else:
                raise ELFError('bad DT_PTREL')

            dt_pltrel = cls._parse_relocs(
                buf, phdrs, dyns, DT_JMPREL, rel_ent_size_tag, DT_PLTRELSZ,
                rel_struct)

            for rel in dt_pltrel:
                sym_idx_max = max(sym_idx_max, rel.r_sym)

        # Parse .gnu.hash
        dt_gnu_hash = cls._get_dynamic(dyns, DT_GNU_HASH)
        if dt_gnu_hash is not None:
            dt_gnu_hash = cls._offset_from_vma(phdrs, dt_gnu_hash)

            num_buckets, sym_idx_base, bloom_size, bloom_shift = \
                struct.unpack_from('IIII', buf, dt_gnu_hash)

            bloom_offset = dt_gnu_hash + 16
            buckets_offset = bloom_offset + elf.WORD_SIZE_IN_BYTES * bloom_size
            chain_offset = buckets_offset + 4 * num_buckets

            # Traverse the buckets
            for i in range(num_buckets):
                sym = struct.unpack_from('I', buf, buckets_offset + i * 4)[0]
                if sym < sym_idx_base:
                    continue

                # Traverse the chain for this bucket
                while True:
                    next_chain = struct.unpack_from(
                        'I', buf, chain_offset + (sym - sym_idx_base) * 4)[0]
                    if next_chain & 1:
                        break
                    sym += 1

                sym_idx_max = max(sym_idx_max, sym)

        # Parse .hash
        dt_hash = cls._get_dynamic(dyns, DT_HASH)
        if dt_hash is not None:
            dt_hash = cls._offset_from_vma(phdrs, dt_hash)
            num_buckets, num_chains = struct.unpack_from('II', buf, dt_hash)
            if num_chains > 0:
                sym_idx_max = max(sym_idx_max, num_chains - 1)

        num_symbols = sym_idx_max + 1

        # Parse .gnu.version_d
        versions = {}
        first_version_name = None

        dt_verdefnum = cls._get_dynamic(dyns, DT_VERDEFNUM)
        if dt_verdefnum is not None:
            dt_verdef = cls._offset_from_vma(
                phdrs, cls._get_dynamic(dyns, DT_VERDEF))

            offset = dt_verdef
            for i in range(dt_verdefnum):
                verdef = Elf_Verdef.unpack_from(buf, offset)

                names = []
                aux = verdef.vd_aux
                aux_offset = offset + aux
                while aux != 0:
                    verdef_aux = Elf_Verdaux.unpack_from(buf, aux_offset)
                    name = cls._extract_c_str(
                        buf, dt_strtab + verdef_aux.vda_name)
                    if first_version_name is None:
                        first_version_name = name
                    names.append(name)
                    aux = verdef_aux.vda_next
                    aux_offset += aux

                versions[verdef.vd_ndx] = names[0]

                offset += verdef.vd_next

        # Parse .gnu.version_r
        dt_verneednum = cls._get_dynamic(dyns, DT_VERNEEDNUM)
        if dt_verneednum is not None:
            dt_verneed = cls._offset_from_vma(
                phdrs, cls._get_dynamic(dyns, DT_VERNEED))

            offset = dt_verneed
            for i in range(dt_verneednum):
                verneed = Elf_Verneed.unpack_from(buf, offset)

                need_file_soname = cls._extract_c_str(
                    buf, dt_strtab + verneed.vn_file)


                aux = verneed.vn_aux
                aux_offset = offset + aux
                while aux != 0:
                    verneed_aux = Elf_Vernaux.unpack_from(buf, aux_offset)
                    name = cls._extract_c_str(
                        buf, dt_strtab + verneed_aux.vna_name)
                    aux = verneed_aux.vna_next
                    aux_offset += aux

                    versions[verneed_aux.vna_other] = name

                offset += verneed.vn_next

        # Parse .dynsym and .gnu.version
        exported = collections.defaultdict(set)
        imported = collections.defaultdict(set)
        dt_symtab = cls._get_dynamic(dyns, DT_SYMTAB)
        dt_versym = cls._get_dynamic(dyns, DT_VERSYM)
        if dt_symtab is not None:
            dt_symtab = cls._offset_from_vma(phdrs, dt_symtab)
            if dt_versym is not None:
                dt_versym = cls._offset_from_vma(phdrs, dt_versym)
            for i in range(1, num_symbols):
                sym = elf.Sym.unpack_from(
                    buf, dt_symtab + elf.Sym.struct_size * i)
                sym_name = cls._extract_c_str(buf, dt_strtab + sym.st_name)

                if dt_versym is None:
                    sym_ver = None
                else:
                    ver = struct.unpack_from('H', buf, dt_versym + 2 * i)[0]
                    sym_ver = versions.get(ver)

                if sym.is_undef:
                    if not sym.is_weak:
                        imported[sym_name].add(sym_ver)
                elif not sym.is_local:
                    exported[sym_name].add(sym_ver)

        # Collect DT_SONAME
        dt_soname = cls._get_dynamic(dyns, DT_SONAME)
        if dt_soname is None:
            dt_soname = file_name
        else:
            dt_soname = cls._extract_c_str(buf, dt_strtab + dt_soname)

        if first_version_name is None:
            first_version_name = dt_soname

        # Collect DT_NEEDED
        dt_needed = []
        for dyn in dyns:
            if dyn.d_tag == DT_NEEDED:
                dt_needed.append(cls._extract_c_str(buf, dt_strtab + dyn.d_val))

        return ParsedELF(dt_soname, dt_needed, imported, exported, header,
                         first_version_name)


    @classmethod
    def _parse_ident(cls, buf):
        try:
            ident = Elf_Ident.unpack_from(buf, 0)
        except struct.error:
            raise ELFInvalidMagicError()

        # Check the ELF magic
        if ident.ei_magic != ELF_MAGIC:
            raise ELFInvalidMagicError()

        # ELF class and endianness
        elf = Elf32 if ident.ei_class == ELFCLASS32 else Elf64
        endianness = '<' if ident.ei_data == ELFDATA2LSB else '>'

        return (elf, endianness)


    @classmethod
    def _parse_program_headers(cls, buf, phoff, phentsize, phnum, elf,
                               endianness):
        phdrs = []
        for offset in range(phoff, phoff + phentsize * phnum, phentsize):
            phdrs.append(elf.Phdr.unpack_from(buf, offset, endianness))
        return phdrs


    @classmethod
    def _parse_dynamics(cls, buf, phdrs, elf, endianness):
        dyns = []
        for phdr in phdrs:
            if phdr.p_type != PT_DYNAMIC:
                continue
            begin = phdr.p_vaddr
            end = begin + phdr.p_memsz
            for vaddr in range(begin, end, elf.Dyn.struct_size):
                file_offset = cls._offset_from_vma(phdrs, vaddr)
                dyns.append(elf.Dyn.unpack_from(buf, file_offset, endianness))
        return dyns


    @classmethod
    def _get_dynamic(cls, dyns, tag):
        for dyn in dyns:
            if dyn.d_tag == tag:
                return dyn.d_val
        return None


    @classmethod
    def _parse_relocs(cls, buf, phdrs, dyns, rel_tag, rel_ent_size_tag,
                      rel_size_tag, rel_struct):
        rel = cls._get_dynamic(dyns, rel_tag)
        rel_ent_size = cls._get_dynamic(dyns, rel_ent_size_tag)
        rel_size = cls._get_dynamic(dyns, rel_size_tag)

        if rel is None and rel_size is None:
            return []

        if rel is None or rel_size is None:
            raise ELFError('bad rel dt_tags')

        if rel_ent_size is None:
            rel_ent_size = rel_struct.struct_size

        if rel_ent_size != rel_struct.struct_size:
            raise ELFError('bad {} size: {}'.format(
                rel_ent_size_tag, rel_ent_size))

        relocs = []
        rel_offset = cls._offset_from_vma(phdrs, rel)
        for offset in range(rel_offset, rel_offset + rel_size, rel_ent_size):
            relocs.append(rel_struct.unpack_from(buf, offset))
        return relocs


    @classmethod
    def _offset_from_vma(cls, phdrs, vaddr):
        for phdr in phdrs:
            if phdr.p_type != PT_LOAD:
                continue
            if (vaddr < phdr.p_vaddr or
                vaddr >= phdr.p_vaddr + phdr.p_memsz or
                vaddr >= phdr.p_vaddr + phdr.p_filesz):
                continue
            return vaddr - phdr.p_vaddr + phdr.p_offset
        raise IndexError('unexpected vaddr ' + str(vaddr))


    @classmethod
    def _extract_c_str(cls, buf, offset):
        pos = buf.find(b'\0', offset)
        return buf[offset:pos] if pos != -1 else buf[offset:]


def _check_dt_needed(filename, main_file, shared_libs):
    """Check whether all DT_NEEDED entries are specified in the build system."""

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
    return version is None or version in lib_sym_vers


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

    known_archs = {EM_386, EM_MIPS, EM_ARM, EM_X86_64, EM_AARCH64}

    # Load the main ELF file (either an executable or a shared lib).
    try:
        main_file = ELF.open(args.filename)
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
            shared_libs.append(ELF.open(name))
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
