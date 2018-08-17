#!/usr/bin/env python
#
# Copyright (C) 2015 The Android Open Source Project
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
Extract package name from AndroidManifest.xml and print it on stdout.
"""

from __future__ import print_function
import os.path
import sys
import xml.etree.ElementTree as ET

def error(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)
    sys.exit(1)

if __name__ == "__main__":
  if len(sys.argv) < 2:
    error("Usage: " + sys.argv[0] + " <manifest>");

  fname = sys.argv[1]
  if (not os.path.isfile(fname)):
    error("Could not open file: " + fname)

  file = ET.parse(fname)
  if not file:
    error("Could not parse file: " + fname)

  if not "package" in file.getroot().attrib:
    error("Xml root does not contain \"package\" attribute")

  print(file.getroot().attrib["package"])
