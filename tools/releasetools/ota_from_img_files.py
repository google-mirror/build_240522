#!/usr/bin/env python
#
# Copyright (C) 2008 The Android Open Source Project
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

"""

"""
import logging
import multiprocessing
import os
import shutil
import sys
import tempfile
import zipfile

import common
import ota_from_target_files

OPTIONS = common.OPTIONS
OPTIONS.package_key = None
OPTIONS.incremental_source = None
OPTIONS.verify = False
OPTIONS.patch_threshold = 0.95
OPTIONS.wipe_user_data = False
OPTIONS.downgrade = False
OPTIONS.timestamp = False
OPTIONS.extra_script = None
OPTIONS.worker_threads = multiprocessing.cpu_count() // 2
if OPTIONS.worker_threads == 0:
  OPTIONS.worker_threads = 1
OPTIONS.two_step = False
OPTIONS.no_signing = False
OPTIONS.block_based = True
OPTIONS.updater_binary = None
OPTIONS.oem_source = None
OPTIONS.oem_no_mount = False
OPTIONS.fallback_to_full = True
OPTIONS.full_radio = False
OPTIONS.full_bootloader = False
# Stash size cannot exceed cache_size * threshold.
OPTIONS.cache_size = None
OPTIONS.stash_threshold = 0.8
OPTIONS.gen_verify = False
OPTIONS.log_diff = None
OPTIONS.payload_signer = None
OPTIONS.payload_signer_args = []
OPTIONS.tgt_info_dir = None
OPTIONS.key_passwords = []

files_location_dict = {
  "system.img" : "IMAGES/system.img",
  "system.map" : "IMAGES/system.map",
  "vendor.map" : "IMAGES/vendor.map",
  "SYSTEM/build.prop" : "SYSTEM/build.prop",
  "VENDOR/build.prop" : "VENDOR/build.prop",
  "misc_info.txt" : "META/misc_info.txt",
  "recovery.fstab" : "RECOVERY/RAMDISK/etc/recovery.fstab",
  "updater" : "OTA/bin/updater"}

def HasVendorPartition(tmp_dir):
  return os.path.exists(os.path.join(tmp_dir, "IMAGES/vendor.img"))

def setup_directory(info_dir, out_dir):
  os.mkdir(os.path.join(out_dir, "IMAGES"))
  for item in os.listdir(out_dir):
    src = os.path.join(out_dir, item)
    dst = os.path.join(out_dir, "IMAGES", item)
    if os.path.isfile(src):
      os.rename(src, dst)

  os.mkdir(os.path.join(out_dir, "META"))
  os.makedirs(os.path.join(out_dir, "RECOVERY/RAMDISK/etc/"))
  os.makedirs(os.path.join(out_dir, "OTA/bin"))
  os.mkdir(os.path.join(out_dir, "SYSTEM"))

  system_prop_path = os.path.join(info_dir, "SYSTEM/build.prop")
  assert os.path.exists(system_prop_path)
  shutil.copy(system_prop_path, os.path.join(out_dir, "SYSTEM/build.prop"))
  if HasVendorPartition(out_dir):
    os.mkdir(os.path.join(out_dir, "VENDOR"))
    vendor_prop_path = os.path.join(info_dir, "VENDOR/build.prop")
    assert os.path.exists(vendor_prop_path)
    shutil.copy(vendor_prop_path, os.path.join(out_dir, "VENDOR/build.prop"))

  for root, _, files in os.walk(info_dir):
    for name in files:
      src = os.path.join(root, name)
      if name in files_location_dict:
        dst = os.path.join(out_dir, files_location_dict[name])
        shutil.copy(src, dst)
      else:
        logging.debug("Found file " + src)


def check_files(prefix):
  has_vendor = HasVendorPartition(prefix)
  for value in files_location_dict.values():
    if not value.lower().startswith("vendor") or has_vendor:
      assert os.path.exists(os.path.join(prefix, value))

def WriteVerifyPackage(output_zip):
  has_vendor = HasVendorPartition(OPTIONS.input_tmp)
  updater_path = OPTIONS.updater_binary
  if updater_path is None:
    updater_path = os.path.join(OPTIONS.input_tmp, "OTA/bin/updater")
  assert os.path.exists(updater_path)
  ota_from_target_files.WriteVerifyPackage(output_zip, has_vendor,
                                           updater_binary=updater_path)

def WriteFullOTAPackage(output_zip):
  has_vendor = HasVendorPartition(OPTIONS.input_tmp)
  updater_path = OPTIONS.updater_binary
  if updater_path is None:
    updater_path = os.path.join(OPTIONS.input_tmp, "OTA/bin/updater")
  assert os.path.exists(updater_path)
  ota_from_target_files.WriteFullOTAPackage(output_zip, has_vendor,
                                            updater_binary=updater_path)

def WriteIncrementalOTAPackage(output_zip):
  has_vendor = HasVendorPartition(OPTIONS.target_tmp)
  if has_vendor and not HasVendorPartition(OPTIONS.source_tmp):
    raise RuntimeError("can't generate incremental that adds /vendor")

  updater_path = OPTIONS.updater_binary
  if updater_path is None:
    if OPTIONS.downgrade:
      updater_path = os.path.join(OPTIONS.source_tmp, "OTA/bin/updater")
    else:
      updater_path = os.path.join(OPTIONS.target_tmp, "OTA/bin/updater")
  assert os.path.exists(updater_path)

  ota_from_target_files.WriteBlockIncrementalOTAPackage(output_zip,
      has_vendor, updater_binary=updater_path)


def main(argv):
  tempfile.tempdir = "/usr/local/google/home/xunchang/1/tmp"
  def option_handler(o, a):
    if o in ("-k", "--package_key"):
      OPTIONS.package_key = a
    elif o in ("-i", "--incremental_from"):
      OPTIONS.incremental_source = a
    elif o == "--full_radio":
      OPTIONS.full_radio = True
    elif o == "--full_bootloader":
      OPTIONS.full_bootloader = True
    elif o in ("-w", "--wipe_user_data"):
      OPTIONS.wipe_user_data = True
    elif o == "--downgrade":
      OPTIONS.downgrade = True
      OPTIONS.wipe_user_data = True
    elif o == "--override_timestamp":
      OPTIONS.timestamp = True
    elif o in ("-o", "--oem_settings"):
      OPTIONS.oem_source = a.split(',')
    elif o == "--oem_no_mount":
      OPTIONS.oem_no_mount = True
    elif o in ("-e", "--extra_script"):
      OPTIONS.extra_script = a
    elif o in ("-t", "--worker_threads"):
      if a.isdigit():
        OPTIONS.worker_threads = int(a)
      else:
        raise ValueError("Cannot parse value %r for option %r - only "
                         "integers are allowed." % (a, o))
    elif o in ("-2", "--two_step"):
      OPTIONS.two_step = True
    elif o == "--no_signing":
      OPTIONS.no_signing = True
    elif o == "--verify":
      OPTIONS.verify = True
    elif o in ("-b", "--binary"):
      OPTIONS.updater_binary = a
    elif o in ("--no_fallback_to_full",):
      OPTIONS.fallback_to_full = False
    elif o == "--stash_threshold":
      try:
        OPTIONS.stash_threshold = float(a)
      except ValueError:
        raise ValueError("Cannot parse value %r for option %r - expecting "
                         "a float" % (a, o))
    elif o == "--gen_verify":
      OPTIONS.gen_verify = True
    elif o == "--log_diff":
      OPTIONS.log_diff = a
    elif o == "--tgt_info_dir":
      OPTIONS.tgt_info_dir = a
    elif o == "--src_info_dir":
      OPTIONS.src_info_dir = a
    else:
      return False
    return True

  args = common.ParseOptions(argv, __doc__,
                             extra_opts="b:k:i:d:we:t:2o:",
                             extra_long_opts=[
                                 "package_key=",
                                 "incremental_from=",
                                 "full_radio",
                                 "full_bootloader",
                                 "wipe_user_data",
                                 "downgrade",
                                 "override_timestamp",
                                 "extra_script=",
                                 "worker_threads=",
                                 "two_step",
                                 "no_signing",
                                 "block",
                                 "binary=",
                                 "oem_settings=",
                                 "oem_no_mount",
                                 "verify",
                                 "no_fallback_to_full",
                                 "stash_threshold=",
                                 "gen_verify",
                                 "log_diff=",
                                 "tgt_info_dir=",
                                 "src_info_dir="
                             ], extra_option_handler=option_handler)

  if len(args) != 2:
    common.Usage(__doc__)
    sys.exit(1)

  if OPTIONS.downgrade:
    # Sanity check to enforce a data wipe.
    if not OPTIONS.wipe_user_data:
      raise ValueError("Cannot downgrade without a data wipe")

    # We should only allow downgrading incrementals (as opposed to full).
    # Otherwise the device may go back from arbitrary build with this full
    # OTA package.
    if OPTIONS.incremental_source is None:
      raise ValueError("Cannot generate downgradable full OTAs")

  assert not (OPTIONS.downgrade and OPTIONS.timestamp), \
      "Cannot have --downgrade AND --override_timestamp both"

  assert OPTIONS.tgt_info_dir is not None
  OPTIONS.info_dict = common.LoadInfoDict(OPTIONS.tgt_info_dir,
                                          OPTIONS.tgt_info_dir)

  assert OPTIONS.info_dict.get("ab_update") != "true"

  # Use the default key to sign the package if not specified with package_key.
  if not OPTIONS.no_signing:
    if OPTIONS.package_key is None:
      OPTIONS.package_key = OPTIONS.info_dict.get(
          "default_system_dev_certificate",
          "build/target/product/security/testkey")
    # Get signing keys
    OPTIONS.key_passwords = common.GetKeyPasswords([OPTIONS.package_key])

  if OPTIONS.extra_script is not None:
    OPTIONS.extra_script = open(OPTIONS.extra_script).read()

  logging.info("unzipping target img-files...")
  OPTIONS.input_tmp, _ = common.UnzipTemp(args[0], ["*.img"])
  setup_directory(OPTIONS.tgt_info_dir, OPTIONS.input_tmp)
  OPTIONS.target_tmp = OPTIONS.input_tmp
  check_files(OPTIONS.target_tmp)

  if OPTIONS.verbose:
    logging.info("--- target info ---")
    common.DumpInfoDict(OPTIONS.info_dict)

  # If the caller explicitly specified the device-specific extensions
  # path via -s/--device_specific, use that.
  # Otherwise, take the path of the file from 'tool_extensions' in the
  # info dict and look for that in the local filesystem, relative to
  # the current directory.

  if OPTIONS.device_specific is None:
    from_input = os.path.join(OPTIONS.input_tmp, "META", "releasetools.py")
    if os.path.exists(from_input):
      logging.info("(using device-specific extensions from target_files)")
      OPTIONS.device_specific = from_input
    else:
      OPTIONS.device_specific = OPTIONS.info_dict.get("tool_extensions", None)

  if OPTIONS.device_specific is not None:
    OPTIONS.device_specific = None
  #  OPTIONS.device_specific = os.path.abspath(OPTIONS.device_specific)

  if OPTIONS.info_dict.get("no_recovery") == "true":
    raise common.ExternalError(
        "--- target build has specified no recovery ---")

  # Set up the output zip. Create a temporary zip file if signing is needed.
  if OPTIONS.no_signing:
    if os.path.exists(args[1]):
      os.unlink(args[1])
    output_zip = zipfile.ZipFile(args[1], "w",
                                 compression=zipfile.ZIP_DEFLATED)
  else:
    temp_zip_file = tempfile.NamedTemporaryFile()
    output_zip = zipfile.ZipFile(temp_zip_file, "w",
                                 compression=zipfile.ZIP_DEFLATED)

  # Non A/B OTAs rely on /cache partition to store temporary files.
  cache_size = OPTIONS.info_dict.get("cache_size", None)
  if cache_size is None:
    logging.info("--- can't determine the cache partition size ---")
  OPTIONS.cache_size = cache_size

  # Generate a verify package.
  if OPTIONS.gen_verify:
    WriteVerifyPackage(output_zip)

  # Generate a full OTA.
  elif OPTIONS.incremental_source is None:
    WriteFullOTAPackage(output_zip)

  # Generate an incremental OTA. It will fall back to generate a full OTA on
  # failure unless no_fallback_to_full is specified.
  else:
    logging.info("unzipping source target-files...")
    # TODO close the zip
    OPTIONS.source_tmp, _ = common.UnzipTemp(OPTIONS.incremental_source,
                                             ["*.img"])
    setup_directory(OPTIONS.src_info_dir, OPTIONS.source_tmp)
    check_files(OPTIONS.source_tmp)

    OPTIONS.target_info_dict = OPTIONS.info_dict
    OPTIONS.source_info_dict = common.LoadInfoDict(OPTIONS.src_info_dir,
                                                   OPTIONS.source_tmp)


    if OPTIONS.verbose:
      logging.info("--- source info ---")
      common.DumpInfoDict(OPTIONS.source_info_dict)
    try:
      WriteIncrementalOTAPackage(output_zip)
      if OPTIONS.log_diff:
        out_file = open(OPTIONS.log_diff, 'w')
        import target_files_diff
        target_files_diff.recursiveDiff('',
                                        OPTIONS.source_tmp,
                                        OPTIONS.input_tmp,
                                        out_file)
        out_file.close()
    except ValueError:
      if not OPTIONS.fallback_to_full:
        raise
      logging.info("--- failed to build incremental; falling back to full ---")
      OPTIONS.incremental_source = None
      WriteFullOTAPackage(output_zip)

  common.ZipClose(output_zip)

  # Sign the generated zip package unless no_signing is specified.
  if not OPTIONS.no_signing:
    ota_from_target_files.SignOutput(temp_zip_file.name, args[1])
    temp_zip_file.close()

  logging.info("done.")


if __name__ == '__main__':
  logging_format = '%(message)s'
  logging.basicConfig(level=logging.INFO, format=logging_format)
  try:
    common.CloseInheritedPipes()
    main(sys.argv[1:])
  except common.ExternalError as e:
    logging.error("\n   ERROR: {}\n".format(e))
    sys.exit(1)
  finally:
    common.Cleanup()
