#
# Copyright (C) 2022 The Android Open Source Project
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

import os.path
import shutil

import common
import merge_target_files
import test_utils
from merge_target_files import (
    DEFAULT_FRAMEWORK_ITEM_LIST,
    DEFAULT_VENDOR_ITEM_LIST,
    DEFAULT_FRAMEWORK_MISC_INFO_KEYS,
    validate_config_lists,
    item_list_to_partition_set,
)


class MergeTargetFilesTest(test_utils.ReleaseToolsTestCase):

  def setUp(self):
    self.testdata_dir = test_utils.get_testdata_dir()
    self.OPTIONS = merge_target_files.OPTIONS
    self.OPTIONS.framework_item_list = DEFAULT_FRAMEWORK_ITEM_LIST
    self.OPTIONS.framework_misc_info_keys = DEFAULT_FRAMEWORK_MISC_INFO_KEYS
    self.OPTIONS.vendor_item_list = DEFAULT_VENDOR_ITEM_LIST

  def test_validate_config_lists_ReturnsFalseIfMissingDefaultItem(self):
    self.OPTIONS.framework_item_list = list(DEFAULT_FRAMEWORK_ITEM_LIST)
    self.OPTIONS.framework_item_list.remove('SYSTEM/*')
    self.assertFalse(validate_config_lists())

  def test_validate_config_lists_ReturnsTrueIfDefaultItemInDifferentList(self):
    self.OPTIONS.framework_item_list = list(DEFAULT_FRAMEWORK_ITEM_LIST)
    self.OPTIONS.framework_item_list.remove('ROOT/*')
    self.OPTIONS.vendor_item_list = list(DEFAULT_VENDOR_ITEM_LIST)
    self.OPTIONS.vendor_item_list.append('ROOT/*')
    self.assertTrue(validate_config_lists())

  def test_validate_config_lists_ReturnsTrueIfExtraItem(self):
    self.OPTIONS.framework_item_list = list(DEFAULT_FRAMEWORK_ITEM_LIST)
    self.OPTIONS.framework_item_list.append('MY_NEW_PARTITION/*')
    self.assertTrue(validate_config_lists())

  def test_validate_config_lists_ReturnsFalseIfSharedExtractedPartition(self):
    self.OPTIONS.vendor_item_list = list(DEFAULT_VENDOR_ITEM_LIST)
    self.OPTIONS.vendor_item_list.append('SYSTEM/my_system_file')
    self.assertFalse(validate_config_lists())

  def test_validate_config_lists_ReturnsFalseIfSharedExtractedPartitionImage(
      self):
    self.OPTIONS.vendor_item_list = list(DEFAULT_VENDOR_ITEM_LIST)
    self.OPTIONS.vendor_item_list.append('IMAGES/system.img')
    self.assertFalse(validate_config_lists())

  def test_validate_config_lists_ReturnsFalseIfBadSystemMiscInfoKeys(self):
    for bad_key in ['dynamic_partition_list', 'super_partition_groups']:
      self.OPTIONS.framework_misc_info_keys = list(
          DEFAULT_FRAMEWORK_MISC_INFO_KEYS)
      self.OPTIONS.framework_misc_info_keys.append(bad_key)
      self.assertFalse(validate_config_lists())

  def test_item_list_to_partition_set(self):
    item_list = [
        'META/apexkeys.txt',
        'META/apkcerts.txt',
        'META/filesystem_config.txt',
        'PRODUCT/*',
        'SYSTEM/*',
        'SYSTEM_EXT/*',
    ]
    partition_set = item_list_to_partition_set(item_list)
    self.assertEqual(set(['product', 'system', 'system_ext']), partition_set)
