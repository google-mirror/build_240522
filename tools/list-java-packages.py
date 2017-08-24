#!/usr/bin/env python
#
# Copyright (C) 2017 The Android Open Source Project
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
#

# Usage:
#   list-java-packages <java-file>*
#
# Prints the unique package names of the java files given on the
# command line
#
# This is done by scanning each <java-file> for a package declaration.
# Each file is scanned only up to the first package declaration
# found, and each directory is only scanned up to the first .java
# file that contains a package declaration.

import re
from sets import Set
import os.path
import sys

def line_match(pattern, path):
  with open(path) as file:
    for line in file.readlines():
      match = pattern.search(line);
      if (match):
        return match.group(1);
  return None

def packages(paths):
  result = Set([])
  dirs_done = Set([])
  pattern = re.compile("^package ([a-zA-Z0-9_\.]+);")
  for path in paths:
    dirname = os.path.dirname(path)
    if dirname in dirs_done:
      continue
    if (path.endswith(".java")):
      p = line_match(pattern, path)
      if (p):
        dirs_done.add(dirname)
        result.add(p)
  return sorted(result)

def create_module_info(module_name, packages):
  print("module " + module_name + ";");
  for package in sorted(packages):
    print("\texports " + package + ";")

def main(args):
  for package in packages(args):
    print(package)

if __name__ == "__main__":
  main(sys.argv)

