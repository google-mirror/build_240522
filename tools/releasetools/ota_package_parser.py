#!/usr/bin/env python
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

from __future__ import print_function

import argparse
import os
import sys
import traceback
import zipfile

from rangelib import RangeSet

class Stash(object):
  def __init__(self):
    self.blocks_stashed = 0
    self.overlap_blocks_stashed = 0
    self.max_stash_needed = 0
    self.current_stash_size = 0
    self.stash_map = {}

  def StashBlocks(self, SHA1, blocks):
    if self.stash_map.has_key(SHA1):
      print("already stashed {}: {}".format(SHA1, blocks))
      return
    self.blocks_stashed += blocks.size()
    self.current_stash_size += blocks.size()
    self.max_stash_needed = max(self.current_stash_size, self.max_stash_needed)
    self.stash_map[SHA1] = blocks

  def FreeBlocks(self, SHA1):
    assert self.stash_map.has_key(SHA1), "stash {} not found".format(SHA1)
    self.current_stash_size -= self.stash_map[SHA1].size()
    del self.stash_map[SHA1]

  def HandleOverlapBlocks(self, SHA1, blocks):
    self.StashBlocks(SHA1, blocks)
    self.overlap_blocks_stashed += blocks.size()
    self.FreeBlocks(SHA1)


class OtaPackageParser(object):
  def __init__(self, package):
    self.package = package
    self.new_data_size = 0
    self.patch_data_size = 0
    self.block_written = 0
    self.block_stashed = 0

  @staticmethod
  def GetSizeString(size):
    assert size >= 0
    base = 1024.0
    if size <= base:
      return "{} bytes".format(size)
    for units in ['K', 'M', 'G']:
      if size <= base * 1024 or units == 'G':
        return "{:.1f}{}".format(size / base, units)
      base *= 1024

  def ParseTransferList(self, name):
    print("\nSimulating commands in '{}':".format(name))
    lines = self.package.read(name).strip().splitlines()
    assert len(lines) >= 4, \
        "{} too short, expect at least 4 lines, has {}".format(name, len(lines))
    assert int(lines[0]) >= 3
    print("(version: {},  expected_write_total: {},  "
          "expected_max_cache_needed: {})".format(lines[0], lines[1], lines[3]))

    blocks_written = 0
    my_stash = Stash()
    for line in lines[4:]:
      cmd_list = line.strip().split(" ")
      cmd_name = cmd_list[0]
      try:
        if cmd_name == "new" or cmd_name == "zero":
          assert len(cmd_list) == 2, "command format error: {}".format(line)
          target_range = RangeSet.parse_raw(cmd_list[1])
          blocks_written += target_range.size()
        elif cmd_name == "move":
          # Example:  move <onehash> <tgt_range> <src_blk_count> <src_range>
          # [<loc_range> <stashed_blocks>]
          assert len(cmd_list) >= 5, "command format error: {}".format(line)
          target_range = RangeSet.parse_raw(cmd_list[2])
          blocks_written += target_range.size()
          if cmd_list[4] == '-':
            continue
          SHA1 = cmd_list[1]
          source_range = RangeSet.parse_raw(cmd_list[4])
          if target_range.overlaps(source_range):
            my_stash.HandleOverlapBlocks(SHA1, source_range)
        elif cmd_name == "bsdiff" or cmd_name == "imgdiff":
          # Example:  bsdiff <offset> <len> <src_hash> <tgt_hash> <tgt_range>
          # <src_blk_count> <src_range> [<loc_range> <stashed_blocks>]
          assert len(cmd_list) >= 8, "command format error: {}".format(line)
          target_range = RangeSet.parse_raw(cmd_list[5])
          blocks_written += target_range.size()
          if cmd_list[7] == '-':
            continue
          source_SHA1 = cmd_list[3]
          source_range = RangeSet.parse_raw(cmd_list[7])
          if target_range.overlaps(source_range):
            my_stash.HandleOverlapBlocks(source_SHA1, source_range)
        elif cmd_name == "stash":
          assert len(cmd_list) == 3, "command format error: {}".format(line)
          SHA1 = cmd_list[1]
          source_range = RangeSet.parse_raw(cmd_list[2])
          my_stash.StashBlocks(SHA1, source_range)
        elif cmd_name == "free":
          assert len(cmd_list) == 2, "command format error: {}".format(line)
          SHA1 = cmd_list[1]
          my_stash.FreeBlocks(SHA1)
      except:
        print("failed to parse command in: " + line)
        raise

    self.block_written += blocks_written
    self.block_stashed += my_stash.blocks_stashed

    print("\nblocks written: " + str(blocks_written))
    print("total blocks stashed: " + str(my_stash.blocks_stashed))
    print("blocks stashed implicitly: " + str(my_stash.overlap_blocks_stashed))
    print("max blocks stashed simultaneously: " + str(
        my_stash.max_stash_needed))

  def PrintDataInfo(self, partition):
    print("\nReading data info for {} partition:".format(partition))
    new_data = self.package.getinfo(partition + ".new.dat")
    patch_data = self.package.getinfo(partition + ".patch.dat")
    print("{:<40}{:<40}".format(new_data.filename, patch_data.filename))
    print("{:<40}{:<40}".format(
          "compress_type: " + str(new_data.compress_type),
          "compress_type: " + str(patch_data.compress_type)))
    print("{:<40}{:<40}".format(
          "compressed_size: " + OtaPackageParser.GetSizeString(
              new_data.compress_size),
          "compressed_size: " + OtaPackageParser.GetSizeString(
              patch_data.compress_size)))
    print("{:<40}{:<40}".format(
        "file_size: " + OtaPackageParser.GetSizeString(new_data.file_size),
        "file_size: " + OtaPackageParser.GetSizeString(patch_data.file_size)))

    self.new_data_size += new_data.file_size
    self.patch_data_size += patch_data.file_size

  def AnalyzePartition(self, partition):
    assert partition in ("system", "vendor")
    assert partition + ".new.dat" in self.package.namelist()
    assert partition + ".patch.dat" in self.package.namelist()
    assert partition + ".transfer.list" in self.package.namelist()

    self.PrintDataInfo(partition)
    self.ParseTransferList(partition + ".transfer.list")

  def PrintMetaData(self):
    meta_path = "META-INF/com/android/metadata"
    print("\nMeta data info:")
    meta_info = {}
    for line in self.package.read(meta_path).strip().splitlines():
      index = line.find("=")
      meta_info[line[0 : index].strip()] = line[index + 1:].strip()
    assert "ota-type" in meta_info and meta_info["ota-type"] == "BLOCK"
    assert "pre-device" in meta_info
    print("device: {}".format(meta_info["pre-device"]))
    if "pre-build" in meta_info:
      print("pre-build: {}".format(meta_info["pre-build"]))
    assert "post-build" in meta_info
    print("post-build: {}".format(meta_info["post-build"]))

  def Analyze(self):
    print("\nAnalyzing ota package: " + self.package.filename)
    self.PrintMetaData()
    assert "system.new.dat" in self.package.namelist()
    self.AnalyzePartition("system")
    if "vendor.new.dat" in self.package.namelist():
      self.AnalyzePartition("vendor")

    BLOCK_SIZE = 4096
    print("\nOTA package analyzed:")
    print("new data size(uncompressed): " + OtaPackageParser.GetSizeString(
        self.new_data_size))
    print("patch data size(uncompressed): " + OtaPackageParser.GetSizeString(
        self.patch_data_size))
    print("total data written: " + OtaPackageParser.GetSizeString(
        self.block_written * BLOCK_SIZE))
    print("total data stashed: " + OtaPackageParser.GetSizeString(
        self.block_stashed * BLOCK_SIZE))


def main(argv):
  parser = argparse.ArgumentParser(description='Analyze an OTA package.')
  parser.add_argument("ota_package", help='Path of the OTA package.')
  args, unknown = parser.parse_known_args(argv)
  try:
    with zipfile.ZipFile(args.ota_package, 'r') as package:
      package_parser = OtaPackageParser(package)
      package_parser.Analyze()
  except:
    print("Failed to read " + args.ota_package)
    traceback.print_exc()
    sys.exit(-1)


if __name__ == '__main__':
  main(sys.argv[1:])









