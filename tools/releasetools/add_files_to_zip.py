#!/usr/bin/env python
#
# Copyright (C) 2015 The Android Open Source Project
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
Add a list of files to a zip file.

Usage:  add_files_to_zip files zip_file.
"""

import sys

if sys.hexversion < 0x02070000:
  print >> sys.stderr, "Python 2.7 or newer is required."
  sys.exit(1)

import os
import zipfile

import common

OPTIONS = common.OPTIONS

OPTIONS.add_missing = False
OPTIONS.rebuild_recovery = False
OPTIONS.replace_verity_public_key = False
OPTIONS.replace_verity_private_key = False
OPTIONS.verity_signer_path = None


def AddFiles(zip_file, source_files):
  """Add the given source file to output_zip.

  Args:
    zip_file: Path to the zip file.
    source_file: Path to the file to be added to the zip file.
  """
  output_zip = zipfile.ZipFile(zip_file, "a",
                               compression=zipfile.ZIP_DEFLATED)
  for source_file in source_files:
    print "Writing %s to %s..." % (source_file, os.path.basename(zip_file))
    common.ZipWrite(output_zip, source_file, os.path.basename(source_file))


def main(argv):
  args = common.ParseOptions(argv, __doc__, extra_opts="ar")

  if len(args) < 2:
    common.Usage(__doc__)
    sys.exit(1)

  AddFiles(args[0], args[1:])
  print "done."

if __name__ == '__main__':
  try:
    common.CloseInheritedPipes()
    main(sys.argv[1:])
  except common.ExternalError as e:
    print
    print "   ERROR: %s" % (e,)
    print
    sys.exit(1)
  finally:
    common.Cleanup()
