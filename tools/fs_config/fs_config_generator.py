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

class AID():

    def __init__(self, name, value, value_norm, homedir, interpreter, file_name):
        self._name = name
        self._value = value
        self._value_norm = value_norm
        self._homedir = homedir
        self._interpreter = interpreter
        self._file_name = file_name

    def get_value(self):
        return self._value

    def get_normalized_value(self):
        return self._value_norm

    def get_name(self):
        return self._name

    def get_file_found(self):
        return self._file_name

    def get_homedir(self):
        return self._homedir

    def get_interpreter(self):
        return self._interpreter

class Utils():
    # from system/core/include/private/android_filesystem_config.h
    AID_OEM_RESERVED_RANGES = [
        (2900, 2999),
        (5000, 5999),
    ]

    # AID's cannot enter this range, tested for base
    AID_APP = 10000

    @staticmethod
    def is_oem_aid(aid):
        return any(lower <= aid <= upper for (lower, upper) in Utils.AID_OEM_RESERVED_RANGES)

    @staticmethod
    def is_app_aid(aid):
        return aid > Utils.AID_APP

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

class FSConfigFileParser():

    _AID_VALID = '^[a-z0-9_]*$'

    def __init__(self, config_files, base=None):

        self._files = []
        self._dirs = []
        self._aids = []
        self._defaults = {}

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

            if is_base:
                self._defaults = {k:v for k,v in config.items('DEFAULT')}

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
                self._aids.sort(key=lambda x: x.get_normalized_value())

    def _handle_aid(self, file_name, section_name, config, is_base):
        value = config.get(section_name, 'value')

        errmsg = '%s for: \"' + section_name + '" file: \"' + file_name + '\"'

        if not re.match(FSConfigFileParser._AID_VALID , section_name):
            raise Exception(errmsg % 'Found bad characters in "%s", expecting within set: "%s"' \
                % (section_name, FSConfigFileParser._AID_VALID))

        if not value:
            raise Exception(errmsg % 'Found specified but unset "value"')

        v = Utils.convert_int(value)
        if v == None:
            raise Exception(errmsg % ('Invalid "value", not a number, got: "%s"' % value))

        # Values must be within OEM range in OEM supplied files
        if not is_base:
            if not Utils.is_oem_aid(v):
                s = '"value" not in valid range %s, got: %s'
                s = s % (str(Utils.AID_OEM_RESERVED_RANGES), value)
                raise Exception(errmsg % s)
        # Base AIDs cannot enter into the app range
        else:
            if Utils.is_app_aid(v):
                raise Exception(errmsg % ('"value", entered into app AID range, got: "%s"' % value))

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

        homedir = '/'
        if config.has_option(section_name, 'homedir'):
            homedir = config.get(section_name, 'homedir')
        elif 'homedir' in self._defaults:
            homedir = self._defaults['homedir']

        interpreter = '/system/bin/sh'
        if config.has_option(section_name, 'interpreter'):
            interpreter = config.get(section_name, 'interpreter')
        elif 'homedir' in self._defaults:
            interpreter = self._defaults['interpreter']

        # Append a tuple of (AID_*, base10(value), str(value))
        # We keep the str version of value so we can print that out in the
        # generated header so investigating parties can identify parts.
        # We store the base10 value for sorting, so everything is ascending
        # later.
        self._aids.append(AID(section_name, value, v, homedir, interpreter, file_name))

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
                    if Utils.convert_int(x) != None:
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
            old_file = None
            for a in aids:
                v = a.get_value()
                x = Utils.convert_int(v)
                # Only need oem AIDs, since base AIDs are
                # define in the included <private/android_filesystem_config.h>
                if not Utils.is_oem_aid(x):
                    continue

                f = a.get_file_found()
                n = a.get_name()

                if f != old_file:
                    # use the preserved str value
                    print FSConfigGen.FILE_COMMENT % f
                    old_file = f
                print FSConfigGen.AID_DEFINE % (n.upper(), v)

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
        if Utils.convert_int(name) != None:
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

    _GENERATED = '''
#
# THIS IS AN AUTOGENERATED FILE! DO NOT MODIFY
#'''

    def __call__(self, files, dirs, aids):

        old_file = None

        print PasswdGen._GENERATED

        for a in aids:
            file_name = a. get_file_found()
            name = a.get_name()
            homedir = a.get_homedir()
            interpreter = a.get_interpreter()

            # use the normalized int value in passwd files.
            aid = a.get_value()

            if file_name != old_file:
                print '# source: "' + file_name + '"'
                old_file = file_name

            # name:passwd:uid:gid:user_name:homedir:intepreter
            print PasswdGen.format_passwd(name, aid, aid, name, homedir, interpreter)

    @staticmethod
    def format_passwd(name, uid, gid, user_name, homedir, interpreter):
        # name:passwd:uid:gid:user_name:homedir:intepreter
        return '%s::%s:%s:%s:%s:%s' % (name, uid, gid, name, homedir, interpreter)

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
