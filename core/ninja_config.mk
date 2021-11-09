ifeq ($(filter address,$(SANITIZE_HOST)),)
NINJA ?= prebuilts/build-tools/$(HOST_PREBUILT_TAG)/bin/ninja
else
NINJA ?= prebuilts/build-tools/$(HOST_PREBUILT_TAG)/asan/bin/ninja
endif

KATI_OUTPUT_PATTERNS := $(OUT_DIR)/build%.ninja $(OUT_DIR)/ninja%.sh

# Modifier goals we don't need to pass to Ninja.
NINJA_EXCLUDE_GOALS := all dist APP-% PRODUCT-%

# A list of goals which affect parsing of makefiles and we need to pass to Kati.
PARSE_TIME_MAKE_GOALS := \
	$(PARSE_TIME_MAKE_GOALS) \
	$(dont_bother_goals) \
	all \
	ECLIPSE-% \
	AUX-% \
	boottarball-nodeps \
	brillo_tests \
	btnod \
	build-art% \
	build_kernel-nodeps \
	clean-oat% \
	continuous_instrumentation_tests \
	continuous_native_tests \
	cts \
	custom_images \
	deps-license \
	dicttool_aosp \
	dist \
	dump-products \
	eng \
	fusion \
	oem_image \
	online-system-api-sdk-docs \
	pdk \
	platform \
	platform-java \
	product-graph \
	samplecode \
	sdk \
	sdk_addon \
	sdk_repo \
	snod \
	stnod \
	systemimage-nodeps \
	systemtarball-nodeps \
	target-files-package \
	test-art% \
	user \
	userdataimage \
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	userdebug \
	valgrind-test-art% \
	vts \
	win_sdk \
	winsdk-tools
=======
	userdebug
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

include $(wildcard vendor/*/build/ninja_config.mk)

# Any Android goals that need to be built.
ANDROID_GOALS := $(filter-out $(KATI_OUTPUT_PATTERNS) $(CKATI) $(MAKEPARALLEL),\
    $(sort $(ORIGINAL_MAKECMDGOALS) $(MAKECMDGOALS)))
# Temporary compatibility support until the build server configs are updated
ANDROID_GOALS := $(patsubst win_sdk,sdk,$(ANDROID_GOALS))
ifneq ($(HOST_OS),linux)
  ANDROID_GOALS := $(filter-out sdk,$(ANDROID_GOALS))
  ANDROID_GOALS := $(patsubst sdk_repo,sdk-repo-build-tools sdk-repo-platform-tools,$(ANDROID_GOALS))
endif
# Goals we need to pass to Ninja.
NINJA_GOALS := $(filter-out $(NINJA_EXCLUDE_GOALS), $(ANDROID_GOALS))
ifndef NINJA_GOALS
  NINJA_GOALS := droid
endif
# Goals we need to pass to Kati.
KATI_GOALS := $(filter $(PARSE_TIME_MAKE_GOALS), $(ANDROID_GOALS))
