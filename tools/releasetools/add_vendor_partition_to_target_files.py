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
Given a target-files zipfile, add all vendor partition files in the given folder
to the zipfile, under IMAGES/vendor subdirectory.

Usage:  add_vendor_partition_to_target_files vendor_partition_folder.
"""

import sys

if sys.hexversion < 0x02070000:
  print >> sys.stderr, "Python 2.7 or newer is required."
  sys.exit(1)

import datetime
import errno
import os
import shutil
import tempfile
import zipfile

import common

OPTIONS = common.OPTIONS

OPTIONS.add_missing = False
OPTIONS.rebuild_recovery = False
OPTIONS.replace_verity_public_key = False
OPTIONS.replace_verity_private_key = False
OPTIONS.verity_signer_path = None

def FindPartitionFiles(vendor_partition_folder):
  """Iterate through bin files in vendor_partition_folder including
  subdirectories.

  Args:
    vendor_partition_folder: Path to the parent folder of vendor partitions.
  Returns:
    An iterator that iterates through all bin files in given path including
    subdirectories.
  """
  if not os.path.exists(vendor_partition_folder):
      yield []
  for root, dirs, files in os.walk(vendor_partition_folder):
    for f in files:
      if f.endswith('.bin'):
        yield os.path.join(root, f)


def AddPartition(output_zip, partition):
  """Add the given partition to output_zip.

  Args:
    output_zip: Path to the zip file.
    partition: Path to the partition file to be added to the zip file.
  """
  patition_path_in_zip = os.path.join("IMAGES", "vendor",
                                      os.path.basename(partition))
  common.ZipWrite(output_zip, partition, patition_path_in_zip)

def main(argv):
  args = common.ParseOptions(argv, __doc__, extra_opts="ar")

  if len(args) != 2:
    common.Usage(__doc__)
    sys.exit(1)

  output_zip = zipfile.ZipFile(args[0], "a",
                               compression=zipfile.ZIP_DEFLATED)
  for partition in FindPartitionFiles(args[1]):
    print "Writing %s to %s..." % (partition, os.path.basename(args[0]))
    AddPartition(output_zip, partition)
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
