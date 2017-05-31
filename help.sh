#!/bin/bash

# locate some directories
cd "$(dirname $0)"
SCRIPT_DIR="${PWD}"
cd ../../..
TOP="${PWD}"

message='The basic Android build process is:

cd '"${TOP}/build"' # Get the source code.
source envsetup.sh  # Add "lunch" (and other utilities and variables) to the shell environment.
lunch               # Choose the device to target.
make                # Execute the configured build.

See '"${SCRIPT_DIR}"'/README.md for more information.
'

echo "$message"
