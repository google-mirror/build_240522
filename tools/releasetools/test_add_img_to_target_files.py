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

import os
import os.path
import unittest
import zipfile

import common
from add_img_to_target_files import AddRadioImagesForAbOta


OPTIONS = common.OPTIONS


class AddImageToTargetFilesTest(unittest.TestCase):

  def setUp(self):
    OPTIONS.input_tmp = common.MakeTempDir()

  def tearDown(self):
    common.Cleanup()

  def test_AddRadioImagesForAbOta_imageExists(self):
    """Tests the case with existing images."""
    image_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    os.mkdir(image_path)

    partitions = ['aboot', 'xbl']
    for partition in partitions:
      partition_image_path = os.path.join(image_path, partition + '.img')
      with open(partition_image_path, 'wb') as partition_fp:
        partition_fp.write(partition.encode())

    AddRadioImagesForAbOta(None, partitions)

    for partition in partitions:
      self.assertTrue(
          os.path.exists(
              os.path.join(OPTIONS.input_tmp, 'IMAGES', partition + '.img')))

  def test_AddRadioImagesForAbOta_copyFromRadio(self):
    """Tests the case with images in RADIO/."""
    radio_path = os.path.join(OPTIONS.input_tmp, 'RADIO')
    os.mkdir(radio_path)

    partitions = ['aboot', 'xbl']
    for partition in partitions:
      partition_image_path = os.path.join(radio_path, partition + '.img')
      with open(partition_image_path, 'wb') as partition_fp:
        partition_fp.write(partition.encode())

    # Set up the output dir.
    image_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    os.mkdir(image_path)

    AddRadioImagesForAbOta(None, partitions)

    for partition in partitions:
      self.assertTrue(
          os.path.exists(os.path.join(image_path, partition + '.img')))

  def test_AddRadioImagesForAbOta_copyFromRadio_zipOutput(self):
    radio_path = os.path.join(OPTIONS.input_tmp, 'RADIO')
    os.mkdir(radio_path)

    partitions = ['aboot', 'xbl']
    for partition in partitions:
      partition_image_path = os.path.join(radio_path, partition + '.img')
      with open(partition_image_path, 'wb') as partition_fp:
        partition_fp.write(partition.encode())

    # Set up the output zip.
    output_file = common.MakeTempFile(suffix='.zip')
    with zipfile.ZipFile(output_file, 'w') as output_zip:
      AddRadioImagesForAbOta(output_zip, partitions)

    with zipfile.ZipFile(output_file, 'r') as verify_zip:
      for partition in partitions:
        self.assertTrue('IMAGES/' + partition + '.img' in verify_zip.namelist())

  def test_AddRadioImagesForAbOta_copyFromVendorImages(self):
    """Tests the case with images in VENDOR_IMAGES/."""
    vendor_images_path = os.path.join(OPTIONS.input_tmp, 'VENDOR_IMAGES')
    os.mkdir(vendor_images_path)

    partitions = ['aboot', 'xbl']
    for index, partition in enumerate(partitions):
      subdir = os.path.join(vendor_images_path, 'subdir-{}'.format(index))
      os.mkdir(subdir)

      partition_image_path = os.path.join(subdir, partition + '.img')
      with open(partition_image_path, 'wb') as partition_fp:
        partition_fp.write(partition.encode())

    # Set up the output dir.
    image_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    os.mkdir(image_path)

    AddRadioImagesForAbOta(None, partitions)

    for partition in partitions:
      self.assertTrue(
          os.path.exists(os.path.join(image_path, partition + '.img')))

  def test_AddRadioImagesForAbOta_missingImages(self):
    radio_path = os.path.join(OPTIONS.input_tmp, 'RADIO')
    os.mkdir(radio_path)

    partitions = ['aboot', 'xbl']
    for partition in partitions:
      partition_image_path = os.path.join(radio_path, partition + '.img')
      with open(partition_image_path, 'wb') as partition_fp:
        partition_fp.write(partition.encode())

    image_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    os.mkdir(image_path)

    self.assertRaises(AssertionError, AddRadioImagesForAbOta, None,
                      partitions + ['foo'])

  def test_AddRadioImagesForAbOta_missingImages_zipOutput(self):
    radio_path = os.path.join(OPTIONS.input_tmp, 'RADIO')
    os.mkdir(radio_path)

    partitions = ['aboot', 'xbl']
    for partition in partitions:
      partition_image_path = os.path.join(radio_path, partition + '.img')
      with open(partition_image_path, 'wb') as partition_fp:
        partition_fp.write(partition.encode())

    # Set up the output zip.
    output_file = common.MakeTempFile(suffix='.zip')
    with zipfile.ZipFile(output_file, 'w') as output_zip:
      self.assertRaises(AssertionError, AddRadioImagesForAbOta, output_zip,
                        partitions + ['foo'])
