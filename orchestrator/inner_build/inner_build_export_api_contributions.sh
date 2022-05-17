#!/bin/bash

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

set -euo pipefail

# TODO: Add doc for this script

# TODO: Build a hardcoded test target for now
BUILD_TARGET=build/make/orchestrator/core/test/apis:cdk


# TODO: Add documentation
# Globals:
#   INNER_TREE_TOP: path to root of the inner tree
#   BUILD_TARGET: TODO add documentation
#   API_DOMAIN: TODO add documentation
# Arguments:
#   $1: bazel subcommand, e.g. build, cquery etc
#   $2-n: additional arguments for bazel subcommand
function bazel() {
  bazel_subcommand="$1"
  shift;
  (cd "${INNER_TREE_TOP}" && STANDALONE_BAZEL=1 tools/bazel "${bazel_subcommand}" "${BUILD_TARGET}" --define api_domain="${API_DOMAIN}" "$@")
}

# Parse input arguments
# Globals:
# Arguments:
#   $@: $@ passed to the script
# Output: Sets the following global variables
#  API_DOMAIN: TODO add documentation
#  INNER_TREE_TOP: path to root of the inner tree
#  OUT_DIR: TODO: add documentation
function parseargs() {
  # Parse arguments
  SHORT_OPTIONS=a:,i:,o:
  LONG_OPTIONS=api_domain:,inner_tree:,out_dir:
  OPTS=$(getopt --options "${SHORT_OPTIONS}" --longoptions "${LONG_OPTIONS}" -- "$@")
  eval set -- "${OPTS}"
  while :
  do
    case "$1" in
      -a | --api_domain )
        API_DOMAIN="$2"
        shift 2;
        ;;
      -i | --inner_tree )
        INNER_TREE_TOP="$2"
        shift 2;
        ;;
      -o | --out_dir )
        OUT_DIR="$2"
        shift 2;
        ;;
      --)
        shift;
        break
        ;;
      *)
        echo Received unrecognized option "$1"
        ;;
    esac
  done
}

# Parse input args and set global variables
parseargs $@

# Build the top level CDK contribution for the api domain
bazel build

# Use bazel query to determine the paths of the generated files
# This is from the root of the workspace, i.e. relative to $INNER_TREE_TOP
contribution_files=$(bazel cquery --output=starlark --starlark:expr="' '.join([f.path for f in target.files.to_list()])")

# Copy the metadata files to the out directory
# TODO: Removing an api contribution should delete the corresponding file from out_dir
# TODO: Copy operation should happen only if content has changed
# TODO: Instead of copy operation, should we write ninja rules instead?
out_dir_api_contributions="${OUT_DIR}"/api_contributions
mkdir -p "${out_dir_api_contributions}"
(cd "${INNER_TREE_TOP}" && cp -f -t "${out_dir_api_contributions}" "${contribution_files}")
