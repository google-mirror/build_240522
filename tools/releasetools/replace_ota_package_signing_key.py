#!/usr/bin/env python
#
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
Replaces the signing key of an OTA package.
"""

from __future__ import print_function

import shlex
import sys

import common
from ota_from_target_files import AbOtaPackage


OPTIONS = common.OPTIONS


def main(args):
  def option_handler(o, a):
    if o in ("-k", "--package_key"):
      OPTIONS.package_key = a
    elif o == "--payload_signer":
      OPTIONS.payload_signer = a
    elif o == "--payload_signer_args":
      OPTIONS.payload_signer_args = shlex.split(a)
    else:
      return False
    return True

  args = common.ParseOptions(args, __doc__,
                             extra_opts="k:",
                             extra_long_opts=[
                                 "package_key=",
                                 "payload_signer=",
                                 "payload_signer_args=",
                             ],
                             extra_option_handler=option_handler)

  if len(args) != 2:
    common.Usage(__doc__)
    sys.exit(1)

  # Get signing keys
  OPTIONS.key_passwords = common.GetKeyPasswords([OPTIONS.package_key])

  package = AbOtaPackage(args[1])
  package.PopulateFromPackageFile(args[0])
  package.Finalize()


if __name__ == '__main__':
  try:
    main(sys.argv[1:])
  except AssertionError as err:
    print('\n    ERROR: %s\n' % (err,))
    sys.exit(1)
  finally:
    common.Cleanup()
