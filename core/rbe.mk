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

# Notice: this works only with Google's RBE service.
ifneq ($(filter-out false,$(USE_RBE)),)
  ifdef RBE_DIR
    rbe_dir := $(RBE_DIR)
  else
    rbe_dir := $(HOME)/rbe
  endif
  RBE_WRAPPER := $(rbe_dir)/rewrapper
  RBE_CC := $(strip $(RBE_WRAPPER) --labels=type=compile,lang=cpp,compiler=clang --)

  # Append rewrapper to existing *_WRAPPER variables so it's possible to
  # use both ccache and rewrapper.
  CC_WRAPPER := $(strip $(CC_WRAPPER) $(RBE_CC))
  CXX_WRAPPER := $(strip $(CXX_WRAPPER) $(RBE_CC))

  rbe_dir :=
endif
