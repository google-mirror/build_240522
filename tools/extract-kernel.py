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

"""
A tool to extract kernel information from a kernel image.
"""

import argparse
import signal
import subprocess
import sys
import re

CONFIG_PREFIX = b'IKCFG_ST'
GZIP_HEADER = b'\037\213\010'
COMPRESSION_ALGO = (
    (["gzip", "-d"], GZIP_HEADER),
    (["xz", "-d"], b'\3757zXZ\000'),
    (["bzip2", "-d"], b'BZh'),
    (["lz4", "-d", "-l"], b'\002\041\114\030'),

    # These are not supported in the build system yet.
    # (["unlzma"], b'\135\0\0\0'),
    # (["lzop", "-d"], b'\211\114\132'),
)
VERSION_PREFIX = b'Linux version '
VERSION_CHARS = b'0123456789.'
VERSION_REGEX = r'[0-9]+[.][0-9]+[.][0-9]+'

def get_version(input_bytes, start_idx):
    end_idx = start_idx
    while end_idx < len(input_bytes) and input_bytes[end_idx] in VERSION_CHARS:
        end_idx += 1
    version = input_bytes[start_idx:end_idx]
    if re.match(VERSION_REGEX, version.decode()):
        return version
    return None

def dump_version(input_bytes):
    idx = 0
    while True:
        idx = input_bytes.find(VERSION_PREFIX, idx)
        if idx < 0:
            return None
        idx += len(VERSION_PREFIX)

        version = get_version(input_bytes, idx)
        if version:
            return version

def dump_configs(input_bytes):
    idx = input_bytes.find(CONFIG_PREFIX + GZIP_HEADER)
    if idx < 0:
        return None
    idx += len(CONFIG_PREFIX)

    sp = subprocess.Popen(["gzip", "-d", "-c"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    o, e = sp.communicate(input=input_bytes[idx:])
    if sp.returncode == 1: # error
        return None
    assert sp.returncode in (0, 2), sp.returncode # success or trailing garbage warning
    return o

def try_decompress(cmd, search_bytes, input_bytes):
    idx = input_bytes.find(search_bytes)
    if idx < 0:
        return None

    idx = 0
    sp = subprocess.Popen(cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    o, e = sp.communicate(input=input_bytes[idx:])
    # ignore errors
    return o

def decompress_dump(func, input_bytes):
    """
    Run func(input_bytes) first; and if that fails (returns value evaluates to False),
    then try different decompression algorithm before running func.
    """
    o = func(input_bytes)
    if o:
        return o
    for cmd, search_bytes in COMPRESSION_ALGO:
        decompressed = try_decompress(cmd, search_bytes, input_bytes)
        if decompressed:
            o = func(decompressed)
            if o:
                return o
        # Force decompress the whole file even if header doesn't match
        decompressed = try_decompress(cmd, b"", input_bytes)
        if decompressed:
            o = func(decompressed)
            if o:
                return o

def main():
    parser = argparse.ArgumentParser(formatter_class=argparse.RawTextHelpFormatter,
        description=__doc__ +
        "\nThese algorithms are tried when decompressing the image:\n    " +
        " ".join(tup[0][0] for tup in COMPRESSION_ALGO))
    parser.add_argument('--input',
                        help='Input kernel image. If not specified, use stdin',
                        metavar='FILE',
                        type=argparse.FileType('rb'),
                        default=sys.stdin)
    parser.add_argument('--output-configs',
                        help='If specified, write configs. Use stdout if no file is specified.',
                        metavar='FILE',
                        nargs='?',
                        type=argparse.FileType('wb'),
                        const=sys.stdout)
    parser.add_argument('--output-version',
                        help='If specified, write version. Use stdout if no file is specified.',
                        metavar='FILE',
                        nargs='?',
                        type=argparse.FileType('wb'),
                        const=sys.stdout)
    parser.add_argument('--tools',
                        help='Decompression tools to use. If not specified, PATH is searched.',
                        metavar='ALGORITHM:EXECUTABLE',
                        nargs='*')
    args = parser.parse_args()

    tools = {pair[0]: pair[1] for pair in (token.split(':') for token in args.tools or [])}
    for cmd, _ in COMPRESSION_ALGO:
        if cmd[0] in tools:
            cmd[0] = tools[cmd[0]]

    input_bytes = args.input.read()

    ret = 0
    if args.output_configs is not None:
        o = decompress_dump(dump_configs, input_bytes)
        if o:
            args.output_configs.write(o)
        else:
            sys.stderr.write("Cannot extract kernel configs in {}".format(args.input.name))
            ret = 1
    if args.output_version is not None:
        o = decompress_dump(dump_version, input_bytes)
        if o:
            args.output_version.write(o)
        else:
            sys.stderr.write("Cannot extract kernel versions in {}".format(args.input.name))
            ret = 1

    return ret


if __name__ == '__main__':
    exit(main())
