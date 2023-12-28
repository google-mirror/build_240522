# Copyright 2023, The Android Open Source Project
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

from test_mapping_parser import getModulesToBuild

import argparse
import subprocess
import os

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

def base_build_command(args: argparse.Namespace) -> list:
  build_command = []
  build_command.append('time')
  build_command.append('./build/soong/soong_ui.bash')
  build_command.append('--make-mode')
  build_command.append('dist')
  build_command.append('DIST_DIR=' + args.dist_dir)
  build_command.append('TARGET_PRODUCT=' + args.target_product)
  build_command.append('TARGET_RELEASE=' + args.target_release)

  return build_command

def get_soong_var(var: str) -> str:
  return run_command(['./build/soong/soong_ui.bash', '--dumpvar-mode', '--abs', var]).strip()

if __name__ == '__main__':
  src_top = os.environ.get("TOP", os.getcwd())

  argparser = argparse.ArgumentParser()
  argparser.add_argument('extra_targets', nargs='*', help='Extra test suites to build.')
  argparser.add_argument('--target_product')
  argparser.add_argument('--target_release')
  argparser.add_argument('--with_dexpreopt_boot_img_and_system_server_only', action='store_true')
  argparser.add_argument('--dist_dir')
  argparser.add_argument('--change_info')

  args = argparser.parse_args()

  if args.change_info is None or not os.path.exists(args.change_info):
    # Not in presubmit, build everything.
    build_command = base_build_command(args)
    build_command.append('general-tests')
    build_command.extend(args.extra_targets)

    run_command(build_command)
    exit

  # Call dumpvars to get the necessary things.
  # TODO(lucafarsi): don't call soong_ui 4 times for this, --dumpvars-mode can
  # do it but it requires parsing.
  host_out_testcases = get_soong_var('HOST_OUT_TESTCASES')
  target_out_testcases = get_soong_var('TARGET_OUT_TESTCASES')
  product_out = get_soong_var('PRODUCT_OUT')
  soong_host_out = get_soong_var('SOONG_HOST_OUT')

  # Call the class to map changed files to modules to build.
  # TODO(lucafarsi): move this into a replaceable class.
  modules_to_build = getModulesToBuild(args.change_info)

  # Call the build command with everything.
  build_command = base_build_command(args)
  build_command.extend(modules_to_build)
  build_command.extend(args.extra_targets)

  run_command(build_command)

  # Call the class to package the outputs.
  # TODO(lucafarsi): move this code into a replaceable class
  host_paths = []
  target_paths = []
  for module in modules_to_build:
    host_path = os.path.join(host_out_testcases, module)
    if os.path.exists(host_path):
      host_paths.append(host_path)

    target_path = os.path.join(target_out_testcases, module)
    if os.path.exists(target_path):
      target_paths.append(target_path)

  build_command = ['time', './out/host/linux-x86/bin/soong_zip']

  # Add host testcases.
  build_command.append('-C')
  build_command.append(os.path.join(src_top, soong_host_out))
  build_command.append('-P')
  build_command.append('host/')
  for path in host_paths:
    build_command.append('-D')
    build_command.append(path)

  # Add target testcases.
  build_command.append('-C')
  build_command.append(os.path.join(src_top, product_out))
  build_command.append('-P')
  build_command.append('target')
  for path in target_paths:
    build_command.append('-D')
    build_command.append(path)

  # Add necessary tools.
  framework_path = os.path.join(soong_host_out, 'framework')

  build_command.append('-C')
  build_command.append(framework_path)
  build_command.append('-P')
  build_command.append('host/tools')
  build_command.append('-f')
  build_command.append(os.path.join(framework_path, 'cts-tradefed.jar'))
  build_command.append('-f')
  build_command.append(os.path.join(framework_path, 'compatibility-host-util.jar'))
  build_command.append('-f')
  build_command.append(os.path.join(framework_path, 'vts-tradefed.jar'))

  # zip to the DIST dir
  build_command.append('-o')
  build_command.append(os.path.join(args.dist_dir, 'general-tests.zip'))

  run_command(build_command)
