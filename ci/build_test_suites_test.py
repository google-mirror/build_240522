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

"""Tests for build_test_suites.py"""

import argparse
import importlib
import logging
import os
import pathlib
import subprocess
import sys
import tempfile
from typing import Any, Sequence, Text
import unittest
from unittest import mock

import build_test_suites
from pyfakefs import fake_filesystem_unittest


class BuildTestSuitesTest(fake_filesystem_unittest.TestCase):

  def setUp(self):
    self.setUpPyfakefs()

    # Disable logging since it breaks the TF Python test output parser.
    # TODO(hzalek): Use TF's `test-output-file` option to re-enable logging.
    logging.getLogger().disabled = True

    os_environ_patcher = mock.patch.dict('os.environ', {})
    self.addCleanup(os_environ_patcher.stop)
    self.mock_os_environ = os_environ_patcher.start()

    subprocess_run_patcher = mock.patch('subprocess.run')
    self.addCleanup(subprocess_run_patcher.stop)
    self.mock_subprocess_run = subprocess_run_patcher.start()

    self.setup_working_build_env()

  def test_missing_target_release_env_var_raises(self):
    del os.environ['TARGET_RELEASE']
    with self.assertRaisesRegex(Exception, 'TARGET_RELEASE'):
      build_test_suites.main([])

  def test_missing_target_product_env_var_raises(self):
    del os.environ['TARGET_PRODUCT']
    with self.assertRaisesRegex(Exception, 'TARGET_PRODUCT'):
      build_test_suites.main([])

  def test_missing_top_env_var_raises(self):
    del os.environ['TOP']
    with self.assertRaisesRegex(Exception, 'TOP'):
      build_test_suites.main([])

  def test_invalid_arg_raises(self):
    with self.assertRaisesRegex(SystemExit, '2'):
      build_test_suites.main(['--invalid_arg'])

  def test_build_failure_returns(self):
    self.mock_subprocess_run.side_effect = subprocess.CalledProcessError(42, None)
    with self.assertRaisesRegex(SystemExit, '42'):
      build_test_suites.main([])

  def test_build_success_returns(self):
    with self.assertRaisesRegex(SystemExit, '0'):
      build_test_suites.main([])

  def setup_working_build_env(self):
    self.fake_top = pathlib.Path('/fake/top')
    self.fake_top.mkdir(parents=True)

    self.soong_ui_dir = self.fake_top.resolve('build/soong')
    self.soong_ui_dir.mkdir(parents=True, exist_ok=True)

    self.soong_ui = self.soong_ui_dir.resolve('soong_ui.bash')
    self.soong_ui.touch()

    self.mock_os_environ.update({
        'TARGET_RELEASE': 'release',
        'TARGET_PRODUCT': 'product',
        'TOP': str(self.fake_top),
    })

    self.mock_subprocess_run.return_value = 0


class RunCommandIntegrationTest(unittest.TestCase):

  def test_streams_stdout(self):
      pass

  def test_streams_stderr(self):
      pass

  def test_propagates_interrupts(self):
      pass


# TODO(lucafarsi): Switch to getting the binary from resources and call it in a
# subprocess so we can do stuff like interrupt
class BuildTestSuitesIntegrationTest(unittest.TestCase):

  def setUp(self):
    self.temp_dir = tempfile.TemporaryDirectory()
    # Logging to stderr breaks the Python test parser.
    logging.getLogger().disabled = True

  def create_build_script(self, contents: list[str] = []):
    soong_dir = pathlib.Path(os.path.join(self.temp_dir.name, 'build', 'soong'))
    soong_dir.mkdir(parents=True)

    build_script = open(
        pathlib.Path(os.path.join(soong_dir, 'soong_ui.bash')), 'w'
    )
    build_script.write('#!/usr/bin/env bash\n')
    build_script.writelines(contents)

    build_script.close()
    os.chmod(pathlib.Path(os.path.join(soong_dir, 'soong_ui.bash')), 0o777)
    return build_script

  def get_valid_env(self):
    return {
        'TARGET_RELEASE': 'release',
        'TARGET_PRODUCT': 'product',
        'TOP': self.temp_dir.name,
    }

  def test_build_script_fails_returns(self):
    self.create_build_script(['exit 1'])
    with mock.patch.dict(os.environ, self.get_valid_env()):
      with self.assertRaisesRegex(SystemExit, '1') as context:
        build_test_suites.main([])

  def DISABLED_test_build_script_interrupt_success(self):
    self.create_build_script(['sleep 1'])
    # bts is an abbreviated package_path i had to use in the Android.bp so that
    # i can find the files in importlib.resources, probably needs a better title
    with importlib.resources.as_file(
        importlib.resources.files('bts').joinpath('build_test_suites')
    ) as build_script:
      with importlib.resources.as_file(
          importlib.resources.files('bts').joinpath(
              'build_test_suites.py.template'
          )
      ) as build_script_2:
        with open(
            os.path.join(self.temp_dir.name, 'build_test_suites'), 'w'
        ) as our_build_script:
          our_build_script.write(build_script.read_text(encoding='utf-8'))
          with open(
              os.path.join(self.temp_dir.name, 'build_test_suites.py'), 'w'
          ) as our_build_script_2:
            our_build_script_2.write(build_script_2.read_text(encoding='utf-8'))

    the_script = os.path.join(self.temp_dir.name, 'build_test_suites')
    os.chmod(os.path.join(self.temp_dir.name, 'build_test_suites'), 0o777)
    os.chmod(os.path.join(self.temp_dir.name, 'build_test_suites.py'), 0o777)
    proc = subprocess.run(args=[the_script], capture_output=True)
    print(proc.stdout)
    print(proc.stderr)


if __name__ == '__main__':
  unittest.main()
