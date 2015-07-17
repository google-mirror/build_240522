KATI ?= $(HOST_OUT_EXECUTABLES)/ckati

KATI_OUTPUTS := $(PRODUCT_OUT)/build.ninja $(PRODUCT_OUT)/ninja.sh
NINJA_GOALS := fastincremental generateonly droid

ANDROID_TARGETS := $(filter-out $(KATI_OUTPUTS) $(NINJA_GOALS),$(MAKECMDGOALS))

ifeq (,$(NINJA_STATUS))
NINJA_STATUS := [%p %s/%t]$(space)
endif

ifeq (,$(filter generateonly,$(MAKECMDGOALS)))
droid fastincremental $(ANDROID_TARGETS): $(KATI_OUTPUTS) ninja.intermediate
	@#empty

.INTERMEDIATE: ninja.intermediate
ninja.intermediate:
	@echo Starting build with ninja
	@PATH=prebuilts/ninja/$(HOST_PREBUILT_TAG)/:$$PATH NINJA_STATUS="$(NINJA_STATUS)" $(PRODUCT_OUT)/ninja.sh -C $(TOP) $(NINJA_ARGS) $(ANDROID_TARGETS)
else
generateonly $(ANDROID_TARGETS): $(KATI_OUTPUTS)
	@#empty
endif

ifeq (,$(filter fastincremental,$(MAKECMDGOALS)))
KATI_FORCE := FORCE
endif

$(KATI_OUTPUTS): $(KATI) kati.intermediate $(KATI_FORCE)

.INTERMEDIATE: kati.intermediate
kati.intermediate:
	@echo Running kati to generate build.ninja...
	@#TODO: use separate ninja file for mm or single target build
	@$(KATI) --ninja --ninja_dir=$(PRODUCT_OUT) --ignore_optional_include=out/%.P --detect_android_echo --use_find_emulator -f build/core/main.mk $(ANDROID_TARGETS) USE_NINJA=false

KATI_CXX := $(CLANG_CXX) $(CLANG_HOST_GLOBAL_CPPFLAGS)
KATI_LD := $(CLANG_CXX) $(CLANG_HOST_GLOBAL_LDFLAGS)
KATI_INTERMEDIATES_PATH := $(HOST_OUT_INTERMEDIATES)/EXECUTABLES/ckati_intermediates
KATI_BIN_PATH := $(HOST_OUT_EXECUTABLES)
include build/kati/Makefile.ckati

.PHONY: FORCE
FORCE:
