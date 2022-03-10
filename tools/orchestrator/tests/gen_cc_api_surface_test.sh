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

MOCK_BUILD_TOP=$(mktemp -t -d st.XXXXX)
REAL_BUILD_TOP="$(readlink -f "$(dirname "$0")"/../../../../../)"
MOCK_OUT_DIR=$(mktemp -t -d st.XXXXX)

# Helper script that creates a minimal bazel workspace in MOCK_BUILD_TOP
function setup(){
  rm -rf "$MOCK_BUILD_TOP"
  # Copy Bazel rules from build/bazel
  mkdir -p "$MOCK_BUILD_TOP"/build/bazel
  cp -R "$REAL_BUILD_TOP"/build/bazel "$MOCK_BUILD_TOP"/build
  # Copy Bazel binary
  mkdir -p "$MOCK_BUILD_TOP"/prebuilts/bazel
  cp -R "$REAL_BUILD_TOP"/prebuilts/bazel "$MOCK_BUILD_TOP"/prebuilts
  # Symlink BUILD and worksapce files
  ln -s "$REAL_BUILD_TOP"/BUILD "$MOCK_BUILD_TOP"/BUILD
  ln -s "$REAL_BUILD_TOP"/WORKSPACE "$MOCK_BUILD_TOP"/WORKSPACE
}

# Helper script that runs bazel build and then runs gen_api_surface.sh
# Arguments
# $1: Bazel target
function gen_cc_api_surface(){
  $(cd "$MOCK_BUILD_TOP" && prebuilts/bazel/linux-x86_64/bazel build "$1")
  $(dirname "$0")/../gen_cc_api_surface.sh "$MOCK_OUT_DIR" "$1" "$MOCK_BUILD_TOP"
}

function fail {
  echo -e "\e[91;1mFAILED:\e[0m" $*
  exit 1
}

# Helper script that writes a mock BUILD file to stdout
# The build file contains
# cc_api_surface called "mysurface"
# cc_api_contributions called "mycontribution1" and "mycontribution2"
function mock_build_file(){
cat << EOF
load("//build/bazel/rules/cc:cc_api_contribution.bzl", "cc_api_contribution")
load("//build/bazel/rules/cc:cc_api_surface.bzl", "cc_api_surface")

cc_api_contribution(
  name="mycontribution1",
  symbol_file="mycontribution1.map.txt",
  headers = ["mycontribution1.h"],
  first_version="29",
)
cc_api_contribution(
  name="mycontribution2",
  symbol_file="mycontribution2.map.txt",
  headers = ["mycontribution2.h"],
  first_version="29",
)
cc_api_surface(
  name="mysurface",
  contributions = [
    ":mycontribution1",
    ":mycontribution2",
  ],
)
EOF
  # Generate the contribution files
  touch $MOCK_BUILD_TOP/mycontribution1.h $MOCK_BUILD_TOP/mycontribution1.map.txt
  touch $MOCK_BUILD_TOP/mycontribution2.h $MOCK_BUILD_TOP/mycontribution2.map.txt
}
# This test checks that the build orchestrator
# 1. copies header files
# 2. copies symbol files
# 3. generates Android.bp
function test_files_are_copied(){
  setup
  mock_build_file > $MOCK_BUILD_TOP/BUILD
  gen_cc_api_surface //:mysurface

  [ -f $MOCK_OUT_DIR/mysurface/mycontribution1/include/mycontribution1.h ] || fail "Header files were not copied"
  [ -f $MOCK_OUT_DIR/mysurface/mycontribution1/mycontribution1.map.txt ] || fail "Symbol file was not copied"
  [ -f $MOCK_OUT_DIR/mysurface/mycontribution1/Android.bp ] || fail "Android.bp file was not generated"
}

# This test checks that the build orchestrator does not copy files if file
# contents have not changed. This ensures that inputs to ninja actions do not
# become dirty after every build orchestrator run
function test_files_are_not_copied_if_no_content_change(){
  setup
  mock_build_file > $MOCK_BUILD_TOP/BUILD
  gen_cc_api_surface //:mysurface
  headersMtime1=$(stat -c "%y" $MOCK_OUT_DIR/mysurface/mycontribution1/include/mycontribution1.h)

  touch $MOCK_BUILD_TOP/mycontribution1.h # touch file, but do no edit
  gen_cc_api_surface //:mysurface
  headersMtime2=$(stat -c "%y" $MOCK_OUT_DIR/mysurface/mycontribution1/include/mycontribution1.h)
  [ "$headersMtime1" == "$headersMtime2" ] || fail "mtime of headers was updated even though file content has not changed"

  echo changes > $MOCK_BUILD_TOP/mycontribution1.h # make edits to file
  gen_cc_api_surface //:mysurface
  headersMtime3=$(stat -c "%y" $MOCK_OUT_DIR/mysurface/mycontribution1/include/mycontribution1.h)
  [ "$headersMtime2" == "$headersMtime3" ] && fail "mtime of headers was not updated even though file content have changed"
}

function test_files_are_deleted_if_contribution_is_removed(){
  setup
  mock_build_file > $MOCK_BUILD_TOP/BUILD

  # Generate cc api surfaces, and check file exists
  gen_cc_api_surface //:mysurface
  [ -f $MOCK_OUT_DIR/mysurface/mycontribution1/include/mycontribution1.h ] || fail "Header files were not copied"

  # Remove contribution, and check files have been removed
  sed -i '/:mycontribution1/d' $MOCK_BUILD_TOP/BUILD
  gen_cc_api_surface //:mysurface
  [ -f $MOCK_OUT_DIR/mysurface/mycontribution1/include/mycontribution1.h ] && fail "Header files were not removed from out directory after contribution was deleted from api surface"
}


test_files_are_copied
test_files_are_not_copied_if_no_content_change
test_files_are_deleted_if_contribution_is_removed
