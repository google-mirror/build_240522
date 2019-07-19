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
""" TODO"""
from __future__ import print_function

import logging
import os
import sys

import common

logger = logging.getLogger(__name__)

OPTIONS = common.OPTIONS
OPTIONS.product_out_framework = None
OPTIONS.product_out_vendor = None
OPTIONS.framework_misc_info_keys = None


def MergeBuilds(temp_dir, product_out_framework, product_out_vendor,
                framework_misc_info_keys):
  framework_dict = common.LoadDictionaryFromFile(
      os.path.join(product_out_framework, "misc_info.txt"))
  vendor_dict = common.LoadDictionaryFromFile(
      os.path.join(product_out_vendor, "misc_info.txt"))
  merged_dict = common.MergeInfoDicts(
      framework_dict, vendor_dict,
      common.LoadListFromFile(framework_misc_info_keys))
  logger.warning(merged_dict)


def main():
  common.InitLogging()

  def option_handler(o, a):
    if o == "--product_out_framework":
      OPTIONS.product_out_framework = a
    elif o == "--product_out_vendor":
      OPTIONS.product_out_vendor = a
    elif o == "--framework_misc_info_keys":
      OPTIONS.framework_misc_info_keys = a
    else:
      return False
    return True

  args = common.ParseOptions(
      sys.argv[1:],
      __doc__,
      extra_long_opts=[
          "product_out_framework=",
          "product_out_vendor=",
          "framework_misc_info_keys=",
      ],
      extra_option_handler=option_handler)

  if (args or OPTIONS.product_out_framework is None or
      OPTIONS.product_out_vendor is None or
      OPTIONS.framework_misc_info_keys is None):
    common.Usage(__doc__)
    sys.exit(1)

  logger.warning("Hi")
  temp_dir = common.MakeTempDir(prefix="merge_builds_")
  try:
    MergeBuilds(temp_dir, OPTIONS.product_out_framework,
                OPTIONS.product_out_vendor, OPTIONS.framework_misc_info_keys)
  finally:
    common.Cleanup()


if __name__ == "__main__":
  main()
