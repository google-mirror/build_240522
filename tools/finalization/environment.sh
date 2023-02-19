#!/bin/bash

set -ex

export FINAL_BUG_ID='0'

export FINAL_PLATFORM_CODENAME='UpsideDownCake'
export CURRENT_PLATFORM_CODENAME='VanillaIceCream'
export FINAL_PLATFORM_CODENAME_JAVA='UPSIDE_DOWN_CAKE'
export FINAL_PLATFORM_SDK_VERSION='34'
export FINAL_PLATFORM_VERSION='14'

export FINAL_BUILD_PREFIX='UP1A'

export FINAL_MAINLINE_EXTENSION='6'
export FINAL_MAINLINE_SDK_COMMIT_MESSAGE=''
# The build has to include both Platform SDK target (sdk) and Modules SDK target (mainline_modules_sdks)
export FINAL_SDK_BUILD_ID=0
