<<<<<<< HEAD   (159713 Merge "Merge empty history for sparse-9081464-L8140000095646)
=======
#!/bin/bash -e
# Copyright (C) 2022 The Android Open Source Project
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

exit 0 #TODO(b/244771600) remove this after CI is enabled

tests=(
 $(dirname $0)/b_tests.sh
)

for test in $tests; do
  bash -x $test
  zsh -x $test
done
>>>>>>> BRANCH (5235f6 Merge "Version bump to TKB1.220921.001.A1 [core/build_id.mk])
