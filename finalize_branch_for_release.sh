#!/bin/bash

set -ex

# This script is WIP and only finalizes part of the Android branch for release.
# The full process can be found at (INTERNAL) go/android-sdk-finalization.

# VNDK snapshot (TODO)
# AIDL
AIDL_TRANSITIVE_FREEZE=true $ANDROID_BUILD_TOP/build/soong/soong_ui.bash --make-mode aidl-freeze-api


# SDK snapshots (TODO)
# Update references in the codebase to new API version (TODO)
# ...

# Test

# TODO(b/229413853): test while simulating 'rel' for more requirements AIDL_FROZEN_REL=true
$ANDROID_BUILD_TOP/build/soong/soong_ui.bash --make-mode nothing
