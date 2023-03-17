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

ARG_SHORT=a:,o:,h
ARG_LONG=artifact:,output:,help
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
    -a | --artifact )
      ARTIFACT_PATH="$2"
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

TARGET_OUT_PATH=$OUT_DIR/artifact

rm -rf $TARGET_OUT_PATH
mkdir -p $TARGET_OUT_PATH
unzip $ARTIFACT_PATH 'SYSTEM/*' 'SYSTEM_EXT/*' -d $TARGET_OUT_PATH

if [ "$?" -ne "0" ] ; then
  echo "Failed to unzip file $ARTIFACT_PATH into $TARGET_OUT_PATH"
  exit 1
fi

# Print information into the output file
echo TIMESTAMP=`date +%m%d%H%M%S` > $OUT_PATH
echo TARGET_OUT_PATH=$TARGET_OUT_PATH >> $OUT_PATH