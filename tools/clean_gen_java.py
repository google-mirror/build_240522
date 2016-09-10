#!/usr/bin/env python
#
# Copyright (C) 2016 The Android Open Source Project
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
Checks the current build configurations against the previous build,
clean artifacts in TARGET_COMMON_OUT_ROOT if necessary.
If the list of generated java files has been changed, the library's generated
source folder needs to be cleaned.

Usage: clean_gen_java.py <intermediates dir>

The intermediates dir is usually $OUT_DIR/target/common/obj, it will contain
current_gen_java.txt and current_modules.txt (and previous_gen_java.txt from a
previous run of this tool). It will also be the root for the relative paths in
current_gen_java.txt.

current_modules.txt contains all java module names in the current build.  This
script modfies current_gen_java.txt if necessary: if there is a directory in
previous_gen_java.txt but absent from current_modules.txt, we copy that line
from previous_gen_java.txt over to current_gen_java.txt before replacing
previous_gen_java.txt. Usually that means we just don't care that module in the
current build disappeared (for example we are switching from a full build to a
partial build with mm/mmm), and we should carry on the previous gen_java config
so that previous_gen_java.txt always reflects the current status of the entire
tree.

Format of current_gen_java.txt and prev_gen_java.txt:
  <directory> <gen_src_file> [gen_src_file ...]
  <directory> <gen_src_file> [gen_src_file ...]
  ...

Format of current_modules.txt:
  <directory>
  <directory>
  ...
"""

from __future__ import print_function
import shutil
import sys
import os

def load_config(filename):
  with open(filename) as f:
    result = {}
    for line in f:
      line = line.strip()
      if not line or line.startswith("#"):
        continue
      words = line.split()
      result[words[0]] = " ".join(words[1:])
    return result

def main(argv):
  if len(argv) != 2:
    print(__doc__, file=sys.stderr)
    sys.exit(1)

  intermediates_dir = argv[1]
  cur_modules_file = os.path.join(intermediates_dir, "current_modules.txt")
  cur_gen_java_file = os.path.join(intermediates_dir, "current_gen_java.txt")
  prev_gen_java_file = os.path.join(intermediates_dir, "previous_gen_java.txt")

  for f in [cur_modules_file, cur_gen_java_file]:
    if not os.path.exists(f):
      print("%s missing" % (os.path.basename(f)), file=sys.stderr)
      sys.exit(1)

  def clean_dir(directory):
    p = os.path.join(intermediates_dir, directory, "src")
    if os.path.exists(p):
      print("Cleaning obsolete generated java files: "+directory)
      shutil.rmtree(p)

  with open(cur_modules_file) as f:
    all_modules = set(f.read().split())

  current_gen_java = load_config(cur_gen_java_file)

  if os.path.exists(prev_gen_java_file):
    previous_gen_java = load_config(prev_gen_java_file)

    carryon = []
    for p in current_gen_java:
      if p not in previous_gen_java:
        clean_dir(p)
      elif current_gen_java[p] != previous_gen_java[p]:
        clean_dir(p)
    for p in previous_gen_java:
      if p not in current_gen_java:
        if p in all_modules:
          clean_dir(p)
        else:
          # we don't build p in the current build.
          carryon.append(p)

    # Add carryon to the current gen_java config file.
    if carryon:
      with open(cur_gen_java_file, "a") as f:
        for p in carryon:
          f.write(p + " " + previous_gen_java[p] + "\n")

  # Move the current gen_java to the previous gen_java
  shutil.move(cur_gen_java_file, prev_gen_java_file)

if __name__ == "__main__":
  main(sys.argv)
