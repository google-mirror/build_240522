#!/usr/bin/env python
#
# Copyright (C) 2022 The Android Open Source Project
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

import unittest

from io import StringIO

from ninja_writer import Writer
from ninja_syntax import Variable, Rule, BuildAction

class TestWriter(unittest.TestCase):

  def test_simple_writer(self):
    with StringIO() as f:
      writer = Writer(f)
      writer.add_variable(Variable(name="cflags", value="-Wall"))
      writer.add_newline()
      cc = Rule(name="ccinputs=["foo.c"])


  def test_comm