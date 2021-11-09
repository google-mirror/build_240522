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
import test_utils
from add_img_to_target_files import (
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    AddCareMapTxtForAbOta, AddPackRadioImages, AddRadioImagesForAbOta,
    GetCareMap)
=======
    AddPackRadioImages,
    CheckAbOtaImages)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
from rangelib import RangeSet
from common import AddCareMapForAbOta, GetCareMap


OPTIONS = common.OPTIONS


class AddImagesToTargetFilesTest(unittest.TestCase):

  def setUp(self):
    OPTIONS.input_tmp = common.MakeTempDir()

  def tearDown(self):
    common.Cleanup()

  @staticmethod
  def _create_images(images, prefix):
    """Creates images under OPTIONS.input_tmp/prefix."""
    path = os.path.join(OPTIONS.input_tmp, prefix)
    if not os.path.exists(path):
      os.mkdir(path)

    for image in images:
      image_path = os.path.join(path, image + '.img')
      with open(image_path, 'wb') as image_fp:
        image_fp.write(image.encode())

    images_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    if not os.path.exists(images_path):
      os.mkdir(images_path)
    return images, images_path

  def test_AddRadioImagesForAbOta_imageExists(self):
    """Tests the case with existing images under IMAGES/."""
    images, images_path = self._create_images(['aboot', 'xbl'], 'IMAGES')
    AddRadioImagesForAbOta(None, images)

    for image in images:
      self.assertTrue(
          os.path.exists(os.path.join(images_path, image + '.img')))

  def test_AddRadioImagesForAbOta_copyFromRadio(self):
    """Tests the case that copies images from RADIO/."""
    images, images_path = self._create_images(['aboot', 'xbl'], 'RADIO')
    AddRadioImagesForAbOta(None, images)

    for image in images:
      self.assertTrue(
          os.path.exists(os.path.join(images_path, image + '.img')))

  def test_AddRadioImagesForAbOta_copyFromRadio_zipOutput(self):
    images, _ = self._create_images(['aboot', 'xbl'], 'RADIO')

    # Set up the output zip.
    output_file = common.MakeTempFile(suffix='.zip')
    with zipfile.ZipFile(output_file, 'w') as output_zip:
      AddRadioImagesForAbOta(output_zip, images)

    with zipfile.ZipFile(output_file, 'r') as verify_zip:
      for image in images:
        self.assertIn('IMAGES/' + image + '.img', verify_zip.namelist())

  def test_AddRadioImagesForAbOta_copyFromVendorImages(self):
    """Tests the case that copies images from VENDOR_IMAGES/."""
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
    images_path = os.path.join(OPTIONS.input_tmp, 'IMAGES')
    os.mkdir(images_path)

    AddRadioImagesForAbOta(None, partitions)

    for partition in partitions:
      self.assertTrue(
          os.path.exists(os.path.join(images_path, partition + '.img')))

  def test_AddRadioImagesForAbOta_missingImages(self):
    images, _ = self._create_images(['aboot', 'xbl'], 'RADIO')
    self.assertRaises(AssertionError, AddRadioImagesForAbOta, None,
                      images + ['baz'])

  def test_AddRadioImagesForAbOta_missingImages_zipOutput(self):
    images, _ = self._create_images(['aboot', 'xbl'], 'RADIO')

    # Set up the output zip.
    output_file = common.MakeTempFile(suffix='.zip')
    with zipfile.ZipFile(output_file, 'w') as output_zip:
      self.assertRaises(AssertionError, AddRadioImagesForAbOta, output_zip,
                        images + ['baz'])

  def test_AddPackRadioImages(self):
    images, images_path = self._create_images(['foo', 'bar'], 'RADIO')
    AddPackRadioImages(None, images)

    for image in images:
      self.assertTrue(
          os.path.exists(os.path.join(images_path, image + '.img')))

  def test_AddPackRadioImages_with_suffix(self):
    images, images_path = self._create_images(['foo', 'bar'], 'RADIO')
    images_with_suffix = [image + '.img' for image in images]
    AddPackRadioImages(None, images_with_suffix)

    for image in images:
      self.assertTrue(
          os.path.exists(os.path.join(images_path, image + '.img')))

  def test_AddPackRadioImages_zipOutput(self):
    images, _ = self._create_images(['foo', 'bar'], 'RADIO')

    # Set up the output zip.
    output_file = common.MakeTempFile(suffix='.zip')
    with zipfile.ZipFile(output_file, 'w') as output_zip:
      AddPackRadioImages(output_zip, images)

    with zipfile.ZipFile(output_file, 'r') as verify_zip:
      for image in images:
        self.assertIn('IMAGES/' + image + '.img', verify_zip.namelist())

  def test_AddPackRadioImages_imageExists(self):
    images, images_path = self._create_images(['foo', 'bar'], 'RADIO')

    # Additionally create images under IMAGES/ so that they should be skipped.
    images, images_path = self._create_images(['foo', 'bar'], 'IMAGES')

    AddPackRadioImages(None, images)

    for image in images:
      self.assertTrue(
          os.path.exists(os.path.join(images_path, image + '.img')))

  def test_AddPackRadioImages_missingImages(self):
    images, _ = self._create_images(['foo', 'bar'], 'RADIO')
    AddPackRadioImages(None, images)

    self.assertRaises(AssertionError, AddPackRadioImages, None,
                      images + ['baz'])

  @staticmethod
  def _test_AddCareMapTxtForAbOta():
    """Helper function to set up the test for test_AddCareMapTxtForAbOta()."""
    OPTIONS.info_dict = {
        'system_verity_block_device' : '/dev/block/system',
        'vendor_verity_block_device' : '/dev/block/vendor',
    }

    # Prepare the META/ folder.
    meta_path = os.path.join(OPTIONS.input_tmp, 'META')
    if not os.path.exists(meta_path):
      os.mkdir(meta_path)

    system_image = test_utils.construct_sparse_image([
        (0xCAC1, 6),
        (0xCAC3, 4),
        (0xCAC1, 6)])
    vendor_image = test_utils.construct_sparse_image([
        (0xCAC2, 10)])

    image_paths = {
        'system' : system_image,
        'vendor' : vendor_image,
    }
    return image_paths

  def test_AddCareMapTxtForAbOta(self):
    image_paths = self._test_AddCareMapTxtForAbOta()

    AddCareMapTxtForAbOta(None, ['system', 'vendor'], image_paths)

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.txt')
    with open(care_map_file, 'r') as verify_fp:
      care_map = verify_fp.read()

    lines = care_map.split('\n')
    self.assertEqual(4, len(lines))
    self.assertEqual('system', lines[0])
    self.assertEqual(RangeSet("0-5 10-15").to_string_raw(), lines[1])
    self.assertEqual('vendor', lines[2])
    self.assertEqual(RangeSet("0-9").to_string_raw(), lines[3])

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  def test_AddCareMapTxtForAbOta_withNonCareMapPartitions(self):
=======
    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(care_map_file, ['system', 'vendor'], image_paths)

    expected = ['system', RangeSet("0-5 10-15").to_string_raw(),
                "ro.system.build.fingerprint",
                "google/sailfish/12345:user/dev-keys",
                'vendor', RangeSet("0-9").to_string_raw(),
                "ro.vendor.build.fingerprint",
                "google/sailfish/678:user/dev-keys"]

    self._verifyCareMap(expected, care_map_file)

  @test_utils.SkipIfExternalToolsUnavailable()
  def test_AddCareMapForAbOta_withNonCareMapPartitions(self):
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
    """Partitions without care_map should be ignored."""
    image_paths = self._test_AddCareMapTxtForAbOta()

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    AddCareMapTxtForAbOta(
        None, ['boot', 'system', 'vendor', 'vbmeta'], image_paths)

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.txt')
    with open(care_map_file, 'r') as verify_fp:
      care_map = verify_fp.read()
=======
    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(
        care_map_file, ['boot', 'system', 'vendor', 'vbmeta'], image_paths)

    expected = ['system', RangeSet("0-5 10-15").to_string_raw(),
                "ro.system.build.fingerprint",
                "google/sailfish/12345:user/dev-keys",
                'vendor', RangeSet("0-9").to_string_raw(),
                "ro.vendor.build.fingerprint",
                "google/sailfish/678:user/dev-keys"]
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

    lines = care_map.split('\n')
    self.assertEqual(4, len(lines))
    self.assertEqual('system', lines[0])
    self.assertEqual(RangeSet("0-5 10-15").to_string_raw(), lines[1])
    self.assertEqual('vendor', lines[2])
    self.assertEqual(RangeSet("0-9").to_string_raw(), lines[3])

  def test_AddCareMapTxtForAbOta_withAvb(self):
    """Tests the case for device using AVB."""
    image_paths = self._test_AddCareMapTxtForAbOta()
    OPTIONS.info_dict = {
        'avb_system_hashtree_enable' : 'true',
        'avb_vendor_hashtree_enable' : 'true',
    }

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    AddCareMapTxtForAbOta(None, ['system', 'vendor'], image_paths)

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.txt')
    with open(care_map_file, 'r') as verify_fp:
      care_map = verify_fp.read()
=======
    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(care_map_file, ['system', 'vendor'], image_paths)

    expected = ['system', RangeSet("0-5 10-15").to_string_raw(),
                "ro.system.build.fingerprint",
                "google/sailfish/12345:user/dev-keys",
                'vendor', RangeSet("0-9").to_string_raw(),
                "ro.vendor.build.fingerprint",
                "google/sailfish/678:user/dev-keys"]
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

    lines = care_map.split('\n')
    self.assertEqual(4, len(lines))
    self.assertEqual('system', lines[0])
    self.assertEqual(RangeSet("0-5 10-15").to_string_raw(), lines[1])
    self.assertEqual('vendor', lines[2])
    self.assertEqual(RangeSet("0-9").to_string_raw(), lines[3])

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  def test_AddCareMapTxtForAbOta_verityNotEnabled(self):
    """No care_map.txt should be generated if verity not enabled."""
    image_paths = self._test_AddCareMapTxtForAbOta()
=======
  @test_utils.SkipIfExternalToolsUnavailable()
  def test_AddCareMapForAbOta_noFingerprint(self):
    """Tests the case for partitions without fingerprint."""
    image_paths = self._test_AddCareMapForAbOta()
    OPTIONS.info_dict = {
        'extfs_sparse_flag' : '-s',
        'system_image_size' : 65536,
        'vendor_image_size' : 40960,
        'system_verity_block_device': '/dev/block/system',
        'vendor_verity_block_device': '/dev/block/vendor',
    }

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(care_map_file, ['system', 'vendor'], image_paths)

    expected = ['system', RangeSet("0-5 10-15").to_string_raw(), "unknown",
                "unknown", 'vendor', RangeSet("0-9").to_string_raw(), "unknown",
                "unknown"]

    self._verifyCareMap(expected, care_map_file)

  @test_utils.SkipIfExternalToolsUnavailable()
  def test_AddCareMapForAbOta_withThumbprint(self):
    """Tests the case for partitions with thumbprint."""
    image_paths = self._test_AddCareMapForAbOta()
    OPTIONS.info_dict = {
        'extfs_sparse_flag': '-s',
        'system_image_size': 65536,
        'vendor_image_size': 40960,
        'system_verity_block_device': '/dev/block/system',
        'vendor_verity_block_device': '/dev/block/vendor',
        'system.build.prop': common.PartitionBuildProps.FromDictionary(
            'system', {
                'ro.system.build.thumbprint':
                'google/sailfish/123:user/dev-keys'}
        ),
        'vendor.build.prop': common.PartitionBuildProps.FromDictionary(
            'vendor', {
                'ro.vendor.build.thumbprint':
                'google/sailfish/456:user/dev-keys'}
        ),
    }

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(care_map_file, ['system', 'vendor'], image_paths)

    expected = ['system', RangeSet("0-5 10-15").to_string_raw(),
                "ro.system.build.thumbprint",
                "google/sailfish/123:user/dev-keys",
                'vendor', RangeSet("0-9").to_string_raw(),
                "ro.vendor.build.thumbprint",
                "google/sailfish/456:user/dev-keys"]

    self._verifyCareMap(expected, care_map_file)

  @test_utils.SkipIfExternalToolsUnavailable()
  def test_AddCareMapForAbOta_skipPartition(self):
    image_paths = self._test_AddCareMapForAbOta()

    # Remove vendor_image_size to invalidate the care_map for vendor.img.
    del OPTIONS.info_dict['vendor_image_size']

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(care_map_file, ['system', 'vendor'], image_paths)

    expected = ['system', RangeSet("0-5 10-15").to_string_raw(),
                "ro.system.build.fingerprint",
                "google/sailfish/12345:user/dev-keys"]

    self._verifyCareMap(expected, care_map_file)

  @test_utils.SkipIfExternalToolsUnavailable()
  def test_AddCareMapForAbOta_skipAllPartitions(self):
    image_paths = self._test_AddCareMapForAbOta()

    # Remove the image_size properties for all the partitions.
    del OPTIONS.info_dict['system_image_size']
    del OPTIONS.info_dict['vendor_image_size']

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(care_map_file, ['system', 'vendor'], image_paths)

    self.assertFalse(os.path.exists(care_map_file))

  def test_AddCareMapForAbOta_verityNotEnabled(self):
    """No care_map.pb should be generated if verity not enabled."""
    image_paths = self._test_AddCareMapForAbOta()
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
    OPTIONS.info_dict = {}
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    AddCareMapTxtForAbOta(None, ['system', 'vendor'], image_paths)

    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.txt')
=======
    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    AddCareMapForAbOta(care_map_file, ['system', 'vendor'], image_paths)

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
    self.assertFalse(os.path.exists(care_map_file))

  def test_AddCareMapTxtForAbOta_missingImageFile(self):
    """Missing image file should be considered fatal."""
    image_paths = self._test_AddCareMapTxtForAbOta()
    image_paths['vendor'] = ''
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    self.assertRaises(AssertionError, AddCareMapTxtForAbOta, None,
=======
    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.pb')
    self.assertRaises(common.ExternalError, AddCareMapForAbOta, care_map_file,
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
                      ['system', 'vendor'], image_paths)

  def test_AddCareMapTxtForAbOta_zipOutput(self):
    """Tests the case with ZIP output."""
    image_paths = self._test_AddCareMapTxtForAbOta()

    output_file = common.MakeTempFile(suffix='.zip')
    with zipfile.ZipFile(output_file, 'w') as output_zip:
      AddCareMapTxtForAbOta(output_zip, ['system', 'vendor'], image_paths)

    with zipfile.ZipFile(output_file, 'r') as verify_zip:
      care_map = verify_zip.read('META/care_map.txt').decode('ascii')

    lines = care_map.split('\n')
    self.assertEqual(4, len(lines))
    self.assertEqual('system', lines[0])
    self.assertEqual(RangeSet("0-5 10-15").to_string_raw(), lines[1])
    self.assertEqual('vendor', lines[2])
    self.assertEqual(RangeSet("0-9").to_string_raw(), lines[3])

  def test_AddCareMapTxtForAbOta_zipOutput_careMapEntryExists(self):
    """Tests the case with ZIP output which already has care_map entry."""
    image_paths = self._test_AddCareMapTxtForAbOta()

    output_file = common.MakeTempFile(suffix='.zip')
    with zipfile.ZipFile(output_file, 'w') as output_zip:
      # Create an existing META/care_map.txt entry.
      common.ZipWriteStr(output_zip, 'META/care_map.txt', 'dummy care_map.txt')

      # Request to add META/care_map.txt again.
      AddCareMapTxtForAbOta(output_zip, ['system', 'vendor'], image_paths)

    # The one under OPTIONS.input_tmp must have been replaced.
    care_map_file = os.path.join(OPTIONS.input_tmp, 'META', 'care_map.txt')
    with open(care_map_file, 'r') as verify_fp:
      care_map = verify_fp.read()

    lines = care_map.split('\n')
    self.assertEqual(4, len(lines))
    self.assertEqual('system', lines[0])
    self.assertEqual(RangeSet("0-5 10-15").to_string_raw(), lines[1])
    self.assertEqual('vendor', lines[2])
    self.assertEqual(RangeSet("0-9").to_string_raw(), lines[3])

    # The existing entry should be scheduled to be replaced.
    self.assertIn('META/care_map.txt', OPTIONS.replace_updated_files_list)

  def test_GetCareMap(self):
    sparse_image = test_utils.construct_sparse_image([
        (0xCAC1, 6),
        (0xCAC3, 4),
        (0xCAC1, 6)])
    OPTIONS.info_dict = {
        'system_adjusted_partition_size' : 12,
    }
    name, care_map = GetCareMap('system', sparse_image)
    self.assertEqual('system', name)
    self.assertEqual(RangeSet("0-5 10-12").to_string_raw(), care_map)

  def test_GetCareMap_invalidPartition(self):
    self.assertRaises(AssertionError, GetCareMap, 'oem', None)

  def test_GetCareMap_invalidAdjustedPartitionSize(self):
    sparse_image = test_utils.construct_sparse_image([
        (0xCAC1, 6),
        (0xCAC3, 4),
        (0xCAC1, 6)])
    OPTIONS.info_dict = {
        'system_adjusted_partition_size' : -12,
    }
    self.assertRaises(AssertionError, GetCareMap, 'system', sparse_image)
