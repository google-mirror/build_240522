#!/usr/bin/env python

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

"""
Validate a given (signed) target_files.zip.

It performs checks to ensure the integrity of the input zip.
 - It verifies the file consistency between the ones in IMAGES/system.img (read
   via IMAGES/system.map) and the ones under unpacked folder of SYSTEM/. The
   same check also applies to the vendor image if present.
"""

import common
import logging
import os.path
import sparse_img
import sys

def _GetImage(which, tmpdir):
  assert which in ('system', 'vendor')

  path = os.path.join(tmpdir, 'IMAGES', which + '.img')
  mappath = os.path.join(tmpdir, 'IMAGES', which + '.map')

  # Map file must exist (allowed to be empty).
  assert os.path.exists(path) and os.path.exists(mappath)

  clobbered_blocks = '0'
  return sparse_img.SparseImage(path, mappath, clobbered_blocks)

def _CalculateFileSha1(file_name, unpacked_name, round_up=False):

  def RoundUpTo4K(value):
    rounded_up = value + 4095
    return rounded_up - (rounded_up % 4096)

  assert os.path.exists(unpacked_name)
  with open(unpacked_name, 'r') as f:
    file_data = f.read()
  file_size = len(file_data)
  if round_up:
    file_size_rounded_up = RoundUpTo4K(file_size)
    file_data += '\0' * (file_size_rounded_up - file_size)
  return common.File(file_name, file_data).sha1

def _ValidateFileAgainstSha1(input_tmp, file_name, file_path, expected_sha1):
  logging.info('Validating the SHA1 of {}'.format(file_name))
  unpacked_name = os.path.join(input_tmp, file_path)
  assert os.path.exists(unpacked_name)
  actual_sha1 = _CalculateFileSha1(file_name, unpacked_name, False)
  assert actual_sha1 == expected_sha1, 'SHA1 mismatches for {}. ' \
      'actual {}, expected {}'.format(file_name, actual_sha1, expected_sha1)

def ValidateFileConsistency(input_zip, input_tmp):
  """Compare the files from image files and unpacked folders."""

  def CheckAllFiles(which):
    logging.info('Checking %s image.', which)
    image = _GetImage(which, input_tmp)
    prefix = '/' + which
    for entry in image.file_map:
      if not entry.startswith(prefix):
        continue

      # Read the blocks that the file resides. Note that it will contain the
      # bytes past the file length, which is expected to be padded with '\0's.
      ranges = image.file_map[entry]
      blocks_sha1 = image.RangeSha1(ranges)

      # The filename under unpacked directory, such as SYSTEM/bin/sh.
      unpacked_name = os.path.join(
          input_tmp, which.upper(), entry[(len(prefix) + 1):])
      file_sha1 = _CalculateFileSha1(entry, unpacked_name, True)

      assert blocks_sha1 == file_sha1, \
          'file: %s, range: %s, blocks_sha1: %s, file_sha1: %s' % (
              entry, ranges, blocks_sha1, file_sha1)

  logging.info('Validating file consistency.')

  # Verify IMAGES/system.img.
  CheckAllFiles('system')

  # Verify IMAGES/vendor.img if applicable.
  if 'VENDOR/' in input_zip.namelist():
    CheckAllFiles('vendor')

  # Not checking IMAGES/system_other.img since it doesn't have the map file.


def ValidateInstallRecoveryScript(input_tmp):
  script_path = 'SYSTEM/bin/install-recovery.sh'
  if not os.path.exists(os.path.join(input_tmp, script_path)):
    logging.info('{} does not exist in input_tmp'.format(script_path))
    return

  logging.info('Checking {}'.format(script_path))
  with open(os.path.join(input_tmp, script_path), 'r') as script:
    lines = script.read().strip().split('\n')
    assert len(lines) >= 6
    trim_index = lines[2].find('&&')
    assert trim_index != -1
    applypatch_argv = lines[2][:trim_index].strip().split(' ')
    assert len(applypatch_argv) >= 5 and applypatch_argv[0] == 'applypatch'

    if len(applypatch_argv) == 5:
      expected_recovery_sha1 = applypatch_argv[3].strip()
      _ValidateFileAgainstSha1(input_tmp, 'recovery.img',
          'SYSTEM/etc/recovery.img', expected_recovery_sha1)
    else:
      if applypatch_argv[1] == "-b":
        assert len(applypatch_argv) == 8
        boot_info_index = 3
      else:
        assert len(applypatch_argv) == 6
        boot_info_index = 1
      boot_info = applypatch_argv[boot_info_index].strip().split(':')
      assert len(boot_info) == 4
      _ValidateFileAgainstSha1(input_tmp, file_name='boot.img',
          file_path='IMAGES/boot.img', expected_sha1=boot_info[3])

      recovery_sha1_index = boot_info_index + 2
      expected_recovery_sha1 = applypatch_argv[recovery_sha1_index]
      _ValidateFileAgainstSha1(input_tmp, file_name='recovery.img',
          file_path='IMAGES/recovery.img',
          expected_sha1=expected_recovery_sha1)

  logging.info('Done checking {}'.format(script_path))


def main(argv):
  def option_handler():
    return True

  args = common.ParseOptions(
      argv, __doc__, extra_opts="",
      extra_long_opts=[],
      extra_option_handler=option_handler)

  if len(args) != 1:
    common.Usage(__doc__)
    sys.exit(1)

  logging_format = '%(asctime)s - %(filename)s - %(levelname)-8s: %(message)s'
  date_format = '%Y/%m/%d %H:%M:%S'
  logging.basicConfig(level=logging.INFO, format=logging_format,
                      datefmt=date_format)

  logging.info("Unzipping the input target_files.zip: %s", args[0])
  input_tmp, input_zip = common.UnzipTemp(args[0])

  #ValidateFileConsistency(input_zip, input_tmp)

  ValidateInstallRecoveryScript(input_tmp)

  # TODO: Check if the OTA keys have been properly updated (the ones on /system,
  # in recovery image).




  logging.info("Done.")


if __name__ == '__main__':
  import tempfile
  tempfile.tempdir = "/usr/local/google/home/xunchang/1"
  try:
    main(sys.argv[1:])
  finally:
    common.Cleanup()
