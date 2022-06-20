#!/bin/bash

# locate some directories
cd "$(dirname $0)"
SCRIPT_DIR="${PWD}"
cd ../..
TOP="${PWD}"

<<<<<<< HEAD   (cbb40d Merge "Merge empty history for sparse-8719481-L1140000095509)
message='usage: tapas [<App1> <App2> ...] [arm|x86|mips|arm64|x86_64|mips64] [eng|userdebug|user]
=======
message='usage: tapas [<App1> <App2> ...] [arm|x86|arm64|x86_64] [eng|userdebug|user] [devkeys]
>>>>>>> BRANCH (eedaac Merge "Version bump to TKB1.220618.001.A1 [core/build_id.mk])

tapas selects individual apps to be built by the Android build system. Unlike
"lunch", "tapas" does not request the building of images for a device.
Additionally, an app built with "tapas" will have its dex file inside its apk,
which should cause it to be suitable for installing on any api-compatible
device. In other words, "tapas" configures the build of unbundled apps.

The names <App1> <App2> ... should match LOCAL_PACKAGE_NAME as defined in an
Android.mk

The usage of the other arguments matches that of the rest of the platform
build system and can be found by running `m help`'

echo "$message"
