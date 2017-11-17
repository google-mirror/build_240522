#!/bin/bash
make_ret=$1
dist_dir=$2
search_dir=$3
find $search_dir -name "*.abidiff" -exec cp {} $dist_dir \;
find_ret=$?
#Give preference to make's return status.
if [ $make_ret != 0 ]
then
  exit $make_ret
fi
if [ $find_ret != 0 ]
then
  exit $find_ret
fi
