#!/bin/bash
set -e

STRIP_PATH="${1}"
CORE="${2}"
VENDOR="${3}"

stripped_core="${CORE}.vndk_lib_check.stripped"
stripped_vendor="${VENDOR}.vndk_lib_check.stripped"

function cleanup() {
  rm -f ${stripped_core} ${stripped_vendor}
}
trap cleanup EXIT

function strip_lib() {
  ${STRIP_PATH} \
    -i ${1} \
    -o ${2} \
    -d /dev/null \
    --remove-build-id
}

strip_lib ${CORE} ${stripped_core}
strip_lib ${VENDOR} ${stripped_vendor}
if ! cmp -s ${stripped_core} ${stripped_vendor}; then
  echo "VNDK library has must_use_vendor_variant=false but has different core and vendor variant: $(basename ${CORE})"
  exit 1
fi
