#
# Copyright (C) 2017 The Android Open Source Project
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

###########################################################
# Basic math functions for positive integers <= 100
#
# (SDK versions for example)
###########################################################
__MATH_NUMBERS :=  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 \
                  21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 \
                  41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 \
                  61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 \
                  81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100

# Returns true if $(1) is a positive integer <= 100, otherwise returns nothing.
define math_is_number
$(strip \
  $(if $(1),,$(error Argument missing)) \
  $(if $(word 2,$(1)),$(error Multiple words in a single argument: $(1))) \
  $(if $(filter $(1),$(__MATH_NUMBERS)),true))
endef

#$(warning true == $(call math_is_number,2))
#$(warning == $(call math_is_number,foo))
#$(call math_is_number,1 2)
#$(call math_is_number,no 2)

define _math_check_valid
$(if $(call math_is_number,$(1)),,$(error Only positive integers <= 100 are supported (not $(1))))
endef

#$(call _math_check_valid,0)
#$(call _math_check_valid,1)
#$(call _math_check_valid,100)
#$(call _math_check_valid,101)
#$(call _math_check_valid,)
#$(call _math_check_valid,1 2)

# Returns the greater of $1 or $2.
# If $1 or $2 is not a positive integer <= 100, then an error is generated.
define math_max
$(strip $(call _math_check_valid,$(1)) $(call _math_check_valid,$(2)) \
  $(lastword $(filter $(1) $(2),$(__MATH_NUMBERS))))
endef

#$(call math_max)
#$(call math_max,1)
#$(call math_max,1 2,3)
#$(warning 1 == $(call math_max,1,1))
#$(warning 42 == $(call math_max,5,42))
#$(warning 42 == $(call math_max,42,5))

define math_gt_or_eq
$(if $(filter $(1),$(call math_max,$(1),$(2))),true)
endef

#$(warning $(call math_gt_or_eq, 2, 1))
#$(warning $(call math_gt_or_eq, 1, 1))
#$(warning $(if $(call math_gt_or_eq, 1, 2),false,true))

# $1 is the variable name to increment
define inc_and_print
$(strip $(eval $(1) := $($(1)) .)$(words $($(1))))
endef

# 0 < $1 <= 8000
define int_check_valid
$(if $(call _int_greater-or-equal,$(call int_encode,8000),$(call int_encode,$(1))),\
  $(if $(call _int_greater-than,$(call int_encode,$(1)),$(call int_encode,0)),,$(call pretty-error,Integer < 0 are not supported!)),\
    $(call pretty-error,Integers > 8000 are not supported!))
endef

_INT_MAX := $(foreach a,x x,$(foreach b,x x x x x x x x x x x x x x x x,\
  $(foreach c,x x x x x x x x x x x x x x x x,x x x x x x x x x x x x x x x x)))

# Max output(integer upper limit) consists of 2^13 (8192) x's built from 16 x's
# 0 < $1 <= 8000 (int_check_valid has to be invoked before calling int_encode)
define int_encode
$(wordlist 1,$(1),$(_INT_MAX))
endef

# _int_max returns the maximum of the two arguments
# input: two (x) lists; output: one (x) list
# integer cannot be passed in directly. It has to be converted using int_encode.
define _int_max
$(subst xx,x,$(join $(1),$(2)))
endef

# first argument is greater than second argument
# output: non-empty if true
# integer cannot be passed in directly. It has to be converted using int_encode.
define _int_greater-than
$(filter-out $(words $(2)),$(words $(call _int_max,$(1),$(2))))
endef

# first argument equals to second argument
# output: non-empty if true
# integer cannot be passed in directly. It has to be converted using int_encode.
define _int_equal
$(filter $(words $(1)),$(words $(2)))
endef

# first argument is greater than or equal to second argument
# output: non-empty if true
# integer cannot be passed in directly. It has to be converted using int_encode.
define _int_greater-or-equal
$(call _int_greater-than,$(1),$(2))$(call _int_equal,$(1),$(2))
endef

# input: two (x) lists; output: one (x) list
# integer cannot be passed in directly. It has to be converted using int_encode.
define int_plus
$(1) $(2)
endef

# first argument subtract second argument if first one is greater or equals to second one.
# input: two (x) lists; output: one (x) list or ERROR
# integer cannot be passed in directly. It has to be converted using int_encode.
define int_subtract
$(if $(call _int_greater-or-equal,$(1),$(2)),$(filter-out xx,$(join $(1),$(2))),$(call pretty-error,$(1) subtract underflow $(2)))
endef

# input: two (x) lists; output: one (x) list
# integer cannot be passed in directly. It has to be converted using int_encode.
define int_multiply
$(foreach a,$(1),$(2))
endef

# caller is responsible that second argument is greater than 0
# input: two (x) lists; output: one (x) list or ERROR
# integer cannot be passed in directly. It has to be converted using int_encode.
define int_divide
$(if $(call _int_greater-or-equal,$(1),$(2)),x $(call int_divide,$(call int_subtract,$(1),$(2)),$(2)),)
endef
