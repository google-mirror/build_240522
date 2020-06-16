#!/bin/bash
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

# Convert Android SDK version to number (unless it already is).
function sdk_version_to_number {
  if [ $# -ne 1 ]; then
    echo "Usage: $0 <sdk-version>"
    exit 1
  fi

  local v="$1"
  case "${v}" in
    "K") v=19 ;;
    "L") v=21 ;;
    "M") v=23 ;;
    "N") v=24 ;;
    "O") v=26 ;;
    "P") v=28 ;;
    "Q") v=29 ;;
    "R") v=30 ;;
    "S") v=31 ;;
    *) ;; # hope it is numeric
  esac
  echo "${v}"
}
