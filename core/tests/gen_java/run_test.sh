#!/bin/bash -e

export TEST_BUILD_SYSTEM=true
export OUT_DIR=out.test
export INTERMEDIATE_SRC=${OUT_DIR}/target/common/obj/JAVA_LIBRARIES/build_system_test_gen_java_intermediates/src

cd ${ANDROID_BUILD_TOP}
rm -rf ${OUT_DIR}
source build/envsetup.sh
lunch aosp_arm-eng
echo "Building aosp_arm-eng for the first time"
mmma -j build/core/tests

if [ -d "${INTERMEDIATE_SRC}" ]; then
  echo "build_system_test_gen_java shouldn't have generated sources for aosp_arm" >&2
  exit 1
fi

AOSP_ARM_KATI_DATE=$(date -r ${OUT_DIR}/build-aosp_arm.ninja +%s)

echo "Rebuilding aosp_arm-eng to re-read previous_build_config.mk and clean_steps.mk"
mmma -j build/core/tests

# Ensure that kati regenerated
if [ "${AOSP_ARM_KATI_DATE}" -eq "$(date -r ${OUT_DIR}/build-aosp_arm.ninja +%s)" ]; then
  echo "error: kati should have regenerated because of previous_build_config.mk or clean_steps.mk" >&2
  exit 1
fi

echo "Rebuilding aosp_arm-eng, shouldn't need to re-read kati"
mmma -j build/core/tests

# Ensure kati didn't regenerate
if [ "${AOSP_ARM_KATI_DATE}" -ne "$(date -r ${OUT_DIR}/build-aosp_arm.ninja +%s)" ]; then
  echo "error: kati should not have regenerated" >&2
  exit 1
fi

lunch aosp_arm64-eng
echo "Building aosp_amr64-eng for the first time"
mmma -j build/core/tests

if [ ! -d "${INTERMEDIATE_SRC}" ]; then
  echo "error: build_system_test_gen_java should have generated sources for aosp_arm64" >&2
  exit 1
fi

AOSP_ARM64_KATI_DATE=$(date -r ${OUT_DIR}/build-aosp_arm64.ninja +%s)

echo "Rebuilding aosp_arm64-eng to re-read previous_build_config.mk and clean_steps.mk"
mmma -j build/core/tests

# Ensure that kati regenerated
if [ "${AOSP_ARM64_KATI_DATE}" -eq "$(date -r ${OUT_DIR}/build-aosp_arm64.ninja +%s)" ]; then
  echo "error: kati should have regenerated because of previous_build_config.mk or clean_steps.mk" >&2
  exit 1
fi

echo "Rebuilding aosp_arm64-eng, shouldn't need to re-read kati"
mmma -j build/core/tests

# Ensure kati didn't regenerate
if [ "${AOSP_ARM64_KATI_DATE}" -ne "$(date -r ${OUT_DIR}/build-aosp_arm64.ninja +%s)" ]; then
  echo "error: kati should not have regenerated" >&2
  exit 1
fi

lunch aosp_arm-eng
echo "Building aosp_arm-eng after aosp_arm64-eng"
mmma -j build/core/tests

if [ -d "${INTERMEDIATE_SRC}" ]; then
  echo "error: build_system_test_gen_java generated sources should have been removed for aosp_arm" >&2
  exit 1
fi

# Ensure kati didn't regenerate
if [ "${AOSP_ARM_KATI_DATE}" -ne "$(date -r ${OUT_DIR}/build-aosp_arm.ninja +%s)" ]; then
  echo "error: kati should not have regenerated" >&2
  exit 1
fi
