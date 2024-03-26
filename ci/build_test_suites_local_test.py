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

"""Integration tests for build_test_suites that require a local build env."""

import ci_test_lib
import os
import pathlib
import shutil
import signal
import subprocess
import tempfile
import time


class BuildTestSuitesLocalTest(ci_test_lib.TestCase):

  KEEP_TEMP_DIR = False  # TODO: Enable this for debugging.

  def setUp(self):
    # We're *not* using `TemporaryDirectory` to temporarily disable cleaning up
    # the directory for debugging and the `delete` keyword was not added until
    # Python 3.12.
    self.temp_dir = pathlib.Path(tempfile.mkdtemp(prefix=self.id()))

    self.top_dir = pathlib.Path(os.environ['ANDROID_BUILD_TOP']).resolve()
    self.executable = self.top_dir.joinpath('build/make/ci/build_test_suites')

  def tearDown(self):
    if not self.KEEP_TEMP_DIR:
      shutil.rmtree(self.temp_dir, ignore_errors=True)

  def build_subprocess_args(self, build_args: list[str]):
    env = os.environ.copy()
    env['TOP'] = str(self.top_dir)

    args = [self.executable] + build_args,
    kwargs = {
      'cwd': self.top_dir,
      'env': env,
      'text': True,
    }

    return (args, kwargs)

  def run_build(self, build_args: list[str]) -> subprocess.CompletedProcess:
    args, kwargs = self.build_subprocess_args(build_args)

    return subprocess.run(
        *args,
        **kwargs,
        check=True,
        capture_output=True,
        timeout=5 * 60,
    )

  def test_fails_for_invalid_arg(self):
    invalid_arg = '--invalid-arg'

    with self.assertRaises(subprocess.CalledProcessError) as cm:
      self.run_build([invalid_arg])

    self.assertIn(invalid_arg, cm.exception.stderr)

  def test_builds_successfully(self):
    self.run_build(['nothing'])

  def test_interrupts_build(self):
    # TODO: Manage process.
    args, kwargs = self.build_subprocess_args(['general-tests'])
    p = subprocess.Popen(*args, **kwargs)
    time.sleep(5)  # Wait for the build to get going.
    self.assertIsNone(p.returncode)  # Check that the process is still alive.
    children = query_child_pids(p.pid)
    # TODO: Add assert utility.
    for c in children:
      self.assertTrue(ci_test_lib.process_alive(c))

    p.send_signal(signal.SIGINT)
    p.wait()

    time.sleep(5)  # Wait for things to die out.
    for c in children:
      self.assertFalse(ci_test_lib.process_alive(c))


# TODO(hzalek): Replace this with `psutils` once available in the tree.
def query_child_pids(parent_pid: int) -> set[int]:
  p = subprocess.run(['pgrep', '-P', str(parent_pid)], check=True, capture_output=True, text=True)
  return { int(pid) for pid in p.stdout.splitlines() }


if __name__ == '__main__':
  ci_test_lib.main()
