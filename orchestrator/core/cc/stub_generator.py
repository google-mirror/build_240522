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

import json
import os
import sys

from typing import NamedTuple

from ninja_tools import Ninja
from ninja_syntax import Variable, Rule, BuildAction

# TODO: This is not hermetic! Make this script into a standalone executable
NDKSTUBGEN="build/build/orchestrator/core/cc/ndkstubgen_runner.sh"

"""A module for writing ninja rules that generate C stubs"""

class GenCcStubsOutput(NamedTuple):
    stub_src: str
    version_script: str
    symbol_list: str

class GenCcStubsInput:
    def __init__(self, arch: str, version: str, api: str, additional_args=""):
        self.arch = arch # target device arch (e.g. x86)
        self.version = version # numeric API level (e.g. 33)
        self.api = api # path to map.txt
        self.additional_args = additional_args # additional args to ndkstubgen (e.g. --llndk)

class StubGenerator:

    def __init__(self):
        self._stubgen_rule = None
        self._version_map_file = None

    def add_stub_gen_rule(self, ninja: Ninja):
        """This adds a ninja rule to run ndkstubgen
        Running ndkstubgen creates C stubs from API .map.txt files"""
        if self._stubgen_rule:
            return
        # Create a variable name for the binary
        ninja.add_variable(Variable("ndkstubgen", NDKSTUBGEN))

        # Add a rule to the ninja file
        rule = Rule("genCcStubsRule")
        rule.add_variable("description", "Generate stub .c files from .map.txt API description files")
        rule.add_variable("command", "${ndkstubgen} --arch ${arch} --api ${apiLevel} --api-map ${apiMap} ${additionalArgs} ${in} ${out}")
        ninja.add_rule(rule)
        self._stubgen_rule = rule

    def add_version_map_file(self, ninja: Ninja, api_surfaces_work_dir: str):
        """This adds a ninja build statement to generate api_levels.json
        This file is an implicit dependency of ndkstubgen"""
        if self._version_map_file:
            return
        version_map_file = os.path.join(api_surfaces_work_dir, "api_levels.json")
        ninja.add_write_file(version_map_file, json.dumps(_get_api_levels()))
        self._version_map_file = version_map_file

    def add_stub_gen_action(self, ninja: Ninja, stub_input: GenCcStubsInput,
                            work_dir: str) -> GenCcStubsOutput:
        """This adds a ninja build action to generate stubs using `genCcStubsRule`"""
        if self._stubgen_rule is None:
            raise Exception("Cannot create ninja build statement to generate C stubs without creating a ninja rule")
        if self._version_map_file is None:
            raise Exception("Cannot create ninja build statement to generate C stubs without a ninja statement to build api_levels.json")

        outputs = GenCcStubsOutput(stub_src=os.path.join(work_dir, stub_input.arch, "stub.c"),
                                   version_script=os.path.join(work_dir, stub_input.arch, "stub.map"),
                                   symbol_list=os.path.join(work_dir, stub_input.arch, "abi_symbol_list.txt"))

        # Create the ninja build action
        gen_stub_build_action = BuildAction(output=list(outputs),
                                            inputs=stub_input.api,
                                            rule=self._stubgen_rule.name,
                                            implicits=[NDKSTUBGEN,
                                                       self._version_map_file]
                                            )
        gen_stub_build_action.add_variable("arch", stub_input.arch)
        gen_stub_build_action.add_variable("apiLevel", stub_input.version)
        gen_stub_build_action.add_variable("apiMap", self._version_map_file)
        gen_stub_build_action.add_variable("additionalArgs", stub_input.additional_args)

        # Add the build action to the ninja file
        ninja.add_build_action(gen_stub_build_action)
        return outputs

# ndkstubgen uses a map for converting Android version codes to a numeric code
# e.g. "R" --> 30
# The map contains active_codenames as well, which get mapped to a preview level
# (9000+)
# TODO: Soong does it in android/api_levels.go, figure out the solution for
# multi-tree
PLATFORM_VERSION_ALL_CODENAMES = ["Tiramisu", "UpsideDownCake"]
PREVIEW_API_LEVEL_BASE = 9000
def _get_api_levels() -> dict:
    api_levels = {
            "G": 9,
            "I": 14,
            "J": 16,
            "J-MR1": 17,
            "J-MR2": 18,
            "K":     19,
            "L":     21,
            "L-MR1": 22,
            "M":     23,
            "N":     24,
            "N-MR1": 25,
            "O":     26,
            "O-MR1": 27,
            "P":     28,
            "Q":     29,
            "R":     30,
            "S":     31,
            "S-V2":  32,
    }
    for index, codename in enumerate(PLATFORM_VERSION_ALL_CODENAMES):
        api_levels[codename] = PREVIEW_API_LEVEL_BASE + index
    return api_levels
