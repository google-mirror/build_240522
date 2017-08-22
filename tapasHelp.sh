#!/bin/bash

# locate some directories
cd "$(dirname $0)"
SCRIPT_DIR="${PWD}"
cd ../..
TOP="${PWD}"

message='usage: tapas [<App1> <App2> ...] [arm|x86|mips|armv5|arm64|x86_64|mips64] [eng|userdebug|user]

tapas selects unbundled apps to be built by the Android build system.

Unbundled apps are distinguished from bundled apps, because the former are
built against the same SDK that is made available to third-party developers,
and avoids the risk of such an app inadvertently accessing an internal API.

The names <App1> <App2> ... should match LOCAL_PACKAGE_NAME as defined in an
Android.mk

The usage of the other arguments matches that of the rest of the platform
buidl system and can be found by running `m help`'

echo "$message"
