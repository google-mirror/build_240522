#!/usr/bin/env python3
#
# Copyright (C) 2023 The Android Open Source Project
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

and manifest from repo information and write it to a file in the top-directory.
This should be run only once per checkout, when sourcing build/envsetup.sh
"""

import argparse
import json
import os
import subprocess


def _get_default_output_file():
  """Returns the filepath for the build output."""
  build_top = os.getenv("ANDROID_BUILD_TOP")
  if not build_top:
    return None
  return os.path.join(build_top, "out", "branchinfo.json")


def _get_default_git_dir():
  """Returns the path of repo git directory."""
  build_top = os.getenv("ANDROID_BUILD_TOP")
  if not build_top:
    return None

  if os.path.exists(os.path.join(build_top, ".git")):
    return os.path.join(build_top, ".git")
  return os.path.join(build_top, ".repo", "manifests.git/")


def _parse_output(git_out):
  """Returns the git_output b string as a standard trimmed string."""
  return git_out.decode("utf-8").strip()


def _execute_git_command(cmd):
  """Executes the given git command and returns the output as a standard string."""

  # Suppress error output - make this a clean operation if it's not run by a git-user
  output = subprocess.check_output(cmd, shell=True, stderr=subprocess.DEVNULL)

  return _parse_output(output)


def _create_dir(path):
  """Creates the directory if it doesn't exist already."""
  path_dir = os.path.dirname(path)
  os.makedirs(path_dir, exist_ok=True)


def _save_data_to_file(filepath, manifest, branch_name):
  """Saves the given data to the given filepath"""
  _create_dir(filepath)

  dictionary = {"manifest": manifest, "branch_name": branch_name}
  output = json.dumps(dictionary, indent=4)
  with open(filepath, "w") as f:
    f.write(output)
    f.close()


def main():
  """Reads branch information from git_dir and writes it to a file.

  Specifically, the manifest and branch name will be written.
  """

  # Parse args
  parser = argparse.ArgumentParser(description="")
  default_output_file = _get_default_output_file()
  parser.add_argument(
      "output_file",
      nargs="?",
      default=default_output_file,
      help=(
          "The filepath to write branch information. Defaults to"
          f" {default_output_file}"
      ),
  )

  default_git_dir = _get_default_git_dir()
  parser.add_argument(
      "git_dir",
      nargs="?",
      default=default_git_dir,
      help=(
          "The git directory to execute commands in. Defaults to"
          f" {default_git_dir}"
      ),
  )

  args = parser.parse_args()
  git_dir = args.git_dir
  output_file = args.output_file

  if git_dir is None or output_file is None:
    print("Skipping branch-name write ")
    return

  # Call git to get pertinent info
  try:
    local_branch_cmd = f"git -C {git_dir} symbolic-ref --short HEAD"
    local_branch = _execute_git_command(local_branch_cmd)

    remote_cmd = f"git -C {git_dir} config --get branch.{local_branch}.remote"
    remote = _execute_git_command(remote_cmd)

    manifest_cmd = f"git -C {git_dir} remote get-url {remote}"
    manifest = _execute_git_command(manifest_cmd)

    branch_name_cmd = (
        f"git -C {git_dir} config --get branch.{local_branch}.merge"
    )
    branch_name = _execute_git_command(branch_name_cmd)
  except Exception as err:
    # Fail silently -
    print("Skipping branch-name write ")
    return
  # Save info

  _save_data_to_file(output_file, manifest, branch_name)


if __name__ == "__main__":
  main()
