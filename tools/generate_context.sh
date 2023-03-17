#!/bin/sh
#
# Copyright (C) 2023 The Android Open Source Project
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

# Use pre-defined location for the artifact
ARTIFACT_PATH=$OUT_DIR/prebuilt_cached/artifacts/ssi

ARG_SHORT=o:,h
ARG_LONG=output:,help
OPTS=$(getopt -n generate_context --options $ARG_SHORT --longoptions $ARG_LONG -- "$@")

eval set -- "$OPTS"

function print_usage(){
  echo "usage: generate_context.sh --output <output_path>"
  exit 2
}

while :
do
  case "$1" in
    -o | --output )
      OUT_PATH="$2"
      shift 2
      ;;
    -h | --help )
      print_usage
      ;;
    -- )
      shift;
      break
      ;;
    * )
      print_usage
      ;;
  esac
done

if [ -z OUT_PATH ] ; then
  print_usage
fi

TARGET_ARCHIVE_FILE=$(find $ARTIFACT_PATH -iname '*target_files*' 2>/dev/null)

if [ "$?" -ne "0" ] ; then
  echo "Cannot find target_files archive file from $ARTIFACT_PATH"
  exit 1
fi

TARGET_OUT_PATH=$OUT_DIR/artifact

rm -rf $TARGET_OUT_PATH
mkdir -p $TARGET_OUT_PATH
unzip $TARGET_ARCHIVE_FILE 'SYSTEM/*' 'SYSTEM_EXT/*' -d $TARGET_OUT_PATH

if [ "$?" -ne "0" ] ; then
  echo "Failed to unzip file $TARGET_ARCHIVE_FILE into $TARGET_OUT_PATH"
  exit 1
fi

# Print information into the output file
echo TIMESTAMP=`date +%m%d%H%M%S` > $OUT_PATH
echo TARGET_OUT_PATH=$TARGET_OUT_PATH >> $OUT_PATH