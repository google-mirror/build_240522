
import copy
import zipfile

import common
import test_utils

from non_ab_ota import NonAbOtaPropertyFiles, WriteFingerprintAssertion
from test_utils import PropertyFilesTestCase


class NonAbOtaPropertyFilesTest(PropertyFilesTestCase):
  """Additional validity checks specialized for NonAbOtaPropertyFiles."""
  def setUp(self):
     common.OPTIONS.no_signing = False
  def test_init(self):
    property_files = NonAbOtaPropertyFiles()
    self.assertEqual('ota-property-files', property_files.name)
    self.assertEqual((), property_files.required)
    self.assertEqual((), property_files.optional)

  def test_Compute(self):
    entries = ()
    zip_file = self.construct_zip_package(entries)
    property_files = NonAbOtaPropertyFiles()
    with zipfile.ZipFile(zip_file) as zip_fp:
      property_files_string = property_files.Compute(zip_fp)

    tokens = self._parse_property_files_string(property_files_string)
    self.assertEqual(1, len(tokens))
    self._verify_entries(zip_file, tokens, entries)

  def test_Finalize(self):
    entries = [
        'META-INF/com/android/metadata',
    ]
    zip_file = self.construct_zip_package(entries)
    property_files = NonAbOtaPropertyFiles()
    with zipfile.ZipFile(zip_file) as zip_fp:
      raw_metadata = property_files.GetPropertyFilesString(
          zip_fp, reserve_space=False)
      property_files_string = property_files.Finalize(zip_fp, len(raw_metadata))
    tokens = self._parse_property_files_string(property_files_string)

    self.assertEqual(1, len(tokens))
    # 'META-INF/com/android/metadata' will be key'd as 'metadata'.
    entries[0] = 'metadata'
    self._verify_entries(zip_file, tokens, entries)

  def test_Verify(self):
    entries = (
        'META-INF/com/android/metadata',
    )
    zip_file = self.construct_zip_package(entries)
    property_files = NonAbOtaPropertyFiles()
    with zipfile.ZipFile(zip_file) as zip_fp:
      raw_metadata = property_files.GetPropertyFilesString(
          zip_fp, reserve_space=False)

      property_files.Verify(zip_fp, raw_metadata)

class NonAbOTATest(test_utils.ReleaseToolsTestCase):
  TEST_TARGET_INFO_DICT = {
      'build.prop': common.PartitionBuildProps.FromDictionary(
          'system', {
              'ro.product.device': 'product-device',
              'ro.build.fingerprint': 'build-fingerprint-target',
              'ro.build.version.incremental': 'build-version-incremental-target',
              'ro.build.version.sdk': '27',
              'ro.build.version.security_patch': '2017-12-01',
              'ro.build.date.utc': '1500000000'}
      )
  }
  TEST_INFO_DICT_USES_OEM_PROPS = {
      'build.prop': common.PartitionBuildProps.FromDictionary(
          'system', {
              'ro.product.name': 'product-name',
              'ro.build.thumbprint': 'build-thumbprint',
              'ro.build.bar': 'build-bar'}
      ),
      'vendor.build.prop': common.PartitionBuildProps.FromDictionary(
          'vendor', {
               'ro.vendor.build.fingerprint': 'vendor-build-fingerprint'}
      ),
      'property1': 'value1',
      'property2': 4096,
      'oem_fingerprint_properties': 'ro.product.device ro.product.brand',
  }
  TEST_OEM_DICTS = [
      {
          'ro.product.brand': 'brand1',
          'ro.product.device': 'device1',
      },
      {
          'ro.product.brand': 'brand2',
          'ro.product.device': 'device2',
      },
      {
          'ro.product.brand': 'brand3',
          'ro.product.device': 'device3',
      },
  ]
  def test_WriteFingerprintAssertion_without_oem_props(self):
    target_info = common.BuildInfo(self.TEST_TARGET_INFO_DICT, None)
    source_info_dict = copy.deepcopy(self.TEST_TARGET_INFO_DICT)
    source_info_dict['build.prop'].build_props['ro.build.fingerprint'] = (
        'source-build-fingerprint')
    source_info = common.BuildInfo(source_info_dict, None)

    script_writer = test_utils.MockScriptWriter()
    WriteFingerprintAssertion(script_writer, target_info, source_info)
    self.assertEqual(
        [('AssertSomeFingerprint', 'source-build-fingerprint',
          'build-fingerprint-target')],
        script_writer.lines)

  def test_WriteFingerprintAssertion_with_source_oem_props(self):
    target_info = common.BuildInfo(self.TEST_TARGET_INFO_DICT, None)
    source_info = common.BuildInfo(self.TEST_INFO_DICT_USES_OEM_PROPS,
                                   self.TEST_OEM_DICTS)

    script_writer = test_utils.MockScriptWriter()
    WriteFingerprintAssertion(script_writer, target_info, source_info)
    self.assertEqual(
        [('AssertFingerprintOrThumbprint', 'build-fingerprint-target',
          'build-thumbprint')],
        script_writer.lines)

  def test_WriteFingerprintAssertion_with_target_oem_props(self):
    target_info = common.BuildInfo(self.TEST_INFO_DICT_USES_OEM_PROPS,
                                   self.TEST_OEM_DICTS)
    source_info = common.BuildInfo(self.TEST_TARGET_INFO_DICT, None)

    script_writer = test_utils.MockScriptWriter()
    WriteFingerprintAssertion(script_writer, target_info, source_info)
    self.assertEqual(
        [('AssertFingerprintOrThumbprint', 'build-fingerprint-target',
          'build-thumbprint')],
        script_writer.lines)

  def test_WriteFingerprintAssertion_with_both_oem_props(self):
    target_info = common.BuildInfo(self.TEST_INFO_DICT_USES_OEM_PROPS,
                                   self.TEST_OEM_DICTS)
    source_info_dict = copy.deepcopy(self.TEST_INFO_DICT_USES_OEM_PROPS)
    source_info_dict['build.prop'].build_props['ro.build.thumbprint'] = (
        'source-build-thumbprint')
    source_info = common.BuildInfo(source_info_dict, self.TEST_OEM_DICTS)

    script_writer = test_utils.MockScriptWriter()
    WriteFingerprintAssertion(script_writer, target_info, source_info)
    self.assertEqual(
        [('AssertSomeThumbprint', 'build-thumbprint',
          'source-build-thumbprint')],
        script_writer.lines)
