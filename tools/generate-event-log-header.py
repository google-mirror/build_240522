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
Usage: generate-event-log-header.py <input-logtag-file> <output-header-file>
"""

import os.path
import sys


if len(sys.argv) != 3:
  print __doc__
  sys.exit(1)

in_logtag = sys.argv[1]
out_header = sys.argv[2]

logtags = []
with open(in_logtag, 'r') as f:
  logtags = f.readlines()

macros = []
for l in logtags:
  l = l.strip()
  if not l:
    continue
  if l.startswith('#'):
    continue

  fields = l.split()
  if len(fields) < 2:
    continue
  num = fields[0]
  tag = fields[1]
  if not fields[0].isdigit():
    continue
  macro_name = tag.upper() + '_LOG_TAG'
  macro = '#define %s %s\n' % (macro_name, num)
  macros.append(macro)

guard_macro = os.path.basename(out_header).upper()
guard_macro = guard_macro.replace('.', '_').replace('-', '_')
guard_macro = '_' + guard_macro

with open(out_header, 'w') as f:
  f.write('/* Auto-generated from %s, DO NOT EDIT! */\n' % in_logtag)
  f.write('#ifndef %s\n' % guard_macro)
  f.write('#define %s\n' % guard_macro)
  f.write('\n')
  f.writelines(macros)
  f.write('\n')
  f.write('#endif')
