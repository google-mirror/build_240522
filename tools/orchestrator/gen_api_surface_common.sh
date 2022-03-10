#!/bin/bash
set -euo pipefail

# Copyright (C) 2022 The Android Open Source Project
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


# This script contains helper functions that will be used to generate api
# surfaces

# Use bazel query to determine output filepath of *_api_surface
# Globals:
# Arguments:
#  $1: inner_tree_root (also the Bazel workspace root)
#  $2: label of the Bazel target
# Outputs:
#  Writes absolute filepath location to stdout, including inner_tree_root prefix
function bazel_output_filepath(){
  # Use prebuilt bazel directly instead of the bazel() function added by
  # build/make/envsetup.sh
  # The shell function runs the prebuilt bazel using --config=bp2build|queryview, and
  # depends on values generated from Soong
  # The api surfaces do not have any dependencies on Soong, and therefore can be
  # built/queried in standalone/pure mode
  output=$(cd "$1" && \
    prebuilts/bazel/linux-x86_64/bazel aquery --output=jsonproto "$2" | \
    jq -r .pathFragments[].label | \
    xargs echo | \
    sed "s/\ /\//g" # Join path fragments using / separator
  )
  echo "$1"/"$output"
}

# Helper function to copy multiple files to out directory
# Globals:
# Arguments:
#  $1: target location in OUT_DIR
#  $2: inner_tree_root
#  $3-n: source files to copy, relative to workspace root
# Outputs:
#  None
function copy_files(){
  out_location="$1"
  shift;
  inner_tree="$1"
  shift;
  mkdir -p "$out_location"
  (cd "$inner_tree" && cp -t "$out_location" "$@")
}
