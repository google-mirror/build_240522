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
"""Merges two non-dist partial builds together.

Given two partial builds, a framework build and a vendor build, merge the builds
together so that the images can be flashed using 'fastboot flashall'.

This script does not require builds to be built with 'make dist'.
This script assumes that images other than super_empty.img do not require
regeneration.

Usage: merge_builds.py [args]

  --framework_images comma_separated_image_list
      Comma-separated list of image names that should come from the framework
      build.

  --product_out_framework product_out_framework_path
      Path to out/target/product/<framework build>.

  --product_out_vendor product_out_vendor_path
      Path to out/target/product/<vendor build>.
"""
from __future__ import print_function

import logging
import os
import sys

import build_super_image
import common

logger = logging.getLogger(__name__)

OPTIONS = common.OPTIONS
OPTIONS.framework_images = ("system",)
OPTIONS.product_out_framework = None
OPTIONS.product_out_vendor = None


def CreateImageSymlinks():
  for image in OPTIONS.framework_images:
    image_path = os.path.join(OPTIONS.product_out_framework, "%s.img" % image)
    symlink_path = os.path.join(OPTIONS.product_out_vendor, "%s.img" % image)
    os.symlink(image_path, symlink_path)


def BuildSuperEmpty():
  framework_dict = common.LoadDictionaryFromFile(
      os.path.join(OPTIONS.product_out_framework, "misc_info.txt"))
  vendor_dict = common.LoadDictionaryFromFile(
      os.path.join(OPTIONS.product_out_vendor, "misc_info.txt"))
  if (framework_dict.get("use_dynamic_partitions") == "true") and (
      vendor_dict.get("use_dynamic_partitions") == "true"):
    merged_dict = common.MergeDynamicPartitionInfoDicts(
        framework_dict=framework_dict,
        vendor_dict=vendor_dict,
        size_prefix="super_",
        size_suffix="_group_size",
        list_prefix="super_",
        list_suffix="_partition_list")
    output_super_empty_path = os.path.join(OPTIONS.product_out_vendor,
                                           "super_empty.img")
    build_super_image.BuildSuperImage(merged_misc_info_dict,
                                      output_super_empty_path)


def MergeBuilds():
  CreateImageSymlinks()
  BuildSuperEmpty()


def main():
  common.InitLogging()

  def option_handler(o, a):
    if o == "--framework_images":
      OPTIONS.framework_images = [i.strip() for i in a.split(",")]
    elif o == "--product_out_framework":
      OPTIONS.product_out_framework = a
    elif o == "--product_out_vendor":
      OPTIONS.product_out_vendor = a
    else:
      return False
    return True

  args = common.ParseOptions(
      sys.argv[1:],
      __doc__,
      extra_long_opts=[
          "framework_images=",
          "product_out_framework=",
          "product_out_vendor=",
      ],
      extra_option_handler=option_handler)

  if (args or OPTIONS.product_out_framework is None or
      OPTIONS.product_out_vendor is None):
    common.Usage(__doc__)
    sys.exit(1)

  MergeBuilds()


if __name__ == "__main__":
  main()
