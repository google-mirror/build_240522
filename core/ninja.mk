KATI=build/core/ckati

KATI_OUTPUTS := build.ninja ninja.sh
NINJA_GOALS := fastincremental generateonly

ANDROID_TARGETS := $(filter-out $(KATI_OUTPUTS) $(NINJA_GOALS),$(MAKECMDGOALS))

.PHONY: all
ifeq (,$(filter generateonly,$(MAKECMDGOALS)))
all fastincremental $(ANDROID_TARGETS): $(KATI_OUTPUTS)
	@echo Starting build with ninja
	@PATH=prebuilts/ninja/linux-x86/:$$PATH ./ninja.sh $(NINJA_ARGS) $(ANDROID_TARGETS)
else
all generateonly $(ANDROID_TARGETS): $(KATI_OUTPUTS)
endif

ifeq (,$(filter fastincremental,$(MAKECMDGOALS)))
KATI_FORCE := FORCE
endif

$(KATI_OUTPUTS): $(KATI) kati.intermediate $(KATI_FORCE)

.INTERMEDIATE: kati.intermediate
kati.intermediate:
	@echo Running kati to generate build.ninja...
	@#TODO: write ninja files to $(PRODUCT_OUT)
	@#TODO: use separate ninja file for mm or single target build
	@$(KATI) --ninja --ignore_optional_include=out/%.P --detect_android_echo --use_find_emulator -f build/core/main.mk $(ANDROID_TARGETS) USE_NINJA=false

.PHONY: FORCE
FORCE:
