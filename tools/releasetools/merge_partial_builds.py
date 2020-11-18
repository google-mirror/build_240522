#!/usr/bin/env python
#
# Copyright (C) 2019 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not
# use this file except in compliance with the License. You may obtain a copy of
# the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations under
# the License.
#
"""Given a set of input target-files, produces any of an image zipfile suitable for use with 'fastboot update', an OTA package, or a super.img.

Usage:  merge_partial_builds --input_file partial_target_file --input_file
partial_target_file  --output_img output_image_zip --output_ota output_ota_zip
--output_super output_super_image

Flags:
  --input_file partial_target_file
      The path to a partial build target files archive that is to be merged.
      This argument can be repeated.

  --output_img output_image_zip
      The path to the output img package. This package is suitable for use with
      'fastboot update'.

  --output_ota output_ota_zip
      The path to the output ota package.

  --output_super output_super_image
      The path to the output super image.
"""

import glob
import itertools
import json
import logging
import os
import shutil
import sys

import build_super_image
import check_target_files_vintf
import common
import find_shareduid_violation
import img_from_target_files
import ota_from_target_files

logger = logging.getLogger(__name__)

OPTIONS = common.OPTIONS

OPTIONS.input_files = []
OPTIONS.output_img = None
OPTIONS.output_ota = None
OPTIONS.output_super = None
OPTIONS.keep_tmp = False


def MergeInfoDictsForReleaseTools(input_files, output_file):
  """Loads and merges the key/value pairs from the given list of input target_files.

  It reads `META/misc_info.txt` file in the target_files input, does validation
  checks and returns the parsed key/value pairs for to the given builds. It's
  usually called early when working on multiple partial input target_files
  files,
  e.g. when generating OTAs. Note that the function may be called against an old
  target_files file (i.e. from past dessert releases). So the property parsing
  needs to be backward compatible.

  Args:
    input_files: The input target_files file, which could be an open
      zipfile.ZipFile instance, or a str for the dir that contains the files
      unzipped from a target_files file.
    output_file: The path to output the merged dictionary of key/value pairs

  Returns:
    A dict that contains the parsed key/value pairs suitable to call the
    following
    tools: build_super_image.py, img_from_target_files.py, and
    ota_from_target_files.py.
  """

  info_dicts = []
  for input_file in input_files:
    info_dicts.append(
        common.LoadDictionaryFromFile(
            os.path.join(input_file, 'META', 'misc_info.txt')))

  merged_dict = common.MergeDynamicPartitionInfoDicts(info_dicts)

  ab_partitions = common.UniqConcat(
      [d.get('ab_partitions', '') for d in info_dicts])
  if ab_partitions:
    merged_dict['ab_partitions'] = ab_partitions

  avb_custom_images_partition_list = common.UniqConcat(
      [d.get('avb_custom_images_partition_list', '') for d in info_dicts])
  if avb_custom_images_partition_list:
    merged_dict[
        'avb_custom_images_partition_list'] = avb_custom_images_partition_list

  for partition_name in common.AVB_VBMETA_PARTITIONS:
    concat = common.UniqConcat(
        [d.get('avb_{}'.format(partition_name), '') for d in info_dicts])
    if concat:
      merged_dict[partition_name] = concat

  for partition_name in common.AVB_VBMETA_PARTITIONS + common.AVB_PARTITIONS:
    for key in ('avb_{}_key_path', 'avb_{}_algorithm', 'avb_{}_args'):
      key = key.format(partition_name)
      is_consistent, value = common.IsValueConsistent(info_dicts, key)
      if is_consistent:
        merged_dict[key] = value

  # Add misc keys needed for release tools
  for key in ('avb_avbtool', 'super_image_in_update_package',
              'bootloader_in_update_package', 'dynamic_partition_retrofit',
              'build_super_partition', 'extfs_sparse_flag', 'vendor.build.prop',
              'vintf_enforce', 'vintf_odm_manifest_skus',
              'vintf_vendor_manifest_skus', 'vintf_include_empty_odm_sku',
              'vintf_include_empty_vendor_sku', 'recovery_api_version',
              'fstab_version', 'allow_non_ab', 'ab_update', 'cache_size',
              'default_system_dev_certificate', 'no_recovery', 'tool_extension',
              'board_uses_vendorimage'):
    is_consistent, value = common.IsValueConsistent(info_dicts, key)
    if is_consistent:
      merged_dict[key] = value

  common.WriteSortedData(merged_dict,
                         os.path.join(output_file, 'META', 'misc_info.txt'))

  # CheckVintf requires the vendor build props
  for input_file in input_files:
    if os.path.exists(os.path.join(
        input_file, 'VENDOR', 'build.prop')) or os.path.exists(
            os.path.join(input_file, 'VENDOR', 'etc', 'build.prop')):
      if 'vendor.build.prop' not in merged_dict:
        merged_dict[
            'vendor.build.prop'] = common.PartitionBuildProps.FromInputFile(
                input_file, 'vendor')
        logger.warning(merged_dict['vendor.build.prop'])
      else:
        raise ValueError(
            'Found two vendor build.prop files! Expected only one in combined input.'
        )

  return merged_dict


def CheckTargetFilesCompatibility(input_files, target_files_dir, info_dict):
  """Run Treble compatibility tests on the merge target files archive.

  Args:
    input_files: The set of extracted input target_files file.
    target_files_dir: The path to the temporary target files directory.
    info_dict: The merged info dict.

  Raises:
    RuntimeError: If CheckVintf fails.
    ValueError: If shared UID violations are present.
  """

  if not check_target_files_vintf.CheckVintf(target_files_dir, info_dict):
    raise RuntimeError('Incompatible VINTF metadata')

  partition_groups = []
  partition_map = {}
  for input_file in input_files:
    image_glob = glob.glob(os.path.join(input_file, '*'))
    partition_groups.append([os.path.basename(image).lower() for image in image_glob])

  partition_map = common.PartitionMapFromTargetFiles(target_files_dir)

  # Generate and check for cross-partition violations of sharedUserId
  # values in APKs. This requires the input target-files packages to contain
  # *.apk files.
  shareduid_violation_modules = os.path.join(
      target_files_dir, 'META', 'shareduid_violation_modules.json')
  with open(shareduid_violation_modules, 'w') as f:
    violation = find_shareduid_violation.FindShareduidViolation(
        target_files_dir, partition_map)

    # Write the output to a file to enable debugging.
    f.write(violation)

    # Check for violations across the input builds' partition groups.
    shareduid_errors = common.SharedUidPartitionViolations(
        json.loads(violation), partition_groups)
    if shareduid_errors:
      for error in shareduid_errors:
        logger.error(error)
      raise ValueError('sharedUserId APK error. See %s' %
                       shareduid_violation_modules)
  # Run host_init_verifier on the combined init rc files.
  filtered_partitions = {
      partition: path
      for partition, path in partition_map.items()
      # host_init_verifier checks only the following partitions:
      if partition in ['system', 'system_ext', 'product', 'vendor', 'odm']
  }
  common.RunHostInitVerifier(
      product_out=target_files_dir,
      partition_map=filtered_partitions)


def MergeTargetFilesForReleaseTools(input_files, temp_dir):
  """Merges a set of partial builds target files archives.

  Creates a temporary, merged target files archive for use with release tools.
  This includes merging any information (e.g. misc_info.txt,
  dynamic_partitions.txt,
  and ab_partitions.txt) needed for the release tools to produce the needed
  build
  artifacts. The merged target files archive will only include the minimum
  information
  needed to allow the release tools to operate. All other information will be
  excluded.

  Args:
    input_files: The set of input target_files.
    temp_dir: The path to the temporary directory used for output.

  Returns:
    A minimally merged target files archive in zip format.

  Raises;
    ValueError: If the union of the contents of the set of target files archives
    is not empty.
  """

  logger.info('In Merge')

  temp_target_files_dir = os.path.join(temp_dir, 'output')
  os.mkdir(temp_target_files_dir)
  os.mkdir(os.path.join(temp_target_files_dir, 'META'))
  os.mkdir(os.path.join(temp_target_files_dir, 'IMAGES'))

  # Extract the needed files and images from each target files archive.
  # The rough huerstic used here is that if the archive has built
  # <partition>.img then we assume it authoritative for that partition and
  # copy both the image itself and the partition directory to the temporary
  # merged target files.
  for input_dir in input_files:
    image_glob = [
        os.path.basename(g)
        for g in glob.glob(os.path.join(input_dir, 'IMAGES', '*.img'))
    ]
    for directory in glob.glob(os.path.join(input_dir, '*')):
      image = '{}.img'.format(os.path.basename(directory)).lower()
      if image in image_glob:
        if os.path.exists(os.path.join(temp_target_files_dir, 'IMAGES', image)):
          raise ValueError('{} already exists in output!'.format(image))
        shutil.copytree(
            directory,
            os.path.join(temp_target_files_dir, os.path.basename(directory)),
            symlinks=True)
        shutil.copyfile(
            os.path.join(input_dir, 'IMAGES', image),
            os.path.join(temp_target_files_dir, 'IMAGES', image))
        if image == 'vendor.img':
          # We assume that the partial build for the vendor image is
          # authortative for android-info.txt and kernel info
          shutil.copytree(
              os.path.join(input_dir, 'OTA'),
              os.path.join(temp_target_files_dir, 'OTA'))
          shutil.copyfile(
              os.path.join(input_dir, 'META', 'kernel_version.txt'),
              os.path.join(temp_target_files_dir, 'META', 'kernel_version.txt'))
          shutil.copyfile(
              os.path.join(input_dir, 'META', 'kernel_configs.txt'),
              os.path.join(temp_target_files_dir, 'META', 'kernel_configs.txt'))
          if os.path.exists(os.path.join(input_dir, 'IMAGES', 'vbmeta_vendor.img')):
            shutil.copyfile(
                os.path.join(input_dir, 'IMAGES', 'vbmeta_system.img'),
                os.path.join(temp_target_files_dir, 'IMAGES', 'vbmeta_vendor.img'))
        if image == 'system.img':
          if os.path.exists(os.path.join(input_dir, 'IMAGES', 'vbmeta_vendor.img')):
            shutil.copyfile(
                os.path.join(input_dir, 'IMAGES', 'vbmeta_system.img'),
                os.path.join(temp_target_files_dir, 'IMAGES', 'vbmeta_system.img'))
    for image in ('userdata.img', 'dtbo.img', 'cache.img'):
      if image in image_glob:
        if os.path.exists(os.path.join(temp_target_files_dir, 'IMAGES', image)):
          raise ValueError('{} already exists in output!'.format(image))
        shutil.copyfile(
            os.path.join(input_dir, 'IMAGES', image),
            os.path.join(temp_target_files_dir, 'IMAGES', image))

  merged_info_dict = MergeInfoDictsForReleaseTools(input_files,
                                                   temp_target_files_dir)

  common.MergeABPartitionsTxt(input_files, temp_target_files_dir)

  common.MergeDynamicPartitionsInfoTxt(input_files, temp_target_files_dir)

  CheckTargetFilesCompatibility(input_files, temp_target_files_dir,
                                merged_info_dict)

  # Build vbmeta.img using partitions in product_out_vendor.
  partitions = {}
  for partition in common.AVB_PARTITIONS + common.AVB_VBMETA_PARTITIONS:
    partition_path = os.path.join(target_dir, 'IMAGES',
                                  "%s.img" % partition)
    if os.path.exists(partition_path):
      partitions[partition] = partition_path

  vbmeta_partitions = common.GetVBMetaPartitions(partitions, merged_info_dict)
  OPTIONS.info_dict = merged_info_dict
  common.BuildVBMeta(
      os.path.join(temp_target_files_dir, 'IMAGES', 'vbmeta.img'), partitions,
      'vbmeta', vbmeta_partitions)

  common.BuildSuperEmpty(temp_target_files_dir)

  temp_target_files_zip = os.path.join(temp_dir, 'temp-target-files.zip')
  common.CreateTargetFilesArchive(temp_target_files_zip, temp_target_files_dir,
                                  temp_dir)

  # Zip up temp target files
  return temp_target_files_zip


def MergePartialBuilds(input_files, output_img, output_ota, output_super):
  """Merges a set of partial builds outputting the requested artifacts.

  Creates a temporary, merged target files archive for use with release tools.
  This includes merging any information (e.g. misc_info.txt,
  dynamic_partitions.txt,
  and ab_partitions.txt) needed for the release tools to produce the needed
  build
  eartifacts.

  Args:
    input_files: The set of input target_files.
    output_img: The path to the output img package.
    output_ota: The path to the output ota package.
    output_super: The path to the output super image.
  """

  logger.info('Starting merge...')
  try:
    extracted_input_files = []
    for input_file in input_files:
      extracted_input_files.append(common.UnzipTemp(input_file))

    temp_dir = common.MakeTempDir()
    temp_target_files_zip = MergeTargetFilesForReleaseTools(
        extracted_input_files, temp_dir)

    if output_img:
      img_from_target_files.main([temp_target_files_zip, output_img])

    if output_ota:
      ota_from_target_files.main([temp_target_files_zip, output_ota])

    if output_super:
      build_super_image.main([temp_target_files_zip, output_super])

  finally:
    logger.info('Cleaning up...')
    if not OPTIONS.keep_tmp:
      common.Cleanup()
  logger.info('Done!')


def main():
  common.InitLogging()

  def OptionHandler(o, a):
    if o == '--input_file':
      OPTIONS.input_files.append(a)
    elif o == '--output_img':
      OPTIONS.output_img = a
    elif o == '--output_ota':
      OPTIONS.output_ota = a
    elif o == '--output_super':
      OPTIONS.output_super = a
    elif o == '--keep_tmp':
      OPTIONS.keep_tmp = True
    else:
      return False
    return True

  args = common.ParseOptions(
      sys.argv[1:],
      __doc__,
      extra_long_opts=['input_file=', 'output_img=', 'output_ota=', 'keep_tmp'],
      extra_option_handler=OptionHandler)

  if (args or OPTIONS.input_files is None or
      (OPTIONS.output_img is None and OPTIONS.output_ota is None)):
    common.Usage(__doc__)
    sys.exit(1)

  MergePartialBuilds(OPTIONS.input_files, OPTIONS.output_img,
                     OPTIONS.output_ota, OPTIONS.output_super)


if __name__ == '__main__':
  main()
