#!/bin/bash
set -eo pipefail

TARGETS=${@:-update-api}

function _print_update_api_help_msg() {
  printf \
"Utililty function to update API in source tree
Usage:
    update_api.sh [targets...]
Examples:
    update_api.sh # update-api is default
    update_api.sh update-api # update all APIs
    update_api.sh ahat-docs-stubs-update-current-api android.car-stubs-docs-update-current-api # update subset of APIs
"
}

function _print_TARGET_PRODUCT_help_msg() {
  printf \
"TARGET_PRODUCT not set. To fix, use one of the following options
1. Rerun lunch in current shell
2. Pass to script as TARGET_PRODUCT=<target_product> update_api.sh [targets...]
"
}

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
  # Run m to generate update_api scripts
  m ${TARGETS}
  # Run update_api scripts to update source tree
  while IFS= read -r update_api_script_name; do
    echo Running script ${update_api_script_name}
    (cd $(gettop) && ${update_api_script_name})
  done < <(_get_update_api_script_names)
}

update_api
