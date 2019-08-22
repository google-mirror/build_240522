#
# Copyright (C) 2019 The Android Open Source Project
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
#

import os.path

import common
import test_utils
from check_target_files_vintf import CheckVintf

class CheckTargetFilesVintfTest(test_utils.ReleaseToolsTestCase):

  def setUp(self):
    self.testdata_dir = test_utils.get_testdata_dir()

  @testutils.SkipIfExternalToolsUnavailable
  def test_CheckVintf_sanity(self):
    self.assertTrue(CheckVintf(os.path.join(self.testdata_dir,
                                            "vintf/sanity")))

  @testutils.SkipIfExternalToolsUnavailable
  def test_CheckVintf_bad_matrix(self):
    self.assertFalse(CheckVintf(os.path.join(self.testdata_dir,
                                             "vintf/bad_matrix")))

  @testutils.SkipIfExternalToolsUnavailable
  def test_CheckVintf_kernel_compat(self):
    self.assertTrue(CheckVintf(os.path.join(self.testdata_dir,
                                            "vintf/kernel_compat")))

  @testutils.SkipIfExternalToolsUnavailable
  def test_CheckVintf_kernel_incompat(self):
    self.assertFalse(CheckVintf(os.path.join(self.testdata_dir,
                                             "vintf/kernel_incompat")))

  @testutils.SkipIfExternalToolsUnavailable
  def test_CheckVintf_sku_compat(self):
    self.assertFalse(CheckVintf(os.path.join(self.testdata_dir,
                                             "vintf/sku_compat")))

  @testutils.SkipIfExternalToolsUnavailable
  def test_CheckVintf_sku_incompat(self):
    self.assertFalse(CheckVintf(os.path.join(self.testdata_dir,
                                             "vintf/sku_incompat")))


