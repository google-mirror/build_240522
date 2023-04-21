#!/bin/bash
#
# Read each package's capability intermediate file in a given directory. Gather
# all capabilities into a metadata XML file, which will be read by the GMS Core
# wearable module at runtime.
#
# Example command:
# 	path/to/merge_wear_capabilities.sh
#		capability/intermediates/path/
# 		output/path/to/static_wear_capabilities.xml

OUTPUT=$1;
shift
INPUTS="$@";

printf '<?xml version="1.0" encoding="utf-8"?>\n<wear_capabilities>\n' \
  > $OUTPUT

cat $INPUTS >> $OUTPUT

printf '</wear_capabilities>' >> $OUTPUT
