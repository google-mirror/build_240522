#!/bin/bash

set -ex

function finalize_main() {
    local top="$(dirname "$0")"/../..

    # default target to modify tree and build SDK
    local m="$top/build/soong/soong_ui.bash --make-mode TARGET_PRODUCT=aosp_arm64 TARGET_BUILD_VARIANT=userdebug"

    # This script is WIP and only finalizes part of the Android branch for release.
    # The full process can be found at (INTERNAL) go/android-sdk-finalization.

    # VNDK snapshot (TODO)
    # SDK snapshots (TODO)
    # Update references in the codebase to new API version (TODO)
    # ...

    # Generate ABI dumps
    ANDROID_BUILD_TOP="$top" \
        $top/development/vndk/tools/header-checker/utils/create_reference_dumps.py \
        -p aosp_arm64 --build-variant user

    AIDL_TRANSITIVE_FREEZE=true $m aidl-freeze-api

    # Update new versions of files. See update-vndk-list.sh (which requires envsetup.sh)
    $m check-vndk-list || \
        { cp $top/out/soong/vndk/vndk.libraries.txt $top/build/make/target/product/gsi/current.txt; }

    # for now, we simulate the release state for AIDL, but in the future, we would want
    # to actually turn the branch into the REL state and test with that
    AIDL_FROZEN_REL=true $m nothing # test build

    # Diff the ABI of user and userdebug builds
    $m droidcore

    # Build SDK (TODO)
    # lunch sdk...
    # m ...
}

finalize_main
