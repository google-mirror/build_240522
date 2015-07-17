KATI ?= build/core/ckati

KATI_OUTPUTS := $(PRODUCT_OUT)/build.ninja $(PRODUCT_OUT)/ninja.sh
NINJA_GOALS := fastincremental generateonly

ANDROID_TARGETS := $(filter-out $(KATI_OUTPUTS) $(NINJA_GOALS),$(MAKECMDGOALS))

.PHONY: all
ifeq (,$(filter generateonly,$(MAKECMDGOALS)))
all droid fastincremental $(ANDROID_TARGETS): $(KATI_OUTPUTS)
	@echo Starting build with ninja
	@PATH=prebuilts/ninja/$(HOST_PREBUILT_TAG)/:$$PATH $(PRODUCT_OUT)/ninja.sh -C $(TOP) $(NINJA_ARGS) $(ANDROID_TARGETS)
else
all droid generateonly $(ANDROID_TARGETS): $(KATI_OUTPUTS)
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

.PHONY: FORCE
FORCE:
