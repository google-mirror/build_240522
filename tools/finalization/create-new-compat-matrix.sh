#!/bin/bash

# Copyright 2024 Google Inc. All rights reserved.
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


set -ex

function create_new_compat_matrix() {

    local top="$(dirname "$0")"/../../../..
    source $top/build/make/tools/finalization/environment.sh

    # create the new file and modify the level
    local current_file=compatibility_matrix."$CURRENT_COMPATIBILITY_MATRIX_LEVEL".xml
    local final_file=compatibility_matrix."$FINAL_COMPATIBILITY_MATRIX_LEVEL".xml
    local src=$top/hardware/interfaces/compatibility_matrices/$current_file
    local dest=$top/hardware/interfaces/compatibility_matrices/$final_file
    sed "s/level=\""$CURRENT_COMPATIBILITY_MATRIX_LEVEL"\"/level=\""$FINAL_COMPATIBILITY_MATRIX_LEVEL"\"/" "$src" > "$dest"

    # add the new module to the end of the Android.bp file
    local bp_file=$top/hardware/interfaces/compatibility_matrices/Android.bp
    echo "
      vintf_compatibility_matrix {
          name: \"framework_$final_file\",
          stem: \"$final_file\",
          srcs: [
              \"$final_file\",
          ],
      }" >> $bp_file

    bpfmt -w $bp_file

    local make_file=$top/hardware/interfaces/compatibility_matrices/Android.mk
    # replace the current compat matrix in the make file with the final one
    # the only place this resides is in the conditional addition
    sed -i "s/$current_file/$final_file/g" $make_file
    # add the current compat matrix to the unconditional addition
    sed -i "/^    framework_compatibility_matrix.device.xml/i \    framework_$current_file \\\\" $make_file
}

create_new_compat_matrix
