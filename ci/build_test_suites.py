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

"""Script to build only the necessary modules for general-tests along

with whatever other targets are passed in.
"""

import argparse
import logging
import os
import pathlib
import subprocess
import sys


class BuildFailure(Exception):
  def __init__(self, message, return_code):
    super().__init__(message)

    self.return_code = return_code


REQUIRED_ENV_VARS = frozenset(['TARGET_PRODUCT', 'TARGET_RELEASE', 'TOP'])

SOONG_UI_EXE_REL_PATH = 'build/soong/soong_ui.bash'

def get_top() -> pathlib.Path:
  return pathlib.Path(os.environ['TOP'])


def build_test_suites(argv):
  check_required_env()
  args = parse_args(argv)

  try:
    build_everything(args)
  except BuildFailure as e:
    logging.error('Build command failed! Check build_log for details.')
    return e.return_code

  return 0


def check_required_env():
  for env_var in REQUIRED_ENV_VARS:
    if env_var not in os.environ:
      raise RuntimeError(f'Required env var {env_var} not found! Aborting.')


def parse_args(argv):
  argparser = argparse.ArgumentParser()
  argparser.add_argument(
      'extra_targets', nargs='*', help='Extra test suites to build.'
  )

  return argparser.parse_args(argv)


def build_everything(args: argparse.Namespace):
  build_command = base_build_command(args, args.extra_targets)
  build_command.append('general-tests')

  try:
    run_command(build_command)
  except subprocess.CalledProcessError as e:
    # TODO(lucafarsi): Don't break the causal chain.
    raise BuildFailure(str(e), e.returncode)


def base_build_command(
    args: argparse.Namespace, extra_targets: set[str]
) -> list:
  build_command = []
  build_command.append(get_top().joinpath(SOONG_UI_EXE_REL_PATH))
  build_command.append('--make-mode')
  build_command.append('dist')
  build_command.extend(extra_targets)

  return build_command


def run_command(args: list[str], print_cmd: bool = True):
  if print_cmd:
    # TODO(lucafarsi): Use join.
    print('+ ' + str(args))
  subprocess.run(args=args, check=True)


def main(argv):
  sys.exit(build_test_suites(argv))
