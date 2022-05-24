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
from cc.library import CompileContext, Compiler, LinkContext, Linker

ARCH = ["arm", "arm64", "x86", "x86_64"]

class CcApiAssemblyContext(object):
    """Context object for managing global state of CC API Assembly"""
    def __init__(self):
        self._stub_generator = StubGenerator()
        self._compiler = Compiler()
        self._linker = Linker()

    def get_cc_api_library_assembler(self):
        """Return a callback function to assemble CC library APIs
        The callback is a member of the context object, and therefore has access to its state"""
        return self.assemble_cc_api_library

    def assemble_cc_api_library(self, context, ninja, build_file, stub_library):
        staging_dir = context.out.api_library_dir(stub_library.api_surface,
                stub_library.api_surface_version, stub_library.name)
        work_dir = context.out.api_library_work_dir(stub_library.api_surface,
                stub_library.api_surface_version, stub_library.name)

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
        self._stub_generator.add_stub_gen_rule(ninja)
        # Generate rule to build api_levels.json
        # This file does not vary across (surface,version,lib) tuples. Write this file to the top-level
        # work directory
        self._stub_generator.add_version_map_file(ninja, context.out.api_library_work_dir("", "", ""))

        for contrib in stub_library.contributions:
            # Copy API file from inner tree to staging directory
            api = contrib.library_contribution["api"]
            api_file_staging_dir = os.path.join(staging_dir, api)
            ninja.add_copy_file(api_file_staging_dir,
                                os.path.join(contrib.inner_tree.root, api))

            # Generate ninja _build actions_ for every single arch
            # The artifacts will not be built unless the arch matches the current lunch config
            #
            # This redundancy is acceptable since this means that we can choose to not
            # recompile the ninja file if we use a different lunch config
            for arch in ARCH:
                # Generate stub .c files using ndkstubgen
                # TODO: Add --llndk or --apex to additonal_args for the correct API surfaces
                inputs = GenCcStubsInput(arch=arch,
                                         version= stub_library.api_surface_version,
                                         api=api_file_staging_dir,
                                         additional_args="",
                                         )
                stub_outputs = self._stub_generator.add_stub_gen_action(ninja, inputs, work_dir)
                # Compile stub .c files to .o files
                object_file = stub_outputs.stub_src + ".o"
                # These compile flags have been plugged in from ndk_library.go
                flags = " ".join(["-Wno-incompatible-library-redeclaration",
                        "-Wno-incomplete-setjmp-declaration",
                        "-Wno-builtin-requires-header",
                        "-Wno-invalid-noreturn",
                        "-Wall",
                        "-Werror",
                        ])
                compile_context = CompileContext(
                        src=stub_outputs.stub_src,
                        flags=flags,
                        out=object_file,
                        frontend=context.tools.clang(),
                )
                self._compiler.compile(ninja, compile_context)

                # Link .o file to .so file
                soname = stub_library.name + ".so"
                shared_library = os.path.join(staging_dir, arch, soname)
                flags = " ".join([
                    "-shared",
                    f"-Wl,-soname,{soname}",
                    f"-Wl,--version-script,{stub_outputs.version_script}",
                    ])
                link_context = LinkContext(
                        objs=[object_file],
                        flags=flags,
                        out=shared_library,
                        frontend=context.tools.clang_cc(),
                )
                self._linker.link(ninja, link_context)


        # Generate phony rule to build the library
        # TODO: This name probably conflictgs with something
        ninja.add_phony("-".join((stub_library.api_surface, str(stub_library.api_surface_version),
                stub_library.name)), includes)

        # Generate build files
