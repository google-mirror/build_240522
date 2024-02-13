#!/bin/bash
# Run filter_teams in CI to create a proto that just has test owners.
# This isn't done as a genrule because module-info.json is not available
# to genrules.

# Get PRODUCT_OUT from make via dumpvars-mode, so we can find the prebuild module-info.json
product_out=$(TARGET_PRODUCT=aosp_x86_64 TARGET_RELEASE=trunk_staging build/soong/soong_ui.bash --dumpvars-mode --vars=PRODUCT_OUT)
eval "${product_out}"
MOD_INFO=${PRODUCT_OUT}/module-info.json

./out/host/linux-x86/bin/filter_teams  --filter_teams \
           --src_teams_file=out/soong/ownership/all_teams.pb \
           --module_info=${MOD_INFO} \
           --out_file=out/soong/ownership/all_test_specs.pb
