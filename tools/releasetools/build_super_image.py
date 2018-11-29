#!/usr/bin/env python
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

"""
Usage: build_super_image dict_file output_dir_or_file

dict_file: a dictionary file containing input arguments to build. Check
    `dump_dynamic_partitions_info' for details.
    In addition, "ab_update" needs to be true for A/B devices.

output_dir_or_file:
    If a single super image is built (for super_empty.img, or super.img for
    launch devices), this argument is the output file.
    If a collection of split images are built (for retrofit devices), this
    argument is the output directory.
"""

from __future__ import print_function

import common
import logging
import sparse_img
import os.path
import shlex
import sys
import zipfile

if sys.hexversion < 0x02070000:
  print("Python 2.7 or newer is required.", file=sys.stderr)
  sys.exit(1)

logger = logging.getLogger(__name__)

OPTIONS = common.OPTIONS
OPTIONS.extracted_input = None


def GetPartitionSizeFromImage(img):
  try:
    simg = sparse_img.SparseImage(img)
    return simg.blocksize * simg.total_blocks
  except ValueError:
    return os.path.getsize(img)


def BuildSuperImage(info_dict, output):

  cmd = [info_dict["lpmake"],
         "--metadata-size", "65536",
         "--super-name", info_dict["super_metadata_device"]]

  ab_update = info_dict.get("ab_update") == "true"
  retrofit = info_dict.get("dynamic_partition_retrofit") == "true"
  block_devices = shlex.split(info_dict.get("super_block_devices", "").strip())
  groups = shlex.split(info_dict.get("super_partition_groups", "").strip())

  if ab_update:
    cmd += ["--metadata-slots", "2"]
  else:
    cmd += ["--metadata-slots", "1"]

  if ab_update and retrofit:
    cmd.append("--auto-slot-suffixing")

  for device in block_devices:
    size = info_dict["super_{}_device_size".format(device)]
    cmd += ["--device", "{}:{}".format(device, size)]

  append_suffix = ab_update and not retrofit
  has_image = False
  for group in groups:
    group_size = info_dict["super_{}_group_size".format(group)]
    if append_suffix:
      cmd += ["--group", "{}_a:{}".format(group, group_size),
              "--group", "{}_b:{}".format(group, group_size)]
    else:
      cmd += ["--group", "{}:{}".format(group, group_size)]

    partition_list = shlex.split(
        info_dict["super_{}_partition_list".format(group)].strip())

    for partition in partition_list:
      image = info_dict.get("{}_image".format(partition))
      image_size = 0
      if image:
        image_size = GetPartitionSizeFromImage(image)
        has_image = True
      if append_suffix:
        cmd += ["--partition",
                "{}_a:readonly:{}:{}_a".format(partition, image_size, group),
                "--partition",
                "{}_b:readonly:0:{}_b".format(partition, group)]
        if image:
          # For A/B devices, super partition always contains sub-partitions in
          # the _a slot, because this image should only be used for
          # bootstrapping / initializing the device. When flashing the image,
          # bootloader fastboot should always mark _a slot as bootable.
          cmd += ["--image", "{}_a={}".format(partition, image)]
      else:
        cmd += ["--partition",
                "{}:readonly:{}:{}".format(partition, image_size, group)]
        if image:
          cmd += ["--image", "{}={}".format(partition, image)]

  if has_image:
    cmd.append("--sparse")

  cmd += ["--output", output]

  common.RunAndCheckOutput(cmd)

  if retrofit and has_image:
    logger.info("Done writing images to directory %s", output)
  else:
    logger.info("Done writing image %s", output)

def main(argv):

  args = common.ParseOptions(argv, __doc__)

  if len(args) != 2:
    common.Usage(__doc__)
    sys.exit(1)

  common.InitLogging()

  def read_all(path):
    with open(path) as f:
      return f.read()

  BuildSuperImage(common.LoadDictionaryFromLines(read_all(args[0]).split("\n")), args[1])

if __name__ == "__main__":
  try:
    common.CloseInheritedPipes()
    main(sys.argv[1:])
  except common.ExternalError:
    logger.exception("\n   ERROR:\n")
    sys.exit(1)
  finally:
    common.Cleanup()
