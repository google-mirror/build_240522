#
# Copyright (C) 2024 The Android Open Source Project
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

""" This script generates jarjar rule files to add a jarjar prefix to all classes, except those
that are API, unsupported API or otherwise excluded."""

import argparse
import io
import re
import subprocess
from zipfile import ZipFile


def parse_arguments(argv):
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '--jars', nargs='+', default=[],
        help='Path to pre-jarjar JAR. Multiple files can be specified.')
    parser.add_argument(
        '--srcjars', nargs='+', default=[],
        help='Path to pre-jarjar src JARS. Multiple files can be specified.')
    parser.add_argument(
        '--container', required=True,
        help='Container of the jarjared library'
             'for example "system".')
    parser.add_argument(
        '--output', required=True, help='Path to output jarjar rules file.')
    return parser.parse_args(argv)


def _list_classes(jar, is_src_jar):
    with ZipFile(jar, 'r') as zip:
        files = zip.namelist()
        file_ext = '.class' if not is_src_jar else '.java'
        class_len = len(file_ext)
        return [f.replace('/', '.')[:-class_len] for f in files
                if f.endswith(file_ext) and not f.endswith('/package-info.class')]


def make_jarjar_rules(args):

    classes = []
    for jar in args.jars:
        classes.extend(_list_classes(jar, False))
    for srcjar in args.srcjars:
        classes.extend(_list_classes(srcjar, True))

    classes.sort()
    with open(args.output, 'w') as f:
        for full_class_name in classes:
            package = full_class_name[:full_class_name.rfind('.')]
            class_name = full_class_name[full_class_name.rfind('.') + 1:]
            result = f'{package}.{args.container}.{class_name}'
            f.write(f'rule {full_class_name} {result}\n')


def _main():
    # Pass in None to use argv
    args = parse_arguments(None)
    make_jarjar_rules(args)


if __name__ == '__main__':
    _main()
