#!/usr/bin/env python
#
# Copyright (C) 2022 The Android Open Source Project
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
"""This script merges two partial target files packages.

One input package contains framework files, and the other contains vendor files.

This script produces a complete, merged target files package:
  - This package can be used to generate a flashable IMG package.
    See --output-img.
  - This package can be used to generate an OTA package. See --output-ota.
  - The merged package is checked for compatibility between the two inputs.

Usage: merge_target_files [args]

  --framework-target-files framework-target-files-zip-archive
      The input target files package containing framework bits. This is a zip
      archive.

  --framework-item-list framework-item-list-file
      The optional path to a newline-separated config file of items that
      are extracted as-is from the framework target files package.

  --framework-misc-info-keys framework-misc-info-keys-file
      The optional path to a newline-separated config file of keys to
      extract from the framework META/misc_info.txt file.

  --vendor-target-files vendor-target-files-zip-archive
      The input target files package containing vendor bits. This is a zip
      archive.

  --vendor-item-list vendor-item-list-file
      The optional path to a newline-separated config file of items that
      are extracted as-is from the vendor target files package.

  --output-target-files output-target-files-package
      If provided, the output merged target files package. Also a zip archive.

  --output-dir output-directory
      If provided, the destination directory for saving merged files. Requires
      the --output-item-list flag.
      Can be provided alongside --output-target-files, or by itself.

  --output-item-list output-item-list-file.
      The optional path to a newline-separated config file that specifies the
      file patterns to copy into the --output-dir. Required if providing
      the --output-dir flag.

  --output-ota output-ota-package
      The output ota package. This is a zip archive. Use of this flag may
      require passing the --path common flag; see common.py.

  --output-img output-img-package
      The output img package, suitable for use with 'fastboot update'. Use of
      this flag may require passing the --path common flag; see common.py.

  --output-super-empty output-super-empty-image
      If provided, creates a super_empty.img file from the merged target
      files package and saves it at this path.

  --rebuild_recovery
      Copy the recovery image used by non-A/B devices, used when
      regenerating vendor images with --rebuild-sepolicy.

  --allow-duplicate-apkapex-keys
      If provided, duplicate APK/APEX keys are ignored and the value from the
      framework is used.

  --rebuild-sepolicy
      If provided, rebuilds odm.img or vendor.img to include merged sepolicy
      files. If odm is present then odm is preferred.

  --vendor-otatools otatools.zip
      If provided, use this otatools.zip when recompiling the odm or vendor
      image to include sepolicy.

  --keep-tmp
      Keep tempoary files for debugging purposes.

  The following only apply when using the VSDK to perform dexopt on vendor apps:

  --framework-dexpreopt-config
      If provided, the location of framwework's dexpreopt_config.zip.

  --framework-dexpreopt-tools
      if provided, the location of framework's dexpreopt_tools.zip.

  --vendor-dexpreopt-config
      If provided, the location of vendor's dexpreopt_config.zip.
"""

import logging
import os
import re
import shutil
import subprocess
import sys
import zipfile

import add_img_to_target_files
import build_image
import build_super_image
import common
import img_from_target_files
import merge_compatibility_checks
import merge_dexopt
import merge_meta
import merge_utils
import ota_from_target_files

from common import ExternalError

logger = logging.getLogger(__name__)

OPTIONS = common.OPTIONS
# Always turn on verbose logging.
OPTIONS.verbose = True
OPTIONS.framework_target_files = None
OPTIONS.framework_item_list = None
OPTIONS.framework_misc_info_keys = None
OPTIONS.vendor_target_files = None
OPTIONS.vendor_item_list = None
OPTIONS.output_target_files = None
OPTIONS.output_dir = None
OPTIONS.output_item_list = None
OPTIONS.output_ota = None
OPTIONS.output_img = None
OPTIONS.output_super_empty = None
OPTIONS.rebuild_recovery = False
# TODO(b/150582573): Remove this option.
OPTIONS.allow_duplicate_apkapex_keys = False
OPTIONS.vendor_otatools = None
OPTIONS.rebuild_sepolicy = False
OPTIONS.keep_tmp = False
OPTIONS.framework_dexpreopt_config = None
OPTIONS.framework_dexpreopt_tools = None
OPTIONS.vendor_dexpreopt_config = None

# In an item list (framework or vendor), we may see entries that select whole
# partitions. Such an entry might look like this 'SYSTEM/*' (e.g., for the
# system partition). The following regex matches this and extracts the
# partition name.

PARTITION_ITEM_PATTERN = re.compile(r'^([A-Z_]+)/\*$')

FRAMEWORK_PARTITIONS = ('system', 'product', 'system_ext', 'system_other',
                        'root', 'system_dlkm')
VENDOR_PARTITIONS = ('vendor', 'odm', 'oem', 'boot', 'vendor_boot',
                     'prebuilt_images', 'radio', 'data', 'vendor_dlkm',
                     'odm_dlkm')


def infer_item_list(input_zip, framework):
  item_list = []

  with zipfile.ZipFile(input_zip, allowZip64=True) as input_zipfile:
    input_namelist = input_zipfile.namelist()

  for partition in (FRAMEWORK_PARTITIONS if framework else VENDOR_PARTITIONS):
    for namelist in input_namelist:
      if namelist.startswith('%s/' % partition.upper()):
        fs_config_prefix = '' if partition == 'system' else '%s_' % partition
        item_list.extend([
            '%s/*' % partition.upper(),
            'IMAGES/%s.img' % partition,
            'IMAGES/%s.map' % partition,
            'META/%sfilesystem_config.txt' % fs_config_prefix,
        ])
        break

  if framework:
    item_list.extend([
        'META/liblz4.so',
        'META/postinstall_config.txt',
        'META/update_engine_config.txt',
        'META/zucchini_config.txt',
    ])
  else:
    item_list.extend([
        'META/kernel_configs.txt',
        'META/kernel_version.txt',
        'META/otakeys.txt',
        'META/releasetools.py',
        'OTA/android-info.txt',
    ])

  return sorted(item_list)


def infer_framework_misc_info_keys():
  keys = [
      'ab_update',
      'avb_vbmeta_system',
      'avb_vbmeta_system_algorithm',
      'avb_vbmeta_system_key_path',
      'avb_vbmeta_system_rollback_index_location',
      'default_system_dev_certificate',
  ]

  for partition in FRAMEWORK_PARTITIONS:
    fs_type_prefix = '' if partition == 'system' else '%s_' % partition
    keys.extend([
        'avb_%s_hashtree_enable' % partition,
        'avb_%s_add_hashtree_footer_args' % partition,
        '%s_disable_sparse' % partition,
        'building_%s_image' % partition,
        '%sfs_type' % fs_type_prefix,
    ])

  return sorted(keys)


def validate_config_lists():
  """Performs validations on the merge config lists.

  Returns:
    False if a validation fails, otherwise true.
  """
  has_error = False

  # Check that partition contents only come from one input per partition.
  for partition in set(FRAMEWORK_PARTITIONS).union(VENDOR_PARTITIONS):
    image_path = 'IMAGES/{}.img'.format(partition.lower().replace('/', ''))
    in_framework = (
        any(item.startswith(partition) for item in OPTIONS.framework_item_list)
        or image_path in OPTIONS.framework_item_list)
    in_vendor = (
        any(item.startswith(partition) for item in OPTIONS.vendor_item_list) or
        image_path in OPTIONS.vendor_item_list)
    if in_framework and in_vendor:
      logger.error(
          'Cannot extract items from %s for both the framework and vendor'
          ' builds. Please ensure only one merge config item list'
          ' includes %s.', partition, partition)
      has_error = True

  if ('dynamic_partition_list' in OPTIONS.framework_misc_info_keys) or (
      'super_partition_groups' in OPTIONS.framework_misc_info_keys):
    logger.error('Dynamic partition misc info keys should come from '
                 'the vendor instance of META/misc_info.txt.')
    has_error = True

  return not has_error


def item_list_to_partition_set(item_list):
  """Converts a target files item list to a partition set.

  The item list contains items that might look like 'SYSTEM/*' or 'VENDOR/*' or
  'OTA/android-info.txt'. Items that end in '/*' are assumed to match entire
  directories where 'SYSTEM' or 'VENDOR' is a directory name that identifies the
  contents of a partition of the same name. Other items in the list, such as the
  'OTA' example contain metadata. This function iterates such a list, returning
  a set that contains the partition entries.

  Args:
    item_list: A list of items in a target files package.

  Returns:
    A set of partitions extracted from the list of items.
  """

  partition_set = set()

  for item in item_list:
    match = PARTITION_ITEM_PATTERN.search(item.strip())
    partition_tag = match.group(1).lower() if match else None

    if partition_tag:
      partition_set.add(partition_tag)

  return partition_set


def create_merged_package(temp_dir):
  """Merges two target files packages into one target files structure.

  Returns:
    Path to merged package under temp directory.
  """
  # Extract "as is" items from the input framework and vendor partial target
  # files packages directly into the output temporary directory, since these items
  # do not need special case processing.

  output_target_files_temp_dir = os.path.join(temp_dir, 'output')
  merge_utils.ExtractItems(
      input_zip=OPTIONS.framework_target_files,
      output_dir=output_target_files_temp_dir,
      extract_item_list=OPTIONS.framework_item_list)
  merge_utils.ExtractItems(
      input_zip=OPTIONS.vendor_target_files,
      output_dir=output_target_files_temp_dir,
      extract_item_list=OPTIONS.vendor_item_list)

  # Perform special case processing on META/* items.
  # After this function completes successfully, all the files we need to create
  # the output target files package are in place.
  merge_meta.MergeMetaFiles(
      temp_dir=temp_dir, merged_dir=output_target_files_temp_dir)

  merge_dexopt.MergeDexopt(
      temp_dir=temp_dir, output_target_files_dir=output_target_files_temp_dir)

  return output_target_files_temp_dir


def generate_missing_images(target_files_dir):
  """Generate any missing images from target files."""

  # Regenerate IMAGES in the target directory.

  add_img_args = [
      '--verbose',
      '--add_missing',
  ]
  if OPTIONS.rebuild_recovery:
    add_img_args.append('--rebuild_recovery')
  add_img_args.append(target_files_dir)

  add_img_to_target_files.main(add_img_args)


def rebuild_image_with_sepolicy(target_files_dir):
  """Rebuilds odm.img or vendor.img to include merged sepolicy files.

  If odm is present then odm is preferred -- otherwise vendor is used.
  """
  partition = 'vendor'
  if os.path.exists(os.path.join(target_files_dir, 'ODM')) or os.path.exists(
      os.path.join(target_files_dir, 'IMAGES/odm.img')):
    partition = 'odm'
  partition_img = '{}.img'.format(partition)
  partition_map = '{}.map'.format(partition)

  logger.info('Recompiling %s using the merged sepolicy files.', partition_img)

  # Copy the combined SEPolicy file and framework hashes to the image that is
  # being rebuilt.
  def copy_selinux_file(input_path, output_filename):
    input_filename = os.path.join(target_files_dir, input_path)
    if not os.path.exists(input_filename):
      input_filename = input_filename.replace('SYSTEM_EXT/', 'SYSTEM/system_ext/') \
          .replace('PRODUCT/', 'SYSTEM/product/')
      if not os.path.exists(input_filename):
        logger.info('Skipping copy_selinux_file for %s', input_filename)
        return
    shutil.copy(
        input_filename,
        os.path.join(target_files_dir, partition.upper(), 'etc/selinux',
                     output_filename))

  copy_selinux_file('META/combined_sepolicy', 'precompiled_sepolicy')
  copy_selinux_file('SYSTEM/etc/selinux/plat_sepolicy_and_mapping.sha256',
                    'precompiled_sepolicy.plat_sepolicy_and_mapping.sha256')
  copy_selinux_file(
      'SYSTEM_EXT/etc/selinux/system_ext_sepolicy_and_mapping.sha256',
      'precompiled_sepolicy.system_ext_sepolicy_and_mapping.sha256')
  copy_selinux_file('PRODUCT/etc/selinux/product_sepolicy_and_mapping.sha256',
                    'precompiled_sepolicy.product_sepolicy_and_mapping.sha256')

  if not OPTIONS.vendor_otatools:
    # Remove the partition from the merged target-files archive. It will be
    # rebuilt later automatically by generate_missing_images().
    os.remove(os.path.join(target_files_dir, 'IMAGES', partition_img))
    return

  # TODO(b/192253131): Remove the need for vendor_otatools by fixing
  # backwards-compatibility issues when compiling images across releases.
  if not OPTIONS.vendor_target_files:
    raise ValueError(
        'Expected vendor_target_files if vendor_otatools is not None.')
  logger.info(
      '%s recompilation will be performed using the vendor otatools.zip',
      partition_img)

  # Unzip the vendor build's otatools.zip and target-files archive.
  vendor_otatools_dir = common.MakeTempDir(
      prefix='merge_target_files_vendor_otatools_')
  vendor_target_files_dir = common.MakeTempDir(
      prefix='merge_target_files_vendor_target_files_')
  common.UnzipToDir(OPTIONS.vendor_otatools, vendor_otatools_dir)
  common.UnzipToDir(OPTIONS.vendor_target_files, vendor_target_files_dir)

  # Copy the partition contents from the merged target-files archive to the
  # vendor target-files archive.
  shutil.rmtree(os.path.join(vendor_target_files_dir, partition.upper()))
  shutil.copytree(
      os.path.join(target_files_dir, partition.upper()),
      os.path.join(vendor_target_files_dir, partition.upper()),
      symlinks=True)

  # Delete then rebuild the partition.
  os.remove(os.path.join(vendor_target_files_dir, 'IMAGES', partition_img))
  rebuild_partition_command = [
      os.path.join(vendor_otatools_dir, 'bin', 'add_img_to_target_files'),
      '--verbose',
      '--add_missing',
  ]
  if OPTIONS.rebuild_recovery:
    rebuild_partition_command.append('--rebuild_recovery')
  rebuild_partition_command.append(vendor_target_files_dir)
  logger.info('Recompiling %s: %s', partition_img,
              ' '.join(rebuild_partition_command))
  common.RunAndCheckOutput(rebuild_partition_command, verbose=True)

  # Move the newly-created image to the merged target files dir.
  if not os.path.exists(os.path.join(target_files_dir, 'IMAGES')):
    os.makedirs(os.path.join(target_files_dir, 'IMAGES'))
  shutil.move(
      os.path.join(vendor_target_files_dir, 'IMAGES', partition_img),
      os.path.join(target_files_dir, 'IMAGES', partition_img))
  shutil.move(
      os.path.join(vendor_target_files_dir, 'IMAGES', partition_map),
      os.path.join(target_files_dir, 'IMAGES', partition_map))

  def copy_recovery_file(filename):
    for subdir in ('VENDOR', 'SYSTEM/vendor'):
      source = os.path.join(vendor_target_files_dir, subdir, filename)
      if os.path.exists(source):
        dest = os.path.join(target_files_dir, subdir, filename)
        shutil.copy(source, dest)
        return
    logger.info('Skipping copy_recovery_file for %s, file not found', filename)

  if OPTIONS.rebuild_recovery:
    copy_recovery_file('etc/recovery.img')
    copy_recovery_file('bin/install-recovery.sh')
    copy_recovery_file('recovery-from-boot.p')


def generate_super_empty_image(target_dir, output_super_empty):
  """Generates super_empty image from target package.

  Args:
    target_dir: Path to the target file package which contains misc_info.txt for
      detailed information for super image.
    output_super_empty: If provided, copies a super_empty.img file from the
      target files package to this path.
  """
  # Create super_empty.img using the merged misc_info.txt.

  misc_info_txt = os.path.join(target_dir, 'META', 'misc_info.txt')

  use_dynamic_partitions = common.LoadDictionaryFromFile(misc_info_txt).get(
      'use_dynamic_partitions')

  if use_dynamic_partitions != 'true' and output_super_empty:
    raise ValueError(
        'Building super_empty.img requires use_dynamic_partitions=true.')
  elif use_dynamic_partitions == 'true':
    super_empty_img = os.path.join(target_dir, 'IMAGES', 'super_empty.img')
    build_super_image_args = [
        misc_info_txt,
        super_empty_img,
    ]
    build_super_image.main(build_super_image_args)

    # Copy super_empty.img to the user-provided output_super_empty location.
    if output_super_empty:
      shutil.copyfile(super_empty_img, output_super_empty)


def create_target_files_archive(output_zip, source_dir, temp_dir):
  """Creates a target_files zip archive from the input source dir.

  Args:
    output_zip: The name of the zip archive target files package.
    source_dir: The target directory contains package to be archived.
    temp_dir: Path to temporary directory for any intermediate files.
  """
  output_target_files_list = os.path.join(temp_dir, 'output.list')
  output_target_files_meta_dir = os.path.join(source_dir, 'META')

  def files_from_path(target_path, extra_args=None):
    """Gets files under the given path and return a sorted list."""
    find_command = ['find', target_path] + (extra_args or [])
    find_process = common.Run(
        find_command, stdout=subprocess.PIPE, verbose=False)
    return common.RunAndCheckOutput(['sort'],
                                    stdin=find_process.stdout,
                                    verbose=False)

  # META content appears first in the zip. This is done by the
  # standard build system for optimized extraction of those files,
  # so we do the same step for merged target_files.zips here too.
  meta_content = files_from_path(output_target_files_meta_dir)
  other_content = files_from_path(
      source_dir,
      ['-path', output_target_files_meta_dir, '-prune', '-o', '-print'])

  with open(output_target_files_list, 'w') as f:
    f.write(meta_content)
    f.write(other_content)

  command = [
      'soong_zip',
      '-d',
      '-o',
      os.path.abspath(output_zip),
      '-C',
      source_dir,
      '-r',
      output_target_files_list,
  ]

  logger.info('creating %s', output_zip)
  common.RunAndCheckOutput(command, verbose=True)
  logger.info('finished creating %s', output_zip)


def merge_target_files(temp_dir):
  """Merges two target files packages together.

  This function uses framework and vendor target files packages as input,
  performs various file extractions, special case processing, and finally
  creates a merged zip archive as output.

  Args:
    temp_dir: The name of a directory we use when we extract items from the
      input target files packages, and also a scratch directory that we use for
      temporary files.
  """

  logger.info('starting: merge framework %s and vendor %s into output %s',
              OPTIONS.framework_target_files, OPTIONS.vendor_target_files,
              OPTIONS.output_target_files)

  output_target_files_temp_dir = create_merged_package(temp_dir)

  partition_map = common.PartitionMapFromTargetFiles(
      output_target_files_temp_dir)

  compatibility_errors = merge_compatibility_checks.CheckCompatibility(
      target_files_dir=output_target_files_temp_dir,
      partition_map=partition_map)
  if compatibility_errors:
    for error in compatibility_errors:
      logger.error(error)
    raise ExternalError(
        'Found incompatibilities in the merged target files package.')

  # Include the compiled policy in an image if requested.
  if OPTIONS.rebuild_sepolicy:
    rebuild_image_with_sepolicy(output_target_files_temp_dir)

  generate_missing_images(output_target_files_temp_dir)

  generate_super_empty_image(output_target_files_temp_dir,
                             OPTIONS.output_super_empty)

  # Finally, create the output target files zip archive and/or copy the
  # output items to the output target files directory.

  if OPTIONS.output_dir:
    merge_utils.CopyItems(output_target_files_temp_dir, OPTIONS.output_dir,
                          OPTIONS.output_item_list)

  if not OPTIONS.output_target_files:
    return

  # Create the merged META/care_map.pb if the device uses A/B updates.
  if OPTIONS.merged_misc_info['ab_update'] == 'true':
    merge_meta.GenerateCareMap(partition_map.keys(),
                               output_target_files_temp_dir)

  create_target_files_archive(OPTIONS.output_target_files,
                              output_target_files_temp_dir, temp_dir)

  # Create the IMG package from the merged target files package.
  if OPTIONS.output_img:
    img_from_target_files.main(
        [OPTIONS.output_target_files, OPTIONS.output_img])

  # Create the OTA package from the merged target files package.

  if OPTIONS.output_ota:
    ota_from_target_files.main(
        [OPTIONS.output_target_files, OPTIONS.output_ota])


def main():
  """The main function.

  Process command line arguments, then call merge_target_files to
  perform the heavy lifting.
  """

  common.InitLogging()

  def option_handler(o, a):
    if o == '--system-target-files':
      logger.warning(
          '--system-target-files has been renamed to --framework-target-files')
      OPTIONS.framework_target_files = a
    elif o == '--framework-target-files':
      OPTIONS.framework_target_files = a
    elif o == '--system-item-list':
      logger.warning(
          '--system-item-list has been renamed to --framework-item-list')
      OPTIONS.framework_item_list = a
    elif o == '--framework-item-list':
      OPTIONS.framework_item_list = a
    elif o == '--system-misc-info-keys':
      logger.warning('--system-misc-info-keys has been renamed to '
                     '--framework-misc-info-keys')
      OPTIONS.framework_misc_info_keys = a
    elif o == '--framework-misc-info-keys':
      OPTIONS.framework_misc_info_keys = a
    elif o == '--other-target-files':
      logger.warning(
          '--other-target-files has been renamed to --vendor-target-files')
      OPTIONS.vendor_target_files = a
    elif o == '--vendor-target-files':
      OPTIONS.vendor_target_files = a
    elif o == '--other-item-list':
      logger.warning('--other-item-list has been renamed to --vendor-item-list')
      OPTIONS.vendor_item_list = a
    elif o == '--vendor-item-list':
      OPTIONS.vendor_item_list = a
    elif o == '--output-target-files':
      OPTIONS.output_target_files = a
    elif o == '--output-dir':
      OPTIONS.output_dir = a
    elif o == '--output-item-list':
      OPTIONS.output_item_list = a
    elif o == '--output-ota':
      OPTIONS.output_ota = a
    elif o == '--output-img':
      OPTIONS.output_img = a
    elif o == '--output-super-empty':
      OPTIONS.output_super_empty = a
    elif o == '--rebuild_recovery' or o == '--rebuild-recovery':
      OPTIONS.rebuild_recovery = True
    elif o == '--allow-duplicate-apkapex-keys':
      OPTIONS.allow_duplicate_apkapex_keys = True
    elif o == '--vendor-otatools':
      OPTIONS.vendor_otatools = a
    elif o == '--rebuild-sepolicy':
      OPTIONS.rebuild_sepolicy = True
    elif o == '--keep-tmp':
      OPTIONS.keep_tmp = True
    elif o == '--framework-dexpreopt-config':
      OPTIONS.framework_dexpreopt_config = a
    elif o == '--framework-dexpreopt-tools':
      OPTIONS.framework_dexpreopt_tools = a
    elif o == '--vendor-dexpreopt-config':
      OPTIONS.vendor_dexpreopt_config = a
    else:
      return False
    return True

  args = common.ParseOptions(
      sys.argv[1:],
      __doc__,
      extra_long_opts=[
          'system-target-files=',
          'framework-target-files=',
          'system-item-list=',
          'framework-item-list=',
          'system-misc-info-keys=',
          'framework-misc-info-keys=',
          'other-target-files=',
          'vendor-target-files=',
          'other-item-list=',
          'vendor-item-list=',
          'output-target-files=',
          'output-dir=',
          'output-item-list=',
          'output-ota=',
          'output-img=',
          'output-super-empty=',
          'framework-dexpreopt-config=',
          'framework-dexpreopt-tools=',
          'vendor-dexpreopt-config=',
          'rebuild_recovery',
          'rebuild-recovery',
          'allow-duplicate-apkapex-keys',
          'vendor-otatools=',
          'rebuild-sepolicy',
          'keep-tmp',
      ],
      extra_option_handler=option_handler)

  # pylint: disable=too-many-boolean-expressions
  if (args or OPTIONS.framework_target_files is None or
      OPTIONS.vendor_target_files is None or
      (OPTIONS.output_target_files is None and OPTIONS.output_dir is None) or
      (OPTIONS.output_dir is not None and OPTIONS.output_item_list is None) or
      (OPTIONS.rebuild_recovery and not OPTIONS.rebuild_sepolicy)):
    common.Usage(__doc__)
    sys.exit(1)

  if OPTIONS.framework_item_list:
    OPTIONS.framework_item_list = common.LoadListFromFile(
        OPTIONS.framework_item_list)
  else:
    OPTIONS.framework_item_list = infer_item_list(
        input_zip=OPTIONS.framework_target_files, framework=True)
  OPTIONS.framework_partition_set = item_list_to_partition_set(
      OPTIONS.framework_item_list)

  if OPTIONS.framework_misc_info_keys:
    OPTIONS.framework_misc_info_keys = common.LoadListFromFile(
        OPTIONS.framework_misc_info_keys)
  else:
    OPTIONS.framework_misc_info_keys = infer_framework_misc_info_keys()

  if OPTIONS.vendor_item_list:
    OPTIONS.vendor_item_list = common.LoadListFromFile(OPTIONS.vendor_item_list)
  else:
    OPTIONS.vendor_item_list = infer_item_list(
        input_zip=OPTIONS.vendor_target_files, framework=False)

  OPTIONS.vendor_partition_set = item_list_to_partition_set(
      OPTIONS.vendor_item_list)

  if OPTIONS.output_item_list:
    OPTIONS.output_item_list = common.LoadListFromFile(OPTIONS.output_item_list)
  else:
    OPTIONS.output_item_list = None

  if not validate_config_lists():
    sys.exit(1)

  temp_dir = common.MakeTempDir(prefix='merge_target_files_')
  try:
    #item_list = infer_item_list(OPTIONS.framework_target_files, True)
    #print('\n'.join(item_list))
    #item_list = infer_item_list(OPTIONS.vendor_target_files, False)
    #print('\n'.join(item_list))
    merge_target_files(temp_dir)
  finally:
    if OPTIONS.keep_tmp:
      logger.info('Keeping temp_dir %s', temp_dir)
    else:
      common.Cleanup()


if __name__ == '__main__':
  main()
