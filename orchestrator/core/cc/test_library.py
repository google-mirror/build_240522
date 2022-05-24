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

from ninja_tools import Ninja
from ninja_syntax import Rule, BuildAction

from cc.library import CompileContext, Compiler, LinkContext, Linker

class TestCompiler(unittest.TestCase):

    def test_compile_rule_created_implicitly(self):
        compiler = Compiler()
        compile_context = CompileContext("src", "flags", "out", frontend="my/clang")
        ninja = Ninja(context=None, file=None)
        compiler.compile(ninja, compile_context)
        assert any([isinstance(node, Rule) and node.name == "cc" for node in ninja.nodes])

    def test_clang_is_implicit_dep(self):
        compiler = Compiler()
        compile_context = CompileContext("src", "flags", "out", frontend="my/clang")
        ninja = Ninja(context=None, file=None)
        compiler.compile(ninja, compile_context)
        compile_action_nodes = [node for node in ninja.nodes if isinstance(node, BuildAction)]
        assert all(["my/clang" in node.implicits for node in compile_action_nodes])

    def test_compile_flags_are_added(self):
        compiler = Compiler()
        compile_context = CompileContext("src", "myflag1 myflag2 myflag3", "out", frontend="my/clang")
        ninja = Ninja(context=None, file=None)
        compiler.compile(ninja, compile_context)
        compile_action_node = [node for node in ninja.nodes if isinstance(node, BuildAction)]
        assert len(compile_action_node) == 1
        compile_action_node = compile_action_node[0]
        variables = compile_action_node.variables
        assert len(variables) == 2
        variables = sorted(variables , key=lambda x: x.name)
        assert variables[0].name == "cFlags"
        assert variables[0].value == compile_context.flags
        assert variables[1].name == "cFrontend"
        assert variables[1].value == compile_context.frontend

class TestLinker(unittest.TestCase):

    def test_link_rule_created_implicitly(self):
        linker = Linker()
        link_context = LinkContext("objs", "flags", "out", frontend="my/clang")
        ninja = Ninja(context=None, file=None)
        linker.link(ninja, link_context)
        assert any([isinstance(node, Rule) and node.name == "ld" for node in ninja.nodes])

    def test_clang_is_implicit_dep(self):
        linker = Linker()
        link_context = LinkContext("objs", "flags", "out", frontend="my/clang")
        ninja = Ninja(context=None, file=None)
        linker.link(ninja, link_context)
        link_action_nodes = [node for node in ninja.nodes if isinstance(node, BuildAction)]
        assert all(["my/clang" in node.implicits for node in link_action_nodes])

    def test_link_flags_are_added(self):
        linker = Linker()
        link_context = LinkContext("src", "myflag1 myflag2 myflag3", "out", frontend="my/clang")
        ninja = Ninja(context=None, file=None)
        linker.link(ninja, link_context)
        link_action_node = [node for node in ninja.nodes if isinstance(node, BuildAction)]
        assert len(link_action_node) == 1
        link_action_node = link_action_node[0]
        variables = link_action_node.variables
        assert len(variables) == 2
        variables = sorted(variables , key=lambda x: x.name)
        assert variables[0].name == "ldFlags"
        assert variables[0].value == link_context.flags
        assert variables[1].name == "ldFrontend"
        assert variables[1].value == link_context.frontend


if __name__ == "__main__":
      unittest.main()
