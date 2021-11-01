#!/bin/bash
set -eo pipefail

TARGETS=${@:-update-api}

function _print_update_api_help_msg() {
  printf \
"Utililty function to update API in source tree
Usage:
    update-api.sh [targets...]
Examples:
    update-api.sh # update-api is default
    update-api.sh update-api # update all APIs
    update-api.sh android.car-stubs-docs-update-current-api android.net.ipsec.ike.stubs.source-update-current-api # update subset of APIs
"
}

function _print_TARGET_PRODUCT_help_msg() {
  printf \
"TARGET_PRODUCT not set. To fix, use one of the following options
1. Rerun lunch in current shell
2. Pass to script as TARGET_PRODUCT=<target_product> update-api.sh [targets...]
"
}

# Returns the *update_current_api.sh input to the *update-current-api phony target
# *update_current_api.sh contains commands to cp the updated API to the source tree
function _get_update_api_script_names() {
  $(gettop)/prebuilts/build-tools/linux-x86/bin/ninja -C $(gettop) -f out/combined-${TARGET_PRODUCT}.ninja -t query ${TARGETS} |
    sed -n '/input/,/output/{/input/b;/output/b;p}' |
    grep -v timestamp$ #filter out timestamp files
}

function update_api() {
  if [[ $(echo "${TARGETS}" | grep -cwe "-h\|--help") -ne 0 ]]; then
    _print_update_api_help_msg
    return 0
  fi

  if [[ -z "${TARGET_PRODUCT}" ]]; then
    _print_TARGET_PRODUCT_help_msg
    return 1
  fi

  # source build/envsetup.sh to get m and gettop
  source "${BASH_SOURCE%/*}/../../../envsetup.sh"
  # get update api script names to build
  UPDATE_API_SCRIPTS=$(_get_update_api_script_names)
  # Run m to generate update_api scripts
  m ${UPDATE_API_SCRIPTS}
  # Run update_api scripts to update source tree
  for update_api_script in ${UPDATE_API_SCRIPTS}; do
    echo Running script ${update_api_script}
    (cd $(gettop) && ${update_api_script})
  done
}

update_api
