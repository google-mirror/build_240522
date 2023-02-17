#!/usr/bin/env python3
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
"""A tool to parse the branch name (e.g., aosp-master-with-phones, master)

from repo information and write it to a file in the top-directory.

This should be run only once per checkout, when sourcing build/envsetup.sh
"""

import argparse
import os
import subprocess


def _get_default_output_file():
  """Returns the filepath for the build output."""
  build_top = os.getenv("ANDROID_BUILD_TOP")
  if not build_top:
    raise Exception(
        "$ANDROID_BUILD_TOP not found in environment. Have you run lunch?"
    )
  return os.path.join(build_top, ".branchinfo", "branchinfo.txt")


def _get_default_git_dir():
  """Returns the path of repo git directory."""
  build_top = os.getenv("ANDROID_BUILD_TOP")
  if not build_top:
    raise Exception(
        "$ANDROID_BUILD_TOP not found in environment. Have you run lunch?"
    )
  return os.path.join(build_top, ".repo", "manifests.git/")


def _parse_output(git_out):
  """Returns the git_output b string as a standard trimmed string."""
  return git_out.decode("utf-8").strip()


def _execute_git_command(cmd):
  """Executes the given git command and returns the output as a standard string."""
  output = subprocess.check_output(cmd, shell=True)

  return _parse_output(output)


def _create_dir(path):
  """Creates the directory if it doesn't exist already."""
  path_dir = os.path.dirname(path)
  subprocess.check_output(f"mkdir -p {path_dir}", shell=True)


def _save_data_to_file(filepath, manifest, branch_name):
  """Saves the given data to the given filepath"""
  _create_dir(filepath)
  file = open(filepath, "w")
  file.write(f"{manifest}\n")
  file.write(f"{branch_name}\n")
  file.close()


def main():
  # Parse args
  parser = argparse.ArgumentParser(description="")
  parser.add_argument(
      "output_file",
      nargs="?",
      default=_get_default_output_file(),
      help="The filepath to write branch information. "
      + "Defaults to .branchinfo/branchinfo.txt",
  )

  parser.add_argument(
      "git_dir",
      nargs="?",
      default=_get_default_git_dir(),
      help=(
          "The git directory to execute commands in. Defaults to"
          " .repo/manifests.git"
      ),
  )
  args = parser.parse_args()
  git_dir = args.git_dir

  # Call git to get pertinent info
  local_branch_cmd = f"git -C {git_dir} symbolic-ref --short HEAD"
  local_branch = _execute_git_command(local_branch_cmd)

  remote_cmd = f"git -C {git_dir} config --get branch.{local_branch}.remote"
  remote = _execute_git_command(remote_cmd)

  manifest_cmd = f"git -C {git_dir} remote get-url {remote}"
  manifest = _execute_git_command(manifest_cmd)

  branch_name_cmd = f"git -C {git_dir} config --get branch.{local_branch}.merge"
  branch_name = _execute_git_command(branch_name_cmd)
  # Save info

  _save_data_to_file(args.output_file, manifest, branch_name)


if __name__ == "__main__":
  main()
