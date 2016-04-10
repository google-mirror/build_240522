#!/usr/bin/env python

import argparse
import ConfigParser
import os
import re
import sys

class generator(object):

    _generators = {}

    def __init__(self, gen):
        self._gen = gen

        if gen in generator._generators:
            raise Exception('Duplicate generator name' + gen)

        generator._generators[gen] = None

    def __call__(self, cls):
        generator._generators[self._gen] = cls
        return cls

    @staticmethod
    def get():
        return generator._generators

    @staticmethod
    def set(gens):
        generator._generators = gens

class FSConfigFileParser():

    # from system/core/include/private/android_filesystem_config.h
    _AID_OEM_RESERVED_RANGES = [
        (2900, 2999),
        (5000, 5999),
    ]

    def __init__(self, config_files, base=None):

        self._files = []
        self._dirs = []
        self._aids = []

        self._seen_paths = {}
        # (name to file, value to aid)
        self._seen_aids = ({}, {})

        self._config_files = config_files
        self._base = base

        if base:
            self._parse(base, is_base=True)

        for f in self._config_files:
            self._parse(f)

    def _parse(self, file_name, is_base=False):

            # Separate config parsers for each file found. If you use read(filenames...) later
            # files can override earlier files which is not what we want. Track state across
            # files and enforce with handle_dup(). Note, strict ConfigParser is set to true in
            # Python >= 3.2, so in previous versions same file sections can override previous
            # sections.
            config = ConfigParser.ConfigParser()
            config.read(file_name)

            for s in config.sections():

                if config.has_option(s, 'value'):
                    FSConfigFileParser._handle_dup('AID', file_name, s, self._seen_aids[0])
                    self._seen_aids[0][s] = file_name
                    self._handle_aid(file_name, s, config, is_base)
                else:
                    FSConfigFileParser._handle_dup('path', file_name, s, self._seen_paths)
                    self._seen_paths[s] = file_name
                    self._handle_path(file_name, s, config)

                # sort entries:
                # * specified path before prefix match
                # ** ie foo before f*
                # * lexicographical less than before other
                # ** ie boo before foo
                # Given these paths:
                # paths=['ac', 'a', 'acd', 'an', 'a*', 'aa', 'ac*']
                # The sort order would be:
                # paths=['a', 'aa', 'ac', 'acd', 'an', 'ac*', 'a*']
                # Thus the fs_config tools will match on specified paths before attempting
                # prefix, and match on the longest matching prefix.
                self._files.sort(key= lambda x: FSConfigFileParser._file_key(x[1]))

                # sort on value of (file_name, name, value, strvalue)
                # This is only cosmetic so AIDS are arranged in ascending order
                # within the generated file.
                self._aids.sort(key=lambda x: x[2])

    def _handle_aid(self, file_name, section_name, config, is_base):
        value = config.get(section_name, 'value')
        section_name = section_name.lower()

        errmsg = '%s for: \"' + section_name + '" file: \"' + file_name + '\"'

        if not value:
            raise Exception(errmsg % 'Found specified but unset "value"')

        v = FSConfigFileParser.convert_int(value)
        if v == None:
            raise Exception(errmsg % ('Invalid "value", not a number, got: \"%s\"' % value))

        # Values must be within OEM range in OEM supplied files
        if not is_base:
            if not any(lower <= v <= upper for (lower, upper) in FSConfigFileParser._AID_OEM_RESERVED_RANGES):
                s = '"value" not in valid range %s, got: %s'
                s = s % (str(FSConfigFileParser._AID_OEM_RESERVED_RANGES), value)
                raise Exception(errmsg % s)

        # use the normalized int value in the dict and detect
        # duplicate definitions of the same value
        v = str(v)
        if v in self._seen_aids[1]:
            # map of value to aid name
            a = self._seen_aids[1][v]

            # aid name to file
            f = self._seen_aids[0][a]

            s = 'Duplicate AID value "%s" found on AID: "%s".' % (value, self._seen_aids[1][v])
            s += ' Previous found in file: "%s."' % f
            raise Exception(errmsg % s)

        self._seen_aids[1][v] = section_name

        # Append a tuple of (AID_*, base10(value), str(value))
        # We keep the str version of value so we can print that out in the
        # generated header so investigating parties can identify parts.
        # We store the base10 value for sorting, so everything is ascending
        # later.
        self._aids.append((file_name, section_name, v, value))

    def _handle_path(self, file_name, section_name, config):

                mode = config.get(section_name, 'mode')
                user = config.get(section_name, 'user')
                group = config.get(section_name, 'group')
                caps = config.get(section_name, 'caps')

                errmsg = 'Found specified but unset option: \"%s" in file: \"' + file_name + '\"'

                if not mode:
                    raise Exception(errmsg % 'mode')

                if not user:
                    raise Exception(errmsg % 'user')

                if not group:
                    raise Exception(errmsg % 'group')

                if not caps:
                    raise Exception(errmsg % 'caps')

                caps = caps.split()

                tmp = []
                for x in caps:
                    if FSConfigFileParser.convert_int(x) != None:
                        tmp.append('(' + x + ')')
                    else:
                        tmp.append('(1ULL << CAP_' + x.upper() + ')')

                caps = tmp

                path = '"' + section_name + '"'

                if len(mode) == 3:
                    mode = '0' + mode

                try:
                    int(mode, 8)
                except:
                    raise Exception('Mode must be octal characters, got: "' + mode + '"')

                if len(mode) != 4:
                    raise Exception('Mode must be 3 or 4 characters, got: "' + mode + '"')


                caps = '|'.join(caps)

                x = [ mode, user, group, caps, section_name ]
                if section_name[-1] == '/':
                    self._dirs.append((file_name, x))
                else:
                    self._files.append((file_name, x))

    def get_files(self):
        return self._files

    def get_dirs(self):
        return self._dirs

    def get_aids(self):
        return self._aids

    @staticmethod
    def _file_key(x):

        # Wrapper class for custom prefix matching strings
        class S(object):
            def __init__(self, str):

                self.orig = str
                self.is_prefix = str[-1] == '*'
                if self.is_prefix:
                    self.str = str[:-1]
                else:
                    self.str = str

            def __lt__(self, other):

                # if were both suffixed the smallest string
                # is 'bigger'
                if self.is_prefix and other.is_prefix:
                    b = len(self.str) > len(other.str)
                # If I am an the suffix match, im bigger
                elif self.is_prefix:
                    b = False
                # If other is the suffix match, he's bigger
                elif other.is_prefix:
                    b = True
                # Alphabetical
                else:
                    b = self.str < other.str
                return b

        return S(x[4])

    @staticmethod
    def _handle_dup(name, file_name, section_name, seen):
            if section_name in seen:
                dups = '"' + seen[section_name] + '" and '
                dups += file_name
                raise Exception('Duplicate ' + name + ' "' + section_name + '" found in files: ' + dups)

    @staticmethod
    def convert_int(num):

            try:
                if num.startswith('0x'):
                    return int(num, 16)
                elif num.startswith('0b'):
                    return int(num, 2)
                elif num.startswith('0'):
                    return int(num, 8)
                else:
                    return int(num, 10)
            except ValueError:
                pass
            return None

@generator("fsconfig")
class FSConfigGen(object):
    '''
    Generates the android_filesystem_config.h file to be used in generating
    fs_config_files and fs_config_dirs.
    '''

    GENERATED = '''
/*
 * THIS IS AN AUTOGENERATED FILE! DO NOT MODIFY
 */
 '''
    INCLUDE = '#include <private/android_filesystem_config.h>'

    DEFINE_NO_DIRS = '#define NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS\n'
    DEFINE_NO_FILES = '#define NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_FILES\n'

    DEFAULT_WARNING = '#warning No device-supplied android_filesystem_config.h, using empty default.'

    NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS_ENTRY = '{ 00000, AID_ROOT,      AID_ROOT,      0, "system/etc/fs_config_dirs" },'
    NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_FILES_ENTRY = '{ 00000, AID_ROOT,      AID_ROOT,      0, "system/etc/fs_config_files" },'

    IFDEF_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS = '#ifdef NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS'
    ENDIF = '#endif'

    OPEN_FILE_STRUCT = 'static const struct fs_path_config android_device_files[] = {'
    OPEN_DIR_STRUCT = 'static const struct fs_path_config android_device_dirs[] = {'
    CLOSE_FILE_STRUCT = '};'

    AID_DEFINE = "#define AID_%s\t%s"

    FILE_COMMENT = '// Defined in file: \"%s\"'

    @staticmethod
    def _generate(files, dirs, aids):
        print FSConfigGen.GENERATED
        print FSConfigGen.INCLUDE
        print

        are_dirs = len(dirs) > 0
        are_files = len(files) > 0
        are_aids = len(aids) > 0

        if are_aids:
            for a in aids:
                # use the preserved str value
                print FSConfigGen.FILE_COMMENT % a[0]
                print FSConfigGen.AID_DEFINE % (a[1].upper(), a[2])

            print

        if not are_dirs:
            print FSConfigGen.DEFINE_NO_DIRS

        if not are_files:
            print FSConfigGen.DEFINE_NO_FILES

        if not are_files and not are_dirs and not are_aids:
            print FSConfigGen.DEFAULT_WARNING
            return

        if are_files:
            print FSConfigGen.OPEN_FILE_STRUCT
            for tup in files:
                f = tup[0]
                c = tup[1]
                # Convert any uses of friendly name uid/gids, ie system, instead
                # of AID_SYSTEM or 1000
                c[1] = FSConfigGen._convert_friendly_name(c[1])
                c[2] = FSConfigGen._convert_friendly_name(c[2])

                c[4] = '"' + c[4] + '"'
                c = '{ ' + '    ,'.join(c) + ' },'
                print FSConfigGen.FILE_COMMENT % f
                print '    ' + c

            if not are_dirs:
                print FSConfigGen.IFDEF_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS
                print '    ' + FSConfigGen.NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS_ENTRY
                print FSConfigGen.ENDIF
            print FSConfigGen.CLOSE_FILE_STRUCT

        if are_dirs:
            print FSConfigGen.OPEN_DIR_STRUCT
            for d in dirs:
                f[4] = '"' + f[4] + '"'
                d = '{ ' + '    ,'.join(d) + ' },'
                print '    ' + d

            print FSConfigGen.CLOSE_FILE_STRUCT

    @staticmethod
    def _convert_friendly_name(name):

        # AID_* just use as is
        if name.startswith('AID'):
            return name

        # numbers, just us as is
        if FSConfigFileParser.convert_int(name) != None:
            return name

        # else use as AID_<name>.upper()
        return 'AID_' + name.upper()

    def __call__(self, files, dirs, aids):
        self._generate(files, dirs, aids)

@generator("passwd")
class PasswdGen(object):
    '''
    Generates a colon delimited passwd file in the format:
      -login name
      - encrypted password
      - numerical uid
      - numerical gid
      - user name
      - home dir
      - interpreter
    Note: Some fields may be blank.
    '''

    def __call__(self, files, dirs, aids):

        for a in aids:
            f = a[0]
            n = a[1]
            # use the normalized int value in passwd files.
            v = a[2]
            # name:passwd:uid:gid:user_name:homedir:intepreter
            print '# source: "' + f + '"'
            print n + '::' + v + ':' + v + ':' + n + ':/:/system/bin/sh'

def main():

    opt_parser = argparse.ArgumentParser(description='A tool for parsing fsconfig config files and producing digestable outputs.')
    opt_parser.add_argument('fsconfig', nargs='*', help='The list of fsconfig files to parse')
    opt_parser.add_argument('-b', '--base', help='Base Android fsconfig file', dest='base')

    generators = generator.get()
    tmp = {}
    opts = []

    opt_group = opt_parser.add_mutually_exclusive_group(required=True)

    # for each generator, instantiate and add them as an option
    for n, g in generators.iteritems():
        o = '--' + n
        opts.append(o)
        opt_group.add_argument(o, help=g.__doc__, action='store_const', const=n, dest='which')
        # Instantiate and save
        tmp[n] = g()

    # reassign constructed generators
    generator.set(tmp)

    args = opt_parser.parse_args()

    files = []
    dirs = []
    aids = []
    seen_paths = {}

    d = vars(args)
    which = d['which']

    parser = FSConfigFileParser(d['fsconfig'], base=d['base'])

    generator.get()[which](parser.get_files(), parser.get_dirs(), parser.get_aids())

if __name__ == '__main__':
    main()
