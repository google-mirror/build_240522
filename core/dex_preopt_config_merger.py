#!/usr/bin/env python
#
# Copyright (C) 2020 The Android Open Source Project
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
"""A tool for merging dexpreopt.config files for <uses-library> dependencies
into the dexpreopt.config file of the library/app that uses them."""

from __future__ import print_function

import json
from collections import OrderedDict
import sys


def main():
  """Program entry point."""
  if len(sys.argv) < 2:
    raise SystemExit('usage: %s <config> [config ...]' % sys.argv[0])

  cfgs = []
  for arg in sys.argv[1:]:
    with open(arg, 'r') as f:
      cfgs.append(json.load(f, object_pairs_hook=OrderedDict))

  cfg0 = cfgs[0]

  uses_libs = {}
  for cfg in cfgs[1:]:
    uses_libs[cfg['Name']] = cfg

  clc_map = cfg0['ClassLoaderContexts']
  clc_map2 = OrderedDict()
  for sdk_ver in clc_map:
    clcs = clc_map[sdk_ver]
    clcs2 = OrderedDict()
    for lib in clcs:
      clc = clcs[lib]
      if lib in uses_libs:
        ulib = uses_libs[lib]
        clc['Host'] = ulib['BuildPath']
        clc['Device'] = ulib['DexLocation']
        clc['Subcontexts'] = ulib['ClassLoaderContexts'].get('any')
        clcs2[ulib['ProvidesUsesLib']] = clc
      else:
        clcs2[lib] = clc
    clc_map2[sdk_ver] = clcs2
  cfg0['ClassLoaderContexts'] = clc_map2

  with open(sys.argv[1], 'w') as f:
    f.write(json.dumps(cfgs[0], indent=4, separators=(',', ': ')))

if __name__ == '__main__':
  main()
