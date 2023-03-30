#!/usr/bin/env python3
#
# Copyright (C) 2023 The Android Open Source Project
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
Generates a list of files that need to be included in the android-xts.zip.
Duplicate files all filter out if this feature is enabled in the given
test suite.

Usage:
  list_xts_file.py
    --xts-test-dir [directory of android-xts]
    --test-suite-name [name of the test suite]
    --output [output file]
"""
import sys
import argparse
import hashlib
import os

from os.path import getsize


# For enabled test suites, this tool will filter out duplicate files.
_REMOVE_DUPLICATE_FILES_ENABLED = ('cts')
_RESOURCE_DIRS = ('/tools', '/testcases', '/lib',  '/lib64')
_RESOURCE_FILES = ('/NOTICE.txt')


def _get_args():
    """Parses input arguments."""
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '--xts-test-dir', required=True,
        help='directory of android-xts')
    parser.add_argument(
        '--test-suite-name', required=True,
        help='name of the test suite')
    parser.add_argument(
        '-o', '--output', required=True,
        help='file path of the output')
    return parser.parse_args()


def _get_checksum(file_path: str) -> bytes:
    """Gets the checksum for the given file."""
    checksum = hashlib.md5()
    with open(file_path, 'rb') as file:
        while True:
            buffer = file.read(8192)
            if not buffer: break
            checksum.update(buffer)
    checksum = checksum.digest()
    return checksum


def _get_xts_module_dirs(files: set[str]) -> set[str]:
    """Gets all xTS test module dirs.

    Each xTS test module contains a xxxxx.config file directly under the
    module dir. This function decides a dir as an xTS test module dir by
    searching for the xxxxx.config file.

    Args:
        files: A set of files under the android-xts dir.

    Returns:
        A set of xTS test module dirs.
    """
    xts_module_dirs = set()
    for file in files:
        if file.endswith('.config'):
            module_dir = file.rsplit("/", 1)[0]
            xts_module_dirs.add(module_dir)
    return xts_module_dirs


def _get_non_duplicate_files(xts_test_dir: str, files: set[str]) -> set[str]:
    """"Gets a set of non-duplicate files to be included in the android-xts.zip.

    This function only remove duplicate files under the android-xts/testcases dir.
    Files under xTS test module dirs are all kept. Duplicate files under non-xTS
    test module dirs are removed.

    Args:
        xts_test_dir: The top android-xts dir.
        files: A set of files under the xts_test_dir.

    Returns:
        A set of files without unnecessary duplications.
    """

    duplicate_files_by_file_key = {}
    xts_module_dirs = _get_xts_module_dirs(files)
    test_cases_dir = os.path.join(xts_test_dir, 'testcases')

    for file in files:
        if not file.startswith(test_cases_dir):
            continue
        file_name = os.path.basename(file)
        file_key = (file_name, getsize(file), _get_checksum(file))
        duplicate_files_by_file_key.setdefault(file_key, set()).add(file)

    for file_key, duplicate_files in duplicate_files_by_file_key.items():
        if len(duplicate_files) <= 1:
            continue
        for duplicate_file in duplicate_files:
            # The path format of each file is {test_cases_dir}/{module_name}/../{file}.
            module_name = duplicate_file.removeprefix(test_cases_dir).split('/', 2)[1]
            module_dir = os.path.join(test_cases_dir, module_name)
            if module_dir in xts_module_dirs:
                continue
            files.remove(duplicate_file)
    return files


def _get_file_and_dir_list(xts_test_dir: str, test_suite_name: str) -> list[str]:
    """Gets a list of files and dirs that need to be included in the android-xts.zip.

    This function first filter out all files under the top android-xts dir based on
    _RESOURCE_DIRS and _RESOURCE_FILES. If the feature of filtering out duplicated
    files is enabled in the given test suite, all unnecessary duplicated files
    are further excluded. The last step in thie function is to collect parent dirs 
    of all filtered files.

    Args:
        xts_test_dir: The top android-xts dir.
        test_suite_name: The name of the test suite.

    Returns:
        A list of files and dirs.
    """
    file_paths = set()
    for root, _, files in os.walk(xts_test_dir):
        for file in files:
            file_path = os.path.join(root, file)
            file_path_without_prefix = file_path.removeprefix(xts_test_dir)
            if (file_path_without_prefix.startswith(_RESOURCE_DIRS)
                or file_path_without_prefix in _RESOURCE_FILES):
                file_paths.add(file_path)

    if test_suite_name in _REMOVE_DUPLICATE_FILES_ENABLED:
        file_paths = _get_non_duplicate_files(xts_test_dir, file_paths)

    dir_paths = set()
    for file_path in file_paths:
        path = file_path
        while True:
            parent_dir = os.path.dirname(path)
            if parent_dir == xts_test_dir or parent_dir in dir_paths:
                break
            dir_paths.add(parent_dir)
            path = parent_dir
    return list(dir_paths | file_paths)


def main(argv):
    args = _get_args()

    xts_test_dir = args.xts_test_dir
    test_suite_name = args.test_suite_name
    output_file = args.output

    with open(output_file, 'w') as file:
        file_paths = _get_file_and_dir_list(xts_test_dir, test_suite_name)
        file_paths.sort()
        file.write('\n'.join(file_paths))

if __name__ == "__main__":
    main(sys.argv)
