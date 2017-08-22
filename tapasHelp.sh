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
This reduces the number of projects that are truly required to be checked out
to do a build of an unbundled app. This makes it easier to validate that an
unbundled app does not use internal APIs: it can be tested by invoking a clean
build using the smaller subset of projects, and confirming that the build
succeeded and that the subset of projects did not expose any internal APIs.
An app built with tapas will also have its dex file inside its apk, which
should cause it to be suitable for installing on any api-compatible device.

The names <App1> <App2> ... should match LOCAL_PACKAGE_NAME as defined in an
Android.mk

The usage of the other arguments matches that of the rest of the platform
build system and can be found by running `m help`'

echo "$message"
