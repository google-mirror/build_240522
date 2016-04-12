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

class Utils():
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

    @staticmethod
    def in_range(value, ranges):
        return any(lower <= value <= upper for (lower, upper) in ranges)

class AID():
    def __init__(self, define, value):
        self.define = define
        self.value = value

class AIDHeaderParser():

    _SKIPWORDS = [ 'UNUSED' ]
    _AID_KW = 'AID_'
    _AID_DEFINE = re.compile('\s*#define\s+%s.*' % _AID_KW)
    _OEM_START_KW = 'START'
    _OEM_END_KW = 'END'
    _OEM_RANGE = re.compile('AID_OEM_RESERVED_[0-9]*_{0,1}(%s|%s)' %(_OEM_START_KW, _OEM_END_KW))

    # Some of the AIDS like AID_MEDIA_EX had names like mediaex
    # list a map of things to fixup until we can correct these
    # at a later data
    _FIXUPS = {
        'media_drm' : 'mediadrm',
        'media_ex' : 'mediaex',
        'media_codec' : 'mediacodec'
    }

    def __init__(self, aid_header):
        self._aid_header = aid_header
        self._aid_name_to_value = {}
        self._aid_value_to_name = {}
        self._oem_ranges = {}

        with open(aid_header) as f:
            self._parse(f)

        self._process_and_check()

    def _parse(self, f):

        lineno = 0

        for l in f:
            lineno = lineno + 1

            errmsg = 'Error "%s" in file: "' + self._aid_header + '" on line: ' + str(lineno)

            if AIDHeaderParser._AID_DEFINE.match(l):
                chunks = l.split()

                if any(x in chunks[1] for x in AIDHeaderParser._SKIPWORDS):
                    continue

                orig = chunks[1]
                value = Utils.convert_int(chunks[2])
                if value == None:
                    raise Exception(errmsg % ('Invalid aid value aid "%s"' % chunks[2]))

                if AIDHeaderParser._is_oem_range(orig):
                    try:
                        self._handle_oem_range(orig, value)
                    except Exception as e:
                        raise Exception(errmsg % (str(e) + ' for "%s"' % orig))

                    continue

                try:
                    self._handle_aid(orig, value)
                except Exception as e:
                    raise Exception(errmsg % (str(e) + ' for "%s"' % orig))

    def _handle_aid(self, orig, value):

        # friendly name
        name = AIDHeaderParser._convert_friendly(orig)

        # duplicate name
        if name in self._aid_name_to_value:
            raise Exception('Duplicate aid "%s"' % orig)

        if value in self._aid_value_to_name:
            raise Exception('Duplicate aid value "%u" for %s' % value, orig)

        self._aid_name_to_value[name] = AID(orig, value)
        self._aid_value_to_name[value] = name

    def _handle_oem_range(self, orig, value):

        # convert AID_OEM_RESERVED_START or AID_OEM_RESERVED_<num>_START
        # to AID_OEM_RESERVED or AID_OEM_RESERVED_<num>
        is_start = orig.endswith(AIDHeaderParser._OEM_START_KW)

        tostrip = len(AIDHeaderParser._OEM_START_KW) if is_start \
                    else len(AIDHeaderParser._OEM_END_KW)

        # ending _
        tostrip = tostrip + 1

        strip = orig[:-tostrip]
        if strip not in self._oem_ranges:
            self._oem_ranges[strip] = []

        if len(self._oem_ranges[strip]) > 2:
            raise Exception('Too many same OEM Ranges "%s"' % orig)

        if len(self._oem_ranges[strip]) == 1:
            x = self._oem_ranges[strip][0]

            if x == value:
                raise Exception('START and END values equal %u' % value)
            elif is_start and x < value:
                raise Exception('END value %u less than START value %u' % (x, value))
            elif not is_start and x > value:
                raise Exception('END value %u less than START value %u' % (value, x))

        self._oem_ranges[strip].append(value)

    def _process_and_check(self):

        # tuplefy the lists since we don't want them mutable
        self._oem_ranges = [
            AIDHeaderParser._convert_lst_to_tup(k, v) for k,v in self._oem_ranges.iteritems()
        ]

        # Check for overlapping ranges
        for r in self._oem_ranges:
            for t in self._oem_ranges:
                if r != t and AIDHeaderParser._get_overlap(r, t):
                    raise Exception("Overlapping OEM Ranges found %s and %s" % (str(r), str(t)))

        # No aids should be in OEM range
        for a in self._aid_value_to_name:

            if Utils.in_range(a, self._oem_ranges):
                name = self._aid_value_to_name[a]
                raise Exception('AID "%s" value: %u within reserved OEM Range: "%s"' % (name, a, str(self._oem_ranges)))

    def get_oem_ranges(self):
        return self._oem_ranges

    def get_aids(self):
        return self._aid_name_to_value

    @staticmethod
    def _convert_lst_to_tup(name, lst):

        if not lst or len(lst) !=2:
            raise Exception('Mismatched range for "%s"' % name)

        return tuple(lst)

    @staticmethod
    def _convert_friendly(aid_define):
        # Translate AID_FOO_BAR to foo_bar (ie name)
        name = aid_define[len(AIDHeaderParser._AID_KW):].lower()

        if name in AIDHeaderParser._FIXUPS:
            return AIDHeaderParser._FIXUPS[name]

        return name

    @staticmethod
    def _is_oem_range(aid):
        return AIDHeaderParser._OEM_RANGE.match(aid)

    @staticmethod
    def _get_overlap(a, b):
        return max(0, min(a[1], b[1]) - max(a[0], b[0]))

class FSConfigFileParser():

    _AID_MATCH = re.compile('AID_[a-zA-Z]+')

    def __init__(self, config_files, oem_ranges):

        self._files = []
        self._dirs = []
        self._aids = []
        self._oem_ranges = oem_ranges

        self._seen_paths = {}
        # (name to file, value to aid)
        self._seen_aids = ({}, {})

        self._config_files = config_files

        for f in self._config_files:
            self._parse(f)

    def _parse(self, file_name):

            # Separate config parsers for each file found. If you use read(filenames...) later
            # files can override earlier files which is not what we want. Track state across
            # files and enforce with handle_dup(). Note, strict ConfigParser is set to true in
            # Python >= 3.2, so in previous versions same file sections can override previous
            # sections.
            config = ConfigParser.ConfigParser()
            config.read(file_name)

            for s in config.sections():

                if FSConfigFileParser._AID_MATCH.match(s) and config.has_option(s, 'value'):
                    FSConfigFileParser._handle_dup('AID', file_name, s, self._seen_aids[0])
                    self._seen_aids[0][s] = file_name
                    self._handle_aid(file_name, s, config)
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

    def _handle_aid(self, file_name, section_name, config):
        value = config.get(section_name, 'value')

        errmsg = '%s for: \"' + section_name + '" file: \"' + file_name + '\"'

        if not value:
            raise Exception(errmsg % 'Found specified but unset "value"')

        v = Utils.convert_int(value)
        if not v:
            raise Exception(errmsg % ('Invalid "value", not a number, got: \"%s\"' % value))

        # Values must be within OEM range
        if self._oem_ranges and not Utils.in_range(v, self._oem_ranges):
            s = '"value" not in valid range %s, got: %s'
            s = s % (str(self._oem_ranges), value)
            raise Exception(errmsg % s)

        # use the normalized int value in the dict and detect
        # duplicate definitions of the same vallue
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
                    if Utils.convert_int(x):
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

    _GENERATED = '''
/*
 * THIS IS AN AUTOGENERATED FILE! DO NOT MODIFY
 */
 '''
    _INCLUDE = '#include <private/android_filesystem_config.h>'

    _DEFINE_NO_DIRS = '#define NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS\n'
    _DEFINE_NO_FILES = '#define NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_FILES\n'

    _DEFAULT_WARNING = '#warning No device-supplied android_filesystem_config.h, using empty default.'

    _NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS_ENTRY = '{ 00000, AID_ROOT,      AID_ROOT,      0, "system/etc/fs_config_dirs" },'
    _NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_FILES_ENTRY = '{ 00000, AID_ROOT,      AID_ROOT,      0, "system/etc/fs_config_files" },'

    _IFDEF_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS = '#ifdef NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS'
    _ENDIF = '#endif'

    _OPEN_FILE_STRUCT = 'static const struct fs_path_config android_device_files[] = {'
    _OPEN_DIR_STRUCT = 'static const struct fs_path_config android_device_dirs[] = {'
    _CLOSE_FILE_STRUCT = '};'

    _GENERIC_DEFINE = "#define %s\t%s"

    _FILE_COMMENT = '// Defined in file: \"%s\"'

    def __init__(self, opt_group):
        opt_group.add_argument('fsconfig', nargs='+', help='The list of fsconfig files to parse')
        opt_group.add_argument('--aidhdr', help='An android_filesystem_config.h file to parse AIDs and OEM Ranges from')

    def __call__(self, args):

        # If supplied, parse the header file for OEM Ranges
        oem_ranges = None
        if 'aidhdr' in args and args['aidhdr']:
            hdr = AIDHeaderParser(args['aidhdr'])
            oem_ranges = hdr.get_oem_ranges()

        parser = FSConfigFileParser(args['fsconfig'], oem_ranges)
        FSConfigGen._generate(parser.get_files(), parser.get_dirs(), parser.get_aids())

    @staticmethod
    def _generate(files, dirs, aids):
        print FSConfigGen._GENERATED
        print FSConfigGen._INCLUDE
        print

        are_dirs = len(dirs) > 0
        are_files = len(files) > 0
        are_aids = len(aids) > 0

        if are_aids:
            for a in aids:
                # use the preserved str value
                print FSConfigGen._FILE_COMMENT % a[0]
                print FSConfigGen._GENERIC_DEFINE % (a[1], a[2])

            print

        if not are_dirs:
            print FSConfigGen._DEFINE_NO_DIRS

        if not are_files:
            print FSConfigGen._DEFINE_NO_FILES

        if not are_files and not are_dirs and not are_aids:
            print FSConfigGen._DEFAULT_WARNING
            return

        if are_files:
            print FSConfigGen._OPEN_FILE_STRUCT
            for tup in files:
                f = tup[0]
                c = tup[1]
                c[4] = '"' + c[4] + '"'
                c = '{ ' + '    ,'.join(c) + ' },'
                print FSConfigGen._FILE_COMMENT % f
                print '    ' + c

            if not are_dirs:
                print FSConfigGen._IFDEF_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS
                print '    ' + FSConfigGen._NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS_ENTRY
                print FSConfigGen._ENDIF
            print FSConfigGen._CLOSE_FILE_STRUCT

        if are_dirs:
            print FSConfigGen._OPEN_DIR_STRUCT
            for d in dirs:
                f[4] = '"' + f[4] + '"'
                d = '{ ' + '    ,'.join(d) + ' },'
                print '    ' + d

            print FSConfigGen._CLOSE_FILE_STRUCT

def main():

    opt_parser = argparse.ArgumentParser(description='A tool for parsing fsconfig config files and producing digestable outputs.')
    subparser = opt_parser.add_subparsers(help='generators')

    generators = generator.get()
    tmp = {}
    opts = []

    # for each generator, instantiate and add them as an option
    for n, g in generators.iteritems():

        p = subparser.add_parser(n, help=g.__doc__)
        p.set_defaults(which=n)

        opt_group = p.add_argument_group(n + ' options')

        # Instantiate and save
        tmp[n] = g(opt_group)

    # reassign constructed generators
    generator.set(tmp)

    args = opt_parser.parse_args()

    files = []
    dirs = []
    aids = []
    seen_paths = {}

    d = vars(args)
    which = d['which']
    del d['which']

    generator.get()[which](d)

if __name__ == '__main__':
    main()
