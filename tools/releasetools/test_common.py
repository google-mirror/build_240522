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
import os
import tempfile
import time
import unittest
import zipfile

import common


def random_string_with_holes(size, block_size, step_size):
  data = ["\0"] * size
  for begin in range(0, size, step_size):
    end = begin + block_size
    data[begin:end] = os.urandom(block_size)
  return "".join(data)

def get_2gb_string():
  kilobytes = 1024
  megabytes = 1024 * kilobytes
  gigabytes = 1024 * megabytes

  size = int(2 * gigabytes + 1)
  block_size = 4 * kilobytes
  step_size = 4 * megabytes
  two_gb_string = random_string_with_holes(
        size, block_size, step_size)
  return two_gb_string


class CommonZipTest(unittest.TestCase):
  def _test_ZipWrite(self, contents, extra_zipwrite_args=None):
    extra_zipwrite_args = dict(extra_zipwrite_args or {})

    test_file = tempfile.NamedTemporaryFile(delete=False)
    zip_file = tempfile.NamedTemporaryFile(delete=False)

    test_file_name = test_file.name
    zip_file_name = zip_file.name

    # File names within an archive strip the leading slash.
    arcname = extra_zipwrite_args.get("arcname", test_file_name)
    if arcname[0] == "/":
      arcname = arcname[1:]

    zip_file.close()
    zip_file = zipfile.ZipFile(zip_file_name, "w")

    try:
      test_file.write(contents)
      test_file.close()

      old_stat = os.stat(test_file_name)
      expected_mode = extra_zipwrite_args.get("perms", 0o644)

      time.sleep(5)  # Make sure the atime/mtime will change measurably.

      common.ZipWrite(zip_file, test_file_name, **extra_zipwrite_args)

      # Verifications.
      new_stat = os.stat(test_file_name)
      self.assertEqual(int(old_stat.st_mode), int(new_stat.st_mode))
      self.assertEqual(int(old_stat.st_mtime), int(new_stat.st_mtime))
      self.assertIsNone(zip_file.testzip())

      common.ZipClose(zip_file)
      zip_file = zipfile.ZipFile(zip_file_name, "r")
      info = zip_file.getinfo(arcname)

      self.assertEqual(info.date_time, (2009, 1, 1, 0, 0, 0))
      mode = (info.external_attr >> 16) & 0o777
      self.assertEqual(mode, expected_mode)
      self.assertEqual(zip_file.read(arcname), contents)
      self.assertIsNone(zip_file.testzip())
    finally:
      os.remove(test_file_name)
      os.remove(zip_file_name)

  def _test_ZipWriteStr(self, zinfo_or_arcname, contents, extra_args=None):
    extra_args = dict(extra_args or {})

    zip_file = tempfile.NamedTemporaryFile(delete=False)
    zip_file_name = zip_file.name
    zip_file.close()

    zip_file = zipfile.ZipFile(zip_file_name, "w")

    try:
      expected_mode = extra_args.get("perms", 0o644)
      time.sleep(5)  # Make sure the atime/mtime will change measurably.

      if not isinstance(zinfo_or_arcname, zipfile.ZipInfo):
        zinfo = zipfile.ZipInfo(filename=zinfo_or_arcname)
      else:
        zinfo = zinfo_or_arcname
      arcname = zinfo.filename

      common.ZipWriteStr(zip_file, zinfo, contents, **extra_args)

      # Verifications.
      self.assertIsNone(zip_file.testzip())
      common.ZipClose(zip_file)

      zip_file = zipfile.ZipFile(zip_file_name, "r")
      info = zip_file.getinfo(arcname)

      self.assertEqual(info.date_time, (2009, 1, 1, 0, 0, 0))
      mode = (info.external_attr >> 16) & 0o777
      self.assertEqual(mode, expected_mode)
      self.assertEqual(zip_file.read(arcname), contents)
      self.assertIsNone(zip_file.testzip())
    finally:
      os.remove(zip_file_name)

  def _test_ZipWriteStr_large_file(self, large, small, extra_args=None):
    extra_args = dict(extra_args or {})

    zip_file = tempfile.NamedTemporaryFile(delete=False)
    zip_file_name = zip_file.name

    test_file = tempfile.NamedTemporaryFile(delete=False)
    test_file_name = test_file.name

    arcname_large = test_file_name
    arcname_small = "bar"

    # File names within an archive strip the leading slash.
    if arcname_large[0] == "/":
      arcname_large = arcname_large[1:]

    zip_file.close()
    zip_file = zipfile.ZipFile(zip_file_name, "w")

    try:
      test_file.write(large)
      test_file.close()

      common.ZipWrite(zip_file, test_file_name, **extra_args)
      common.ZipWriteStr(zip_file, arcname_small, small, **extra_args)

      self.assertIsNone(zip_file.testzip())
      common.ZipClose(zip_file)

      expected_mode = 0o644

      # Check the contents written by ZipWrite()
      zip_file = zipfile.ZipFile(zip_file_name, "r")
      info = zip_file.getinfo(arcname_large)

      self.assertEqual(info.date_time, (2009, 1, 1, 0, 0, 0))
      mode = (info.external_attr >> 16) & 0o777
      self.assertEqual(mode, expected_mode)
      self.assertEqual(zip_file.read(arcname_large), large)
      self.assertIsNone(zip_file.testzip())

      # Check the contents written by ZipWriteStr()
      info = zip_file.getinfo(arcname_small)

      self.assertEqual(info.date_time, (2009, 1, 1, 0, 0, 0))
      mode = (info.external_attr >> 16) & 0o777
      self.assertEqual(mode, expected_mode)
      self.assertEqual(zip_file.read(arcname_small), small)
      self.assertIsNone(zip_file.testzip())
    finally:
      os.remove(zip_file_name)
      os.remove(test_file_name)

  def _test_reset_ZIP64_LIMIT(self, func, *args):
    default_limit = (1 << 31) - 1
    self.assertEqual(default_limit, zipfile.ZIP64_LIMIT)
    func(*args)
    self.assertEqual(default_limit, zipfile.ZIP64_LIMIT)

  def test_ZipWrite(self):
    file_contents = os.urandom(1024)
    self._test_ZipWrite(file_contents)

  def test_ZipWrite_with_opts(self):
    file_contents = os.urandom(1024)
    self._test_ZipWrite(file_contents, {
        "arcname": "foobar",
        "perms": 0o777,
        "compress_type": zipfile.ZIP_DEFLATED,
    })
    self._test_ZipWrite(file_contents, {
        "arcname": "foobar",
        "perms": 0o700,
        "compress_type": zipfile.ZIP_STORED,
    })

  def test_ZipWrite_large_file(self):
    file_contents = get_2gb_string()
    self._test_ZipWrite(file_contents, {
        "compress_type": zipfile.ZIP_DEFLATED,
    })

  def test_ZipWrite_resets_ZIP64_LIMIT(self):
    self._test_reset_ZIP64_LIMIT(self._test_ZipWrite, "")

  def test_ZipWriteStr(self):
    random_string = os.urandom(1024)
    # Passing arcname
    self._test_ZipWriteStr("foo", random_string)

    # Passing zinfo
    zinfo = zipfile.ZipInfo(filename="foo")
    self._test_ZipWriteStr(zinfo, random_string)

    # Timestamp in the zinfo should be overwritten.
    zinfo.date_time = (2015, 3, 1, 15, 30, 0)
    self._test_ZipWriteStr(zinfo, random_string)

  def test_ZipWriteStr_with_opts(self):
    random_string = os.urandom(1024)
    # Passing arcname
    self._test_ZipWriteStr("foo", random_string, {
        "perms": 0o777,
        "compress_type": zipfile.ZIP_DEFLATED,
    })
    self._test_ZipWriteStr("foo", random_string, {
        "perms": 0o777,
        "compress_type": zipfile.ZIP_STORED,
    })

    # Passing zinfo
    zinfo = zipfile.ZipInfo(filename="foo")
    self._test_ZipWriteStr(zinfo, random_string, {
        "perms": 0o755,
        "compress_type": zipfile.ZIP_DEFLATED,
    })
    self._test_ZipWriteStr(zinfo, random_string, {
        "perms": 0o755,
        "compress_type": zipfile.ZIP_STORED,
    })

  def test_ZipWriteStr_large_file(self):
    # zipfile.writestr() doesn't work when the str size is over 2GiB even with
    # the workaround. We will only test the case of writing a string into a
    # large archive.
    long_string = get_2gb_string()
    short_string = os.urandom(1024)
    self._test_ZipWriteStr_large_file(long_string, short_string, {
        "compress_type": zipfile.ZIP_DEFLATED,
    })

  def test_ZipWriteStr_resets_ZIP64_LIMIT(self):
    self._test_reset_ZIP64_LIMIT(self._test_ZipWriteStr, "foo", "")
    zinfo = zipfile.ZipInfo(filename="foo")
    self._test_reset_ZIP64_LIMIT(self._test_ZipWriteStr, zinfo, "")
