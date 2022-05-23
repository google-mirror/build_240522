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
# ibuted under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import unittest

import os
import sys

from ninja_tools import Ninja
from ninja_syntax import Variable, BuildAction
from cc.stub_generator import StubGenerator, GenCcStubsInput, NDKSTUBGEN

class TestStubGenerator(unittest.TestCase):

    def test_ndkstubgen_var_is_added(self):
        ninja = Ninja(context=None, file=None)
        stub_generator = StubGenerator()
        stub_generator.add_stub_gen_rule(ninja)
        variables = [node for node in ninja.nodes if isinstance(node, Variable)]
        assert variables
        assert any([var.value == NDKSTUBGEN for var in variables])

    # stubgen ninja build actions without a stubgen ninja rule is an exception
    def test_gen_rule_is_required(self):
        ninja = Ninja(context=None, file=None)
        stub_generator = StubGenerator()
        stub_inputs = GenCcStubsInput("x86", "33", "libfoo.map.txt")
        with self.assertRaises(Exception):
            stub_generator.add_stub_gen_action(ninja, stub_inputs, "out")

    # the ndkstubgen binary is an implicit deps
    # ninja should recompile stubs if it changes
    def test_implicit_deps(self):
        ninja = Ninja(context=None, file=None)
        stub_generator = StubGenerator()
        stub_generator.add_stub_gen_rule(ninja)
        stub_inputs = GenCcStubsInput("x86", "33", "libfoo.map.txt")
        stub_generator.add_stub_gen_action(ninja, stub_inputs, "out")
        build_actions = [node for node in ninja.nodes if isinstance(node,
                                                                    BuildAction)]
        assert build_actions
        assert all([NDKSTUBGEN in build_action.implicits for build_action in
                    build_actions])

    def test_output_contains_c_stubs(self):
        ninja = Ninja(context=None, file=None)
        stub_generator = StubGenerator()
        stub_generator.add_stub_gen_rule(ninja)
        stub_inputs = GenCcStubsInput("x86", "33", "libfoo.map.txt")
        outputs = stub_generator.add_stub_gen_action(ninja, stub_inputs, "out")
        assert len(outputs) > 0
        assert "stub.c" in outputs.stub_src
        assert "stub.map" in outputs.version_script

if __name__ == "__main__":
      unittest.main()
