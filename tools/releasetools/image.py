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
# See the License for the specific

from hashlib import sha1

from rangelib import RangeSet

__all__ = ["EmptyImage", "DataImage"]


class Image(object):
  def RangeSha1(self, ranges):
    raise NotImplementedError

  def ReadRangeSet(self, ranges):
    raise NotImplementedError

  def TotalSha1(self, include_clobbered_blocks=False):
    raise NotImplementedError

  def WriteRangeDataToFd(self, ranges, fd):
    raise NotImplementedError

  def GetFileMap(self):
    raise NotImplementedError

  @staticmethod
  def LoadFileBlockMap(care_map, map_path, blocksize, clobbered_blocks,
                       hashtree_info, block_read_func, allow_shared_blocks):
    """Loads the given block map file.

    If the map file presents, constructs a dictionary based on its content. For
    the remaining blocks in the care map, categorizes them into clobbered
    blocks, hashtree blocks, zero blocks, and nonzero blocks.

    Args:
      care_map: A RangeSet of all cared blocks on the device.
      map_path: The filename of the block map file.
      blocksize: The size in bytes for each block.
      clobbered_blocks: A RangeSet instance for the clobbered blocks.
      hashtree_info: An object contains the information of hashtree.
      block_read_func: A generator function that returns the block data within
          a given range.
      allow_shared_blocks: Whether having shared blocks is allowed.

    Returns:
      A dictionary that contains the filename and its ranges.
    """

    if map_path:
      out, remaining = Image.ParseMapFile(care_map, map_path, clobbered_blocks,
                                          allow_shared_blocks)
    else:
      out = {}
      remaining = care_map

    remaining = remaining.subtract(clobbered_blocks)
    if hashtree_info:
      remaining = remaining.subtract(hashtree_info.hashtree_range)

    # For all the remaining blocks in the care_map (ie, those that
    # aren't part of the data for any file nor part of the clobbered_blocks),
    # divide them into blocks that are all zero and blocks that aren't.
    # (Zero blocks are handled specially because (1) there are usually
    # a lot of them and (2) bsdiff handles files with long sequences of
    # repeated bytes especially poorly.)
    zero_blocks, nonzero_groups = Image.FindNonZeroBlocks(
        remaining, blocksize, block_read_func)

    assert zero_blocks or nonzero_groups or clobbered_blocks

    if zero_blocks:
      out["__ZERO"] = RangeSet(data=zero_blocks)
    if nonzero_groups:
      for i, blocks in enumerate(nonzero_groups):
        out["__NONZERO-%d" % i] = RangeSet(data=blocks)
    if clobbered_blocks:
      out["__COPY"] = clobbered_blocks
    if hashtree_info:
      out["__HASHTREE"] = hashtree_info.hashtree_range

    return out

  @staticmethod
  def ParseMapFile(care_map, map_path, clobbered_blocks, allow_shared_blocks):
    """Parses the map file and constructs a dict from filename to its ranges."""

    out = {}
    remaining = care_map
    with open(map_path) as f:
      for line in f:
        fn, ranges = line.split(None, 1)
        ranges = RangeSet.parse(ranges)

        if allow_shared_blocks:
          # Find the shared blocks that have been claimed by others. If so, tag
          # the entry so that we can skip applying imgdiff on this file.
          shared_blocks = ranges.subtract(remaining)
          if shared_blocks:
            non_shared = ranges.subtract(shared_blocks)
            if not non_shared:
              continue

            # There shouldn't anything in the extra dict yet.
            assert not ranges.extra, "Non-empty RangeSet.extra"

            # Put the non-shared RangeSet as the value in the block map, which
            # has a copy of the original RangeSet.
            non_shared.extra['uses_shared_blocks'] = ranges
            ranges = non_shared

        out[fn] = ranges
        assert ranges.size() == ranges.intersect(remaining).size()

        # Currently we assume that blocks in clobbered_blocks are not part of
        # any file.
        assert not clobbered_blocks.overlaps(ranges)
        remaining = remaining.subtract(ranges)

    return out, remaining

  @staticmethod
  def FindNonZeroBlocks(remaining, blocksize, read_func):
    """Parses the remaining blocks to separate out the zero & nonzero blocks.

    Returns:
      A tuple of (zero_blocks, nonzero_groups), the zero blocks is just a list
      of block indices; while the nonzero_groups is a list of list, with the
      size of each sublist smaller than MAX_BLOCKS_PER_GROUP.
    """

    zero_blocks = []
    nonzero_blocks = []
    reference = '\0' * blocksize

    # Workaround for bug 23227672. For squashfs, we don't have a system.map. So
    # the whole system image will be treated as a single file. But for some
    # unknown bug, the updater will be killed due to OOM when writing back the
    # patched image to flash (observed on lenok-userdebug MEA49). Prior to
    # getting a real fix, we evenly divide the non-zero blocks into smaller
    # groups (currently 1024 blocks or 4MB per group).
    # Bug: 23227672
    MAX_BLOCKS_PER_GROUP = 1024
    nonzero_groups = []

    for index, data in read_func(remaining):
      if data == reference:
        zero_blocks.append(index)
        zero_blocks.append(index + 1)
      else:
        nonzero_blocks.append(index)
        nonzero_blocks.append(index + 1)

        if len(nonzero_blocks) * 2 > MAX_BLOCKS_PER_GROUP:
          nonzero_groups.append(nonzero_blocks)
          nonzero_blocks = []

    if nonzero_blocks:
      nonzero_groups.append(nonzero_blocks)

    return zero_blocks, nonzero_groups


class EmptyImage(Image):
  """A zero-length image."""

  def __init__(self):
    self.blocksize = 4096
    self.care_map = RangeSet()
    self.clobbered_blocks = RangeSet()
    self.extended = RangeSet()
    self.total_blocks = 0
    self.file_map = {}
    self.hashtree_info = None

  def RangeSha1(self, ranges):
    return sha1().hexdigest()

  def ReadRangeSet(self, ranges):
    return ()

  def TotalSha1(self, include_clobbered_blocks=False):
    # EmptyImage always carries empty clobbered_blocks, so
    # include_clobbered_blocks can be ignored.
    assert self.clobbered_blocks.size() == 0
    return sha1().hexdigest()

  def WriteRangeDataToFd(self, ranges, fd):
    raise ValueError("Can't write data from EmptyImage to file")

  def GetFileMap(self):
    return self.file_map


class DataImage(Image):
  """An image wrapped around a single string of data."""

  def __init__(self, data, trim=False, pad=False, file_map_fn=None):
    self.data = data
    self.blocksize = 4096

    assert not (trim and pad)

    partial = len(self.data) % self.blocksize
    padded = False
    if partial > 0:
      if trim:
        self.data = self.data[:-partial]
      elif pad:
        self.data += '\0' * (self.blocksize - partial)
        padded = True
      else:
        raise ValueError(("data for DataImage must be multiple of %d bytes "
                          "unless trim or pad is specified") %
                         (self.blocksize,))

    assert len(self.data) % self.blocksize == 0

    self.total_blocks = len(self.data) / self.blocksize
    self.care_map = RangeSet(data=(0, self.total_blocks))
    # When the last block is padded, we always write the whole block even for
    # incremental OTAs. Because otherwise the last block may get skipped if
    # unchanged for an incremental, but would fail the post-install
    # verification if it has non-zero contents in the padding bytes.
    # Bug: 23828506
    if padded:
      clobbered_blocks = [self.total_blocks-1, self.total_blocks]
    else:
      clobbered_blocks = []
    self.clobbered_blocks = RangeSet(data=clobbered_blocks)
    self.extended = RangeSet()

    self.file_map = self.BuildFileMap(file_map_fn)

  def BuildFileMap(self, file_map_fn):
    """Initializes the file_map field for data images.

    Loads the file map if it presents. And divides the remaining ranges into
    __ZERO, __NONZERO, and clobbered_blocks.
    """

    # The block data generator for DataImage.
    def DataImageReader(block_ranges):
      for s, e in block_ranges:
        for i in range(s, e):
          yield i, self.data[i * self.blocksize: (i + 1) * self.blocksize]

    return Image.LoadFileBlockMap(self.care_map, file_map_fn, self.blocksize,
                                  self.clobbered_blocks, hashtree_info=None,
                                  block_read_func=DataImageReader,
                                  allow_shared_blocks=True)

  def GetFileMap(self):
    return self.file_map

  def _GetRangeData(self, ranges):
    for s, e in ranges:
      yield self.data[s*self.blocksize:e*self.blocksize]

  def RangeSha1(self, ranges):
    h = sha1()
    for data in self._GetRangeData(ranges):  # pylint: disable=not-an-iterable
      h.update(data)
    return h.hexdigest()

  def ReadRangeSet(self, ranges):
    return [self._GetRangeData(ranges)]

  def TotalSha1(self, include_clobbered_blocks=False):
    if not include_clobbered_blocks:
      return self.RangeSha1(self.care_map.subtract(self.clobbered_blocks))
    else:
      return sha1(self.data).hexdigest()

  def WriteRangeDataToFd(self, ranges, fd):
    for data in self._GetRangeData(ranges):  # pylint: disable=not-an-iterable
      fd.write(data)
