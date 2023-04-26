#!/bin/bash

set -ex

export FINAL_BUG_ID='0'

export FINAL_PLATFORM_CODENAME='VanillaIceCream'
export CURRENT_PLATFORM_CODENAME='VanillaIceCream'
export FINAL_PLATFORM_CODENAME_JAVA='VANILLA_ICE_CREAM'
export FINAL_PLATFORM_SDK_VERSION='35'
export FINAL_PLATFORM_VERSION='15'

export FINAL_BUILD_PREFIX='VP1A'

export FINAL_MAINLINE_EXTENSION='8'

# Options:
# 'unfinalized' - branch is in development state,
# 'sdk' - SDK/API is finalized
# 'rel' - branch is finalized, switched to REL
export FINAL_STATE='unfinalized'
