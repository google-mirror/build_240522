#!/bin/bash

set -ex

function finalize_vintf_resources() {
    local top="$(dirname "$0")"/../../../..
    source $top/build/make/tools/finalization/environment.sh

    # TODO(b/314010764): finalize LL_NDK

    # system/sepolicy
    board_api_level=$(TARGET_RELEASE=next $top/build/soong/soong_ui.bash --dumpvar-mode BOARD_API_LEVEL)
    system/sepolicy/tools/finalize-vintf-resources.sh "$top" "$board_api_level"

    # pre-finalization build target (trunk)
    local aidl_m="$top/build/soong/soong_ui.bash --make-mode TARGET_PRODUCT=aosp_arm64 TARGET_RELEASE=trunk TARGET_BUILD_VARIANT=userdebug DIST_DIR=out/dist"
    AIDL_TRANSITIVE_FREEZE=true $aidl_m aidl-freeze-api

    # build/make
    cp "$top/build/make/target/product/gsi/current.txt" "$top/build/make/target/product/gsi/$FINAL_PLATFORM_SDK_VERSION.txt"
}

finalize_vintf_resources

