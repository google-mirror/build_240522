#!/bin/bash

# Late in Android Q, libhwbinder and libhidltransport were combined into
# libhidlbase. This is because of the memory overhead associated with extra
# libraries. However, there is still an overhead of ~4kb per empty library. So,
# we want to remove unnecessary usages of the now empty
# libhwbinder/libhidltransport.

echo "Removing libhwbinder/libhidltransport deps from $(pwd)."
echo "Then, this script will list remaining locations that might have to be manually cleaned up."
echo "ENTER or ctrl+C"
read || exit 1

echo "Starting..."
find . -type f -not -path "*/.git/*" -not -path "*/.repo/*" -not -path "*/out/*" -name "Android.bp" \
        -exec sed -i -e "/^ \+\"\(libhidltransport\|libhwbinder\)\",\?$/d" {} \;

find . -type f -not -path "*/.git/*" -not -path "*/.repo/*" -not -path "*/out/*" -name "*.mk" \
        -exec sed -i -e "/^ \+\(libhidltransport\|libhwbinder\) \+\\\\$/d" {} \;

echo "Locations unsure about"
find . -type f -not -path "*/.git/*" -not -path "*/.repo/*" -not -path "*/out/*" -name "*.mk" \
        -exec grep -nP "\b(libhidltransport|libhwbinder)\b" {} /dev/null \;

find . -type f -not -path "*/.git/*" -not -path "*/.repo/*" -not -path "*/out/*" \
        -not -path "*/prebuilts/vndk/*" -name "Android.bp" \
        -exec grep -nP "\b(libhidltransport|libhwbinder)\b" {} /dev/null \;
