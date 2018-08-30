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

from __future__ import print_function

import struct
import subprocess

import common
from rangelib import RangeSet


FIXED_SALT = "aee087a5be3b982978c923f566a94613496b417f2af592639bc80d141e34dfe7"


def GetVerityFECSize(partition_size):
  cmd = ["fec", "-s", str(partition_size)]
  p = common.Run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                 verbose=False)
  output, _ = p.communicate()
  if p.returncode != 0:
    return False, 0
  return True, int(output)


def GetVerityTreeSize(partition_size):
  cmd = ["build_verity_tree", "-s", str(partition_size)]
  p = common.Run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                 verbose=False)
  output, _ = p.communicate()
  if p.returncode != 0:
    return False, 0
  return True, int(output)


def GetVerityMetadataSize(partition_size):
  cmd = ["build_verity_metadata.py", "size", str(partition_size)]
  p = common.Run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                 verbose=False)
  output, _ = p.communicate()
  if p.returncode != 0:
    return False, 0
  return True, int(output)


def GetVeritySize(partition_size, fec_supported):
  success, verity_tree_size = GetVerityTreeSize(partition_size)
  if not success:
    return 0
  success, verity_metadata_size = GetVerityMetadataSize(partition_size)
  if not success:
    return 0
  verity_size = verity_tree_size + verity_metadata_size
  if fec_supported:
    success, fec_size = GetVerityFECSize(partition_size + verity_size)
    if not success:
      return 0
    return verity_size + fec_size
  return verity_size


def AdjustPartitionSizeForVerity(partition_size, fec_supported, block_size=4096,
                                 verbose=False):
  """Modifies the provided partition size to account for the verity metadata.

  This information is used to size the created image appropriately.

  Args:
    partition_size: the size of the partition to be verified.

  Returns:
    A tuple of the size of the partition adjusted for verity metadata, and
    the size of verity metadata.
  """
  key = "%d %d" % (partition_size, fec_supported)
  if key in AdjustPartitionSizeForVerity.results:
    return AdjustPartitionSizeForVerity.results[key]

  hi = partition_size
  if hi % block_size != 0:
    hi = (hi // block_size) * block_size

  # verity tree and fec sizes depend on the partition size, which
  # means this estimate is always going to be unnecessarily small
  verity_size = GetVeritySize(hi, fec_supported)
  lo = partition_size - verity_size
  result = lo

  # do a binary search for the optimal size
  while lo < hi:
    i = ((lo + hi) // (2 * block_size)) * block_size
    v = GetVeritySize(i, fec_supported)
    if i + v <= partition_size:
      if result < i:
        result = i
        verity_size = v
      lo = i + block_size
    else:
      hi = i

  if verbose:
    print("Adjusted partition size for verity, partition_size: {},"
          " verity_size: {}".format(result, verity_size))
  AdjustPartitionSizeForVerity.results[key] = (result, verity_size)
  return result, verity_size


AdjustPartitionSizeForVerity.results = {}


def BuildVerityFEC(sparse_image_path, verity_path, verity_fec_path,
                   padding_size):
  cmd = ["fec", "-e", "-p", str(padding_size), sparse_image_path,
         verity_path, verity_fec_path]
  p = common.Run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
  output, _ = p.communicate()
  if p.returncode != 0:
    print("Could not build FEC data! Error: %s" % output)
    return False
  return True


def BuildVerityTree(sparse_image_path, verity_image_path, prop_dict):
  cmd = ["build_verity_tree", "-A", FIXED_SALT, sparse_image_path,
         verity_image_path]
  p = common.Run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
  output, _ = p.communicate()
  if p.returncode != 0:
    print("Could not build verity tree! Error: %s" % output)
    return False
  root, salt = output.split()
  prop_dict["verity_root_hash"] = root
  prop_dict["verity_salt"] = salt
  return True


def BuildVerityMetadata(image_size, verity_metadata_path, root_hash, salt,
                        block_device, signer_path, key, signer_args,
                        verity_disable):
  cmd = ["build_verity_metadata.py", "build", str(image_size),
         verity_metadata_path, root_hash, salt, block_device, signer_path, key]
  if signer_args:
    cmd.append("--signer_args=\"%s\"" % (' '.join(signer_args),))
  if verity_disable:
    cmd.append("--verity_disable")
  p = common.Run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
  output, _ = p.communicate()
  if p.returncode != 0:
    print("Could not build verity metadata! Error: %s" % output)
    return False
  return True


class VerityTreeInfo(object):
  """A class that parses the metadata of hash_tree for a given partition."""

  def __init__(self, partition_name, tgt, info_dict):
    """Initialize VerityTreeInfo with the sparse image and input property.

    Arguments:
      partition_name: name of the target partition to compute the hash_tree
      tgt: sparse image of the target partition
      info_dict: The build-time info dict.
    """

    verity_partition = (
        info_dict.get("verity") == "true" and
        "{}_verity_block_device".format(partition_name) in info_dict)

    if not verity_partition:
      self.enabled = False
      return

    self.enabled = True
    self.tgt = tgt
    self.block_size = tgt.blocksize
    self.partition_size = partition_size = info_dict.get(
        "{}_size".format(partition_name))
    image_size = tgt.total_blocks * tgt.blocksize
    assert image_size == partition_size, "partition size {} doesn't match" \
        " with the calculated image size {}".format(partition_size, image_size)
    self.fec_supported = info_dict.get("verity_fec") == "true"

    self.verity_tree_range = None
    self.file_system_range = None
    self.hash_algorithm = None
    self.salt = None
    self.root_hash = None

  def GetVerityDataSize(self):
    """Calculate the verity size based on the size of the input image.

    Since we already know the structure of a verity enabled image to be:
    [file_system, verity_hash_tree, verity_metadata, fec_data]. We can then
    calculate the size and offset of each section.

    Returns:
      A tuple of filesystem_size, hash_tree_size, metadata_size
    """

    adjusted_size, _ = AdjustPartitionSizeForVerity(
        self.partition_size, self.fec_supported, self.block_size)
    assert adjusted_size % self.block_size == 0

    result, verity_tree_size = GetVerityTreeSize(adjusted_size)
    if not result:
      return None
    assert verity_tree_size % self.block_size == 0

    result, metadata_size = GetVerityMetadataSize(adjusted_size)
    if not result:
      return None
    assert metadata_size % self.block_size == 0

    return adjusted_size, verity_tree_size, metadata_size

  def ParseVerityMetadata(self, meta_data, file_system_size):
    """Parses the hash_algorithm, root_hash, salt from the metadata block."""

    # More info about the metadata structure available in:
    # system/extras/verity/build_verity_metadata.py
    META_HEADER_SIZE = 268
    header_bin = meta_data[0:META_HEADER_SIZE]
    header = struct.unpack("II256sI", header_bin)

    # header: magic_number, version, signature, table_len
    assert header[0] == 0xb001b001, header[0]
    table_len = header[3]
    verity_table = meta_data[META_HEADER_SIZE: META_HEADER_SIZE + table_len]
    table_entries = verity_table.rstrip().split()

    # Expected verity table format: "1 block_device block_device BLOCK_SIZE
    # BLOCK_SIZE data_blocks data_blocks hash_algorithm root_hash salt"
    assert len(table_entries) == 10, "Unexpected verity table size {}".format(
        len(table_entries))
    assert (int(table_entries[3]) == self.block_size and
            int(table_entries[4]) == self.block_size)
    assert (int(table_entries[5]) * self.block_size == file_system_size and
            int(table_entries[6]) * self.block_size == file_system_size)

    self.hash_algorithm = table_entries[7]
    self.root_hash = table_entries[8]
    self.salt = table_entries[9]

  def ValidateVerityTree(self):
    """Checks that we can reconstruct the verity hash tree."""

    # Writes the file system section to a temp file; and calls the executable
    # build_verity_tree to construct the hash tree.
    adjusted_partition = common.MakeTempFile(prefix="adjusted_partition")
    with open(adjusted_partition, "wb") as fd:
      self.tgt.WriteRangeDataToFd(self.file_system_range, fd)

    generated_verity_tree = common.MakeTempFile(prefix="verity")
    prop_dict = {}
    if not BuildVerityTree(adjusted_partition, generated_verity_tree,
                           prop_dict):
      return False

    assert prop_dict["verity_salt"] == self.salt
    if prop_dict["verity_root_hash"] != self.root_hash:
      print("Calculated verty root hash {} doesn't match the one in metadata"
            " {}".format(prop_dict["verity_root_hash"], self.root_hash))
      return False

    # Reads the generated hash tree and checks if it has the exact same bytes
    # as the one in the sparse image.
    with open(generated_verity_tree, "rb") as fd:
      return fd.read() == ''.join(self.tgt.ReadRangeSet(self.verity_tree_range))

  def GetVerityTreeInfo(self):
    """Parses and validates the hash_tree info in a sparse image.

    Returns:
      True: If we successfully constructs the hash_tree from the embedded
          metadata.
      False: If a sanity check fails or we fail to get the exact bytes of the
          hash tree.
    """

    size_tuple = self.GetVerityDataSize()
    if not size_tuple:
      return False
    file_system_size, verity_tree_size, meta_size = size_tuple
    self.verity_tree_range = RangeSet(
        data=[file_system_size / self.block_size,
              (file_system_size + verity_tree_size) / self.block_size])
    self.file_system_range = RangeSet(
        data=[0, file_system_size / self.block_size])

    meta_start = file_system_size + verity_tree_size
    meta_range = RangeSet(data=[meta_start / self.block_size,
                                (meta_start + meta_size) / self.block_size])
    meta_data = ''.join(self.tgt.ReadRangeSet(meta_range))
    self.ParseVerityMetadata(meta_data, file_system_size)

    if not self.ValidateVerityTree():
      print("Failed to reconstruct the verity tree")
      return False

    return True
