#!/bin/bash

#ensure that OUT_DIR is set by soong_ui
if [ -z "${OUT_DIR}" ]; then
  #get repo root
  cd "$(dirname $0)/../../.."
  TOP="${PWD}"
  echo "usage: cd ${TOP} && make help"
  exit 1
fi

#make OUT_DIR absolute
OUT_DIR="$(echo ${OUT_DIR} | sed "s|^|${TOP}/|")"

message='
usage: make [-j] [<targets>] [<variable>=<value>...]

Builds the specified Android artifacts.

Ways to specify what to build:
  The general way to specify what to build is to set the environment variables $TARGET_PRODUCT and
    $TARGET_BUILD_VARIANT .

    $TARGET_PRODUCT refers to the device that the created image is intended to be run on. A list of
      valid values is given by "lunch", explained below
    $TARGET_BUILD_VARIANT is one of "user", "userdebug", or "eng", and controls the amount of debugging
      to be added into the generated image

  Usually, users set the product and variant using the following shorthand:

    cd '"${TOP}/build"'
    source envsetup.sh
    lunch [<product>-<variant>] #If you omit the argument, "lunch" will list the choices and prompt.
   
    This also has the side-effect of updating some other variables, including adding to PATH the
      paths of some generated tools, such as adb (adb is a tool that communicates with the device)

    An alternative to setting $TARGET_PRODUCT and $TARGET_BUILD_VARIANT (which you may see be used by
      build servers), is to include a target of PRODUCT-<product>-<variant> .

  A target may also be a file path. For example, out/host/linux-x86/bin/adb .
    Note that when giving a relative file path as a target, that path is interpreted relative to the
    root of the source tree (rather than relative to the current working directory)

  A target may also be any other target defined within a Makefile. Common targets include:
    clean                   (aka clobber) equivalent to rm -rf out/
    checkbuild              Builds every module defined in the source tree
    droid                   Default target

    java                    Build all the java code in the source tree
    native                  Build all the native code in the source tree

    host                    Build all the host code (not to be run on a device) in the source tree
    target                  Build all the target code (to be run on the device) in the source tree

    (java|native)-(host|target)
    (host|target)-(java|native)
                            Build the intersection of the two given arguments

    snod                    Quickly rebuild the system image from built packages
    vnod                    Quickly rebuild the vendor image from built packages
    offline-sdk-docs        Generate the HTML for the developer SDK docs
    doc-comment-check-docs  Check HTML doc links & validity, without generating HTML
    libandroid_runtime      All the JNI framework stuff
    framework               All the java framework stuff
    services                The system server (Java) and friends
    help                    Shows this message

  To view the modules and targets defined in a particular directory, look for:
    files named *.mk (most commonly Android.mk)
      these files are defined in Make syntax
    files named Android.bp
      these files are defined in Blueprint syntax

  To obtain the full (extremely large) compiled list of targets, run the build once, and then,
    see '"${OUT_DIR}/build-<product>*.ninja"'
    and '"${OUT_DIR}/soong/build.ninja"'

Targets that adjust an existing build:
  showcommands              Display the individual commands run to implement the build
  dist                      Copy into ${DIST_DIR} the portion of the build that must be distributed

Flags
  -j <N>                    Runs <N> processes at once
  -j                        Autodetects the number of processes to run at once, and runs that many

Variables
  Variables can either be set in the surrounding shell environment or can be passed as command-line
    arguments. For example:
      I_AM_A_SHELL_VAR=1
      I_AM_ANOTHER_SHELL_VAR=2 make droid I_AM_AN_ARGUMENT=3
  Here are some common variables and their meanings:
    TARGET_PRODUCT          The <product> to build
    TARGET_BUILD_VARIANT    The <variant> to build
    DIST_DIR                The directory in which to place the distribution artifacts.
    OUT_DIR                 The directory in which to place non-distribution artifacts.


  To see a fuller list of variables, run the build once, and then,
    see '"${OUT_DIR}"'/.kati_stamp-*
    
'

echo "$message"
