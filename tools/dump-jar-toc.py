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

import os.path
import subprocess
import sys


if len(sys.argv) != 2:
  sys.stderr.write('Usage: %s <jar-file>\n' % sys.argv[0])
  sys.exit(1)
jar_name = sys.argv[1]

files_in_jar = subprocess.check_output(['jar', '-tf', jar_name])
classes = []
for file_in_jar in files_in_jar.splitlines():
  base, ext = os.path.splitext(file_in_jar)
  if ext == '.class':
    classes.append(base)

print subprocess.check_output(
    ['javap', '-constants', '-cp', jar_name] + classes)
