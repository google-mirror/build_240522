#!/bin/bash

# locate some directories
cd "$(dirname $0)"
SCRIPT_DIR="${PWD}"
cd ../..
TOP="${PWD}"

message='usage: tapas [<App1> <App2> ...] [arm|x86|mips|armv5|arm64|x86_64|mips64] [eng|userdebug|user]

tapas selects unbundled apps to be built by the Android build system.

Unbundled apps are distinguished from bundled apps, because the former are
built against the same SDK that is made available to third-party developers.
Because a successfully built unbundled app does not need to be built against
the platform internals directly, it makes it easier to avoid using internal
APIs. An unbundled app will also have its dex file in its APK, which is
required to run on an API-compatible device.

The names <App1> <App2> ... should match LOCAL_PACKAGE_NAME as defined in an
Android.mk

The usage of the other arguments matches that of the rest of the platform
build system and can be found by running `m help`'

echo "$message"
