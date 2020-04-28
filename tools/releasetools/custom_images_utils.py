#!/usr/bin/env python
#
# Copyright (C) 2020 The Android Open Source Project
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

import logging
import os.path
import re
import shutil
import sys
import zipfile

import common

logger = logging.getLogger(__name__)

OPTIONS = common.OPTIONS

class CustomInfoError(Exception):
  """An Exception raised during Custom Information command."""

  def __init__(self, message):
    Exception.__init__(self, message)


class CustomSigningError(Exception):
  """An Exception raised during Custom image signing."""

  def __init__(self, message):
    Exception.__init__(self, message)


def SignCustomImage(raw_image, key_path, algorithm, extra_args=None):
  """Signs a given raw_image with the key and algorithm."""
  temp_raw_image = common.MakeTempFile(prefix='custom-', suffix='.img')
  shutil.copyfile(raw_image, temp_raw_image)

  image_info = ParseCustomImageInfo(temp_raw_image)
  if image_info is None:
    return None;

  # Add the new footer. Old footer, if any, will be replaced by avbtool.
  cmd = ['avbtool', 'add_hashtree_footer',
         '--do_not_generate_fec',
         '--hash_algorithm', image_info['hash_algorithm'],
         '--partition_name', image_info['partition_name'],
         '--partition_size', image_info['partition_size'],
         '--salt', image_info['Salt'],
         '--algorithm', algorithm,
         '--key', key_path,
         '--image', temp_raw_image]
  if extra_args is not None:
    cmd.append(extra_args)

  try:
    common.RunAndCheckOutput(cmd)
  except common.ExternalError as e:
    raise CustomSigningError(
        'Failed to sign Custom Image {} with {}:\n{}'.format(
            raw_image, key_path, e))

  # Verify the signed Custom image with specified public key.
  logger.info('Verifying %s', temp_raw_image)
  VerifyCustomImage(temp_raw_image, image_info['partition_name'], key_path)

  return temp_raw_image


def VerifyCustomImage(custom_image, partition_name, key_path):
  """Verifies the Custom image signature with the given key."""

  # To verify the image, image name must be same as partition name
  verify_dir = common.MakeTempDir(prefix='custom-verify-')
  verifying_custom_image = os.path.join(verify_dir, partition_name + '.img')
  shutil.copyfile(custom_image, verifying_custom_image)

  cmd = ['avbtool', 'verify_image', '--image', verifying_custom_image,
         '--key', key_path]
  try:
    common.RunAndCheckOutput(cmd)
  except common.ExternalError as e:
    raise CustomSigningError(
        'Failed to validate Custom signing for {} with {}:\n{}'.format(
            custom_image, key_path, e))


def ParseCustomImageInfo(custom_image):
  """Parses the Custom image info."""
  if not os.path.exists(custom_image):
    raise CustomInfoError('Failed to find image: {}'.format(custom_image))

  cmd = ['avbtool', 'info_image', '--image', custom_image]
  try:
    output = common.RunAndCheckOutput(cmd)
  except common.ExternalError as e:
    raise CustomInfoError(
        'Failed to get Custom image info for {}:\n{}'.format(
            custom_image, e))

  # Extract the Image size / Hash Algorithm/ Partition Name / Salt info from
  # Custom Image (i.e. an image signed with avbtool). For example,
  # Image size:               5242880 bytes
  #       Hash Algorithm:        sha1
  #       Partition Name:        oem
  #       Salt:                  9999cad66543df71c361c5203cd6e025f2529bce
  IMAGE_INFO_PATTERN = (
      r'^\s*(?P<key>Image size|Hash Algorithm|Partition Name|Salt)' \
          '\:\s*(?P<value>.*?)$')
  image_info_matcher = re.compile(IMAGE_INFO_PATTERN)

  image_info = {}
  for line in output.split('\n'):
    line_info = image_info_matcher.match(line)
    if not line_info:
      continue

    key, value = line_info.group('key'), line_info.group('value')

    if key == 'Image size':
      image_info['partition_size'] = value.split(" ")[0]
    elif key == 'Hash Algorithm':
      image_info['hash_algorithm'] = value
    elif key == 'Partition Name':
      image_info['partition_name'] = value
    else:
      image_info[key] = value

  # Sanity check.
  if len(image_info) != 4:
    return None
  for key in ('partition_size', 'hash_algorithm', 'partition_name', 'Salt'):
    if key not in image_info:
      raise CustomInfoError(
          'Failed to find {} prop in {}'.format(key, custom_image))

  return image_info
