#!/bin/bash

set -ex

function create_new_compat_matrix() {

    local top="$(dirname "$0")"/../../../..
    source $top/build/make/tools/finalization/environment.sh

    # create the new file and modify the level
    local current="$CURRENT_COMPATIBILITY_MATRIX_LEVEL"
    local final="$FINAL_COMPATIBILITY_MATRIX_LEVEL"
    local current_file=compatibility_matrix."$current".xml
    local final_file=compatibility_matrix."$final".xml
    local src=$top/hardware/interfaces/compatibility_matrices/compatibility_matrix."$current".xml
    local dest=$top/hardware/interfaces/compatibility_matrices/compatibility_matrix."$final".xml
    sed "s/level=\""$current"\"/level=\""$final"\"/" "$src" > "$dest"

    # add the new module to the end of the Android.bp file
    local bp_file=$top/hardware/interfaces/compatibility_matrices/Android.bp
    echo "" >> $bp_file
    echo "vintf_compatibility_matrix {" >> $bp_file
    echo "    name: \"framework_compatibility_matrix.10.xml\"," >> $bp_file
    echo "    stem: \"compatibility_matrix.10.xml\"," >> $bp_file
    echo "    srcs: [" >> $bp_file
    echo "        \"compatibility_matrix.10.xml\"," >> $bp_file
    echo "    ]," >> $bp_file
    echo "    kernel_configs: [" >> $bp_file
    echo "    ]," >> $bp_file
    echo "}" >> $bp_file

    # hack to get bpmodify to display the previous kernel_config properties by
    # removing the property in the current module and displaying the diff.
    local kernel_configs=$(bpmodify -m framework_$current_file -remove-property --property kernel_configs -d $bp_file | grep -Poz "\"kernel.*\N" | sed 's/\"//g')
    # add the preious kernel_configs to the new module
    bpmodify -m framework_$final_file -property kernel_configs -a $kernel_configs -w $bp_file
    bpfmt -w $bp_file

    local make_file=$top/hardware/interfaces/compatibility_matrices/Android.mk
    # replace the current compat matrix in the make file with the final one
    # the only place this resides is in the conditional addition
    sed -i "s/$current_file/$final_file/g" $make_file
    # add the current compat matrix to the unconditional addition
    sed -i "/^    framework_compatibility_matrix.device.xml/i \    framework_$current_file \\\\" $make_file
}

create_new_compat_matrix
