#
# Copyright (C) 2009 The Android Open Source Project
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

<<<<<<< HEAD   (5f91bd Merge "Merge empty history for sparse-8435393-L9030000095399)
# the foreach and the if remove the single space entries that creep in because of the evals
define gather-all-products
$(sort $(foreach p, \
	$(eval _all_products_visited := )
  $(call all-products-inner, $(ALL_PRODUCTS)) \
	, $(if $(strip $(p)),$(strip $(p)),)) \
)
=======
# the sort also acts as a strip to remove the single space entries that creep in because of the evals
define gather-all-makefiles-for-current-product
$(eval _all_products_visited := )\
$(sort $(call gather-all-makefiles-for-current-product-inner,$(INTERNAL_PRODUCT)))
>>>>>>> BRANCH (436489 Merge "Version bump to TKB1.220417.001.A1 [core/build_id.mk])
endef

define gather-all-makefiles-for-current-product-inner
	$(foreach p,$(1),\
		$(if $(filter $(p),$(_all_products_visited)),, \
			$(p) \
			$(eval _all_products_visited += $(p)) \
			$(call all-products-inner, $(PRODUCTS.$(strip $(p)).INHERITS_FROM))
		) \
	)
endef

<<<<<<< HEAD   (5f91bd Merge "Merge empty history for sparse-8435393-L9030000095399)

this_makefile := build/make/core/product-graph.mk

products_svg := $(OUT_DIR)/products.svg
products_pdf := $(OUT_DIR)/products.pdf
products_graph := $(OUT_DIR)/products.dot
ifeq ($(strip $(ANDROID_PRODUCT_GRAPH)),)
products_list := $(INTERNAL_PRODUCT)
else
ifeq ($(strip $(ANDROID_PRODUCT_GRAPH)),--all)
products_list := --all
else
products_list := $(foreach prod,$(ANDROID_PRODUCT_GRAPH),$(call resolve-short-product-name,$(prod)))
endif
endif

really_all_products := $(call gather-all-products)
=======
node_color_target := orange
node_color_common := beige
node_color_vendor := lavenderblush
node_color_default := white
define node-color
$(if $(filter $(1),$(PRIVATE_TOP_LEVEL_MAKEFILE)),\
  $(node_color_target),\
  $(if $(filter build/make/target/product/%,$(1)),\
    $(node_color_common),\
    $(if $(filter vendor/%,$(1)),$(node_color_vendor),$(node_color_default))\
  )\
)
endef
>>>>>>> BRANCH (436489 Merge "Version bump to TKB1.220417.001.A1 [core/build_id.mk])

open_parethesis := (
close_parenthesis := )

# Emit properties of a product node to a file.
# $(1) the product
# $(2) the output file
define emit-product-node-props
$(hide) echo \"$(1)\" [ \
<<<<<<< HEAD   (5f91bd Merge "Merge empty history for sparse-8435393-L9030000095399)
label=\"$(dir $(1))\\n$(notdir $(1))\\n\\n$(subst $(close_parenthesis),,$(subst $(open_parethesis),,$(PRODUCTS.$(strip $(1)).PRODUCT_MODEL)))\\n$(PRODUCTS.$(strip $(1)).PRODUCT_DEVICE)\" \
$(if $(filter $(1),$(PRIVATE_PRODUCTS_FILTER)), style=\"filled\" fillcolor=\"#FFFDB0\",) \
colorscheme=\"svg\" fontcolor=\"darkblue\" href=\"products/$(1).html\" \
=======
label=\"$(dir $(1))\\n$(notdir $(1))$(if $(filter $(1),$(PRIVATE_TOP_LEVEL_MAKEFILE)),$(subst $(open_parethesis),,$(subst $(close_parenthesis),,\\n\\n$(PRODUCT_MODEL)\\n$(PRODUCT_DEVICE))))\" \
style=\"filled\" fillcolor=\"$(strip $(call node-color,$(1)))\" \
colorscheme=\"svg\" fontcolor=\"darkblue\" \
>>>>>>> BRANCH (436489 Merge "Version bump to TKB1.220417.001.A1 [core/build_id.mk])
] >> $(2)

endef

<<<<<<< HEAD   (5f91bd Merge "Merge empty history for sparse-8435393-L9030000095399)
$(products_graph): PRIVATE_PRODUCTS := $(really_all_products)
$(products_graph): PRIVATE_PRODUCTS_FILTER := $(products_list)
=======
products_graph := $(OUT_DIR)/products.dot
>>>>>>> BRANCH (436489 Merge "Version bump to TKB1.220417.001.A1 [core/build_id.mk])

$(products_graph): PRIVATE_ALL_MAKEFILES_FOR_THIS_PRODUCT := $(call gather-all-makefiles-for-current-product)
$(products_graph): PRIVATE_TOP_LEVEL_MAKEFILE := $(INTERNAL_PRODUCT)
$(products_graph):
	@echo Product graph DOT: $@ for $(PRIVATE_TOP_LEVEL_MAKEFILE)
	$(hide) echo 'digraph {' > $@
	$(hide) echo 'graph [ ratio=.5 ];' >> $@
	$(hide) $(foreach p,$(PRIVATE_ALL_MAKEFILES_FOR_THIS_PRODUCT), \
	  $(foreach d,$(PRODUCTS.$(strip $(p)).INHERITS_FROM), echo \"$(d)\" -\> \"$(p)\" >> $@;))
	$(foreach p,$(PRIVATE_ALL_MAKEFILES_FOR_THIS_PRODUCT),$(call emit-product-node-props,$(p),$@))
	$(hide) echo '}' >> $@

# Evaluates to the name of the product file
# $(1) product file
define product-debug-filename
$(OUT_DIR)/products/$(strip $(1)).html
endef

# Makes a rule for the product debug info
# $(1) product file
define transform-product-debug
$(OUT_DIR)/products/$(strip $(1)).txt: $(this_makefile)
	@echo Product debug info file: $$@
	$(hide) rm -f $$@
	$(hide) mkdir -p $$(dir $$@)
	$(hide) echo 'FILE=$(strip $(1))' >> $$@
	$(hide) echo 'PRODUCT_NAME=$$(PRODUCTS.$(strip $(1)).PRODUCT_NAME)' >> $$@
	$(hide) echo 'PRODUCT_MODEL=$$(PRODUCTS.$(strip $(1)).PRODUCT_MODEL)' >> $$@
	$(hide) echo 'PRODUCT_LOCALES=$$(PRODUCTS.$(strip $(1)).PRODUCT_LOCALES)' >> $$@
	$(hide) echo 'PRODUCT_AAPT_CONFIG=$$(PRODUCTS.$(strip $(1)).PRODUCT_AAPT_CONFIG)' >> $$@
	$(hide) echo 'PRODUCT_AAPT_PREF_CONFIG=$$(PRODUCTS.$(strip $(1)).PRODUCT_AAPT_PREF_CONFIG)' >> $$@
	$(hide) echo 'PRODUCT_PACKAGES=$$(PRODUCTS.$(strip $(1)).PRODUCT_PACKAGES)' >> $$@
	$(hide) echo 'PRODUCT_DEVICE=$$(PRODUCTS.$(strip $(1)).PRODUCT_DEVICE)' >> $$@
	$(hide) echo 'PRODUCT_MANUFACTURER=$$(PRODUCTS.$(strip $(1)).PRODUCT_MANUFACTURER)' >> $$@
	$(hide) echo 'PRODUCT_PROPERTY_OVERRIDES=$$(PRODUCTS.$(strip $(1)).PRODUCT_PROPERTY_OVERRIDES)' >> $$@
	$(hide) echo 'PRODUCT_DEFAULT_PROPERTY_OVERRIDES=$$(PRODUCTS.$(strip $(1)).PRODUCT_DEFAULT_PROPERTY_OVERRIDES)' >> $$@
	$(hide) echo 'PRODUCT_SYSTEM_DEFAULT_PROPERTIES=$$(PRODUCTS.$(strip $(1)).PRODUCT_SYSTEM_DEFAULT_PROPERTIES)' >> $$@
	$(hide) echo 'PRODUCT_PRODUCT_PROPERTIES=$$(PRODUCTS.$(strip $(1)).PRODUCT_PRODUCT_PROPERTIES)' >> $$@
	$(hide) echo 'PRODUCT_CHARACTERISTICS=$$(PRODUCTS.$(strip $(1)).PRODUCT_CHARACTERISTICS)' >> $$@
	$(hide) echo 'PRODUCT_COPY_FILES=$$(PRODUCTS.$(strip $(1)).PRODUCT_COPY_FILES)' >> $$@
	$(hide) echo 'PRODUCT_OTA_PUBLIC_KEYS=$$(PRODUCTS.$(strip $(1)).PRODUCT_OTA_PUBLIC_KEYS)' >> $$@
	$(hide) echo 'PRODUCT_EXTRA_RECOVERY_KEYS=$$(PRODUCTS.$(strip $(1)).PRODUCT_EXTRA_RECOVERY_KEYS)' >> $$@
	$(hide) echo 'PRODUCT_PACKAGE_OVERLAYS=$$(PRODUCTS.$(strip $(1)).PRODUCT_PACKAGE_OVERLAYS)' >> $$@
	$(hide) echo 'DEVICE_PACKAGE_OVERLAYS=$$(PRODUCTS.$(strip $(1)).DEVICE_PACKAGE_OVERLAYS)' >> $$@
	$(hide) echo 'PRODUCT_SDK_ADDON_NAME=$$(PRODUCTS.$(strip $(1)).PRODUCT_SDK_ADDON_NAME)' >> $$@
	$(hide) echo 'PRODUCT_SDK_ADDON_COPY_FILES=$$(PRODUCTS.$(strip $(1)).PRODUCT_SDK_ADDON_COPY_FILES)' >> $$@
	$(hide) echo 'PRODUCT_SDK_ADDON_COPY_MODULES=$$(PRODUCTS.$(strip $(1)).PRODUCT_SDK_ADDON_COPY_MODULES)' >> $$@
	$(hide) echo 'PRODUCT_SDK_ADDON_DOC_MODULES=$$(PRODUCTS.$(strip $(1)).PRODUCT_SDK_ADDON_DOC_MODULES)' >> $$@
	$(hide) echo 'PRODUCT_DEFAULT_WIFI_CHANNELS=$$(PRODUCTS.$(strip $(1)).PRODUCT_DEFAULT_WIFI_CHANNELS)' >> $$@
	$(hide) echo 'PRODUCT_DEFAULT_DEV_CERTIFICATE=$$(PRODUCTS.$(strip $(1)).PRODUCT_DEFAULT_DEV_CERTIFICATE)' >> $$@
	$(hide) echo 'PRODUCT_RESTRICT_VENDOR_FILES=$$(PRODUCTS.$(strip $(1)).PRODUCT_RESTRICT_VENDOR_FILES)' >> $$@
	$(hide) echo 'PRODUCT_VENDOR_KERNEL_HEADERS=$$(PRODUCTS.$(strip $(1)).PRODUCT_VENDOR_KERNEL_HEADERS)' >> $$@

$(call product-debug-filename, $(p)): \
			$(OUT_DIR)/products/$(strip $(1)).txt \
			build/make/tools/product_debug.py \
			$(this_makefile)
	@echo Product debug html file: $$@
	$(hide) mkdir -p $$(dir $$@)
	$(hide) cat $$< | build/make/tools/product_debug.py > $$@
endef

product_debug_files:=
$(foreach p,$(really_all_products), \
			$(eval $(call transform-product-debug, $(p))) \
			$(eval product_debug_files += $(call product-debug-filename, $(p))) \
   )

$(products_pdf): $(products_graph)
	@echo Product graph PDF: $@
	dot -Tpdf -Nshape=box -o $@ $<

$(products_svg): $(products_graph) $(product_debug_files)
	@echo Product graph SVG: $@
	dot -Tsvg -Nshape=box -o $@ $<

product-graph: $(products_pdf) $(products_svg)
