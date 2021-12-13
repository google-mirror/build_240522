/*
 * Copyright (C) 2021 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "android-base/file.h"

namespace zipalign {
std::string GetExecutablePath() {
#if defined(__linux__)
  // Bazel test runner will set TEST_SRCDIR. This is a clue that test is run
  // under Bazel.
  std::string test_srcdir = getenv("TEST_SRCDIR");
  if (test_srcdir.empty()) {
    return android::base::GetExecutablePath();
  } else {
    // When test is run under Bazel, use the real execution path
    // (program_invocation_name is initialized from argv[0]) instead of
    // following symlink to find the executable path (which could live outside
    // of Bazel sandbox).
    std::string path = program_invocation_name;
    return path;
  }
#else
  return android::base::GetExecutablePath();
#endif
}
} // namespace zipalign
