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

"""
Signs a standalone Custom image.

Usage:  sign_custom_image [flags] input_custom_image output_custom_image

  --key <key>
      Mandatory flag that specifies the signing key.

  --algorithm <algorithm>
      Mandatory flag that specifies the signing algorithm.
"""

import logging
import shutil
import sys

import common
import custom_utils

logger = logging.getLogger(__name__)


def main(argv):

  options = {}

  def option_handler(o, a):
    if o == '--key':
      options['key'] = a
    elif o == '--algorithm':
      options['algorithm'] = a
    else:
      return False
    return True

  args = common.ParseOptions(
      argv, __doc__,
      extra_opts='',
      extra_long_opts=[
          'key=',
          'algorithm=',
      ],
      extra_option_handler=option_handler)

  if (len(args) != 2 or 'key' not in options or 'algorithm' not in options):
    common.Usage(__doc__)
    sys.exit(1)

  common.InitLogging()

  input_image = args[0]
  output_image = args[1]
  signed_image = common.MakeTempFile(prefix='custom-', suffix='.img')

  signed_image = custom_utils.SignCustomImage(input_image,
                                              options['key'],
                                              options['algorithm'])

  if signed_image is None:
    print('%s is not an AVB image, Skip ...' % (input_image,))
    sys.exit(1)

  shutil.copyfile(signed_image, output_image)
  logger.info("done.")


if __name__ == '__main__':
  try:
    main(sys.argv[1:])
  except common.ExternalError:
    logger.exception("\n   ERROR:\n")
    sys.exit(1)
  finally:
    common.Cleanup()
