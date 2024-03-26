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

"""Testing utilities for tests in the CI package."""

import logging
import os
import unittest


# Export the TestCase class to reduce the number of imports tests have to list.
TestCase = unittest.TestCase


def process_alive(pid):
  """Check For the existence of a pid."""

  try:
    os.kill(pid, 0)
  except OSError:
    return False

  return True


def main():

  # Disable logging since it breaks the TF Python test output parser.
  # TODO(hzalek): Use TF's `test-output-file` option to re-enable logging.
  logging.getLogger().disabled = True

  unittest.main()
