<<<<<<< HEAD   (cddf34 Merge "Merge empty history for sparse-9088035-L9400000095648)
=======
#!/bin/bash

set -ex

function finalize_main() {
    local top="$(dirname "$0")"/../..

    # default target to modify tree and build SDK
    local m="$top/build/soong/soong_ui.bash --make-mode TARGET_PRODUCT=aosp_arm64 TARGET_BUILD_VARIANT=userdebug"

    # Build finalization artifacts.
    source $top/build/make/finalize-aidl-vndk-sdk-resources.sh

    # This command tests:
    #   The release state for AIDL.
    #   ABI difference between user and userdebug builds.
    #   Resource/SDK finalization.
    # In the future, we would want to actually turn the branch into the REL
    # state and test with that.
    AIDL_FROZEN_REL=true $m droidcore

    # Build SDK (TODO)
    # lunch sdk...
    # m ...
}

finalize_main

>>>>>>> BRANCH (68e5c5 Merge "Version bump to TKB1.220925.001.A1 [core/build_id.mk])
