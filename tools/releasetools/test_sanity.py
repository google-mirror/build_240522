#!/usr/bin/env python
#
# Copyright (C) 2015 The Android Open Source Project
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
#

"""
Usage:
line1
line2
line3
"""

from __future__ import print_function

import argparse
import datetime
import hashlib
import inspect
import json
import os
import random
import string
import subprocess
import sys
import tempfile
#import time
import unittest
import zipfile

#import ota_from_target_files
import add_img_to_target_files
import common

# Discard outputs by redirecting them to DEVNULL.
DEVNULL = open(os.devnull, 'wb')

# The file that contains the testcases.
TEST_DATA_FILE = 'build/tools/releasetools/checksums.json'
TEST_FILENAME_LEN_MAX = 15
TEST_FILE_LEN_MAX = 4096

update = False
reserve = False

top = os.getenv('ANDROID_BUILD_TOP')
if top is None:
  print('ANDROID_BUILD_TOP not set.\n', file=sys.stderr)
  sys.exit(1)
else:
  os.chdir(top)


def hash_file(source):
  ctx = hashlib.sha1()
  blocksize = 4096
  with open(source, 'rb') as f:
    buf = f.read(blocksize)
    while len(buf) > 0:
      ctx.update(buf)
      buf = f.read(blocksize)
  return ctx.hexdigest()


def dump_test(testcase):
  print('cmd:\t%s' % (testcase['cmd'],))
  print('args:\t%s' % (testcase['args'],))
  print('hash:\t%s' % (testcase['checksum']))
  print('-' * 80)


def getid_generator():
  """Get a generator that returns an identifier.

  It returns an generator that gives an identifier according to caller's class
  and method names with a suffix number.
  """
  stack = inspect.stack()
  parent = stack[1][0]
  class_name = parent.f_locals['self'].__class__.__name__
  method_name = stack[1][3]
  prefix = class_name + ':' + method_name

  i = 0
  while True:
    i += 1
    yield prefix + ':' + str(i)

target_files_gen = [('SYSTEM/', 0),
                    ('SYSTEM/app/', 10),
                    ('SYSTEM/app/app1', 1),
                    ('SYSTEM/app/app2', 1),
                    ('SYSTEM/app/app3', 1),
                    ('SYSTEM/bin/', 5),
                    ('SYSTEM/etc/', 8),
                    ('SYSTEM/fonts/', 3),
                    ('SYSTEM/framework/', 20),
                    ('SYSTEM/lib/', 3),
                    ('SYSTEM/media/', 2),
                    ('SYSTEM/priv-app/', 5),
                    ('SYSTEM/tts/', 3),
                    ('SYSTEM/usr/', 6),
                    ('SYSTEM/usr/prog1', 1),
                    ('SYSTEM/usr/prog2', 1),
                    ('SYSTEM/xbin/', 7),
                    ('BOOT/', 0),
                    ('BOOT/RAMDISK/', 10),
                    ('BOOT/kernel', 1),
                    ('RECOVERY/', 0),
                    ('RECOVERY/RAMDISK/', 10),
                    ('RECOVERY/kernel', 1),
                    # ('RECOVERY/RAMDISK/etc/', 0),
                    # ('RECOVERY/RAMDISK/etc/recovery.fstab', 15),
                    ('OTA/', 0),
                    ('OTA/bin/updater', 1)]

target_files_stub = [('META/misc_info.txt'),
                     ('SYSTEM/build.prop')]

def generate_target_files(add_images=False, seed=2015):
  def random_data(size):
    len = random.randint(1, size)
    return ''.join(random.choice(string.ascii_letters+string.digits) for _ in range(len))

  def set_timestamp(filename):
    # Use a fixed timestamp so the output is repeatable.
    epoch = datetime.datetime.fromtimestamp(0)
    timestamp = (datetime.datetime(2015, 4, 1) - epoch).total_seconds()
    os.utime(filename, (timestamp, timestamp))

  random.seed(seed)
  target_files_temp = tempfile.NamedTemporaryFile(delete=False, suffix='.zip')
  # target_files_name = 'fake1_target-files.zip'
  target_files_name = target_files_temp.name
  archive = zipfile.ZipFile(target_files_temp, 'w', compression=zipfile.ZIP_DEFLATED)

  for (path, num) in target_files_gen:
    # Directory
    if path.endswith('/'):
      zip_dir = zipfile.ZipInfo(path)
      zip_dir.external_attr = 0o40750 << 16
      zip_dir.date_time = (2015, 4, 1, 0, 0, 0)
      archive.writestr(zip_dir, '')

      if num > 0:
        for i in range(random.randint(1, num)):
          filename = os.path.join(path, random_data(TEST_FILENAME_LEN_MAX))
          temp_file = tempfile.NamedTemporaryFile(delete=False)
          temp_file.write(random_data(TEST_FILE_LEN_MAX))
          temp_file.close()
          set_timestamp(temp_file.name)
          archive.write(temp_file.name, arcname=filename)
          os.remove(temp_file.name)

    # File
    else:
      filename = path
      temp_file = tempfile.NamedTemporaryFile(delete=False)
      temp_file.write(random_data(TEST_FILE_LEN_MAX))
      temp_file.close()
      set_timestamp(temp_file.name)
      archive.write(temp_file.name, arcname=filename)
      os.remove(temp_file.name)

  # TODO: handle the prefix
  stub_prefix = 'build/tools/releasetools/testdata'
  for filename in target_files_stub:
    full_name = os.path.join(top, stub_prefix, filename)
    set_timestamp(full_name)
    archive.write(full_name, arcname=filename)

  # We need to close the file to get the full zip.
  archive.close()
  target_files_temp.close()

  # Generate META/filesystem_config.txt.
  cmdline = ['zipinfo', '-1', target_files_name]
  zipinfo = subprocess.Popen(cmdline, stdout=subprocess.PIPE)

  # import shlex
  # cmdline = 'awk \'BEGIN { FS="SYSTEM/" } /^SYSTEM\// {print "system/" $2}\''
  # cmdline = shlex.split(cmdline)
  cmdline = ['awk', 'BEGIN { FS="SYSTEM/" } /^SYSTEM\// {print "system/" $2}']
  # print(cmdline)
  zipinfo2 = subprocess.Popen(cmdline, stdin=zipinfo.stdout, stdout=subprocess.PIPE)

  cmdline = [os.path.join(top, 'out/host/linux-x86/bin', 'fs_config'), '-C', '-S', 'build/tools/releasetools/file_contexts']
  output = subprocess.check_output(cmdline, stdin=zipinfo2.stdout)
  zipinfo.wait()
  assert zipinfo.returncode == 0, 'command failed'

  temp_file = tempfile.NamedTemporaryFile(delete=False)
  temp_file.write(output)
  temp_file.close()
  set_timestamp(temp_file.name)

  archive = zipfile.ZipFile(target_files_name, 'a', compression=zipfile.ZIP_DEFLATED)
  archive.write(temp_file.name, arcname='META/filesystem_config.txt')
  archive.close()

  os.remove(temp_file.name)

  if add_images:
    add_img_to_target_files.main([target_files_name])

  return target_files_name


class OtaScriptTest(unittest.TestCase):

  def __init__(self, *args, **kwargs):
    super(OtaScriptTest, self).__init__(*args, **kwargs)

  @classmethod
  def setUpClass(cls, add_images, need_src_tf=False):
    print('Hello from OtaScriptTest')
    # Generate testdata
    cls.target_files = generate_target_files(add_images=add_images, seed=2015)
    if need_src_tf:
      cls.src_tf = generate_target_files(add_images=add_images, seed=2016)
    else:
      cls.src_tf = None

    with open(TEST_DATA_FILE, "r") as json_file:
      cls.checksums = json.load(json_file)

  @classmethod
  def tearDownClass(cls):
    # Clean up
    if not reserve:
      os.remove(cls.target_files)
      if cls.src_tf is not None:
        os.remove(cls.src_tf)
    else:
      print('Target TF is reserved at \'%s\'.' % (cls.target_files))
      if cls.src_tf is not None:
        print('Source TF is reserved at \'%s\'.' % (cls.src_tf))

    if update:
      with open(TEST_DATA_FILE, "w") as json_file:
        json.dump(cls.checksums, json_file, sort_keys=True, indent=4,
                  separators=(',', ':'))

  def _compare(self, test_id, checksum):
    if update:
      self.checksums[test_id] = checksum
    else:
      expected = self.checksums[test_id]
      self.assertEqual(checksum, expected)

  def _run_test(self, cmdline):
    outfile_temp = tempfile.NamedTemporaryFile()
    outfile_name = outfile_temp.name
    cmdline = cmdline + [outfile_name]
    p = subprocess.Popen(cmdline, cwd=top, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = p.communicate()
    p.wait()
    assert p.returncode == 0, 'Command \'%s\' failed.\nout: %s\nerr: %s' % (
        ' '.join(cmdline), out, err)

    checksum = hash_file(outfile_name)
    outfile_temp.close()
    return checksum


class OtaFromTargetFilesTest(OtaScriptTest):
  """Test the functions in ota_from_target_files.py"""

  def __init__(self, *args, **kwargs):
    super(OtaFromTargetFilesTest, self).__init__(*args, **kwargs)

  @classmethod
  def setUpClass(cls):
    print('Hello from OtaFromTargetFilesTest')
    super(OtaFromTargetFilesTest, cls).setUpClass(add_images=True, need_src_tf=True)

  def test_WriteFullOTAPackage(self):
    test_id_gen = getid_generator()
    # Test generating full OTA
    checksum = self._run_test(
        ['build/tools/releasetools/ota_from_target_files.py',
         self.target_files])
    self._compare(next(test_id_gen), checksum)

    # Option --no_prereq: don't check build time
    checksum = self._run_test(
        ['build/tools/releasetools/ota_from_target_files.py', '--no_prereq',
         self.target_files])
    self._compare(next(test_id_gen), checksum)

    checksum = self._run_test(
        ['build/tools/releasetools/ota_from_target_files.py', '-n',
         self.target_files])
    self._compare(next(test_id_gen), checksum)

    # File-based OTA
    # Block-based OTA

  def test_Integration(self):
    test_id_gen = getid_generator()
    # Test generating full OTA
    checksum = self._run_test(
        ['build/tools/releasetools/ota_from_target_files.py',
         self.target_files])
    self._compare(next(test_id_gen), checksum)

    # Test generating file-based incremental OTA
    # LRX22L -> LMY47P userdebug
    checksum = self._run_test(
        ['build/tools/releasetools/ota_from_target_files.py',
         '-k', 'build/target/product/security/testkey',
         '-i',
         'build/tools/releasetools/testdata/volantis-userdebug-target_files-1816899.zip',
         'build/tools/releasetools/testdata/volantis-userdebug-target_files-1792736.zip'])
    self._compare(next(test_id_gen), checksum)

    # Test generate block-based incremental OTA
    # LRX22L -> LMY47P userdebug, but with test keys
    checksum = self._run_test(
        ['build/tools/releasetools/ota_from_target_files.py',
         '-k', 'build/target/product/security/testkey',
         '--block', '-i',
         'build/tools/releasetools/testdata/volantis-userdebug-target_files-1816899.zip',
         'build/tools/releasetools/testdata/volantis-userdebug-target_files-1792736.zip'])
    self._compare(next(test_id_gen), checksum)


class AddImgToTargetFilesTest(OtaScriptTest):

  def __init__(self, *args, **kwargs):
    super(AddImgToTargetFilesTest, self).__init__(*args, **kwargs)

  @classmethod
  def setUpClass(cls):
    print('Hello from AddImgToTargetFilesTest')
    super(AddImgToTargetFilesTest, cls).setUpClass(add_images=False)

  def test_AddSystem(self):
    test_id_gen = getid_generator()

    OPTIONS = common.OPTIONS
    OPTIONS.input_tmp, input_zip = common.UnzipTemp(self.target_files)
    OPTIONS.info_dict = common.LoadInfoDict(input_zip)
    input_zip.close()

    # AddSystem
    output_zip = zipfile.ZipFile(self.target_files, 'a', compression=zipfile.ZIP_DEFLATED)
    add_img_to_target_files.AddSystem(output_zip)
    output_zip.close()
    checksum = hash_file(self.target_files)
    self._compare(next(test_id_gen), checksum)

    # Generate again, which should be a no-op since it already exists.
    OPTIONS.input_tmp, input_zip = common.UnzipTemp(self.target_files)
    OPTIONS.info_dict = common.LoadInfoDict(input_zip)
    input_zip.close()

    output_zip = zipfile.ZipFile(self.target_files, 'a', compression=zipfile.ZIP_DEFLATED)
    add_img_to_target_files.AddSystem(output_zip)
    output_zip.close()
    checksum = hash_file(self.target_files)
    self._compare(next(test_id_gen), checksum)

    # AddSystem and rebuild the recovery patch file
    """
    OPTIONS.rebuild_recovery = True
    output_zip_temp = tempfile.NamedTemporaryFile(suffix='.zip', delete=False)
    output_zip = zipfile.ZipFile(output_zip_temp.name, 'w', compression=zipfile.ZIP_DEFLATED)
    add_img_to_target_files.AddSystem(output_zip)
    checksum = hash_file(output_zip_temp.name)
    output_zip.close()
    self._compare(next(test_id_gen), checksum)
    """

    # AddCache
    OPTIONS.input_tmp, input_zip = common.UnzipTemp(self.target_files)
    OPTIONS.info_dict = common.LoadInfoDict(input_zip)
    input_zip.close()

    output_zip = zipfile.ZipFile(self.target_files, 'a', compression=zipfile.ZIP_DEFLATED)
    add_img_to_target_files.AddCache(output_zip)
    output_zip.close()
    checksum = hash_file(self.target_files)
    self._compare(next(test_id_gen), checksum)


class ImgFromTargetFilesTest(OtaScriptTest):

  def __init__(self, *args, **kwargs):
    super(ImgFromTargetFilesTest, self).__init__(*args, **kwargs)

  def _test3(self):
    test_id_gen = getid_generator()

    # Test generate block-based incremental OTA
    checksum = self._run_test(
        ['build/tools/releasetools/img_from_target_files.py',
         'build/tools/releasetools/testdata/aosp_hammerhead-target_files-eng.tbao.zip'])
    self._compare(next(test_id_gen), checksum)


def main(argv):
  parser = argparse.ArgumentParser(description='Parse')
  parser.add_argument('-u', '--update', help='Update the checksums', action='store_true')
  parser.add_argument('-e', '--reserve', help='Reserve the testdata', action='store_true')
  # parser.add_argument('-d', '--dump', help='dump the testcases', action='store_true')
  parser.add_argument('unittest_args', nargs='*')
  args = parser.parse_args(argv)
  global update
  update = args.update
  global reserve
  reserve = args.reserve

  print('unittest_args: [%s]' % args.unittest_args)

  # unittest.main(argv=['./test_sanity.py'] + args.unittest_args)
  unittest.main(argv=args.unittest_args)


if __name__ == '__main__':
  main(sys.argv)
