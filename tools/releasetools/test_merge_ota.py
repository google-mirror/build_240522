# Copyright (C) 2008 The Android Open Source Project
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


import os
import tempfile
import test_utils
import merge_ota


class MergeOtaTest(test_utils.ReleaseToolsTestCase):
  def setUp(self) -> None:
    self.testdata_dir = test_utils.get_testdata_dir()
    return super().setUp()

  def test_MergeThreeOtas(self):
    ota1 = os.path.join(self.testdata_dir, "oriole_vbmeta.zip")
    ota2 = os.path.join(self.testdata_dir, "oriole_vbmeta_system.zip")
    ota3 = os.path.join(self.testdata_dir, "oriole_vbmeta_vendor.zip")
    with tempfile.NamedTemporaryFile() as output_file:
      merge_ota.main([ota1, ota2, ota3, "--output", output_file.name])
