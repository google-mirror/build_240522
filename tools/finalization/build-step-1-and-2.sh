#!/bin/bash

set -ex

function finalize_main_step12() {
    local top="$(dirname "$0")"/../..

    # default target to modify tree and build SDK
    local m="$top/build/soong/soong_ui.bash --make-mode TARGET_PRODUCT=aosp_arm64 TARGET_BUILD_VARIANT=userdebug"

    # SDK codename -> int
    source $top/build/make/tools/finalization/finalize-aidl-vndk-sdk-resources.sh

    # Platform/Mainline SDKs build and move to prebuilts
    source $top/build/make/tools/finalization/localonly-finalize-mainline-sdk.sh

    # REL
    source $top/build/make/tools/finalization/finalize-sdk-rel.sh

    # This command tests:
    #   The release state for AIDL.
    #   ABI difference between user and userdebug builds.
    #   Resource/SDK finalization.
    AIDL_FROZEN_REL=true $m
}

finalize_main_step12

