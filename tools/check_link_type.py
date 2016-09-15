#!/usr/bin/env python
#
# Copyright (C) 2012 The Android Open Source Project
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

from __future__ import print_function
import argparse
import os
import sys

def main(args):
    failed = False
    message = '\033[1;31merror:\033[0m\033[1m'
    if args.warn:
        message = '\033[1;35mwarning:\033[0m\033[1m'

    for dep in args.deps:
        dep_name = os.path.basename(os.path.dirname(dep))
        if dep_name.endswith('_intermediates'):
            dep_name = dep_name[:len(dep_name)-len('_intermediates')]

        dep_type = ''
        with open(dep, 'r') as dep_file:
            dep_type = dep_file.read().strip()

        if dep_type in args.allowed:
            continue

        if not failed:
            print('', file=sys.stderr)
        print('\033[1m%s: %s %s (%s) should not link to %s (%s)\033[0m' %
              (args.makefile, message, args.module, args.type, dep_name,
               dep_type), file=sys.stderr)
        failed = True

    if failed:
        print('', file=sys.stderr)

        if not args.warn:
            sys.exit(1)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Check link types')
    parser.add_argument('--warn', help='Warn instead of error',
                        action='store_true')
    parser.add_argument('--makefile', help='Makefile defining module')
    parser.add_argument('--module', help='The module being checked')
    parser.add_argument('--type', help='The link type of module')
    parser.add_argument('--allowed', help='The allowed types for deps',
                        action='append', metavar='TYPE')
    parser.add_argument('deps', help='The dependencies to check',
                        metavar='DEP', nargs='*')
    main(parser.parse_args())
