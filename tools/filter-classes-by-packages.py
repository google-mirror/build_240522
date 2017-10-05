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
#   filter-class-paths-by-packages [include|exclude] <java-package>*
#
# Given a list of paths (e.g. java/util/ArrayList.class) of class
# files, prints the ones that are [in|not in] one of the given packages.

import re
import os.path
from sets import Set
import sys

def main(args):
  mode = args[1]
  if (not(mode in ("include", "exclude"))):
    raise Exception("invalid mode " + mode + ", should be either include or exclude")
  packages = Set(args[2:])
  pattern = re.compile("^(.+)/[^/]+\.class")
  for path in sys.stdin:
    path = os.path.normpath(path)
    match = pattern.search(path)
    if (match):
      package = match.group(1).replace("/", ".")
      if (bool(package in packages) == bool(mode == "include")):
        print(path)

if __name__ == "__main__":
  main(sys.argv)

