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

from importlib import resources
import logging
import multiprocessing
import os
import pathlib
import shutil
import signal
import stat
import subprocess
import sys
import tempfile
import textwrap
import time
from typing import Callable
import unittest
from unittest import mock
import build_test_suites
from pyfakefs import fake_filesystem_unittest
from timeout_decorator import timeout_decorator


class BuildTestSuitesTest(fake_filesystem_unittest.TestCase):

  def setUp(self):
    self.setUpPyfakefs()

    os_environ_patcher = mock.patch.dict('os.environ', {})
    self.addCleanup(os_environ_patcher.stop)
    self.mock_os_environ = os_environ_patcher.start()

    subprocess_run_patcher = mock.patch('subprocess.run')
    self.addCleanup(subprocess_run_patcher.stop)
    self.mock_subprocess_run = subprocess_run_patcher.start()

    self._setup_working_build_env()

  def test_missing_target_release_env_var_raises(self):
    del os.environ['TARGET_RELEASE']

    with self.assertRaisesRegex(build_test_suites.Error, 'TARGET_RELEASE'):
      build_test_suites.main([])

  def test_missing_target_product_env_var_raises(self):
    del os.environ['TARGET_PRODUCT']

    with self.assertRaisesRegex(build_test_suites.Error, 'TARGET_PRODUCT'):
      build_test_suites.main([])

  def test_missing_top_env_var_raises(self):
    del os.environ['TOP']

    with self.assertRaisesRegex(build_test_suites.Error, 'TOP'):
      build_test_suites.main([])

  def test_invalid_arg_raises(self):
    invalid_args = ['--invalid_arg']

    with self.assertRaisesRegex(SystemExit, '2'):
      build_test_suites.main(invalid_args)

  def test_build_failure_returns(self):
    self.mock_subprocess_run.side_effect = subprocess.CalledProcessError(
        42, None
    )

    with self.assertRaisesRegex(SystemExit, '42'):
      build_test_suites.main([])

  def test_build_success_returns(self):
    with self.assertRaisesRegex(SystemExit, '0'):
      build_test_suites.main([])

  def _setup_working_build_env(self):
    self.fake_top = pathlib.Path('/fake/top')
    self.fake_top.mkdir(parents=True)

    self.soong_ui_dir = self.fake_top.joinpath('build/soong')
    self.soong_ui_dir.mkdir(parents=True, exist_ok=True)

    self.soong_ui = self.soong_ui_dir.joinpath('soong_ui.bash')
    self.soong_ui.touch()

    self.mock_os_environ.update({
        'TARGET_RELEASE': 'release',
        'TARGET_PRODUCT': 'product',
        'TOP': str(self.fake_top),
    })

    self.mock_subprocess_run.return_value = 0


class RunCommandIntegrationTest(unittest.TestCase):

  DEFAULT_TEST_TIMEOUT_SECS = 5

  def setUp(self):
    # self.temp_dir_obj = tempfile.TemporaryDirectory(prefix=self.id())
    # self.temp_dir = pathlib.Path(self.temp_dir_obj.name)
    self.temp_dir = pathlib.Path(tempfile.mkdtemp(prefix=self.id()))

    # Copy the Python executable from 'non-code' resources and make it
    # executable for use by tests that launch a subprocess. Note that we don't
    # use Python's native `sys.executable` property since that is not set when
    # running via the embedded launcher.
    base_name = 'py3-cmd'
    dest_file = self.temp_dir.joinpath(base_name)
    with resources.as_file(
        resources.files('testdata').joinpath(base_name)
    ) as p:
      shutil.copy(p, dest_file)
    dest_file.chmod(dest_file.stat().st_mode | stat.S_IEXEC)
    self.python_executable = dest_file

  # TODO: Not working.
  @timeout_decorator.timeout(seconds=DEFAULT_TEST_TIMEOUT_SECS)
  def test_streams_stdout(self):
    marker = 'Spinach'
    stdout_file = self.temp_dir.joinpath('stdout.txt')

    def run(stdout_file):
      with open(stdout_file, 'w') as f:
        build_test_suites.run_command(
            [
                self.python_executable,
                '-c',
                textwrap.dedent(f"""\
                  import time

                  print('{marker}', end='', flush=True)

                  # Keep process alive until we check stdout.
                  time.sleep(10)
                  """),
            ],
            stdout=f,
        )

    p = multiprocessing.Process(target=run, args=[stdout_file])
    p.start()
    wait_until(p.is_alive)

    # TODO: Encoding and size.
    try:
      wait_until(
          lambda: stdout_file.is_file() and stdout_file.stat().st_size > 0
      )
      self.assertEquals(marker, read_file_contents(stdout_file))
    finally:
      os.kill(p.pid, signal.SIGINT)  # TODO: Kill.

  @timeout_decorator.timeout(seconds=DEFAULT_TEST_TIMEOUT_SECS)
  def test_propagates_interrupts(self):
    pid_file = self.temp_dir.joinpath('pid.txt')

    def run(pid_file):
      build_test_suites.run_command([
          self.python_executable,
          '-c',
          textwrap.dedent(f"""\
              import os
              import pathlib
              import time

              pathlib.Path('{pid_file}').write_text(str(os.getpid()))

              # Keep the process alive for us to explicitly interrupt it.
              time.sleep(10)
              """),
      ])

    p = multiprocessing.Process(target=run, args=[pid_file])
    p.start()
    wait_until(p.is_alive)

    wait_until(lambda: pid_file.is_file() and pid_file.stat().st_size > 0)

    process_pid = int(read_file_contents(pid_file))
    self.assertTrue(process_exists(process_pid))
    os.kill(p.pid, signal.SIGINT)
    # TODO: Assert 130.

    # TODO: Replace with assert.
    wait_until(lambda: not process_exists(process_pid))

  def assert_file_contains(self, file: pathlib.Path, substring: str):
    self.assertIn(substring, read_file_contents(f))


def process_exists(pid):
  """Check For the existence of a unix pid."""

  try:
    os.kill(pid, 0)
  except OSError:
    return False

  return True


def wait_until(
    condition_function: Callable[[], bool],
    timeout_secs: float = 3.0,
    polling_interval_secs: float = 0.1,
):
  """Waits until a condition function returns True."""

  start_time_secs = time.time()

  while not condition_function():
    if time.time() - start_time_secs > timeout_secs:
      raise TimeoutError(
          f'Condition not met within timeout: {timeout_secs} seconds'
      )

    time.sleep(polling_interval_secs)


def read_file_contents(file: pathlib.Path) -> str:
  with open(file, 'r') as f:
    return f.read()


if __name__ == '__main__':

  # Disable logging since it breaks the TF Python test output parser.
  # TODO(hzalek): Use TF's `test-output-file` option to re-enable logging.
  logging.getLogger().disabled = True

  unittest.main()
