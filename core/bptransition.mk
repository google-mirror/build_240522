# If a rule fails, delete $@.
.DELETE_ON_ERROR:

# Default rule
.PHONY: androidmk
androidmk:

ANDROIDBP := $(OUT_DIR)/.bootstrap/bin/androidbp
BPFIND := $(OUT_DIR)/.bootstrap/bin/bpfind

# Run bootstrap.bash if necessary
ifeq ($(wildcard $(OUT_DIR)/soong),)
result := $(shell cd $(OUT_DIR) >/dev/null; ../bootstrap.bash >&2 || echo -n "error")
ifeq ($(result),error)
$(error bootstrap.bash failed)
endif
endif

# Build the tools
result := $(shell $(OUT_DIR)/soong .bootstrap/bin/bpfind .bootstrap/bin/androidbp >&2 || echo -n "error")
ifeq ($(result),error)
$(error androidbp/bpfind build failed)
endif

ALL_MAKEFILES :=

# Rule to translate path of Android.bp file to the name of the generated Android_*.mk file
# $1: Android.bp file
define mk-for-bp
$(MK_OUT_DIR)/Android_$(patsubst %_,%,$(subst /,_,$(dir $(1)))).mk
endef

# Rule to translate Android.bp files to Android.mk files
# $1: Android.bp file
define translate-bp-to-mk
  $(eval $(call translate-bp-to-mk-internal,$(1),$(call mk-for-bp,$(1))))
endef
define translate-bp-to-mk-internal
ALL_MAKEFILES += $2
$2 : $1 $(ANDROIDBP) $(BPFIND)
	@echo "Translating $$<..." >&2
	@mkdir -p $$(dir $$@)
	@$(ANDROIDBP) $$< $$@ >&2
endef

# Android.bp files that we expect no matching Android.mk files
missing_blacklist := \
	Android.bp \
	art/build/Android.bp \
	build/soong/Android.bp

# Find Android.bp files that don't have corresponding Android.mk files
missing_makefiles := $(shell $(BPFIND) -name=Android.bp -without=Android.mk Android.bp)
missing_makefiles := $(filter-out $(missing_blacklist),$(missing_makefiles))

# Translate the Android.bp files to Android.mk files
$(foreach mk, $(missing_makefiles), $(call translate-bp-to-mk,$(mk)))

androidmk: $(ALL_MAKEFILES)

# Delete all stale files
all_generated_makefiles := $(wildcard $(MK_OUT_DIR)/Android_*.mk)
stale_makefiles := $(filter-out $(ALL_MAKEFILES),$(all_generated_makefiles))
$(foreach m,$(stale_makefiles), \
   $(info Removing stale makefile $(m))$(shell rm -f $(m)))
