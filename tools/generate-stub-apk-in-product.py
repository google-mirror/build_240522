#!/usr/bin/env python

import argparse
import hashlib
import os
import re
import subprocess
import sys
import tempfile

def ParseArgs(argv):
  parser = argparse.ArgumentParser(description='Create a Stub Apk')
  parser.add_argument('-v', '--verbose', help='verbose execution')
  parser.add_argument('--key', help='path to the private key file')
  parser.add_argument('--cert', help='path to the certificate file')
  parser.add_argument('--path', help='path to the src of stub apk')
  parser.add_argument('--outpath', help='path to stub apk')
  return parser.parse_args(argv)

def RunCommand(cmd, verbose=False, env=None):
  print("Running: " + " ".join(cmd))
  p = subprocess.Popen(
      cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, env=env)
  output, _ = p.communicate()

  if verbose or p.returncode is not 0:
    print(output.rstrip())

  assert p.returncode is 0, "Failed to execute: " + " ".join(cmd)
  return (output, p.returncode)

def PrepareAndroidManifest(package, sharedUserId):
  template = """\
<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
  package="{package}" android:sharedUserId="{sharedUserId}">
  <application android:hasCode="false" />
</manifest>
"""
  return template.format(package=package, sharedUserId=sharedUserId)

def CreateStubApk(args):
  src_manifest_path = os.path.join(args.path, "AndroidManifest.xml")
  if not os.path.exists(src_manifest_path):
    print("Manifest file '" + src_manifest_path + "' does not exist")
    return False
  shared_user_id = ReadManifest(src_manifest_path)
  package_name = None
  if shared_user_id is None:
    return False

  package_name = shared_user_id + ".stub"
  temp_dir = tempfile.mkdtemp()
  content_dir = os.path.join(temp_dir, "content")
  os.mkdir(content_dir)

  android_manifest_file = os.path.join(content_dir, 'AndroidManifest.xml')
  with open(android_manifest_file, 'w+') as f:
    f.write(PrepareAndroidManifest(package_name, shared_user_id))

  apk_file_name = package_name + ".apk"
  apk_dir = os.path.join(args.outpath, package_name)
  if not os.path.exists(apk_dir):
    os.mkdir(apk_dir)
    
  apk_file = os.path.join(apk_dir, apk_file_name)
  if os.path.exists(apk_file):
    return False

  apk_temp_file = os.path.join(temp_dir, "temp.apk")
  cmd = ['aapt2']
  cmd.append('link')
  cmd.extend(['--manifest', android_manifest_file])
  cmd.extend(['-o',  apk_temp_file])
  cmd.extend(['-I', "prebuilts/sdk/current/android.jar"])
  RunCommand(cmd, args.verbose)

  cmd = ['zipalign']
  cmd.append('-f')
  cmd.append('4096')
  cmd.append(apk_temp_file)
  cmd.append(apk_file)
  RunCommand(cmd, args.verbose)

  cmd = ['apksigner']
  cmd.append('sign')
  cmd.append('--v1-signing-enabled')
  cmd.append('false')
  cmd.append('--v2-signing-enabled')
  cmd.append('false')
  cmd.append('--key')
  cmd.append(args.key)
  cmd.append('--cert')
  cmd.append(args.cert)
  cmd.append('--lineage')
  cmd.append(args.cert.replace("x509.pem","lineage"))
  cmd.append(apk_file)
  RunCommand(cmd, args.verbose)

  return True

def ReadManifest(manifest_filename):
  shared_uid = None
  manifest_raw = None

  with open(manifest_filename, "r") as f:
    manifest_raw = f.read()

  for line in manifest_raw.split("\n"):
    line = line.strip()
    m = re.search(r'(\S*|\s*)android:sharedUserId(\s*=\s*)"(\S*)"', line)
    if m:
      shared_uid = m.group(3)
      if shared_uid is not None:
        break
  return shared_uid

def main(argv):
  args = ParseArgs(argv)
  success = CreateStubApk(args)

if __name__ == '__main__':
  main(sys.argv[1:])
