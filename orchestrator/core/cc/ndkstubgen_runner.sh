#!/bin/bash -eu

# Copyright 2017 Google Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# A wrapper script that sets the python path for running ndkstubgen
# This is necessary since orchestrator uses symlinks to ndkstubgen and its deps

# The module search path for symlinks in Python is calculated after following the symlinks
# https://docs.python.org/3/tutorial/modules.html#the-module-search-path

dir=$(dirname -- "$0")
PYTHONPATH="${PYTHONPATH}:${dir}" python "${dir}"/ndkstubgen.py $@
