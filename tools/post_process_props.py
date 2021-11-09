#!/usr/bin/env python
#
# Copyright (C) 2009 The Android Open Source Project
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

import sys

# Usage: post_process_props.py file.prop [blacklist_key, ...]
# Blacklisted keys are removed from the property file, if present

# See PROP_VALUE_MAX in system_properties.h.
# The constant in system_properties.h includes the terminating NUL,
# so we decrease the value by 1 here.
PROP_VALUE_MAX = 91

# Put the modifications that you need to make into the /system/build.prop into this
# function. The prop object has get(name) and put(name,value) methods.
def mangle_build_prop(prop):
  pass

# Put the modifications that you need to make into /vendor/default.prop and
# /odm/default.prop into this function. The prop object has get(name) and
# put(name,value) methods.
def mangle_default_prop_override(prop):
  pass

# Put the modifications that you need to make into the /system/etc/prop.default into this
# function. The prop object has get(name) and put(name,value) methods.
def mangle_default_prop(prop):
  # If ro.debuggable is 1, then enable adb on USB by default
  # (this is for userdebug builds)
  if prop.get("ro.debuggable") == "1":
    val = prop.get("persist.sys.usb.config")
    if "adb" not in val:
      if val == "":
        val = "adb"
      else:
        val = val + ",adb"
      prop.put("persist.sys.usb.config", val)
  # UsbDeviceManager expects a value here.  If it doesn't get it, it will
  # default to "adb". That might not the right policy there, but it's better
  # to be explicit.
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  if not prop.get("persist.sys.usb.config"):
    prop.put("persist.sys.usb.config", "none");
=======
  if not prop_list.get_value("persist.sys.usb.config"):
    prop_list.put("persist.sys.usb.config", "none")

def validate_grf_props(prop_list, sdk_version):
  """Validate GRF properties if exist.

  If ro.board.first_api_level is defined, check if its value is valid for the
  sdk version.
  Also, validate the value of ro.board.api_level if defined.

  Returns:
    True if the GRF properties are valid.
  """
  grf_api_level = prop_list.get_value("ro.board.first_api_level")
  board_api_level = prop_list.get_value("ro.board.api_level")

  if not grf_api_level:
    if board_api_level:
      sys.stderr.write("error: non-GRF device must not define "
                       "ro.board.api_level\n")
      return False
    # non-GRF device skips the GRF validation test
    return True

  grf_api_level = int(grf_api_level)
  if grf_api_level > sdk_version:
    sys.stderr.write("error: ro.board.first_api_level(%d) must be less than "
                     "or equal to ro.build.version.sdk(%d)\n"
                     % (grf_api_level, sdk_version))
    return False

  if board_api_level:
    board_api_level = int(board_api_level)
    if board_api_level < grf_api_level or board_api_level > sdk_version:
      sys.stderr.write("error: ro.board.api_level(%d) must be neither less "
                       "than ro.board.first_api_level(%d) nor greater than "
                       "ro.build.version.sdk(%d)\n"
                       % (board_api_level, grf_api_level, sdk_version))
      return False

  return True
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

def validate(prop):
  """Validate the properties.

  Returns:
    True if nothing is wrong.
  """
  check_pass = True
  buildprops = prop.to_dict()
  for key, value in buildprops.iteritems():
    # Check build properties' length.
    if len(value) > PROP_VALUE_MAX and not key.startswith("ro."):
      check_pass = False
      sys.stderr.write("error: %s cannot exceed %d bytes: " %
                       (key, PROP_VALUE_MAX))
      sys.stderr.write("%s (%d)\n" % (value, len(value)))
  return check_pass

class PropFile:

  def __init__(self, lines):
    self.lines = [s.strip() for s in lines]

  def to_dict(self):
    props = {}
    for line in self.lines:
      if not line or line.startswith("#"):
        continue
      if "=" in line:
        key, value = line.split("=", 1)
        props[key] = value
    return props

  def get(self, name):
    key = name + "="
    for line in self.lines:
      if line.startswith(key):
        return line[len(key):]
    return ""

  def put(self, name, value):
    key = name + "="
    for i in range(0,len(self.lines)):
      if self.lines[i].startswith(key):
        self.lines[i] = key + value
        return
    self.lines.append(key + value)

  def delete(self, name):
    key = name + "="
    self.lines = [ line for line in self.lines if not line.startswith(key) ]

  def write(self, f):
    f.write("\n".join(self.lines))
    f.write("\n")

def main(argv):
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  filename = argv[1]
  f = open(filename)
  lines = f.readlines()
  f.close()
=======
  parser = argparse.ArgumentParser(description="Post-process build.prop file")
  parser.add_argument("--allow-dup", dest="allow_dup", action="store_true",
                      default=False)
  parser.add_argument("filename")
  parser.add_argument("disallowed_keys", metavar="KEY", type=str, nargs="*")
  parser.add_argument("--sdk-version", type=int, required=True)
  args = parser.parse_args()
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

  properties = PropFile(lines)

  if filename.endswith("/build.prop"):
    mangle_build_prop(properties)
  elif (filename.endswith("/vendor/default.prop") or
        filename.endswith("/odm/default.prop")):
    mangle_default_prop_override(properties)
  elif (filename.endswith("/default.prop") or # legacy
        filename.endswith("/prop.default")):
    mangle_default_prop(properties)
  else:
    sys.stderr.write("bad command line: " + str(argv) + "\n")
    sys.exit(1)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  if not validate(properties):
=======
  props = PropList(args.filename)
  mangle_build_prop(props)
  if not override_optional_props(props, args.allow_dup):
    sys.exit(1)
  if not validate_grf_props(props, args.sdk_version):
    sys.exit(1)
  if not validate(props):
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
    sys.exit(1)

  # Drop any blacklisted keys
  for key in argv[2:]:
    properties.delete(key)

  f = open(filename, 'w+')
  properties.write(f)
  f.close()

if __name__ == "__main__":
  main(sys.argv)
