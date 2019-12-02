#!/bin/bash
for f in $(git diff HEAD^ --name-only); do
  sed -i -e 's|sdkinfo|sdkext|g' $f
  sed -i -e 's|module_sdk|module_sdkext|g' $f
done
