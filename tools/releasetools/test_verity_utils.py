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

"""Unittests for verity_utils.py."""

from __future__ import print_function

import os
import os.path
import unittest

import build_image
import common
import sparse_img
import test_utils
import verity_utils
from rangelib import RangeSet


class VerityUtilsTest(unittest.TestCase):
  def setUp(self):
    self.testdata_dir = test_utils.get_testdata_dir()
    self.prop_dict = {
        'verity': 'true',
        'verity_fec': 'true',
        'system_verity_block_device': '/dev/block/system',
        'system_size': 1024 * 1024
    }

    self.hash_algorithm = "sha256"
    self.fixed_salt = \
        "aee087a5be3b982978c923f566a94613496b417f2af592639bc80d141e34dfe7"
    self.expected_root_hash = \
        "0b7c4565e87b1026e11fbab91c0bc29e185c847a5b44d40e6e86e461e8adf80d"

  def tearDown(self):
    common.Cleanup()

  def _create_simg(self, raw_data):
    output_file = common.MakeTempFile()
    raw_image = common.MakeTempFile()
    with open(raw_image, 'wb') as f:
      f.write(raw_data)

    cmd = ["img2simg", raw_image, output_file, '4096']
    p = common.Run(cmd)
    p.communicate()
    self.assertEqual(0, p.returncode)

    return output_file

  def _generate_image(self):
    partition_size = 1024 * 1024
    adjusted_size, verity_size = verity_utils.AdjustPartitionSizeForVerity(
        partition_size, True)

    raw_image = ""
    for i in range(adjusted_size):
      raw_image += str(i % 10)

    output_file = self._create_simg(raw_image)

    # Append the verity metadata.
    prop_dict = {
        'original_partition_size' : str(partition_size),
        'partition_size' : str(adjusted_size),
        'verity_block_device' : '/dev/block/system',
        'verity_key' : os.path.join(self.testdata_dir, 'testkey'),
        'verity_signer_cmd' : 'verity_signer',
        'verity_size' : str(verity_size),
    }
    self.assertTrue(
        build_image.MakeVerityEnabledImage(output_file, True, prop_dict))

    return output_file

  def test_VerityTreeInfo_init(self):
    partition_size = 1024 * 1024
    image_file = sparse_img.SparseImage(self._generate_image())

    info1 = verity_utils.VerityTreeInfo('system', image_file, self.prop_dict)
    self.assertTrue(info1.enabled)
    self.assertEqual(partition_size, info1.partition_size)
    self.assertTrue(info1.fec_supported)

    info2 = verity_utils.VerityTreeInfo('system', image_file, {})
    self.assertFalse(info2.enabled)

    self.prop_dict['system_size'] = 0
    with self.assertRaises(AssertionError):
      verity_utils.VerityTreeInfo('system', image_file, self.prop_dict)

  def test_VerityTreeInfo_getVeritySize(self):
    image_file = sparse_img.SparseImage(self._generate_image())

    info = verity_utils.VerityTreeInfo('system', image_file, self.prop_dict)
    fs_size, tree_size, meta_size = info.GetVerityDataSize()
    self.assertEqual(991232, fs_size)
    self.assertEqual(12288, tree_size)
    self.assertEqual(32768, meta_size)

  def test_VerityTreeInfo_parseMetadata(self):
    image_path = self._generate_image()
    result, unsparse_image = build_image.UnsparseImage(image_path, False)
    self.assertTrue(result)
    with open(unsparse_image, 'r') as f:
      f.seek(1003520)
      metadata = f.read(32768)

    image_file = sparse_img.SparseImage(image_path)
    info = verity_utils.VerityTreeInfo('system', image_file, self.prop_dict)
    info.ParseVerityMetadata(metadata, 991232)

    self.assertEqual(self.hash_algorithm, info.hash_algorithm)
    self.assertEqual(self.fixed_salt, info.salt)
    self.assertEqual(self.expected_root_hash, info.root_hash)

  def test_VerityTreeInfo_validateHashTree_smoke(self):
    image_file = sparse_img.SparseImage(self._generate_image())

    info = verity_utils.VerityTreeInfo('system', image_file, self.prop_dict)
    info.hash_algorithm = self.hash_algorithm
    info.salt = self.fixed_salt
    info.root_hash = self.expected_root_hash

    info.file_system_range = RangeSet(data=[0, 991232 / 4096])
    info.verity_tree_range = RangeSet(
        data=[991232 / 4096, (991232 + 12288) / 4096])
    self.assertTrue(info.ValidateVerityTree())

  def test_VerityTreeInfo_validateHashTree_failure(self):
    image_file = sparse_img.SparseImage(self._generate_image())

    info = verity_utils.VerityTreeInfo('system', image_file, self.prop_dict)
    info.hash_algorithm = self.hash_algorithm
    info.salt = self.fixed_salt
    info.root_hash = "a" + self.expected_root_hash[1:]

    info.file_system_range = RangeSet(data=[0, 991232 / 4096])
    info.verity_tree_range = RangeSet(
        data=[991232 / 4096, (991232 + 12288) / 4096])
    self.assertFalse(info.ValidateVerityTree())

  def test_VerityTreeInfo_getVerityInfo(self):
    image_file = sparse_img.SparseImage(self._generate_image())

    info = verity_utils.VerityTreeInfo('system', image_file, self.prop_dict)
    info.GetVerityTreeInfo()

    self.assertEqual(RangeSet(data=[0, 991232 / 4096]), info.file_system_range)
    self.assertEqual(RangeSet(data=[991232 / 4096, (991232 + 12288) / 4096]),
                     info.verity_tree_range)
    self.assertEqual(self.hash_algorithm, info.hash_algorithm)
    self.assertEqual(self.fixed_salt, info.salt)
    self.assertEqual(self.expected_root_hash, info.root_hash)
