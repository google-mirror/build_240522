#
# Copyright (C) 2007 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# When specifying "dist", the user has asked that we copy the important
# files from this build into DIST_DIR.

ifdef dist_goal

# $(1): source file
# $(2): destination file
# $(3): goals that should copy the file
#
define copy-one-dist-file
$(3): $(2)
$(2): $(1)
	@echo "Dist: $$@"
	$$(copy-file-to-new-target-with-cp)
endef

# A global variable to remember all dist'ed src:dst pairs.
# So if a src:dst is already dist'ed by another goal,
# we should just establish the dependency and don't really call the
# copy-one-dist-file to avoid multiple rules for the same target.
_all_dist_src_dst_pairs :=
# Other parts of the system should use this function to associate
# certain files with certain goals.  When those goals are built
# and "dist" is specified, the marked files will be copied to DIST_DIR.
#
# $(1): a list of goals  (e.g. droid, sdk, pdk, ndk)
# $(2): the dist files to add to those goals.  If the file contains ':',
#       the text following the colon is the name that the file is copied
#       to under the dist directory.  Subdirs are ok, and will be created
#       at copy time if necessary.
define dist-for-goals
$(foreach file,$(2), \
<<<<<<< HEAD   (326d62 Merge "Merge empty history for sparse-8747889-L1870000095520)
  $(eval fw := $(subst :,$(space),$(file))) \
  $(eval src := $(word 1,$(fw))) \
  $(eval dst := $(word 2,$(fw))) \
  $(eval dst := $(if $(dst),$(dst),$(notdir $(src)))) \
  $(if $(filter $(_all_dist_src_dst_pairs),$(src):$(dst)),\
    $(eval $(call add-dependency,$(1),$(DIST_DIR)/$(dst))),\
    $(eval $(call copy-one-dist-file,\
      $(src),$(DIST_DIR)/$(dst),$(1)))\
      $(eval _all_dist_src_dst_pairs += $(src):$(dst))\
  )\
=======
  $(eval src := $(call word-colon,1,$(file))) \
  $(eval dst := $(call word-colon,2,$(file))) \
  $(if $(dst),,$(eval dst := $$(notdir $$(src)))) \
  $(eval _all_dist_src_dst_pairs += $$(src):$$(dst)) \
  $(foreach goal,$(1), \
    $(eval _all_dist_goal_output_pairs += $$(goal):$$(dst))))
endef

.PHONY: shareprojects

define __share-projects-rule
$(1) : PRIVATE_TARGETS := $(2)
$(1) : PRIVATE_ARGUMENT_FILE := $(call intermediates-dir-for,METAPACKAGING,codesharing)/$(1)/arguments
$(1): $(2) $(COMPLIANCE_LISTSHARE)
	$(hide) rm -f $$@
	mkdir -p $$(dir $$@)
	mkdir -p $$(dir $$(PRIVATE_ARGUMENT_FILE))
	$$(if $$(strip $$(PRIVATE_TARGETS)),$$(call dump-words-to-file,$$(PRIVATE_TARGETS),$$(PRIVATE_ARGUMENT_FILE)))
	$$(if $$(strip $$(PRIVATE_TARGETS)),OUT_DIR=$(OUT_DIR) $(COMPLIANCE_LISTSHARE) -o $$@ @$$(PRIVATE_ARGUMENT_FILE),touch $$@)
endef

# build list of projects to share in $(1) for dist targets in $(2)
#
# $(1): the intermediate project sharing file
# $(2): the dist files to base the sharing on
define _share-projects-rule
$(eval $(call __share-projects-rule,$(1),$(call corresponding-license-metadata,$(2))))
endef

.PHONY: alllicensetexts

define __license-texts-rule
$(2) : PRIVATE_GOAL := $(1)
$(2) : PRIVATE_TARGETS := $(3)
$(2) : PRIVATE_ROOTS := $(4)
$(2) : PRIVATE_ARGUMENT_FILE := $(call intermediates-dir-for,METAPACKAGING,licensetexts)/$(2)/arguments
$(2): $(3) $(TEXTNOTICE)
	$(hide) rm -f $$@
	mkdir -p $$(dir $$@)
	mkdir -p $$(dir $$(PRIVATE_ARGUMENT_FILE))
	$$(if $$(strip $$(PRIVATE_TARGETS)),$$(call dump-words-to-file,\
            -product="$$(PRIVATE_GOAL)" -title="$$(PRIVATE_GOAL)" \
            $$(addprefix -strip_prefix ,$$(PRIVATE_ROOTS)) \
            -strip_prefix=$(PRODUCT_OUT)/ -strip_prefix=$(HOST_OUT)/\
            $$(PRIVATE_TARGETS),\
            $$(PRIVATE_ARGUMENT_FILE)))
	$$(if $$(strip $$(PRIVATE_TARGETS)),OUT_DIR=$(OUT_DIR) $(TEXTNOTICE) -o $$@ @$$(PRIVATE_ARGUMENT_FILE),touch $$@)
endef

# build list of projects to share in $(2) for dist targets in $(3) for dist goal $(1)
#
# $(1): the name of the dist goal
# $(2): the intermediate project sharing file
# $(3): the dist files to base the sharing on
define _license-texts-rule
$(eval $(call __license-texts-rule,$(1),$(2),$(call corresponding-license-metadata,$(3)),$(sort $(dir $(3)))))
endef

define _add_projects_to_share
$(strip $(eval _idir := $(call intermediates-dir-for,METAPACKAGING,shareprojects))) \
$(strip $(eval _tdir := $(call intermediates-dir-for,METAPACKAGING,licensetexts))) \
$(strip $(eval _goals := $(sort $(_all_dist_goals)))) \
$(strip $(eval _opairs := $(sort $(_all_dist_goal_output_pairs)))) \
$(strip $(eval _dpairs := $(sort $(_all_dist_src_dst_pairs)))) \
$(strip $(eval _allt :=)) \
$(foreach goal,$(_goals), \
  $(eval _f := $(_idir)/$(goal).shareprojects) \
  $(eval _n := $(_tdir)/$(goal).txt) \
  $(call dist-for-goals,$(goal),$(_f):shareprojects/$(basename $(notdir $(_f)))) \
  $(call dist-for-goals,$(goal),$(_n):licensetexts/$(basename $(notdir $(_n)))) \
  $(eval _targets :=) \
  $(foreach op,$(filter $(goal):%,$(_opairs)),$(foreach p,$(filter %:$(call word-colon,2,$(op)),$(_dpairs)),$(eval _targets += $(call word-colon,1,$(p))))) \
  $(eval _allt += $(_targets)) \
  $(eval $(call _share-projects-rule,$(_f),$(_targets))) \
  $(eval $(call _license-texts-rule,$(goal),$(_n),$(_targets))) \
>>>>>>> BRANCH (a6db6e Merge "Version bump to TKB1.220626.001.A1 [core/build_id.mk])
)
endef

else # !dist_goal

# empty definition when not building dist
define dist-for-goals
endef

endif # !dist_goal
