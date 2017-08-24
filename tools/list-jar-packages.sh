#!/bin/bash

# Copyright (C) 2017 The Android Open Source Project
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

# Prints out the Java package names of all classes in the .jar files
# whose paths are given on the command line.
#
# Usage:
# jar-packages.sh path/to/foo.jar path/to/bar.jar [...]

for j in "$@"; do jar tf $j ; done \
          | grep -E '[a-zA-Z0-9_\$]+\.class$' \
          | sed 's/\/[a-zA-Z0-9_\$]\+\.class$//g' \
          | sed 's/\//./g' \
          | sort \
          | uniq

