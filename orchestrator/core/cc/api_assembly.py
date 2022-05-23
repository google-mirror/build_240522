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

import os

from cc.stub_generator import StubGenerator, GenCcStubsInput

# TODO: Move this from global scope to some context object
stub_generator = StubGenerator()

def assemble_cc_api_library(context, ninja, build_file, stub_library):
    print("\nassembling cc_api_library %s-%s %s from:" % (stub_library.api_surface,
        stub_library.api_surface_version, stub_library.name))
    for contrib in stub_library.contributions:
        print("  %s %s" % (contrib.api_domain, contrib.library_contribution))

    staging_dir = context.out.api_library_dir(stub_library.api_surface,
            stub_library.api_surface_version, stub_library.name)
    work_dir = context.out.api_library_work_dir(stub_library.api_surface,
            stub_library.api_surface_version, stub_library.name)
    print("staging_dir=%s" % (staging_dir))
    print("work_dir=%s" % (work_dir))

    # Generate rules to copy headers
    includes = []
    include_dir = os.path.join(staging_dir, "include")
    for contrib in stub_library.contributions:
        for headers in contrib.library_contribution["headers"]:
            root = headers["root"]
            for file in headers["files"]:
                # TODO: Deal with collisions of the same name from multiple contributions
                include = os.path.join(include_dir, file)
                ninja.add_copy_file(include, os.path.join(contrib.inner_tree.root, root, file))
                includes.append(include)

    # Generate rule to run ndkstubgen
    stub_generator.add_stub_gen_rule(ninja)
    # Generate rule to build api_levels.json
    # This file does not vary across (surface,version,lib) tuples. Write this file to the top-level
    # work directory
    stub_generator.add_version_map_file(ninja, context.out.api_library_work_dir("", "", ""))

    arch = _get_device_arch()
    for contrib in stub_library.contributions:
        # Copy API file from inner tree to staging directory
        api = contrib.library_contribution["api"]
        # TODO: This should be a single element, list check should not be required
        # TODO: Update configs/bazel rules to reflect this
        api = api[0] if isinstance(api, list) else api
        api_file_staging_dir = os.path.join(staging_dir, api)
        ninja.add_copy_file(api_file_staging_dir,
                            os.path.join(contrib.inner_tree.root, api))

        # Generate stub .c files using ndkstubgen
        inputs = GenCcStubsInput(arch=arch,
                                 version= stub_library.api_surface_version,
                                 api=api_file_staging_dir,
                                 )
        stub_outputs = stub_generator.add_stub_gen_action(ninja, inputs, work_dir)

        # Compile stub .c files to .o files

        # Link .o file to .so file

    # Generate phony rule to build the library
    # TODO: This name probably conflictgs with something
    ninja.add_phony("-".join((stub_library.api_surface, str(stub_library.api_surface_version),
            stub_library.name)), includes)

    # Generate build files

def _get_device_arch() -> str:
    """Returns architecture of the target device
    This is used by ndkstubgen to gate arch-specific APIs defined in .map.txt
    files"""
    # TODO: This should be configured somewhere (probably lunch?)
    # TODO: Return "arm64" for now
    return "arm64"
