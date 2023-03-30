#!/usr/bin/env python
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

from glob import glob
from pathlib import Path
from os import remove
from os.path import join
import argparse

class FileLister:
    def __init__(self, args) -> None:
        self.out_file = args.out_file
        sanitize(self.out_file)

        self.folder_dir = args.dir
        self.extensions = [extension if
                           len(extension) > 0 and extension[0] == "." else
                           "." + extension for extension in args.extensions]
        self.files_list = list()

    def get_files(self) -> None:
        """Get all files directory in the input directory including the files in the subdirectories

        Recursively finds all files in the input directory.
        Set file_list as a list of file directory strings,
        which do not include directories but only files.
        List is sorted in alphabetical order of the file directories.

        Args:
            dir: Directory to get the files. String.

        Raises:
            FileNotFoundError: An error occurred accessing the non-existing directory
        """

        if not dir_exists(self.folder_dir):
            raise FileNotFoundError(f"Directory {self.folder_dir} does not exist")

        if self.folder_dir[:-2] != "**":
            self.folder_dir = join(self.folder_dir, "**")

        self.files_list = [file for file in sorted(glob(self.folder_dir, recursive=True)) if Path(file).is_file()]

    def list(self) -> None:
        self.files_list.reverse()
        while self.files_list:
            file = self.files_list.pop()
            if not self.extensions or Path(file).suffix in self.extensions:
                self.write(file)

    def write(self, line: str) -> None:
        if self.out_file == "":
            print(line)
        else:
            write_line(self.out_file, line)

###
# Helper functions
###
def dir_exists(dir: str) -> bool:
    return Path(dir).exists()

def write_line(out_file: str, line: str) -> None:
    with open(out_file, "a") as f:
        f.write(line + "\n")
        f.close()

def sanitize(dir: str) -> None:
    if dir_exists(dir):
        remove(dir)

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('dir', action='store', type=str,
                        help="directory to list all subdirectory files")
    parser.add_argument('--out', dest='out_file',
                        action='store', default="", type=str,
                        help="optional directory to write subdirectory files. If not set, will print to console")
    parser.add_argument('--extensions', nargs='*', default=list(), dest='extensions',
                        help="Extensions to include in the output. If not set, all files are included")

    args = parser.parse_args()

    file_lister = FileLister(args)
    file_lister.get_files()
    file_lister.list()
