#!/bin/bash

# locate some directories
cd "$(dirname $0)"
SCRIPT_DIR="${PWD}"
cd ../../..
TOP="${PWD}"

message='The basic Android build process is:

cd '"${TOP}"'
source build/envsetup.sh    # Add "lunch" (and other utilities and variables)
                            # to the shell environment.
lunch [<product>-<variant>] # Choose the device to target.
make -j [<targets>]         # Execute the configured build.

See '"${SCRIPT_DIR}"'/README.txt for more info about build usage and concepts.

See '"${SCRIPT_DIR}"'/common_targets.txt for a list of common build targets.
'

echo "$message"
