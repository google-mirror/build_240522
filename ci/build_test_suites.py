# Copyright 2024, The Android Open Source Project
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

"""
Script to build only the necessary modules for general-tests along
with whatever other targets are passed in.
"""

from collections.abc import Sequence
from typing import Set, Text, Dict, Any

from test_mapping_module_retriever import GetTestMappings

import argparse
import json
import pathlib
import re
import subprocess
import sys
import os

# List of modules that are always required to be in general-tests.zip
REQUIRED_MODULES = frozenset(['ats-tradefed-tests', 'DsuGsiIntegrationTest', 'adb-remount-sh', 'mediaroutertest', 'CompanionDeviceManagerMultiDeviceTestCases', 'soong_zip'])

def build_test_suites(argv):
  args = parse_args(argv)

  if args.change_info is None or not os.path.exists(args.change_info):
    # Not in presubmit, build everything.
    build_command = base_build_command(args)
    build_command.append('general-tests')

    run_command(build_command)
    return

  # Call the class to map changed files to modules to build.
  # TODO(lucafarsi): Move this into a replaceable class.
  modules_to_build = find_modules_to_build(pathlib.Path(args.change_info))

  # Call the build command with everything.
  build_command = base_build_command(args)
  build_command.extend(modules_to_build)

  run_command(build_command)

  zip_build_outputs(modules_to_build)

def parse_args(argv):
  argparser = argparse.ArgumentParser()
  argparser.add_argument('extra_targets', nargs='*', help='Extra test suites to build.')
  argparser.add_argument('--target_product')
  argparser.add_argument('--target_release')
  argparser.add_argument('--with_dexpreopt_boot_img_and_system_server_only', action='store_true')
  argparser.add_argument('--dist_dir')
  argparser.add_argument('--change_info')

  return argparser.parse_args()

def base_build_command(args: argparse.Namespace) -> list:
  build_command = []
  build_command.append('time')
  build_command.append('./build/soong/soong_ui.bash')
  build_command.append('--make-mode')
  build_command.append('dist')
  build_command.append('DIST_DIR=' + args.dist_dir)
  build_command.append('TARGET_PRODUCT=' + args.target_product)
  build_command.append('TARGET_RELEASE=' + args.target_release)
  build_command.extend(args.extra_targets)

  return build_command

def run_command(args: list[str]) -> str:
  result = subprocess.run(
      args=args,
      text=True,
      capture_output=True,
      check=False,
  )
  # If the process failed, print its stdout and propagate the exception.
  if not result.returncode == 0:
    print('Build command failed! output:')
    print('stdout: ' + result.stdout)
    print('stderr: ' + result.stderr)

  result.check_returncode()
  return result.stdout

def get_soong_var(var: str) -> str:
  return run_command(['./build/soong/soong_ui.bash', '--dumpvar-mode', '--abs', var]).strip()

def find_modules_to_build(change_info: pathlib.Path) -> Set[Text]:
  changed_files = find_changed_files(change_info)

  test_mappings = GetTestMappings(changed_files, set())

  modules_to_build = set(REQUIRED_MODULES)

  modules_to_build.update(read_modules_from_test_mappings(test_mappings, changed_files))

  return modules_to_build

def find_changed_files(change_info: pathlib.Path) -> Set[Text]:
  with open(change_info) as change_info_file:
    change_info_contents = json.load(change_info_file)

  changed_files = set()

  for change in change_info_contents['changes']:
    project_path = change.get('projectPath') + '/'

    for revision in change.get('revisions'):
      for file_info in revision.get('fileInfos'):
        changed_files.add(project_path + file_info.get('path'))

  return changed_files

def read_modules_from_test_mappings(test_mappings: Dict[str, Any], changed_files: Set[Text]) -> Set[Text]:
  module_list = set()

  # The test_mappings object returned by GetTestMappings is organized as
  # follows:
  # {
  #   'test_mapping_file_path': {
  #     'group_name' : [
  #       'name': 'module_name',
  #     ],
  #   }
  # }
  for test_mapping in test_mappings.values():
    for group in test_mapping.values():
      for element in group:
        module = element.get('name', None)
        if module and match_file_pattern(test_mapping, element, changed_files):
          module_list.add(module)

  return module_list

def match_file_pattern(test_mapping: Text, element: Dict[str, Any], changed_files: Set[Text]) -> bool:
  file_patterns = element.get('file_patterns', None)
  if not file_patterns:
    return True
  for changed_file in changed_files:
    for pattern in file_patterns:
      if re.search(pattern, changed_file):
        return True

  return False

def zip_build_outputs(modules_to_build: Set[Text]):
  src_top = os.environ.get("TOP", os.getcwd())

  # Call dumpvars to get the necessary things.
  # TODO(lucafarsi): Don't call soong_ui 4 times for this, --dumpvars-mode can
  # do it but it requires parsing.
  host_out_testcases = get_soong_var('HOST_OUT_TESTCASES')
  target_out_testcases = get_soong_var('TARGET_OUT_TESTCASES')
  product_out = get_soong_var('PRODUCT_OUT')
  soong_host_out = get_soong_var('SOONG_HOST_OUT')
  host_out = get_soong_var('HOST_OUT')

  # Call the class to package the outputs.
  # TODO(lucafarsi): Move this code into a replaceable class.
  host_paths = []
  target_paths = []
  for module in modules_to_build:
    host_path = os.path.join(host_out_testcases, module)
    if os.path.exists(host_path):
      host_paths.append(host_path)

    target_path = os.path.join(target_out_testcases, module)
    if os.path.exists(target_path):
      target_paths.append(target_path)

  zip_command = ['time', os.path.join(host_out, 'soong_zip')]

  # Add host testcases.
  zip_command.append('-C')
  zip_command.append(os.path.join(src_top, soong_host_out))
  zip_command.append('-P')
  zip_command.append('host/')
  for path in host_paths:
    zip_command.append('-D')
    zip_command.append(path)

  # Add target testcases.
  zip_command.append('-C')
  zip_command.append(os.path.join(src_top, product_out))
  zip_command.append('-P')
  zip_command.append('target')
  for path in target_paths:
    zip_command.append('-D')
    zip_command.append(path)

  # Add necessary tools. These are also hardcoded in general-tests.mk.
  framework_path = os.path.join(soong_host_out, 'framework')

  zip_command.append('-C')
  zip_command.append(framework_path)
  zip_command.append('-P')
  zip_command.append('host/tools')
  zip_command.append('-f')
  zip_command.append(os.path.join(framework_path, 'cts-tradefed.jar'))
  zip_command.append('-f')
  zip_command.append(os.path.join(framework_path, 'compatibility-host-util.jar'))
  zip_command.append('-f')
  zip_command.append(os.path.join(framework_path, 'vts-tradefed.jar'))

  # Zip to the DIST dir.
  zip_command.append('-o')
  zip_command.append(os.path.join(dist_dir, 'general-tests.zip'))

  run_command(zip_command)

if __name__ == '__main__':
  build_test_suites(sys.argv)
