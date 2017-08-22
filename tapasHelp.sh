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
Unbundled apps therefore have been validated to be suitable for installing on
any API-compatible device. A successfully built unbundled app will avoid
inadvertently using internal APIs, and will also have its dex file in its APK.

The names <App1> <App2> ... should match LOCAL_PACKAGE_NAME as defined in an
Android.mk

The usage of the other arguments matches that of the rest of the platform
build system and can be found by running `m help`'

echo "$message"
