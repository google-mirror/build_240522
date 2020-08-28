#!/usr/bin/env python
#
# Copyright (C) 2016 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the 'License');
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an 'AS IS' BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import getopt, json, sys

def PrintFileNames(path):
  with open(path) as jf:
    data = json.load(jf)
  for line in data:
    print(line["Name"])

def PrintCanonicalList(path):
  with open(path) as jf:
    data = json.load(jf)
  for line in data:
    print "{0:12d}  {1}".format(line["Size"], line["Name"])

def CompareFileNames(path, snapshot, mode):
  with open(path) as jf:
    data = json.load(jf)
  new_files = []
  for line in data:
    new_files.append(str(line["Name"]))
  with open(snapshot) as f:
    ori_files = f.read().splitlines()
  error = False
  if '+' in mode:
    added = frozenset(new_files).difference(ori_files)
    for d in added:
      print("  +" + d)
      error = True
  if '-' in mode:
    deled = frozenset(ori_files).difference(new_files)
    for d in deled:
      print("  -" + d)
      error = True
  if error:
    print("Fatal differences are detected. Or you may want to update the snapshot with:")
    print(" %s -n %s > %s " % (sys.argv[0], path, snapshot))
    sys.exit(2)

def PrintUsage(name):
  print("""
Usage: %s -[nc] json_files_list | -f jason_files_list snapshot_flies_list:[+-]*
 -n produces list of files only
 -c produces classic installed-files.txt
 -f compare list of files with the snapshot and the specified mode
""" % (name))

def main(argv):
  try:
    opts, args = getopt.getopt(argv[1:], "ncf", "")
  except getopt.GetoptError, err:
    print(err)
    PrintUsage(argv[0])
    sys.exit(2)

  if len(opts) == 0:
    print("No conversion option specified")
    PrintUsage(argv[0])
    sys.exit(2)

  if len(args) == 0:
    print("No input file specified")
    PrintUsage(argv[0])
    sys.exit(2)

  for o, a in opts:
    if o == ("-n"):
      PrintFileNames(args[0])
      sys.exit()
    elif o == ("-c"):
      PrintCanonicalList(args[0])
      sys.exit()
    elif o == ("-f"):
      if len(args) < 2:
        print("PRODUCT_SYSTEM_FREEZE is not specified")
        return 
      t = args[1].split(':')
      if len(t) < 2:
        print("No compare mode provided")
        PrintUsage(argv[0])
        sys.exit(2)
      CompareFileNames(args[0], t[0], t[1])
    else:
      assert False, "Unsupported option"

if __name__ == '__main__':
  main(sys.argv)
