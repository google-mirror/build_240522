#!/usr/bin/env python
#
# Copyright (C) 2008 The Android Open Source Project
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
Given a target-files zipfile, produces an image zipfile suitable for
use with 'fastboot update'.

Usage:  img_from_target_files [flags] input_target_files output_image_zip

  -z  (--bootable_zip)
      Include only the bootable images (eg 'boot' and 'recovery') in
      the output.

"""

from __future__ import print_function

import logging
import os
import shutil
import sys
import zipfile

import common
from build_super_image import BuildSuperImage

if sys.hexversion < 0x02070000:
  print("Python 2.7 or newer is required.", file=sys.stderr)
  sys.exit(1)

logger = logging.getLogger(__name__)

OPTIONS = common.OPTIONS


def CopyInfo(output_zip):
  """Copy the android-info.txt file from the input to the output."""
  common.ZipWrite(
      output_zip, os.path.join(OPTIONS.input_tmp, "OTA", "android-info.txt"),
      "android-info.txt")


def main(argv):
  # This allows modifying the value from inner function.
  bootable_only_array = [False]

  def option_handler(o, _):
    if o in ("-z", "--bootable_zip"):
      bootable_only_array[0] = True
    else:
      return False
    return True

  args = common.ParseOptions(argv, __doc__,
                             extra_opts="z",
                             extra_long_opts=["bootable_zip"],
                             extra_option_handler=option_handler)

  bootable_only = bootable_only_array[0]

  if len(args) != 2:
    common.Usage(__doc__)
    sys.exit(1)

  common.InitLogging()

  OPTIONS.input_zip = args[0]
  OPTIONS.info_dict = common.LoadInfoDict(OPTIONS.input_zip)

  put_super = OPTIONS.info_dict.get("super_image_in_update_package") == "true"
  dynamic_partition_list = OPTIONS.info_dct.get("dynamic_partition_list",
                                                "").strip().split()
  super_device_list = OPTIONS.info_dict.get("super_block_devices",
                                            "").strip().split()
  retrofit_dap = OPTIONS.info_dict.get("dynamic_partition_retrofit") == "true"
  should_build_super = OPTIONS.info_dict.get("build_super_partition") == "true"

  unzip_pattern = ["IMAGES/*", "OTA/android-info.txt"]
  if put_super and should_build_super:
    if retrofit_dap:
      unzip_pattern.append("OTA/super_*.img")
    else:
      unzip_pattern.append("META/*")

  OPTIONS.input_tmp = common.UnzipTemp(OPTIONS.input_zip, unzip_pattern)
  output_zip = zipfile.ZipFile(args[1], "w", compression=zipfile.ZIP_DEFLATED)
  CopyInfo(output_zip)

  try:
    dynamic_images = [p + ".img" for p in dynamic_partition_list]

    images_path = os.path.join(OPTIONS.input_tmp, "IMAGES")
    # A target-files zip must contain the images since Lollipop.
    assert os.path.exists(images_path)
    for image in sorted(os.listdir(images_path)):
      if bootable_only and image not in ("boot.img", "recovery.img"):
        continue
      if not image.endswith(".img"):
        continue
      if image == "recovery-two-step.img":
        continue
      if put_super:
        if image == "super_empty.img":
          continue
        if image in dynamic_images:
          continue
      common.ZipWrite(output_zip, os.path.join(images_path, image), image)

    if put_super and should_build_super:
      if retrofit_dap:
        # retrofit devices already have super split images under OTA/
        images_path = os.path.join(OPTIONS.input_tmp, "OTA")
        super_split_images = ["super_" + p + ".img" for p in super_device_list]
        for image in super_split_images:
          common.ZipWrite(output_zip, os.path.join(images_path, image), image)
      else:
        # super image for non-retrofit devices aren't in target files package,
        # so build it.
        super_file = common.MakeTempFile("super", ".img")
        BuildSuperImage(OPTIONS.input_tmp, super_file)
        common.ZipWrite(output_zip, super_file, "super.img")

  finally:
    logger.info("cleaning up...")
    common.ZipClose(output_zip)
    shutil.rmtree(OPTIONS.input_tmp)

  logger.info("done.")


if __name__ == '__main__':
  try:
    common.CloseInheritedPipes()
    main(sys.argv[1:])
  except common.ExternalError as e:
    logger.exception("\n   ERROR:\n")
    sys.exit(1)
