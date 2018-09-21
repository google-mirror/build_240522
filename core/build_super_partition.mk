#
# Copyright (C) 2018 The Android Open Source Project
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

# Require:
# LOCAL_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_NAME: name of the super partition. Example is
# "super", "super_foo", etc.

ifeq ($(LOCAL_SUPER_PARTITION_NAME),)
$(error LOCAL_SUPER_PARTITION_NAME must not be empty.)
endif

LOCAL_SUPER_PARTITION_NAME_UPPER := $(call to-upper,$(LOCAL_SUPER_PARTITION_NAME))
LOCAL_CHECK_RULE_NAME := check_$(LOCAL_SUPER_PARTITION_NAME)_partition_sizes

# -----------------------------------------------------------------
# super partition image

ifeq (true,$(PRODUCT_BUILD_SUPER_PARTITION))

# BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_SIZE must be defined to build super image.
ifdef BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_SIZE

INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET := $(PRODUCT_OUT)/$(LOCAL_SUPER_PARTITION_NAME).img
INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET := $(PRODUCT_OUT)/$(LOCAL_SUPER_PARTITION_NAME)_empty.img

$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET): $(call images-for-partitions,$(BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_PARTITION_LIST))

# For A/B devices, super partition always contains sub-partitions in the _a slot, because this
# image should only be used for bootstrapping / initializing the device. When flashing the image,
# bootloader fastboot should always mark _a slot as bootable.
ifeq ($(AB_OTA_UPDATER),true)
$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET) $(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET): PRIVATE_PARTITION_SUFFIX=_a
endif # AB_OTA_UPDATER

$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET) $(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET): $(LPMAKE)

$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET) $(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET): PRIVATE_NAME := $(LOCAL_SUPER_PARTITION_NAME)
$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET) $(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET): PRIVATE_NAME_UPPER := $(LOCAL_SUPER_PARTITION_NAME_UPPER)

$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET):
	$(call pretty,"Target $(PRIVATE_NAME) fs image: $@")
	$(call build-superimage-target, \
		$(BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_SIZE), \
		$(BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_PARTITION_LIST), \
		$@,$(PRIVATE_PARTITION_SUFFIX),true)

$(call dist-for-goals,dist_files,$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET))

$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET):
	$(call pretty,"Target empty $(PRIVATE_NAME) fs image: $@")
	$(call build-superimage-target, \
		$(BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_SIZE), \
		$(BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_PARTITION_LIST), \
		$@,$(PRIVATE_PARTITION_SUFFIX))

$(call dist-for-goals,dist_files,$(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET))

endif # BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_SIZE
endif # PRODUCT_BUILD_SUPER_PARTITION

# -----------------------------------------------------------------
# Check image sizes <= size of super partition

ifeq (,$(TARGET_BUILD_APPS))
# Do not check for apps-only build

ifeq (true,$(PRODUCT_BUILD_SUPER_PARTITION))
ifdef BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_SIZE

ifdef BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_PARTITION_LIST

.PHONY: $(LOCAL_CHECK_RULE_NAME)
check_android_partition_sizes: $(LOCAL_CHECK_RULE_NAME)

# Add image dependencies so that generated_*_image_info.txt are written before checking.
$(LOCAL_CHECK_RULE_NAME): $(call images-for-partitions,$(BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_PARTITION_LIST))

$(LOCAL_CHECK_RULE_NAME): PRIVATE_NAME_UPPER := $(LOCAL_SUPER_PARTITION_NAME_UPPER)

$(LOCAL_CHECK_RULE_NAME):
	partition_size_list="$(call read-size-of-partitions,$(BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_PARTITION_LIST))"; \
	sum_sizes_expr=$$(sed -e 's/ /+/g' <<< "$${partition_size_list}"); \
	max_size_tail=$(if $(filter true,$(AB_OTA_UPDATER))," / 2"); \
	max_size_expr=$(BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_SIZE)$${max_size_tail}; \
	if [ $$(( $${sum_sizes_expr} )) -gt $$(( $${max_size_expr} )) ]; then \
		echo "The sum of sizes of all logical partitions is larger than BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_SIZE$${max_size_tail}:"; \
		echo $${sum_sizes_expr} '==' $$(( $${sum_sizes_expr} )) '>' $${max_size_expr} '==' $$(( $${max_size_expr} )); \
		exit 1; \
	else \
		echo "The sum of sizes of all logical partitions is within BOARD_$(PRIVATE_NAME_UPPER)_PARTITION_SIZE$${max_size_tail}:"; \
		echo $${sum_sizes_expr} '==' $$(( $${sum_sizes_expr} )) '<=' $${max_size_expr} '==' $$(( $${max_size_expr} )); \
	fi

endif # BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_PARTITION_LIST
endif # BOARD_$(LOCAL_SUPER_PARTITION_NAME_UPPER)_PARTITION_SIZE
endif # PRODUCT_BUILD_SUPER_PARTITION
endif # TARGET_BUILD_APPS

.PHONY: $(LOCAL_SUPER_PARTITION_NAME)image
$(LOCAL_SUPER_PARTITION_NAME)image: $(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_TARGET)

.PHONY: $(LOCAL_SUPER_PARTITION_NAME)image_empty
$(LOCAL_SUPER_PARTITION_NAME)image_empty: $(INSTALLED_$(LOCAL_SUPER_PARTITION_NAME_UPPER)IMAGE_EMPTY_TARGET)

LOCAL_CHECK_RULE_NAME :=
LOCAL_SUPER_PARTITION_NAME :=
LOCAL_SUPER_PARTITION_NAME_UPPER :=
