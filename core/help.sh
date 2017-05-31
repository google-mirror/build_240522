#!/bin/bash

#get repo root
cd "$(dirname $0)/../../.."
TOP="$PWD"

#find out dir
OUT_DIR="${OUT_DIR}"
if [ "${OUT_DIR}" == "" -o "${OUT_DIR}" == "out" ]; then
  OUT_DIR="${TOP}/out"
fi

message='
usage: make [-j] [<targets>] [<variable>=<value>...]

Builds the specified Android artifacts.

Targets that specify what to build:
  The most basic targets are of the form PRODUCT-<product>-<variant> .
    <product> refers to the device that the created image is intended to be run on.
    <variant> is one of "user", "userdebug", or "eng", and controls the amount of debugging to be added
      into the generated image
    For example, "PRODUCT-aosp_arm-userdebug" specifies to do a userdebug build to run on an
      aosp_arm device. Note that "aosp_arm" is special in that it does not correspond to a real
      physical device and is instead a generic hypothetical device for testing.
    In practice, few users directly type "PRODUCT-<product>-<variant>" as an argument.
      Most users save their <product> and <variant> selections into alternate environment variables,
      to save typing. The way to do this is:

      cd '"${TOP}/build"'
      source envsetup.sh
      lunch <product>-<variant>

      This also has the side-effect of updating your PATH to include the paths of some generated
      tools, such as adb, which can communicate with the device

  A target may also be a filename. For example, out/host/linux-x86/bin/adb

  A target may also be any other target defined within a Makefile. Common targets include:
    clean                   (aka clobber) equivalent to rm -rf out/
    checkbuild              Builds every module defined in the source tree
    droid                   Default target
    snod                    Quickly rebuild the system image from built packages
    vnod                    Quickly rebuild the vendor image from built packages
    offline-sdk-docs        Generate the HTML for the developer SDK docs
    doc-comment-check-docs  Check HTML doc links & validity, without generating HTML
    libandroid_runtime      All the JNI framework stuff
    framework               All the java framework stuff
    services                The system server (Java) and friends
    help                    Shows this message

  To see a full list of targets, run the build once, and then,
    see '"${OUT_DIR}/build-<product>-<variant>.ninja"'

Targets that adjust an existing build (these targets should not be the only target) :
  showcommands              Display the individual commands run to implement the build
  dist                      Copy into ${DIST_DIR} (or ${TOP}/out if unset) the portion of the build that
                            is needed for distributing

Flags
  Flags adjust the build
    -j                      Enable running multiple processes at once

Variables
  Make variables can either be set in the surrounding shell environment or can be passed as
  command-line arguments. Here are some common variables and their meanings:
    TARGET_PRODUCT          The <product> to build
    TARGET_BUILD_VARIANT    The <variant> to build
    DIST_DIR                The directory in which to place the distribution artifacts.
    OUT_DIR                 The directory in which to place non-distribution artifacts. Defaults to
                            ${TOP}/out if not set


  To see a fuller list of variables, run the build once, and then,
    see '"${OUT_DIR}"'/soong/make_vars-<product>-<variant>.mk
'

echo "$message"
