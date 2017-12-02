#!/usr/bin/env python3
#
# Copyright (C) 2017 The Android Open Source Project
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

import os
import os.path
import tempfile
import unittest
from unittest.mock import patch, call

import add_img_to_target_files
import common


OPTIONS = common.OPTIONS


class AddImageToTargetFilesTest(unittest.TestCase):
  def setUp(self):
    OPTIONS.input_tmp = tempfile.mkdtemp()
    OPTIONS.tempfiles = [OPTIONS.input_tmp]

  def tearDown(self):
    common.Cleanup()

  def test_AddRadioImagesForAbOta_ImageExists(self):
    """Tests the case with existing images."""
    image_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    os.mkdir(image_path)

    boot_image_path = os.path.join(image_path, 'boot.img')
    with open(boot_image_path, 'w') as boot_image_file:
      boot_image_file.write('boot')
    system_image_path = os.path.join(image_path, 'system.img')
    with open(system_image_path, 'w') as system_image_file:
      system_image_file.write('system')

    add_img_to_target_files.AddRadioImagesForAbOta(None, ['boot', 'system'])

    self.assertTrue(os.path.exists(boot_image_path))
    self.assertTrue(os.path.exists(system_image_path))

  def test_AddRadioImagesForAbOta_CopyImagesFromRadio(self):
    pass

  @patch('add_img_to_target_files.GetCareMap')
  def test_AddCareMapTxtForAbOta(self, mock_GetCareMap):
    OPTIONS.info_dict = dict()
    OPTIONS.info_dict['system_verity_block_device'] = '/dev/block/system'
    OPTIONS.info_dict['vendor_verity_block_device'] = '/dev/block/vendor'

    image_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    os.mkdir(image_path)

    meta_path = os.path.join(OPTIONS.input_tmp, 'META')
    os.mkdir(meta_path)

    image_paths = dict()

    system_image_path = os.path.join(image_path, 'system.img')
    with open(system_image_path, 'w') as system_image_file:
      system_image_file.write('system')
    image_paths['system'] = system_image_path

    vendor_image_path = os.path.join(image_path, 'vendor.img')
    with open(vendor_image_path, 'w') as vendor_image_file:
      vendor_image_file.write('system')
    image_paths['vendor'] = vendor_image_path

    mock_GetCareMap.return_value = ['foo', 'bar']

    add_img_to_target_files.AddCareMapTxtForAbOta(None, ['system', 'vendor'],
                                                  image_paths)

    mock_GetCareMap.assert_has_calls([call('system', system_image_path),
                                      call('vendor', vendor_image_path)])
    self.assertTrue(os.path.exists(os.path.join(meta_path, 'care_map.txt')))

  def test_AddCareMapTxtForAbOta_ReplaceExistingFiles(self):
    pass
