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

import common
from image import DataImage
from rangelib import RangeSet
from sparse_img import SparseImage
from test_utils import ReleaseToolsTestCase


class ImageLoadFileMapTest(ReleaseToolsTestCase):
  """Checks the loaded block file map function for DataImage and SparseImage."""

  @staticmethod
  def _build_image_data():
    # block 0 :    zeros
    # blocks 1-3:  filled with 'a'
    # block 4:     first 10 bytes 'b', followed by zeros
    # blocks 5-8:  filled with 'c'
    # blocks 9-10: zeros
    data = ['\0' * 4096, 'a' * 4096 * 3, 'b' * 10 + '\0' * 4086, 'c' * 4096 * 4,
            '\0' * 4096 * 2]
    file_map_data = ["file_a 1-3", "file_c 5-8"]
    file_map_fn = common.MakeTempFile(suffix=".map")
    with open(file_map_fn, "wb") as f:
      f.write('\n'.join(file_map_data))

    return ''.join(data), file_map_fn

  def test_LoadFileMap_dataImage(self):
    image_data, file_map = self._build_image_data()

    image = DataImage(image_data, pad=True, file_map_fn=file_map)

    expected_file_map = {
        "file_a": RangeSet.parse("1-3"),
        "file_c": RangeSet.parse("5-8"),
        "__ZERO": RangeSet.parse("0 9-10"),
        "__NONZERO-0": RangeSet.parse("4")
    }
    self.assertEqual(expected_file_map, image.GetFileMap())

  def test_LoadFileMap_sparseImage(self):
    image_data, file_map = self._build_image_data()

    raw_image_file = common.MakeTempFile(prefix="raw-", suffix=".img")
    with open(raw_image_file, "wb") as f:
      f.write(image_data)

    simg_file = common.MakeTempFile(prefix="sparse-", suffix=".img")
    common.RunAndCheckOutput(["img2simg", raw_image_file, simg_file])
    image = SparseImage(simg_file, file_map_fn=file_map,
                        clobbered_blocks=[0, 1])

    expected_file_map = {
        "file_a": RangeSet.parse("1-3"),
        "file_c": RangeSet.parse("5-8"),
        "__COPY": RangeSet.parse("0"),
        "__ZERO": RangeSet.parse("9-10"),
        "__NONZERO-0": RangeSet.parse("4")
    }
    self.assertEqual(expected_file_map, image.GetFileMap())
