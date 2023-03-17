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

ARG_SHORT=c:,t:,o:,h
ARG_LONG=context:,target:,output:,help
OPTS=$(getopt -n import_artifact --options $ARG_SHORT --longoptions $ARG_LONG -- "$@")

eval set -- "$OPTS"

function print_usage(){
  echo "usage: import_artifact.sh --context <context_file> --target <target_file_path> [--output <output_file_path>]"
  exit 2
}

while :
do
  case "$1" in
    -c | --context )
      CONTEXT_FILE="$2"
      shift 2
      ;;
    -t | --target )
      TARGET_PATH="$2"
      shift 2
      ;;
    -o | --output )
      OUTPUT_PATH="$2"
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

if [ -z $CONTEXT_FILE ] || [ -z $TARGET_PATH ] ; then
  print_usage
fi

if [ -z $OUTPUT_PATH ]; then
  OUTPUT_PATH=$TARGET_PATH
fi

while IFS='=' read -ra line; do
  if [ ${line[0]} = 'TARGET_OUT_PATH' ]; then
    ARTIFACT_BASE_PATH=${line[1]}
  fi
done < $CONTEXT_FILE

if [ -z $ARTIFACT_BASE_PATH ] ; then
  echo "Cannot find a target out path from " $CONTEXT_FILE
  exit 2
fi

if [[ "$TARGET_PATH" == *"/system/"* ]] ; then
  SRC_PATH=$ARTIFACT_BASE_PATH$(sed 's/^.*\/system\//\/SYSTEM\//g' <<< $TARGET_PATH)
fi

if [[ "$TARGET_PATH" == *"/system_ext/"* ]] ; then
  SRC_PATH=$ARTIFACT_BASE_PATH$(sed 's/^.*\/system_ext\//\/SYSTEM_EXT\//g' <<< $TARGET_PATH)
fi

cp -f $SRC_PATH $OUTPUT_PATH