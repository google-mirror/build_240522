#!/usr/bin/python3
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


from ninja_tools import Ninja
from ninja_syntax import BuildAction, Rule

from typing import List

"""A module for compiling and linking cc artifacts"""


class CompileContext():

    def __init__(self, src: str, flags:str, out: str, frontend: str):
        self.src = src
        self.flags = flags
        self.out = out
        self.frontend = frontend

class Compiler():

    def __init__(self):
        self._compile_rule = None

    def create_compile_rule(self, ninja: Ninja) -> Rule:
        rule = Rule("cc")
        rule.add_variable("description", "compile source to object file using clang/clang++")
        rule.add_variable("command", "${cFrontend} -c ${cFlags} -o ${out} ${in}")
        ninja.add_rule(rule)
        return rule

    def compile(self, ninja: Ninja, compile_context: CompileContext) -> None:
        if not self._compile_rule:
            self._compile_rule = self.create_compile_rule(ninja)

        compile_action = BuildAction(output=compile_context.out,
                    inputs=compile_context.src,
                    rule=self._compile_rule.name,
                    implicits=[compile_context.frontend]
                    )
        compile_action.add_variable("cFrontend", compile_context.frontend)
        compile_action.add_variable("cFlags", compile_context.flags)
        ninja.add_build_action(compile_action)

class LinkContext():

    def __init__(self, objs: List[str], flags: str, out: str, frontend: str):
        self.objs = objs
        self.flags = flags
        self.out = out
        self.frontend = frontend

class Linker():

    def __init__(self):
        self._link_rule = None

    def create_link_rule(self, ninja: Ninja) -> Rule:
        rule = Rule("ld")
        rule.add_variable("description", "link object files using clang/clang++")
        rule.add_variable("command", "${ldFrontend} ${ldFlags} -o ${out} ${in}")
        ninja.add_rule(rule)
        return rule

    def link(self, ninja: Ninja, link_context: LinkContext) -> None:
        if not self._link_rule:
            self._link_rule = self.create_link_rule(ninja)

        link_action = BuildAction(output=link_context.out,
                inputs=link_context.objs,
                rule=self._link_rule.name,
                implicits=[link_context.frontend]
                )
        link_action.add_variable("ldFrontend", link_context.frontend)
        link_action.add_variable("ldFlags", link_context.flags)
        ninja.add_build_action(link_action)
