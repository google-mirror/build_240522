#!/usr/bin/env bash
# Copyright (C) 2022 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

source $(dirname $0)/../envsetup.sh

unset OUT_DIR OUT_DIR_COMMON_BASE

function check_getoutdir
(
    OUT_DIR=$1
    OUT_DIR_COMMON_BASE=$2
    [ "$3" = "$(getoutdir)" ] || ( echo "OUT_DIR=$1 OUT_DIR_COMMON_BASE=$2:  expected getoutdir to return '$3', got '$(getoutdir)'" && exit 1 )
)

# default
check_getoutdir "" "" $(gettop)/out

# OUT_DIR
check_getoutdir out       "" $(gettop)/out
check_getoutdir out1      "" $(gettop)/out1
check_getoutdir /tmp/out1 "" /tmp/out1
check_getoutdir ../out    "" $(gettop)/../out

# OUT_DIR_COMMON_BASE
check_getoutdir "" base     $(gettop)/base/$(basename $(gettop))
check_getoutdir "" /mnt     /mnt/$(basename $(gettop))
check_getoutdir "" /mnt/out /mnt/out/$(basename $(gettop))
check_getoutdir "" /mnt/..  /mnt/../$(basename $(gettop))

# both, OUT_DIR_COMMON_BASE is ignored
check_getoutdir out1 /mnt $(gettop)/out1
