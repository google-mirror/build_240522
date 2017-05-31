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
  The common way to specify what to build is to set that information in the environment via:

    cd '"${TOP}/build"'
    source envsetup.sh #Sets up the shell environment. Run "hmm" after sourcing it for more info.
    lunch [<product>-<variant>] #Selects the build target.
    make [<options>] [<targets>] [<variable>=<value>...] #Invokes the configured build.

      <product> refers to the device that the created image is intended to be run on.
        This gets saved in the shell environment as $TARGET_PRODUCT by "lunch".
      <variant> is one of "user", "userdebug", or "eng", and controls the amount of debugging to be
        added into the generated image.
        This gets saved in the shell environment as $TARGET_BUILD_VARIANT by "lunch".

      If you omit the argument to "lunch", it will list the choices and prompt.

    An alternative to setting $TARGET_PRODUCT and $TARGET_BUILD_VARIANT is to execute:

      make PRODUCT-<product>-<variant> # you may see this usage in build servers

    The <options>, <targets>, and <variable assignments> are all optional.
      If no targets are specified, the build system will build a system image for the configured
        product and variant.

  A target may be a file path. For example, out/host/linux-x86/bin/adb .
    Note that when giving a relative file path as a target, that path is interpreted relative to the
    root of the source tree (rather than relative to the current working directory).

  A target may also be any other target defined within a Makefile. Common targets include:
    clean                   (aka clobber) equivalent to rm -rf out/
    checkbuild              Build every module defined in the source tree
    droid                   Default target
    nothing                 Do not build anything, just parse and validate the build structure

    java                    Build all the java code in the source tree
    native                  Build all the native code in the source tree

    host                    Build all the host code (not to be run on a device) in the source tree
    target                  Build all the target code (to be run on the device) in the source tree

    (java|native)-(host|target)
    (host|target)-(java|native)
                            Build the intersection of the two given arguments

    snod                    Quickly rebuild the system image from built packages
                            Stands for "System, NO Dependencies"
    vnod                    Quickly rebuild the vendor image from built packages
                            Stands for "Vendor, NO Dependencies"
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

  For now, the full (extremely large) compiled list of targets can be found (after running the build
    once), split among these two files:

    '"${OUT_DIR}/build-<product>*.ninja"'
    '"${OUT_DIR}/soong/build.ninja"'

    If you find yourself interacting with these files, you are encouraged to provide a tool to browse
    them and to mention the tool here.

Targets that adjust an existing build:
  showcommands              Display the individual commands run to implement the build
  dist                      Copy into ${DIST_DIR} the portion of the build that must be distributed

Flags
  -j <N>                    Run <N> processes at once
  -j                        Autodetect the number of processes to run at once, and run that many

Variables
  Variables can either be set in the surrounding shell environment or can be passed as command-line
    arguments. For example:
      export I_AM_A_SHELL_VAR=1
      I_AM_ANOTHER_SHELL_VAR=2 make droid I_AM_AN_ARGUMENT=3
  Here are some common variables and their meanings:
    TARGET_PRODUCT          The <product> to build #as described above
    TARGET_BUILD_VARIANT    The <variant> to build #as described above
    DIST_DIR                The directory in which to place the distribution artifacts.
    OUT_DIR                 The directory in which to place non-distribution artifacts.

  There does not yet exist a good method by which to discover the full list of supported variables.
  Please mention it here when it does.

'

echo "$message"
