#!/bin/bash
#
# Get package name and wear capability resource files of a single
# pre-installed APK in the current build and store as an XML file.
#
# Example command:
# 	path/to/get_wear_capabilities_per_apk.sh
# 		output/intermediates/path/ path/to/apk

OUTPUT=$1
INPUT=$2
AAPT=$3

# Find wear capabilities information of a single app package.
function find_capabilities_for_package() {
  $AAPT dump resources $package \
    | grep -Poz '(?s)android_wear_capabilities.*?\]\n' \
    | grep -Poz '(?s)\[.*\]' | tr -d '[" ' | tr ',]' '\n'
}

# Write package name and wear capabilities to the output file.
package=$2;
printf '\t<package name="%s">\n' $($AAPT dump packagename $package) > $OUTPUT
find_capabilities_for_package $package | while read capability ; do
  if [[ -n $capability ]]; then
    printf '\t\t<capability name="%s"/>\n' $capability \
      >> $OUTPUT
  fi
done
printf '\t</package>\n' >> $OUTPUT
