function hmm() {
cat <<EOF

Run "m help" for help with the build system itself.

Invoke ". build/envsetup.sh" from your shell to add the following functions to your environment:
- lunch:      lunch <product_name>-<build_variant>
              Selects <product_name> as the product to build, and <build_variant> as the variant to
              build, and stores those selections in the environment to be read by subsequent
              invocations of 'm' etc.
- tapas:      tapas [<App1> <App2> ...] [arm|x86|mips|arm64|x86_64|mips64] [eng|userdebug|user]
- croot:      Changes directory to the top of the tree.
- m:          Makes from the top of the tree.
- mm:         Builds all of the modules in the current directory, but not their dependencies.
- mmm:        Builds all of the modules in the supplied directories, but not their dependencies.
              To limit the modules being built use the syntax: mmm dir/:target1,target2.
- mma:        Builds all of the modules in the current directory, and their dependencies.
- mmma:       Builds all of the modules in the supplied directories, and their dependencies.
- provision:  Flash device with all required partitions. Options will be passed on to fastboot.
- cgrep:      Greps on all local C/C++ files.
- ggrep:      Greps on all local Gradle files.
- jgrep:      Greps on all local Java files.
- resgrep:    Greps on all local res/*.xml files.
- mangrep:    Greps on all local AndroidManifest.xml files.
- mgrep:      Greps on all local Makefiles files.
- sepgrep:    Greps on all local sepolicy files.
- sgrep:      Greps on all local source files.
- godir:      Go to the directory containing a file.
- allmod:     List all modules.
- gomod:      Go to the directory containing a module.
- refreshmod: Refresh list of modules for allmod/gomod.

Environment options:
- SANITIZE_HOST: Set to 'true' to use ASAN for all host modules. Note that
                 ASAN_OPTIONS=detect_leaks=0 will be set by default until the
                 build is leak-check clean.

Look at the source to view more functions. The complete list is:
EOF
    local T=$(gettop)
    local A=""
    local i
    for i in `cat $T/build/envsetup.sh | sed -n "/^[[:blank:]]*function /s/function \([a-z_]*\).*/\1/p" | sort | uniq`; do
      A="$A $i"
    done
    echo $A
}

# Get all the build variables needed by this script in a single call to the build system.
function build_build_var_cache()
{
    local T=$(gettop)
    # Grep out the variable names from the script.
    cached_vars=(`cat $T/build/envsetup.sh | tr '()' '  ' | awk '{for(i=1;i<=NF;i++) if($i~/get_build_var/) print $(i+1)}' | sort -u | tr '\n' ' '`)
    cached_abs_vars=(`cat $T/build/envsetup.sh | tr '()' '  ' | awk '{for(i=1;i<=NF;i++) if($i~/get_abs_build_var/) print $(i+1)}' | sort -u | tr '\n' ' '`)
    # Call the build system to dump the "<val>=<value>" pairs as a shell script.
    build_dicts_script=`\builtin cd $T; build/soong/soong_ui.bash --dumpvars-mode \
                        --vars="${cached_vars[*]}" \
                        --abs-vars="${cached_abs_vars[*]}" \
                        --var-prefix=var_cache_ \
                        --abs-var-prefix=abs_var_cache_`
    local ret=$?
    if [ $ret -ne 0 ]
    then
        unset build_dicts_script
        return $ret
    fi
    # Execute the script to store the "<val>=<value>" pairs as shell variables.
    eval "$build_dicts_script"
    ret=$?
    unset build_dicts_script
    if [ $ret -ne 0 ]
    then
        return $ret
    fi
    BUILD_VAR_CACHE_READY="true"
}

# Delete the build var cache, so that we can still call into the build system
# to get build variables not listed in this script.
function destroy_build_var_cache()
{
    unset BUILD_VAR_CACHE_READY
    local v
    for v in $cached_vars; do
      unset var_cache_$v
    done
    unset cached_vars
    for v in $cached_abs_vars; do
      unset abs_var_cache_$v
    done
    unset cached_abs_vars
}

# Get the value of a build variable as an absolute path.
function get_abs_build_var()
{
    if [ "$BUILD_VAR_CACHE_READY" = "true" ]
    then
        eval "echo \"\${abs_var_cache_$1}\""
    return
    fi

    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP." >&2
        return
    fi
    (\cd $T; build/soong/soong_ui.bash --dumpvar-mode --abs $1)
}

# Get the exact value of a build variable.
function get_build_var()
{
    if [ "$BUILD_VAR_CACHE_READY" = "true" ]
    then
        eval "echo \"\${var_cache_$1}\""
    return
    fi

    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP." >&2
        return
    fi
    (\cd $T; build/soong/soong_ui.bash --dumpvar-mode $1)
}

# check to see if the supplied product is one we can build
function check_product()
{
    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP." >&2
        return
    fi
        TARGET_PRODUCT=$1 \
        TARGET_BUILD_VARIANT= \
        TARGET_BUILD_TYPE= \
        TARGET_BUILD_APPS= \
        get_build_var TARGET_DEVICE > /dev/null
    # hide successful answers, but allow the errors to show
}

VARIANT_CHOICES=(user userdebug eng)

# check to see if the supplied variant is valid
function check_variant()
{
    local v
    for v in ${VARIANT_CHOICES[@]}
    do
        if [ "$v" = "$1" ]
        then
            return 0
        fi
    done
    return 1
}

function setpaths()
{
    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP."
        return
    fi

    ##################################################################
    #                                                                #
    #              Read me before you modify this code               #
    #                                                                #
    #   This function sets ANDROID_BUILD_PATHS to what it is adding  #
    #   to PATH, and the next time it is run, it removes that from   #
    #   PATH.  This is required so lunch can be run more than once   #
    #   and still have working paths.                                #
    #                                                                #
    ##################################################################

    # Note: on windows/cygwin, ANDROID_BUILD_PATHS will contain spaces
    # due to "C:\Program Files" being in the path.

    # out with the old
    if [ -n "$ANDROID_BUILD_PATHS" ] ; then
        export PATH=${PATH/$ANDROID_BUILD_PATHS/}
    fi
    if [ -n "$ANDROID_PRE_BUILD_PATHS" ] ; then
        export PATH=${PATH/$ANDROID_PRE_BUILD_PATHS/}
        # strip leading ':', if any
        export PATH=${PATH/:%/}
    fi

    # and in with the new
    local prebuiltdir=$(getprebuilt)
    local gccprebuiltdir=$(get_abs_build_var ANDROID_GCC_PREBUILTS)

    # defined in core/config.mk
    local targetgccversion=$(get_build_var TARGET_GCC_VERSION)
    local targetgccversion2=$(get_build_var 2ND_TARGET_GCC_VERSION)
    export TARGET_GCC_VERSION=$targetgccversion

    # The gcc toolchain does not exists for windows/cygwin. In this case, do not reference it.
    export ANDROID_TOOLCHAIN=
    export ANDROID_TOOLCHAIN_2ND_ARCH=
    local ARCH=$(get_build_var TARGET_ARCH)
    local toolchaindir toolchaindir2=
    case $ARCH in
        x86) toolchaindir=x86/x86_64-linux-android-$targetgccversion/bin
            ;;
        x86_64) toolchaindir=x86/x86_64-linux-android-$targetgccversion/bin
            ;;
        arm) toolchaindir=arm/arm-linux-androideabi-$targetgccversion/bin
            ;;
        arm64) toolchaindir=aarch64/aarch64-linux-android-$targetgccversion/bin;
               toolchaindir2=arm/arm-linux-androideabi-$targetgccversion2/bin
            ;;
        mips|mips64) toolchaindir=mips/mips64el-linux-android-$targetgccversion/bin
            ;;
        *)
            echo "Can't find toolchain for unknown architecture: $ARCH"
            toolchaindir=xxxxxxxxx
            ;;
    esac
    if [ -d "$gccprebuiltdir/$toolchaindir" ]; then
        export ANDROID_TOOLCHAIN=$gccprebuiltdir/$toolchaindir
    fi

    if [ "$toolchaindir2" -a -d "$gccprebuiltdir/$toolchaindir2" ]; then
        export ANDROID_TOOLCHAIN_2ND_ARCH=$gccprebuiltdir/$toolchaindir2
    fi

    export ANDROID_DEV_SCRIPTS=$T/development/scripts:$T/prebuilts/devtools/tools:$T/external/selinux/prebuilts/bin

    # add kernel specific binaries
    case $(uname -s) in
        Linux)
            export ANDROID_DEV_SCRIPTS=$ANDROID_DEV_SCRIPTS:$T/prebuilts/misc/linux-x86/dtc:$T/prebuilts/misc/linux-x86/libufdt
            ;;
        *)
            ;;
    esac

    ANDROID_BUILD_PATHS=$(get_build_var ANDROID_BUILD_PATHS):$ANDROID_TOOLCHAIN
    if [ -n "$ANDROID_TOOLCHAIN_2ND_ARCH" ]; then
        ANDROID_BUILD_PATHS=$ANDROID_BUILD_PATHS:$ANDROID_TOOLCHAIN_2ND_ARCH
    fi
    ANDROID_BUILD_PATHS=$ANDROID_BUILD_PATHS:$ANDROID_DEV_SCRIPTS:
    export ANDROID_BUILD_PATHS

    # If prebuilts/android-emulator/<system>/ exists, prepend it to our PATH
    # to ensure that the corresponding 'emulator' binaries are used.
    case $(uname -s) in
        Darwin)
            ANDROID_EMULATOR_PREBUILTS=$T/prebuilts/android-emulator/darwin-x86_64
            ;;
        Linux)
            ANDROID_EMULATOR_PREBUILTS=$T/prebuilts/android-emulator/linux-x86_64
            ;;
        *)
            ANDROID_EMULATOR_PREBUILTS=
            ;;
    esac
    if [ -n "$ANDROID_EMULATOR_PREBUILTS" -a -d "$ANDROID_EMULATOR_PREBUILTS" ]; then
        ANDROID_BUILD_PATHS=$ANDROID_BUILD_PATHS$ANDROID_EMULATOR_PREBUILTS:
        export ANDROID_EMULATOR_PREBUILTS
    fi

    export PATH=$ANDROID_BUILD_PATHS$PATH

    # out with the duplicate old
    if [ -n $ANDROID_PYTHONPATH ]; then
        export PYTHONPATH=${PYTHONPATH//$ANDROID_PYTHONPATH/}
    fi
    # and in with the new
    export ANDROID_PYTHONPATH=$T/development/python-packages:
    export PYTHONPATH=$ANDROID_PYTHONPATH$PYTHONPATH

    export ANDROID_JAVA_HOME=$(get_abs_build_var ANDROID_JAVA_HOME)
    export JAVA_HOME=$ANDROID_JAVA_HOME
    export ANDROID_JAVA_TOOLCHAIN=$(get_abs_build_var ANDROID_JAVA_TOOLCHAIN)
    export ANDROID_PRE_BUILD_PATHS=$ANDROID_JAVA_TOOLCHAIN:
    export PATH=$ANDROID_PRE_BUILD_PATHS$PATH

    unset ANDROID_PRODUCT_OUT
    export ANDROID_PRODUCT_OUT=$(get_abs_build_var PRODUCT_OUT)
    export OUT=$ANDROID_PRODUCT_OUT

    unset ANDROID_HOST_OUT
    export ANDROID_HOST_OUT=$(get_abs_build_var HOST_OUT)

    unset ANDROID_HOST_OUT_TESTCASES
    export ANDROID_HOST_OUT_TESTCASES=$(get_abs_build_var HOST_OUT_TESTCASES)

    unset ANDROID_TARGET_OUT_TESTCASES
    export ANDROID_TARGET_OUT_TESTCASES=$(get_abs_build_var TARGET_OUT_TESTCASES)

    # needed for building linux on MacOS
    # TODO: fix the path
    #export HOST_EXTRACFLAGS="-I "$T/system/kernel_headers/host_include
}

function printconfig()
{
    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP." >&2
        return
    fi
    get_build_var report_config
}

function set_stuff_for_environment()
{
    setpaths
    set_sequence_number

    export ANDROID_BUILD_TOP=$(gettop)
    # With this environment variable new GCC can apply colors to warnings/errors
    export GCC_COLORS='error=01;31:warning=01;35:note=01;36:caret=01;32:locus=01:quote=01'
    export ASAN_OPTIONS=detect_leaks=0
}

function set_sequence_number()
{
    export BUILD_ENV_SEQUENCE_NUMBER=13
}

# Takes a command name, and check if it's in ENVSETUP_NO_COMPLETION or not.
function should_add_completion() {
    local cmd="$(basename $1| sed 's/_completion//' |sed 's/\.\(.*\)*sh$//')"
    case :"$ENVSETUP_NO_COMPLETION": in
        *:"$cmd":*)
            return 1
            ;;
    esac
    return 0
}

function addcompletions()
{
    local T dir f

    # Keep us from trying to run in something that's neither bash nor zsh.
    if [ -z "$BASH_VERSION" -a -z "$ZSH_VERSION" ]; then
        return
    fi

    # Keep us from trying to run in bash that's too old.
    if [ -n "$BASH_VERSION" -a ${BASH_VERSINFO[0]} -lt 3 ]; then
        return
    fi

    local completion_files=(
      system/core/adb/adb.bash
      system/core/fastboot/fastboot.bash
      tools/tradefederation/core/atest/atest_completion.sh
    )
    # Completion can be disabled selectively to allow users to use non-standard completion.
    # e.g.
    # ENVSETUP_NO_COMPLETION=adb # -> disable adb completion
    # ENVSETUP_NO_COMPLETION=adb:bit # -> disable adb and bit completion
    for f in ${completion_files[*]}; do
        if [ -f "$f" ] && should_add_completion "$f"; then
            . $f
        fi
    done

    if should_add_completion bit ; then
        complete -C "bit --tab" bit
    fi
    complete -F _lunch lunch

    complete -F _complete_android_module_names gomod
    complete -F _complete_android_module_names m
}

function choosetype()
{
    echo "Build type choices are:"
    echo "     1. release"
    echo "     2. debug"
    echo

    local DEFAULT_NUM DEFAULT_VALUE
    DEFAULT_NUM=1
    DEFAULT_VALUE=release

    export TARGET_BUILD_TYPE=
    local ANSWER
    while [ -z $TARGET_BUILD_TYPE ]
    do
        echo -n "Which would you like? ["$DEFAULT_NUM"] "
        if [ -z "$1" ] ; then
            read ANSWER
        else
            echo $1
            ANSWER=$1
        fi
        case $ANSWER in
        "")
            export TARGET_BUILD_TYPE=$DEFAULT_VALUE
            ;;
        1)
            export TARGET_BUILD_TYPE=release
            ;;
        release)
            export TARGET_BUILD_TYPE=release
            ;;
        2)
            export TARGET_BUILD_TYPE=debug
            ;;
        debug)
            export TARGET_BUILD_TYPE=debug
            ;;
        *)
            echo
            echo "I didn't understand your response.  Please try again."
            echo
            ;;
        esac
        if [ -n "$1" ] ; then
            break
        fi
    done

    build_build_var_cache
    set_stuff_for_environment
    destroy_build_var_cache
}

#
# This function isn't really right:  It chooses a TARGET_PRODUCT
# based on the list of boards.  Usually, that gets you something
# that kinda works with a generic product, but really, you should
# pick a product by name.
#
function chooseproduct()
{
    local default_value
    if [ "x$TARGET_PRODUCT" != x ] ; then
        default_value=$TARGET_PRODUCT
    else
        default_value=aosp_arm
    fi

    export TARGET_BUILD_APPS=
    export TARGET_PRODUCT=
    local ANSWER
    while [ -z "$TARGET_PRODUCT" ]
    do
        echo -n "Which product would you like? [$default_value] "
        if [ -z "$1" ] ; then
            read ANSWER
        else
            echo $1
            ANSWER=$1
        fi

        if [ -z "$ANSWER" ] ; then
            export TARGET_PRODUCT=$default_value
        else
            if check_product $ANSWER
            then
                export TARGET_PRODUCT=$ANSWER
            else
                echo "** Not a valid product: $ANSWER"
            fi
        fi
        if [ -n "$1" ] ; then
            break
        fi
    done

    build_build_var_cache
    set_stuff_for_environment
    destroy_build_var_cache
}

function choosevariant()
{
    echo "Variant choices are:"
    local index=1
    local v
    for v in ${VARIANT_CHOICES[@]}
    do
        # The product name is the name of the directory containing
        # the makefile we found, above.
        echo "     $index. $v"
        index=$(($index+1))
    done

    local default_value=eng
    local ANSWER

    export TARGET_BUILD_VARIANT=
    while [ -z "$TARGET_BUILD_VARIANT" ]
    do
        echo -n "Which would you like? [$default_value] "
        if [ -z "$1" ] ; then
            read ANSWER
        else
            echo $1
            ANSWER=$1
        fi

        if [ -z "$ANSWER" ] ; then
            export TARGET_BUILD_VARIANT=$default_value
        elif (echo -n $ANSWER | grep -q -e "^[0-9][0-9]*$") ; then
            if [ "$ANSWER" -le "${#VARIANT_CHOICES[@]}" ] ; then
                export TARGET_BUILD_VARIANT=${VARIANT_CHOICES[$(($ANSWER-1))]}
            fi
        else
            if check_variant $ANSWER
            then
                export TARGET_BUILD_VARIANT=$ANSWER
            else
                echo "** Not a valid variant: $ANSWER"
            fi
        fi
        if [ -n "$1" ] ; then
            break
        fi
    done
}

function choosecombo()
{
    choosetype $1

    echo
    echo
    chooseproduct $2

    echo
    echo
    choosevariant $3

    echo
    build_build_var_cache
    set_stuff_for_environment
    printconfig
    destroy_build_var_cache
}

# Clear this variable.  It will be built up again when the vendorsetup.sh
# files are included at the end of this file.
unset LUNCH_MENU_CHOICES
function add_lunch_combo()
{
    local new_combo=$1
    local c
    for c in ${LUNCH_MENU_CHOICES[@]} ; do
        if [ "$new_combo" = "$c" ] ; then
            return
        fi
    done
    LUNCH_MENU_CHOICES=(${LUNCH_MENU_CHOICES[@]} $new_combo)
}

function print_lunch_menu()
{
    local uname=$(uname)
    echo
    echo "You're building on" $uname
    echo
    echo "Lunch menu... pick a combo:"

    local i=1
    local choice
    for choice in $(TARGET_BUILD_APPS= LUNCH_MENU_CHOICES="${LUNCH_MENU_CHOICES[@]}" get_build_var COMMON_LUNCH_CHOICES)
    do
        echo "     $i. $choice"
        i=$(($i+1))
    done

    echo
}

function lunch()
{
    local answer

    if [ "$1" ] ; then
        answer=$1
    else
        print_lunch_menu
        echo -n "Which would you like? [aosp_arm-eng] "
        read answer
    fi

    local selection=

    if [ -z "$answer" ]
    then
        selection=aosp_arm-eng
    elif (echo -n $answer | grep -q -e "^[0-9][0-9]*$")
    then
        local choices=($(TARGET_BUILD_APPS= LUNCH_MENU_CHOICES="${LUNCH_MENU_CHOICES[@]}" get_build_var COMMON_LUNCH_CHOICES))
        if [ $answer -le ${#choices[@]} ]
        then
            # array in zsh starts from 1 instead of 0.
            if [ -n "$ZSH_VERSION" ]
            then
                selection=${choices[$(($answer))]}
            else
                selection=${choices[$(($answer-1))]}
            fi
        fi
    else
        selection=$answer
    fi

    export TARGET_BUILD_APPS=

    local product variant_and_version variant version

    product=${selection%%-*} # Trim everything after first dash
    variant_and_version=${selection#*-} # Trim everything up to first dash
    if [ "$variant_and_version" != "$selection" ]; then
        variant=${variant_and_version%%-*}
        if [ "$variant" != "$variant_and_version" ]; then
            version=${variant_and_version#*-}
        fi
    fi

    if [ -z "$product" ]
    then
        echo
        echo "Invalid lunch combo: $selection"
        return 1
    fi

    TARGET_PRODUCT=$product \
    TARGET_BUILD_VARIANT=$variant \
    TARGET_PLATFORM_VERSION=$version \
    build_build_var_cache
    if [ $? -ne 0 ]
    then
        return 1
    fi

    export TARGET_PRODUCT=$(get_build_var TARGET_PRODUCT)
    export TARGET_BUILD_VARIANT=$(get_build_var TARGET_BUILD_VARIANT)
    if [ -n "$version" ]; then
      export TARGET_PLATFORM_VERSION=$(get_build_var TARGET_PLATFORM_VERSION)
    else
      unset TARGET_PLATFORM_VERSION
    fi
    export TARGET_BUILD_TYPE=release

    echo

    set_stuff_for_environment
    printconfig
    destroy_build_var_cache
}

unset COMMON_LUNCH_CHOICES_CACHE
# Tab completion for lunch.
function _lunch()
{
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    if [ -z "$COMMON_LUNCH_CHOICES_CACHE" ]; then
        COMMON_LUNCH_CHOICES_CACHE=$(TARGET_BUILD_APPS= LUNCH_MENU_CHOICES="${LUNCH_MENU_CHOICES[@]}" get_build_var COMMON_LUNCH_CHOICES)
    fi

    COMPREPLY=( $(compgen -W "${COMMON_LUNCH_CHOICES_CACHE}" -- ${cur}) )
    return 0
}

# Configures the build to build unbundled apps.
# Run tapas with one or more app names (from LOCAL_PACKAGE_NAME)
function tapas()
{
    local showHelp="$(echo $* | xargs -n 1 echo | \grep -E '^(help)$' | xargs)"
    local arch="$(echo $* | xargs -n 1 echo | \grep -E '^(arm|x86|mips|arm64|x86_64|mips64)$' | xargs)"
    local variant="$(echo $* | xargs -n 1 echo | \grep -E '^(user|userdebug|eng)$' | xargs)"
    local density="$(echo $* | xargs -n 1 echo | \grep -E '^(ldpi|mdpi|tvdpi|hdpi|xhdpi|xxhdpi|xxxhdpi|alldpi)$' | xargs)"
    local apps="$(echo $* | xargs -n 1 echo | \grep -E -v '^(user|userdebug|eng|arm|x86|mips|arm64|x86_64|mips64|ldpi|mdpi|tvdpi|hdpi|xhdpi|xxhdpi|xxxhdpi|alldpi)$' | xargs)"

    if [ "$showHelp" != "" ]; then
      $(gettop)/build/make/tapasHelp.sh
      return
    fi

    if [ $(echo $arch | wc -w) -gt 1 ]; then
        echo "tapas: Error: Multiple build archs supplied: $arch"
        return
    fi
    if [ $(echo $variant | wc -w) -gt 1 ]; then
        echo "tapas: Error: Multiple build variants supplied: $variant"
        return
    fi
    if [ $(echo $density | wc -w) -gt 1 ]; then
        echo "tapas: Error: Multiple densities supplied: $density"
        return
    fi

    local product=aosp_arm
    case $arch in
      x86)    product=aosp_x86;;
      mips)   product=aosp_mips;;
      arm64)  product=aosp_arm64;;
      x86_64) product=aosp_x86_64;;
      mips64)  product=aosp_mips64;;
    esac
    if [ -z "$variant" ]; then
        variant=eng
    fi
    if [ -z "$apps" ]; then
        apps=all
    fi
    if [ -z "$density" ]; then
        density=alldpi
    fi

    export TARGET_PRODUCT=$product
    export TARGET_BUILD_VARIANT=$variant
    export TARGET_BUILD_DENSITY=$density
    export TARGET_BUILD_TYPE=release
    export TARGET_BUILD_APPS=$apps

    build_build_var_cache
    set_stuff_for_environment
    printconfig
    destroy_build_var_cache
}

function gettop
{
    local TOPFILE=build/make/core/envsetup.mk
    if [ -n "$TOP" -a -f "$TOP/$TOPFILE" ] ; then
        # The following circumlocution ensures we remove symlinks from TOP.
        (cd $TOP; PWD= /bin/pwd)
    else
        if [ -f $TOPFILE ] ; then
            # The following circumlocution (repeated below as well) ensures
            # that we record the true directory name and not one that is
            # faked up with symlink names.
            PWD= /bin/pwd
        else
            local HERE=$PWD
            local T=
            while [ \( ! \( -f $TOPFILE \) \) -a \( $PWD != "/" \) ]; do
                \cd ..
                T=`PWD= /bin/pwd -P`
            done
            \cd $HERE
            if [ -f "$T/$TOPFILE" ]; then
                echo $T
            fi
        fi
    fi
}

function m()
{
    local T=$(gettop)
    if [ "$T" ]; then
        _wrap_build $T/build/soong/soong_ui.bash --make-mode $@
    else
        echo "Couldn't locate the top of the tree.  Try setting TOP."
        return 1
    fi
}

function findmakefile()
{
    local TOPFILE=build/make/core/envsetup.mk
    local HERE=$PWD
    if [ "$1" ]; then
        \cd $1
    fi;
    local T=
    while [ \( ! \( -f $TOPFILE \) \) -a \( $PWD != "/" \) ]; do
        T=`PWD= /bin/pwd`
        if [ -f "$T/Android.mk" -o -f "$T/Android.bp" ]; then
            echo $T/Android.mk
            \cd $HERE
            return
        fi
        \cd ..
    done
    \cd $HERE
    return 1
}

function mm()
{
    local T=$(gettop)
    # If we're sitting in the root of the build tree, just do a
    # normal build.
    if [ -f build/soong/soong_ui.bash ]; then
        _wrap_build $T/build/soong/soong_ui.bash --make-mode $@
    else
        # Find the closest Android.mk file.
        local M=$(findmakefile)
        local MODULES=
        local GET_INSTALL_PATH=
        local ARGS=
        # Remove the path to top as the makefilepath needs to be relative
        local M=`echo $M|sed 's:'$T'/::'`
        if [ ! "$T" ]; then
            echo "Couldn't locate the top of the tree.  Try setting TOP."
            return 1
        elif [ ! "$M" ]; then
            echo "Couldn't locate a makefile from the current directory."
            return 1
        else
            local ARG
            for ARG in $@; do
                case $ARG in
                  GET-INSTALL-PATH) GET_INSTALL_PATH=$ARG;;
                esac
            done
            if [ -n "$GET_INSTALL_PATH" ]; then
              MODULES=
              ARGS=GET-INSTALL-PATH-IN-$(dirname ${M})
              ARGS=${ARGS//\//-}
            else
              MODULES=MODULES-IN-$(dirname ${M})
              # Convert "/" to "-".
              MODULES=${MODULES//\//-}
              ARGS=$@
            fi
            if [ "1" = "${WITH_TIDY_ONLY}" -o "true" = "${WITH_TIDY_ONLY}" ]; then
              MODULES=tidy_only
            fi
            ONE_SHOT_MAKEFILE=$M _wrap_build $T/build/soong/soong_ui.bash --make-mode $MODULES $ARGS
        fi
    fi
}

function mmm()
{
    local T=$(gettop)
    if [ "$T" ]; then
        local MAKEFILE=
        local MODULES=
        local MODULES_IN_PATHS=
        local ARGS=
        local DIR TO_CHOP
        local DIR_MODULES
        local GET_INSTALL_PATH=
        local GET_INSTALL_PATHS=
        local DASH_ARGS=$(echo "$@" | awk -v RS=" " -v ORS=" " '/^-.*$/')
        local DIRS=$(echo "$@" | awk -v RS=" " -v ORS=" " '/^[^-].*$/')
        for DIR in $DIRS ; do
            DIR_MODULES=`echo $DIR | sed -n -e 's/.*:\(.*$\)/\1/p' | sed 's/,/ /'`
            DIR=`echo $DIR | sed -e 's/:.*//' -e 's:/$::'`
            # Remove the leading ./ and trailing / if any exists.
            DIR=${DIR#./}
            DIR=${DIR%/}
            local M
            if [ "$DIR_MODULES" = "" ]; then
                M=$(findmakefile $DIR)
            else
                # Only check the target directory if a module is specified.
                if [ -f $DIR/Android.mk -o -f $DIR/Android.bp ]; then
                    local HERE=$PWD
                    cd $DIR
                    M=`PWD= /bin/pwd`
                    M=$M/Android.mk
                    cd $HERE
                fi
            fi
            if [ "$M" ]; then
                # Remove the path to top as the makefilepath needs to be relative
                local M=`echo $M|sed 's:'$T'/::'`
                if [ "$DIR_MODULES" = "" ]; then
                    MODULES_IN_PATHS="$MODULES_IN_PATHS MODULES-IN-$(dirname ${M})"
                    GET_INSTALL_PATHS="$GET_INSTALL_PATHS GET-INSTALL-PATH-IN-$(dirname ${M})"
                else
                    MODULES="$MODULES $DIR_MODULES"
                fi
                MAKEFILE="$MAKEFILE $M"
            else
                case $DIR in
                  showcommands | snod | dist | *=*) ARGS="$ARGS $DIR";;
                  GET-INSTALL-PATH) GET_INSTALL_PATH=$DIR;;
                  *) if [ -d $DIR ]; then
                         echo "No Android.mk in $DIR.";
                     else
                         echo "Couldn't locate the directory $DIR";
                     fi
                     return 1;;
                esac
            fi
        done
        if [ -n "$GET_INSTALL_PATH" ]; then
          ARGS=${GET_INSTALL_PATHS//\//-}
          MODULES=
          MODULES_IN_PATHS=
        fi
        if [ "1" = "${WITH_TIDY_ONLY}" -o "true" = "${WITH_TIDY_ONLY}" ]; then
          MODULES=tidy_only
          MODULES_IN_PATHS=
        fi
        # Convert "/" to "-".
        MODULES_IN_PATHS=${MODULES_IN_PATHS//\//-}
        ONE_SHOT_MAKEFILE="$MAKEFILE" _wrap_build $T/build/soong/soong_ui.bash --make-mode $DASH_ARGS $MODULES $MODULES_IN_PATHS $ARGS
    else
        echo "Couldn't locate the top of the tree.  Try setting TOP."
        return 1
    fi
}

function mma()
{
  local T=$(gettop)
  if [ -f build/soong/soong_ui.bash ]; then
    _wrap_build $T/build/soong/soong_ui.bash --make-mode $@
  else
    if [ ! "$T" ]; then
      echo "Couldn't locate the top of the tree.  Try setting TOP."
      return 1
    fi
    local M=$(findmakefile || echo $(realpath $PWD)/Android.mk)
    # Remove the path to top as the makefilepath needs to be relative
    local M=`echo $M|sed 's:'$T'/::'`
    local MODULES_IN_PATHS=MODULES-IN-$(dirname ${M})
    # Convert "/" to "-".
    MODULES_IN_PATHS=${MODULES_IN_PATHS//\//-}
    _wrap_build $T/build/soong/soong_ui.bash --make-mode $@ $MODULES_IN_PATHS
  fi
}

function mmma()
{
  local T=$(gettop)
  if [ "$T" ]; then
    local DASH_ARGS=$(echo "$@" | awk -v RS=" " -v ORS=" " '/^-.*$/')
    local DIRS=$(echo "$@" | awk -v RS=" " -v ORS=" " '/^[^-].*$/')
    local MY_PWD=`PWD= /bin/pwd`
    if [ "$MY_PWD" = "$T" ]; then
      MY_PWD=
    else
      MY_PWD=`echo $MY_PWD|sed 's:'$T'/::'`
    fi
    local DIR=
    local MODULES_IN_PATHS=
    local ARGS=
    for DIR in $DIRS ; do
      if [ -d $DIR ]; then
        # Remove the leading ./ and trailing / if any exists.
        DIR=${DIR#./}
        DIR=${DIR%/}
        if [ "$MY_PWD" != "" ]; then
          DIR=$MY_PWD/$DIR
        fi
        MODULES_IN_PATHS="$MODULES_IN_PATHS MODULES-IN-$DIR"
      else
        case $DIR in
          showcommands | snod | dist | *=*) ARGS="$ARGS $DIR";;
          *) echo "Couldn't find directory $DIR"; return 1;;
        esac
      fi
    done
    # Convert "/" to "-".
    MODULES_IN_PATHS=${MODULES_IN_PATHS//\//-}
    _wrap_build $T/build/soong/soong_ui.bash --make-mode $DASH_ARGS $ARGS $MODULES_IN_PATHS
  else
    echo "Couldn't locate the top of the tree.  Try setting TOP."
    return 1
  fi
}

function croot()
{
    local T=$(gettop)
    if [ "$T" ]; then
        if [ "$1" ]; then
            \cd $(gettop)/$1
        else
            \cd $(gettop)
        fi
    else
        echo "Couldn't locate the top of the tree.  Try setting TOP."
    fi
}

function cproj()
{
    local TOPFILE=build/make/core/envsetup.mk
    local HERE=$PWD
    local T=
    while [ \( ! \( -f $TOPFILE \) \) -a \( $PWD != "/" \) ]; do
        T=$PWD
        if [ -f "$T/Android.mk" ]; then
            \cd $T
            return
        fi
        \cd ..
    done
    \cd $HERE
    echo "can't find Android.mk"
}

# simplified version of ps; output in the form
# <pid> <procname>
function qpid() {
    local prepend=''
    local append=''
    if [ "$1" = "--exact" ]; then
        prepend=' '
        append='$'
        shift
    elif [ "$1" = "--help" -o "$1" = "-h" ]; then
        echo "usage: qpid [[--exact] <process name|pid>"
        return 255
    fi

    local EXE="$1"
    if [ "$EXE" ] ; then
        qpid | \grep "$prepend$EXE$append"
    else
        adb shell ps \
            | tr -d '\r' \
            | sed -e 1d -e 's/^[^ ]* *\([0-9]*\).* \([^ ]*\)$/\1 \2/'
    fi
}

# coredump_setup - enable core dumps globally for any process
#                  that has the core-file-size limit set correctly
#
# NOTE: You must call also coredump_enable for a specific process
#       if its core-file-size limit is not set already.
# NOTE: Core dumps are written to ramdisk; they will not survive a reboot!

function coredump_setup()
{
    echo "Getting root...";
    adb root;
    adb wait-for-device;

    echo "Remounting root partition read-write...";
    adb shell mount -w -o remount -t rootfs rootfs;
    sleep 1;
    adb wait-for-device;
    adb shell mkdir -p /cores;
    adb shell mount -t tmpfs tmpfs /cores;
    adb shell chmod 0777 /cores;

    echo "Granting SELinux permission to dump in /cores...";
    adb shell restorecon -R /cores;

    echo "Set core pattern.";
    adb shell 'echo /cores/core.%p > /proc/sys/kernel/core_pattern';

    echo "Done."
}

# coredump_enable - enable core dumps for the specified process
# $1 = PID of process (e.g., $(pid mediaserver))
#
# NOTE: coredump_setup must have been called as well for a core
#       dump to actually be generated.

function coredump_enable()
{
    local PID=$1;
    if [ -z "$PID" ]; then
        printf "Expecting a PID!\n";
        return;
    fi;
    echo "Setting core limit for $PID to infinite...";
    adb shell /system/bin/ulimit -p $PID -c unlimited
}

# core - send SIGV and pull the core for process
# $1 = PID of process (e.g., $(pid mediaserver))
#
# NOTE: coredump_setup must be called once per boot for core dumps to be
#       enabled globally.

function core()
{
    local PID=$1;

    if [ -z "$PID" ]; then
        printf "Expecting a PID!\n";
        return;
    fi;

    local CORENAME=core.$PID;
    local COREPATH=/cores/$CORENAME;
    local SIG=SEGV;

    coredump_enable $1;

    local done=0;
    while [ $(adb shell "[ -d /proc/$PID ] && echo -n yes") ]; do
        printf "\tSending SIG%s to %d...\n" $SIG $PID;
        adb shell kill -$SIG $PID;
        sleep 1;
    done;

    adb shell "while [ ! -f $COREPATH ] ; do echo waiting for $COREPATH to be generated; sleep 1; done"
    echo "Done: core is under $COREPATH on device.";
}

# systemstack - dump the current stack trace of all threads in the system process
# to the usual ANR traces file
function systemstack()
{
    stacks system_server
}

# Read the ELF header from /proc/$PID/exe to determine if the process is
# 64-bit.
function is64bit()
{
    local PID="$1"
    if [ "$PID" ] ; then
        if [[ "$(adb shell cat /proc/$PID/exe | xxd -l 1 -s 4 -p)" -eq "02" ]] ; then
            echo "64"
        else
            echo ""
        fi
    else
        echo ""
    fi
}

case `uname -s` in
    Darwin)
        function sgrep()
        {
            find -E . -name .repo -prune -o -name .git -prune -o  -type f -iregex '.*\.(c|h|cc|cpp|hpp|S|java|xml|sh|mk|aidl|vts)' \
                -exec grep --color -n "$@" {} +
        }

        ;;
    *)
        function sgrep()
        {
            find . -name .repo -prune -o -name .git -prune -o  -type f -iregex '.*\.\(c\|h\|cc\|cpp\|hpp\|S\|java\|xml\|sh\|mk\|aidl\|vts\)' \
                -exec grep --color -n "$@" {} +
        }
        ;;
esac

function gettargetarch
{
    get_build_var TARGET_ARCH
}

function ggrep()
{
    find . -name .repo -prune -o -name .git -prune -o -name out -prune -o -type f -name "*\.gradle" \
        -exec grep --color -n "$@" {} +
}

function jgrep()
{
    find . -name .repo -prune -o -name .git -prune -o -name out -prune -o -type f -name "*\.java" \
        -exec grep --color -n "$@" {} +
}

function cgrep()
{
    find . -name .repo -prune -o -name .git -prune -o -name out -prune -o -type f \( -name '*.c' -o -name '*.cc' -o -name '*.cpp' -o -name '*.h' -o -name '*.hpp' \) \
        -exec grep --color -n "$@" {} +
}

function resgrep()
{
    local dir
    for dir in `find . -name .repo -prune -o -name .git -prune -o -name out -prune -o -name res -type d`; do
        find $dir -type f -name '*\.xml' -exec grep --color -n "$@" {} +
    done
}

function mangrep()
{
    find . -name .repo -prune -o -name .git -prune -o -path ./out -prune -o -type f -name 'AndroidManifest.xml' \
        -exec grep --color -n "$@" {} +
}

function sepgrep()
{
    find . -name .repo -prune -o -name .git -prune -o -path ./out -prune -o -name sepolicy -type d \
        -exec grep --color -n -r --exclude-dir=\.git "$@" {} +
}

function rcgrep()
{
    find . -name .repo -prune -o -name .git -prune -o -name out -prune -o -type f -name "*\.rc*" \
        -exec grep --color -n "$@" {} +
}

case `uname -s` in
    Darwin)
        function mgrep()
        {
            find -E . -name .repo -prune -o -name .git -prune -o -path ./out -prune -o \( -iregex '.*/(Makefile|Makefile\..*|.*\.make|.*\.mak|.*\.mk|.*\.bp)' -o -regex '(.*/)?soong/[^/]*.go' \) -type f \
                -exec grep --color -n "$@" {} +
        }

        function treegrep()
        {
            find -E . -name .repo -prune -o -name .git -prune -o -type f -iregex '.*\.(c|h|cpp|hpp|S|java|xml)' \
                -exec grep --color -n -i "$@" {} +
        }

        ;;
    *)
        function mgrep()
        {
            find . -name .repo -prune -o -name .git -prune -o -path ./out -prune -o \( -regextype posix-egrep -iregex '(.*\/Makefile|.*\/Makefile\..*|.*\.make|.*\.mak|.*\.mk|.*\.bp)' -o -regextype posix-extended -regex '(.*/)?soong/[^/]*.go' \) -type f \
                -exec grep --color -n "$@" {} +
        }

        function treegrep()
        {
            find . -name .repo -prune -o -name .git -prune -o -regextype posix-egrep -iregex '.*\.(c|h|cpp|hpp|S|java|xml)' -type f \
                -exec grep --color -n -i "$@" {} +
        }

        ;;
esac

function getprebuilt
{
    get_abs_build_var ANDROID_PREBUILTS
}

function tracedmdump()
{
    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP."
        return
    fi
    local prebuiltdir=$(getprebuilt)
    local arch=$(gettargetarch)
    local KERNEL=$T/prebuilts/qemu-kernel/$arch/vmlinux-qemu

    local TRACE=$1
    if [ ! "$TRACE" ] ; then
        echo "usage:  tracedmdump  tracename"
        return
    fi

    if [ ! -r "$KERNEL" ] ; then
        echo "Error: cannot find kernel: '$KERNEL'"
        return
    fi

    local BASETRACE=$(basename $TRACE)
    if [ "$BASETRACE" = "$TRACE" ] ; then
        TRACE=$ANDROID_PRODUCT_OUT/traces/$TRACE
    fi

    echo "post-processing traces..."
    rm -f $TRACE/qtrace.dexlist
    post_trace $TRACE
    if [ $? -ne 0 ]; then
        echo "***"
        echo "*** Error: malformed trace.  Did you remember to exit the emulator?"
        echo "***"
        return
    fi
    echo "generating dexlist output..."
    /bin/ls $ANDROID_PRODUCT_OUT/system/framework/*.jar $ANDROID_PRODUCT_OUT/system/app/*.apk $ANDROID_PRODUCT_OUT/data/app/*.apk 2>/dev/null | xargs dexlist > $TRACE/qtrace.dexlist
    echo "generating dmtrace data..."
    q2dm -r $ANDROID_PRODUCT_OUT/symbols $TRACE $KERNEL $TRACE/dmtrace || return
    echo "generating html file..."
    dmtracedump -h $TRACE/dmtrace >| $TRACE/dmtrace.html || return
    echo "done, see $TRACE/dmtrace.html for details"
    echo "or run:"
    echo "    traceview $TRACE/dmtrace"
}

# communicate with a running device or emulator, set up necessary state,
# and run the hat command.
function runhat()
{
    # process standard adb options
    local adbTarget=""
    if [ "$1" = "-d" -o "$1" = "-e" ]; then
        adbTarget=$1
        shift 1
    elif [ "$1" = "-s" ]; then
        adbTarget="$1 $2"
        shift 2
    fi
    local adbOptions=${adbTarget}
    #echo adbOptions = ${adbOptions}

    # runhat options
    local targetPid=$1

    if [ "$targetPid" = "" ]; then
        echo "Usage: runhat [ -d | -e | -s serial ] target-pid"
        return
    fi

    # confirm hat is available
    if [ -z $(which hat) ]; then
        echo "hat is not available in this configuration."
        return
    fi

    # issue "am" command to cause the hprof dump
    local devFile=/data/local/tmp/hprof-$targetPid
    echo "Poking $targetPid and waiting for data..."
    echo "Storing data at $devFile"
    adb ${adbOptions} shell am dumpheap $targetPid $devFile
    echo "Press enter when logcat shows \"hprof: heap dump completed\""
    echo -n "> "
    read

    local localFile=/tmp/$$-hprof

    echo "Retrieving file $devFile..."
    adb ${adbOptions} pull $devFile $localFile

    adb ${adbOptions} shell rm $devFile

    echo "Running hat on $localFile"
    echo "View the output by pointing your browser at http://localhost:7000/"
    echo ""
    hat -JXmx512m $localFile
}

function getbugreports()
{
    local reports=(`adb shell ls /sdcard/bugreports | tr -d '\r'`)

    if [ ! "$reports" ]; then
        echo "Could not locate any bugreports."
        return
    fi

    local report
    for report in ${reports[@]}
    do
        echo "/sdcard/bugreports/${report}"
        adb pull /sdcard/bugreports/${report} ${report}
        gunzip ${report}
    done
}

function getsdcardpath()
{
    adb ${adbOptions} shell echo -n \$\{EXTERNAL_STORAGE\}
}

function getscreenshotpath()
{
    echo "$(getsdcardpath)/Pictures/Screenshots"
}

function getlastscreenshot()
{
    local screenshot_path=$(getscreenshotpath)
    local screenshot=`adb ${adbOptions} ls ${screenshot_path} | grep Screenshot_[0-9-]*.*\.png | sort -rk 3 | cut -d " " -f 4 | head -n 1`
    if [ "$screenshot" = "" ]; then
        echo "No screenshots found."
        return
    fi
    echo "${screenshot}"
    adb ${adbOptions} pull ${screenshot_path}/${screenshot}
}

function startviewserver()
{
    local port=4939
    if [ $# -gt 0 ]; then
            port=$1
    fi
    adb shell service call window 1 i32 $port
}

function stopviewserver()
{
    adb shell service call window 2
}

function isviewserverstarted()
{
    adb shell service call window 3
}

function key_home()
{
    adb shell input keyevent 3
}

function key_back()
{
    adb shell input keyevent 4
}

function key_menu()
{
    adb shell input keyevent 82
}

function smoketest()
{
    if [ ! "$ANDROID_PRODUCT_OUT" ]; then
        echo "Couldn't locate output files.  Try running 'lunch' first." >&2
        return
    fi
    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP." >&2
        return
    fi

    (\cd "$T" && mmm tests/SmokeTest) &&
      adb uninstall com.android.smoketest > /dev/null &&
      adb uninstall com.android.smoketest.tests > /dev/null &&
      adb install $ANDROID_PRODUCT_OUT/data/app/SmokeTestApp.apk &&
      adb install $ANDROID_PRODUCT_OUT/data/app/SmokeTest.apk &&
      adb shell am instrument -w com.android.smoketest.tests/android.test.InstrumentationTestRunner
}

# simple shortcut to the runtest command
function runtest()
{
    local T=$(gettop)
    if [ ! "$T" ]; then
        echo "Couldn't locate the top of the tree.  Try setting TOP." >&2
        return
    fi
    ("$T"/development/testrunner/runtest.py $@)
}

function godir () {
    if [[ -z "$1" ]]; then
        echo "Usage: godir <regex>"
        return
    fi
    local T=$(gettop)
    local FILELIST
    if [ ! "$OUT_DIR" = "" ]; then
        mkdir -p $OUT_DIR
        FILELIST=$OUT_DIR/filelist
    else
        FILELIST=$T/filelist
    fi
    if [[ ! -f $FILELIST ]]; then
        echo -n "Creating index..."
        (\cd $T; find . -wholename ./out -prune -o -wholename ./.repo -prune -o -type f > $FILELIST)
        echo " Done"
        echo ""
    fi
    local lines
    lines=($(\grep "$1" $FILELIST | sed -e 's/\/[^/]*$//' | sort | uniq))
    if [[ ${#lines[@]} = 0 ]]; then
        echo "Not found"
        return
    fi
    local pathname
    local choice
    if [[ ${#lines[@]} > 1 ]]; then
        while [[ -z "$pathname" ]]; do
            local index=1
            local line
            for line in ${lines[@]}; do
                printf "%6s %s\n" "[$index]" $line
                index=$(($index + 1))
            done
            echo
            echo -n "Select one: "
            unset choice
            read choice
            if [[ $choice -gt ${#lines[@]} || $choice -lt 1 ]]; then
                echo "Invalid choice"
                continue
            fi
            pathname=${lines[$(($choice-1))]}
        done
    else
        pathname=${lines[0]}
    fi
    \cd $T/$pathname
}

# Update module-info.json in out.
function refreshmod() {
    if [ ! "$ANDROID_PRODUCT_OUT" ]; then
        echo "No ANDROID_PRODUCT_OUT. Try running 'lunch' first." >&2
        return 1
    fi

    echo "Refreshing modules (building module-info.json). Log at $ANDROID_PRODUCT_OUT/module-info.json.build.log." >&2

    # for the output of the next command
    mkdir -p $ANDROID_PRODUCT_OUT || return 1

    # Note, can't use absolute path because of the way make works.
    m out/target/product/$(get_build_var TARGET_DEVICE)/module-info.json \
        > $ANDROID_PRODUCT_OUT/module-info.json.build.log 2>&1
}

# List all modules for the current device, as cached in module-info.json. If any build change is
# made and it should be reflected in the output, you should run 'refreshmod' first.
function allmod() {
    if [ ! "$ANDROID_PRODUCT_OUT" ]; then
        echo "No ANDROID_PRODUCT_OUT. Try running 'lunch' first." >&2
        return 1
    fi

    if [ ! -f "$ANDROID_PRODUCT_OUT/module-info.json" ]; then
        echo "Could not find module-info.json. It will only be built once, and it can be updated with 'refreshmod'" >&2
        refreshmod || return 1
    fi

    python -c "import json; print '\n'.join(sorted(json.load(open('$ANDROID_PRODUCT_OUT/module-info.json')).keys()))"
}

# Go to a specific module in the android tree, as cached in module-info.json. If any build change
# is made, and it should be reflected in the output, you should run 'refreshmod' first.
function gomod() {
    if [ ! "$ANDROID_PRODUCT_OUT" ]; then
        echo "No ANDROID_PRODUCT_OUT. Try running 'lunch' first." >&2
        return 1
    fi

    if [[ $# -ne 1 ]]; then
        echo "usage: gomod <module>" >&2
        return 1
    fi

    if [ ! -f "$ANDROID_PRODUCT_OUT/module-info.json" ]; then
        echo "Could not find module-info.json. It will only be built once, and it can be updated with 'refreshmod'" >&2
        refreshmod || return 1
    fi

    local relpath=$(python -c "import json, os
module = '$1'
module_info = json.load(open('$ANDROID_PRODUCT_OUT/module-info.json'))
if module not in module_info:
    exit(1)
print module_info[module]['path'][0]" 2>/dev/null)

    if [ -z "$relpath" ]; then
        echo "Could not find module '$1' (try 'refreshmod' if there have been build changes?)." >&2
        return 1
    else
        cd $ANDROID_BUILD_TOP/$relpath
    fi
}

function _complete_android_module_names() {
    local word=${COMP_WORDS[COMP_CWORD]}
    COMPREPLY=( $(allmod | grep -E "^$word") )
}

# Print colored exit condition
function pez {
    "$@"
    local retval=$?
    if [ $retval -ne 0 ]
    then
        echo $'\E'"[0;31mFAILURE\e[00m"
    else
        echo $'\E'"[0;32mSUCCESS\e[00m"
    fi
    return $retval
}

function get_make_command()
{
    # If we're in the top of an Android tree, use soong_ui.bash instead of make
    if [ -f build/soong/soong_ui.bash ]; then
        # Always use the real make if -C is passed in
        for arg in "$@"; do
            if [[ $arg == -C* ]]; then
                echo command make
                return
            fi
        done
        echo build/soong/soong_ui.bash --make-mode
    else
        echo command make
    fi
}

function _wrap_build()
{
    local start_time=$(date +"%s")
    "$@"
    local ret=$?
    local end_time=$(date +"%s")
    local tdiff=$(($end_time-$start_time))
    local hours=$(($tdiff / 3600 ))
    local mins=$((($tdiff % 3600) / 60))
    local secs=$(($tdiff % 60))
    local ncolors=$(tput colors 2>/dev/null)
    if [ -n "$ncolors" ] && [ $ncolors -ge 8 ]; then
        color_failed=$'\E'"[0;31m"
        color_success=$'\E'"[0;32m"
        color_reset=$'\E'"[00m"
    else
        color_failed=""
        color_success=""
        color_reset=""
    fi
    echo
    if [ $ret -eq 0 ] ; then
        echo -n "${color_success}#### build completed successfully "
    else
        echo -n "${color_failed}#### failed to build some targets "
    fi
    if [ $hours -gt 0 ] ; then
        printf "(%02g:%02g:%02g (hh:mm:ss))" $hours $mins $secs
    elif [ $mins -gt 0 ] ; then
        printf "(%02g:%02g (mm:ss))" $mins $secs
    elif [ $secs -gt 0 ] ; then
        printf "(%s seconds)" $secs
    fi
    echo " ####${color_reset}"
    echo
    return $ret
}

function make()
{
    _wrap_build $(get_make_command "$@") "$@"
}

function provision()
{
    if [ ! "$ANDROID_PRODUCT_OUT" ]; then
        echo "Couldn't locate output files.  Try running 'lunch' first." >&2
        return 1
    fi
    if [ ! -e "$ANDROID_PRODUCT_OUT/provision-device" ]; then
        echo "There is no provisioning script for the device." >&2
        return 1
    fi

    # Check if user really wants to do this.
    if [ "$1" = "--no-confirmation" ]; then
        shift 1
    else
        echo "This action will reflash your device."
        echo ""
        echo "ALL DATA ON THE DEVICE WILL BE IRREVOCABLY ERASED."
        echo ""
        echo -n "Are you sure you want to do this (yes/no)? "
        read
        if [[ "${REPLY}" != "yes" ]] ; then
            echo "Not taking any action. Exiting." >&2
            return 1
        fi
    fi
    "$ANDROID_PRODUCT_OUT/provision-device" "$@"
}

function atest()
{
    # Let's use the built version over the prebuilt, then source code.
    local os_arch=$(get_build_var HOST_PREBUILT_TAG)
    local built_atest=${ANDROID_HOST_OUT}/bin/atest
    local prebuilt_atest="$(gettop)"/prebuilts/asuite/atest/$os_arch/atest
    if [[ -x $built_atest ]]; then
        $built_atest "$@"
    elif [[ -x $prebuilt_atest ]]; then
        $prebuilt_atest "$@"
    else
        # TODO: once prebuilt atest released, remove the source code section
        # and change the location of atest_completion.sh in addcompletions().
        "$(gettop)"/tools/tradefederation/core/atest/atest.py "$@"
    fi
}

# Zsh needs bashcompinit called to support bash-style completion.
function enable_zsh_completion() {
    # Don't override user's options if bash-style completion is already enabled.
    if ! declare -f complete >/dev/null; then
        autoload -U compinit && compinit
        autoload -U bashcompinit && bashcompinit
    fi
}

function validate_current_shell() {
    local current_sh="$(ps -o command -p $$)"
    case "$current_sh" in
        *bash*)
            function check_type() { type -t "$1"; }
            ;;
        *zsh*)
            function check_type() { type "$1"; }
            enable_zsh_completion ;;
        *)
            echo -e "WARNING: Only bash and zsh are supported.\nUse of other shell would lead to erroneous results."
            ;;
    esac
}

function acloud()
{
    # Let's use the built version over the prebuilt.
    local built_acloud=${ANDROID_HOST_OUT}/bin/acloud
    if [ -f $built_acloud ]; then
        $built_acloud "$@"
        return $?
    fi

    local host_os_arch=$(get_build_var HOST_PREBUILT_TAG)
    case $host_os_arch in
        linux-x86) "$(gettop)"/prebuilts/asuite/acloud/linux-x86/acloud "$@"
        ;;
        darwin-x86) "$(gettop)"/prebuilts/asuite/acloud/darwin-x86/acloud "$@"
        ;;
    *)
        echo "acloud is not supported on your host arch: $host_os_arch"
        ;;
    esac
}

function aidegen()
{
    # Always use the prebuilt version.
    local host_os_arch=$(get_build_var HOST_PREBUILT_TAG)
    case $host_os_arch in
        linux-x86) "$(gettop)"/prebuilts/asuite/aidegen/linux-x86/aidegen "$@"
        ;;
    *)
        echo "aidegen is not supported on your host arch: $host_os_arch"
        ;;
    esac
}

# Execute the contents of any vendorsetup.sh files we can find.
function source_vendorsetup() {
    for dir in device vendor product; do
        for f in $(test -d $dir && \
            find -L $dir -maxdepth 4 -name 'vendorsetup.sh' 2>/dev/null | sort); do
            echo "including $f"; . $f
        done
    done
}

validate_current_shell
source_vendorsetup
addcompletions
