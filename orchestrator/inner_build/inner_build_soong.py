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

import argparse
import subprocess
import sys

import os

sys.dont_write_bytecode = True
import common

class InnerBuildSoong(common.Commands):
    def describe(self, args):
        pass


    def export_api_contributions(self, args):
        # TODO: Add a verbose and error log to the Context object?
        current_dir = os.path.dirname(__file__)
        export_api_script_path = os.path.join(current_dir, "inner_build_export_api_contributions.sh")
        for api_domain in args.api_domain:
          cmd = [
              export_api_script_path,
              "--out_dir",
              args.out_dir,
              "--api_domain",
              api_domain,
              "--inner_tree",
              args.inner_tree,
          ]
          # TODO: write to log provided by context object instead of stdout
          print(f"Exporting the contributions of api_domain={api_domain} to out_dir={args.out_dir}")
          proc = subprocess.run(cmd, shell=False, capture_output=True)
          # TODO: Add to verbose log
          if proc.returncode:
            # TODO: Add to error log
            sys.stderr.write("export_api_contribution failed with error message:\n")
            sys.stderr.write(proc.stderr.decode())
            sys.exit(proc.returncode)



def main(argv):
    return InnerBuildSoong().Run(argv)


if __name__ == "__main__":
    sys.exit(main(sys.argv))
