KATI=build/core/ckati

KATI_OUTPUTS := build.ninja ninja.sh

ANDROID_TARGETS := $(filter-out $(KATI_OUTPUTS),$(MAKECMDGOALS))

.PHONY: all
all $(ANDROID_TARGETS): $(KATI_OUTPUTS)
	@echo Starting build with ninja
	@PATH=prebuilts/ninja/linux-x86/:$$PATH ./ninja.sh $(ANDROID_TARGETS)


$(KATI_OUTPUTS): $(KATI) kati.intermediate FORCE

.INTERMEDIATE: kati.intermediate
kati.intermediate:
	@echo Running kati to generate build.ninja...
	@$(KATI) --ninja --ignore_optional_include=out/%.P --use_find_emulator -f build/core/main.mk $(ANDROID_TARGETS) USE_NINJA=false

.PHONY: FORCE
FORCE:
