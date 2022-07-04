#!/usr/bin/env python3
#
# Copyright 2022 Google Inc. All rights reserved.
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
`build_manifest_generator` generates build manifest APK file containing digests
of target files.
"""

import argparse
import common
import os
import subprocess
import sys

HASH_ALGORITHM = 'sha256'

def _digest(fsverity_path, input_file):
  cmd = [fsverity_path, 'digest', input_file]
  cmd.extend(['--compact'])
  cmd.extend(['--hash-alg', HASH_ALGORITHM])
  out = subprocess.check_output(cmd, universal_newlines=True).strip()
  return bytes(bytearray.fromhex(out))

if __name__ == '__main__':
  p = argparse.ArgumentParser()
  p.add_argument(
      '--output',
      help='Path to the output manifest APK',
      required=True)
  p.add_argument(
      '--fsverity-path',
      help='path to the fsverity program',
      required=True)
  p.add_argument(
      '--fsverity-manifest-generator-path',
      help='path to the fsverity_manifest_generator program',
      required=True)
  p.add_argument(
      '--aapt2-path',
      help='path to the aapt2 program',
      required=True)
  p.add_argument(
      '--min-sdk-version',
      help='minimum supported sdk version of the generated manifest apk',
      required=True)
  p.add_argument(
      '--version-code',
      help='version code for the generated manifest apk',
      required=True)
  p.add_argument(
      '--version-name',
      help='version name for the generated manifest apk',
      required=True)
  p.add_argument(
      '--framework-res',
      help='path to framework-res.apk',
      required=True)
  p.add_argument(
      '--apksigner-path',
      help='path to the apksigner program',
      required=True)
  p.add_argument(
      '--apk-key-path',
      help='path to the apk key',
      required=True)
  p.add_argument(
      '--apk-manifest-path',
      help='path to AndroidManifest.xml',
      required=True)
  p.add_argument(
      '--base-dir',
      help='directory to use as a relative root for the inputs',
      required=True)
  p.add_argument(
      'inputs',
      nargs='+',
      help='input file for the build manifest')
  args = p.parse_args(sys.argv[1:])

  temp_dir = common.MakeTempDir()

  os.mkdir(os.path.join(temp_dir, "assets"))
  manifest_path = os.path.join(temp_dir, "assets", "build_manifest.pb")
  common.RunAndCheckOutput([args.fsverity_manifest_generator_path,
      "--fsverity-path", args.fsverity_path,
      "--base-dir", args.base_dir,
      "--output" manifest_path].append(args.inputs))

  common.RunAndCheckOutput([args.aapt2_path, "link",
      "-A", os.path.join(temp_dir, "assets"),
      "-o", args.output,
      "--min-sdk-version", args.min_sdk_version,
      "--version-code", args.version_code,
      "--version-name", args.version_name,
      "-I", args.framework_res,
      "--manifest", args.apk_manifest_path])
  common.RunAndCheckOutput([args.apksigner_path, "sign", "--in", args.output,
      "--cert", args.apk_key_path + ".x509.pem",
      "--key", args.apk_key_path + ".pk8"])
