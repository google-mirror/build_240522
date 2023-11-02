# This file contains logic to enforce artifact path requirements
# defined in product makefiles.

# Fakes don't get installed, and NDK stubs aren't installed to device.
static_allowed_patterns := $(TARGET_OUT_FAKE)/% $(SOONG_OUT_DIR)/ndk/%
# RROs become REQUIRED by the source module, but are always placed on the vendor partition.
static_allowed_patterns += %__auto_generated_rro_product.apk
static_allowed_patterns += %__auto_generated_rro_vendor.apk
# Auto-included targets are not considered
static_allowed_patterns += $(call product-installed-files,)
# $(PRODUCT_OUT)/apex is where shared libraries in APEXes get installed.
# The path can be considered as a fake path, as the shared libraries
# are installed there just to have symbols files for them under
# $(PRODUCT_OUT)/symbols/apex for debugging purpose. The /apex directory
# is never compiled into a filesystem image.
static_allowed_patterns += $(PRODUCT_OUT)/apex/%
ifeq (true,$(BOARD_USES_SYSTEM_OTHER_ODEX))
  # Allow system_other odex space optimization.
  static_allowed_patterns += \
    $(TARGET_OUT_SYSTEM_OTHER)/%.odex \
    $(TARGET_OUT_SYSTEM_OTHER)/%.vdex \
    $(TARGET_OUT_SYSTEM_OTHER)/%.art
endif

ifneq (,$(filter-out true false relaxed strict,$(PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS))$(filter-out 1 0,$(words $(PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS))))
  $(error PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS must be one of [true, false, relaxed, strict], found: $(PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS))
endif

define is-bounded-product
$(filter $(ARTIFACT_PATH_REQUIREMENT_PRODUCTS),$(1))
endef

define get-enforcement
$(filter-out false,\
  $(call get-product-var,$(1),PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS))
endef

_nmp_cache_key = $(strip $(1)).$(strip $(2))
_nmp_cache_var = APR_NEAREST_MAKEFILES.$(_nmp_cache_key)
APR_NEAREST_MAKEFILES_CACHE_KEYS :=

define nearest-matching-products
$(eval _nmp_all_visited_parents := ) \
$(if $(filter $(APR_NEAREST_MAKEFILES_CACHE_KEYS),$(_nmp_cache_key)),,\
  $(eval $(_nmp_cache_var) := \
    $(if $(strip $(call $(1),$(2))),$(2),\
      $(foreach p,$(PRODUCTS.$(strip $(2)).INHERITS_FROM),\
        $(if $(filter $(_nmp_all_visited_parents),$(p)),,\
          $(eval _nmp_all_visited_parents := $(sort $(_nmp_all_visited_parents) $(p))) \
          $(call nearest-matching-products,$(1),$(p)))))) \
  $(eval APR_NEAREST_MAKEFILES_CACHE_KEYS := $(sort $(APR_NEAREST_MAKEFILES_CACHE_KEYS) $(_nmp_cache_key)))) \
$($(_nmp_cache_var))
endef

define nearest-bounded-products
$(foreach mk,$(1),$(call nearest-matching-products,is-bounded-product,$(mk)))
endef

define nearest-enforced-products
$(foreach mk,$(1),$(call nearest-matching-products,get-enforcement,$(mk)))
endef

define _cpmb_strict_exn_var
_CPMB_INHERITED_STRICT_EXN.$(strip $(1))
endef

define _cpmb_relaxed_exn_var
_CPMB_INHERITED_RELAXED_EXN.$(strip $(1))
endef

_checked_mk_boundaries :=
define check-makefile-path-boundary
$(if $(filter $(_checked_mk_boundaries),$(1)),,\
  $(eval _checked_mk_boundaries := $(sort $(_checked_mk_boundaries) $(1))) \
  $(eval $(_cpmb_strict_exn_var) := ) \
  $(eval $(_cpmb_relaxed_exn_var) := ) \
  $(foreach mk,\
    $(call nearest-bounded-products,$(PRODUCTS.$(strip $(1)).INHERITS_FROM)),\
    $(call check-makefile-path-boundary,$(mk)) \
    $(eval $(_cpmb_strict_exn_var) += $(filter-out $($(_cpmb_strict_exn_var)),$($(call _cpmb_strict_exn_var,$(mk))))) \
    $(eval $(_cpmb_relaxed_exn_var) += $(filter-out $($(_cpmb_relaxed_exn_var)),$($(call _cpmb_relaxed_exn_var,$(mk))))) \
    $(if $(PRODUCTS.$(1).ARTIFACT_PATH_REQUIREMENT_ALLOW_INHERITED),\
      $(eval # These must be relaxed because derived products may override them away) \
      $(eval $(_cpmb_relaxed_exn_var) += $(filter-out $($(_cpmb_relaxed_exn_var)),\
        $(patsubst $(PRODUCT_OUT)/%,%,$(call product-installed-files,$(mk))))))) \
  $(eval allowed_var := $(if $(PRODUCTS.$(strip $(1)).ARTIFACT_PATH_REQUIREMENT_IS_RELAXED),\
    $(_cpmb_relaxed_exn_var),\
    $(_cpmb_strict_exn_var))) \
  $(eval $(allowed_var) += \
    $(filter-out $($(allowed_var)),$(PRODUCTS.$(strip $(1)).ARTIFACT_PATH_EXCEPTION))) \
  $(eval strict_allowed := $($(_cpmb_strict_exn_var))) \
  $(eval relaxed_allowed := $($(_cpmb_relaxed_exn_var))) \
  $(eval boundary := $(PRODUCTS.$(strip $(1)).ARTIFACT_PATH_BOUNDARY)) \
  $(if $(boundary),\
    $(eval ### Verify that the product only produces files inside its path boundary.) \
    $(eval path_patterns := $(call resolve-product-relative-paths,$(boundary),%)) \
    $(eval strict_allowed_patterns := $(call resolve-product-relative-paths,$(strict_allowed))) \
    $(eval relaxed_allowed_patterns := $(call resolve-product-relative-paths,$(relaxed_allowed))) \
    $(eval allowed_patterns := \
      $(path_patterns) \
      $(strict_allowed_patterns) \
      $(relaxed_allowed_patterns) \
      $(static_allowed_patterns)) \
    $(eval files := $(call product-installed-files, $(1))) \
    $(eval offending_files := \
      $(filter-out $(allowed_patterns),$(files))) \
    $(call maybe-print-list-and-error,$(offending_files),\
      $(1) produces files outside its artifact path requirement. \
      Allowed paths are $(subst $(space),$(comma)$(space),$(addsuffix *,$(boundary)))) \
    $(eval unused_allowed := \
      $(foreach pat,$(filter-out $(files),$(strict_allowed_patterns)),\
        $(if $(filter $(pat),$(files)),,$(pat)))) \
    $(call maybe-print-list-and-error,$(unused_allowed),$(1) includes redundant allowed entries in its artifact path requirement.)))
endef

all_offending_files :=

_epb_enforcement_id = $(strip $(1)):$(strip $(2))
_checked_mk_enforcement :=
define enforce-product-bounds
$(if $(filter $(_checked_mk_enforcement),$(_epb_enforcement_id)),,\
  $(eval _checked_mk_enforcement := $(sort $(_checked_mk_enforcement) $(_epb_enforcement_id))) \
  $(foreach mk,\
    $(call nearest-bounded-products,$(PRODUCTS.$(strip $(2)).INHERITS_FROM)),\
    $(call enforce-product-bounds,$(1),$(mk))) \
  $(eval files := $(call product-installed-files,$(2))) \
  $(eval extra_files := $(filter-out $(files) $(HOST_OUT)/%,$(installed_files))) \
  $(eval ### Optionally verify that nothing else produces files inside this artifact path requirement.) \
  $(eval ownership := $(PRODUCTS.$(2).ARTIFACT_PATH_OWNERSHIP)) \
  $(eval path_patterns := $(call resolve-product-relative-paths,$(ownership),%)) \
  $(eval files_in_requirement := $(filter $(path_patterns),$(extra_files))) \
  $(eval all_offending_files += $(filter-out $(all_offending_files),$(files_in_requirement))) \
  $(eval offending_files := $(filter-out $(allowed_patterns),$(files_in_requirement))) \
  $(eval used_allowed := $(foreach p,$(allowed_patterns),$(if $(filter $(p),$(files_in_requirement)),$(p)))) \
  $(eval unused_allowed := $(filter-out $(used_allowed),$(unused_allowed))) \
  $(call maybe-print-list-and-error,$(offending_files),\
    $(1) produces files inside $(2)'s artifact path requirement. \
    $(call get-product-var,$(1),PRODUCT_ARTIFACT_PATH_REQUIREMENT_HINT)))
endef

define enforce-all-product-bounds
$(if $(filter $(_checked_mk_enforcement),$(1)),,\
  $(eval _checked_mk_enforcement := $(sort $(_checked_mk_enforcement) $(1))) \
  $(foreach mk,\
    $(call nearest-enforced-products,$(PRODUCTS.$(strip $(1)).INHERITS_FROM)),\
    $(call enforce-all-product-bounds,$(mk))) \
  $(eval installed_files := $(call product-installed-files,$(1))) \
  $(eval allowed := $(call get-product-var,$(1),PRODUCT_ARTIFACT_PATH_REQUIREMENT_ALLOWED_LIST)) \
  $(eval allowed_patterns := $(call resolve-product-relative-paths,$(allowed))) \
  $(eval unused_allowed := $(allowed_patterns)) \
  $(foreach mk,$(call nearest-bounded-products,$(1)),\
    $(call enforce-product-bounds,$(1),$(mk))) \
  $(if $(filter true strict,$(call get-enforcement,$(1))),\
    $(call maybe-print-list-and-error,$(unused_allowed),$(1) includes redundant artifact path requirement allowed list entries.)))
endef

$(foreach mk,$(call nearest-bounded-products,$(INTERNAL_PRODUCT)),\
  $(call check-makefile-path-boundary,$(mk)))

$(foreach mk,$(call nearest-enforced-products,$(INTERNAL_PRODUCT)),\
  $(call enforce-all-product-bounds,$(mk)))

$(PRODUCT_OUT)/offending_artifacts.txt:
	rm -f $@
	$(foreach f,$(sort $(all_offending_files)),echo $(f) >> $@;)
