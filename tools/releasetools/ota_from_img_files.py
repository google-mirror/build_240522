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
import edify_generator
import ota_from_target_files
import sparse_img


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


def GetImage(which, tmpdir):
  """Returns an image object suitable for passing to BlockImageDiff.

  'which' partition must be "system" or "vendor". A prebuilt image and file
  map must already exist in tmpdir.
  """

  assert which in ("system", "vendor")

  path = os.path.join(tmpdir,  "IMAGES", which + ".img")
  mappath = os.path.join(tmpdir, which + ".map")

  # The image and map files must have been created prior to calling
  # ota_from_target_files.py (since LMP).
  print(path)
  assert os.path.exists(path) and os.path.exists(mappath)

  # Bug: http://b/20939131
  # In ext4 filesystems, block 0 might be changed even being mounted
  # R/O. We add it to clobbered_blocks so that it will be written to the
  # target unconditionally. Note that they are still part of care_map.
  clobbered_blocks = "0"

  return sparse_img.SparseImage(path, mappath, clobbered_blocks)


def HasVendorPartition(tmp_dir):
  print(os.path.join(tmp_dir, "IMAGES/vendor.img"))
  return os.path.exists(os.path.join(tmp_dir, "IMAGES/vendor.img"))


def WriteVerifyPackage():
  print ""



def WriteFullOTAPackage(output_zip):
  script = edify_generator.EdifyGenerator(3, OPTIONS.info_dict)

  recovery_mount_options = OPTIONS.info_dict.get("recovery_mount_options")
  oem_props = OPTIONS.info_dict.get("oem_fingerprint_properties")
  oem_dicts = None
  if oem_props:
    oem_dicts = ota_from_target_files._LoadOemDicts(script, recovery_mount_options)

  target_fp = ota_from_target_files.CalculateFingerprint(oem_props, oem_dicts and oem_dicts[0],
                                   OPTIONS.info_dict)
  metadata = {
      "post-build": target_fp,
      "pre-device": ota_from_target_files.GetOemProperty("ro.product.device", oem_props,
                                   oem_dicts and oem_dicts[0],
                                   OPTIONS.info_dict),
      "post-timestamp": ota_from_target_files.GetBuildProp("ro.build.date.utc", OPTIONS.info_dict),
  }

  device_specific = common.DeviceSpecificParams(
      input_version=OPTIONS.info_dict["recovery_api_version"],
      output_zip=output_zip,
      script=script,
      metadata=metadata,
      info_dict=OPTIONS.info_dict)

  #assert HasRecoveryPatch(input_zip)

  metadata["ota-type"] = "BLOCK"

  ts = ota_from_target_files.GetBuildProp("ro.build.date.utc", OPTIONS.info_dict)
  ts_text = ota_from_target_files.GetBuildProp("ro.build.date", OPTIONS.info_dict)
  script.AssertOlderBuild(ts, ts_text)

  ota_from_target_files.AppendAssertions(script, OPTIONS.info_dict, oem_dicts)
  device_specific.FullOTA_Assertions()

  # Two-step package strategy (in chronological order, which is *not*
  # the order in which the generated script has things):
  #
  # if stage is not "2/3" or "3/3":
  #    write recovery image to boot partition
  #    set stage to "2/3"
  #    reboot to boot partition and restart recovery
  # else if stage is "2/3":
  #    write recovery image to recovery partition
  #    set stage to "3/3"
  #    reboot to recovery partition and restart recovery
  # else:
  #    (stage must be "3/3")
  #    set stage to ""
  #    do normal full package installation:
  #       wipe and install system, boot image, etc.
  #       set up system to update recovery partition on first boot
  #    complete script normally
  #    (allow recovery to mark itself finished and reboot)

  recovery_img = common.GetBootableImage("recovery.img", "recovery.img",
                                         OPTIONS.input_tmp, "RECOVERY")
  if OPTIONS.two_step:
    if not OPTIONS.info_dict.get("multistage_support", None):
      assert False, "two-step packages not supported by this build"
    fs = OPTIONS.info_dict["fstab"]["/misc"]
    assert fs.fs_type.upper() == "EMMC", \
        "two-step packages only supported on devices with EMMC /misc partitions"
    bcb_dev = {"bcb_dev": fs.device}
    common.ZipWriteStr(output_zip, "recovery.img", recovery_img.data)
    script.AppendExtra("""
if get_stage("%(bcb_dev)s") == "2/3" then
""" % bcb_dev)

    # Stage 2/3: Write recovery image to /recovery (currently running /boot).
    script.Comment("Stage 2/3")
    script.WriteRawImage("/recovery", "recovery.img")
    script.AppendExtra("""
set_stage("%(bcb_dev)s", "3/3");
reboot_now("%(bcb_dev)s", "recovery");
else if get_stage("%(bcb_dev)s") == "3/3" then
""" % bcb_dev)

    # Stage 3/3: Make changes.
    script.Comment("Stage 3/3")

  # Dump fingerprints
  script.Print("Target: %s" % target_fp)

  device_specific.FullOTA_InstallBegin()

  system_progress = 0.75

  if OPTIONS.wipe_user_data:
    system_progress -= 0.1
  #if HasVendorPartition(input_zip):
  #  system_progress -= 0.1

  # Place a copy of file_contexts.bin into the OTA package which will be used
  # by the recovery program.
  if "selinux_fc" in OPTIONS.info_dict:
    ota_from_target_files.WritePolicyConfig(OPTIONS.info_dict["selinux_fc"], output_zip)

  recovery_mount_options = OPTIONS.info_dict.get("recovery_mount_options")

  script.ShowProgress(system_progress, 0)

  # Full OTA is done as an "incremental" against an empty source image. This
  # has the effect of writing new data from the package to the entire
  # partition, but lets us reuse the updater code that writes incrementals to
  # do it.
  system_tgt = GetImage("system", OPTIONS.input_tmp)
  system_tgt.ResetFileMap()
  system_diff = common.BlockDifference("system", system_tgt, src=None)
  system_diff.WriteScript(script, output_zip)

  if os.path.exists(os.path.join(os.path.join(OPTIONS.input_tmp, "IMAGES", "vendor.img"))):
    script.ShowProgress(0.1, 0)

    vendor_tgt = GetImage("vendor", OPTIONS.input_tmp)
    vendor_tgt.ResetFileMap()
    vendor_diff = common.BlockDifference("vendor", vendor_tgt)
    vendor_diff.WriteScript(script, output_zip)

  boot_path = os.path.join(OPTIONS.input_tmp, "boot.img");
  assert os.path.exists(boot_path)
  boot_img = common.File.FromLocalFile("boot.img", boot_path)
  common.CheckSize(boot_img.data, "boot.img", OPTIONS.info_dict)
  common.ZipWriteStr(output_zip, "boot.img", boot_img.data)

  script.ShowProgress(0.05, 5)
  script.WriteRawImage("/boot", "boot.img")

  script.ShowProgress(0.2, 10)
  #device_specific.FullOTA_InstallEnd()

  if OPTIONS.extra_script is not None:
    script.AppendExtra(OPTIONS.extra_script)

  script.UnmountAll()

  if OPTIONS.wipe_user_data:
    script.ShowProgress(0.1, 10)
    script.FormatPartition("/data")

  if OPTIONS.two_step:
    script.AppendExtra("""
set_stage("%(bcb_dev)s", "");
""" % bcb_dev)
    script.AppendExtra("else\n")

    # Stage 1/3: Nothing to verify for full OTA. Write recovery image to /boot.
    script.Comment("Stage 1/3")
    ota_from_target_files._WriteRecoveryImageToBoot(script, output_zip)

    script.AppendExtra("""
set_stage("%(bcb_dev)s", "2/3");
reboot_now("%(bcb_dev)s", "");
endif;
endif;
""" % bcb_dev)

  script.SetProgress(1)
  updater_path = OPTIONS.updater_binary
  if updater_path is None:
    updater_path = os.path.join(OPTIONS.input_tmp, "OTA/bin/updater")
  script.AddToZip(None, output_zip, input_path=updater_path)
  metadata["ota-required-cache"] = str(script.required_cache)
  ota_from_target_files.WriteMetadata(metadata, output_zip)


def HandleDowngradeMetadata(metadata):
  # Only incremental OTAs are allowed to reach here.
  assert OPTIONS.incremental_source is not None

  post_timestamp = ota_from_target_files.GetBuildProp("ro.build.date.utc", OPTIONS.target_info_dict)
  pre_timestamp = ota_from_target_files.GetBuildProp("ro.build.date.utc", OPTIONS.source_info_dict)
  is_downgrade = long(post_timestamp) < long(pre_timestamp)

  if OPTIONS.downgrade:
    if not is_downgrade:
      raise RuntimeError("--downgrade specified but no downgrade detected: "
                         "pre: %s, post: %s" % (pre_timestamp, post_timestamp))
    metadata["ota-downgrade"] = "yes"
  elif OPTIONS.timestamp:
    if not is_downgrade:
      raise RuntimeError("--timestamp specified but no timestamp hack needed: "
                         "pre: %s, post: %s" % (pre_timestamp, post_timestamp))
    metadata["post-timestamp"] = str(long(pre_timestamp) + 1)
  else:
    if is_downgrade:
      raise RuntimeError("Downgrade detected based on timestamp check: "
                         "pre: %s, post: %s. Need to specify --timestamp OR "
                         "--downgrade to allow building the incremental." % (
                             pre_timestamp, post_timestamp))
    metadata["post-timestamp"] = post_timestamp


def WriteIncrementalOTAPackage(output_zip):
  source_version = OPTIONS.source_info_dict["recovery_api_version"]
  target_version = OPTIONS.target_info_dict["recovery_api_version"]

  if source_version == 0:
    print("WARNING: generating edify script for a source that "
          "can't install it.")
  script = edify_generator.EdifyGenerator(
      source_version, OPTIONS.target_info_dict,
      fstab=OPTIONS.source_info_dict["fstab"])

  recovery_mount_options = OPTIONS.source_info_dict.get(
      "recovery_mount_options")
  source_oem_props = OPTIONS.source_info_dict.get("oem_fingerprint_properties")
  target_oem_props = OPTIONS.target_info_dict.get("oem_fingerprint_properties")
  oem_dicts = None
  if source_oem_props and target_oem_props:
    oem_dicts = ota_from_target_files._LoadOemDicts(script, recovery_mount_options)

  metadata = {
      "pre-device": ota_from_target_files.GetOemProperty("ro.product.device", source_oem_props,
                                   oem_dicts and oem_dicts[0],
                                   OPTIONS.source_info_dict),
      "ota-type": "BLOCK",
  }

  HandleDowngradeMetadata(metadata)

  device_specific = common.DeviceSpecificParams(
      source_version=source_version,
      target_version=target_version,
      script=script,
      metadata=metadata,
      info_dict=OPTIONS.source_info_dict)

  source_fp = ota_from_target_files.CalculateFingerprint(source_oem_props, oem_dicts and oem_dicts[0],
                                   OPTIONS.source_info_dict)
  target_fp = ota_from_target_files.CalculateFingerprint(target_oem_props, oem_dicts and oem_dicts[0],
                                   OPTIONS.target_info_dict)
  metadata["pre-build"] = source_fp
  metadata["post-build"] = target_fp
  metadata["pre-build-incremental"] = ota_from_target_files.GetBuildProp(
      "ro.build.version.incremental", OPTIONS.source_info_dict)
  metadata["post-build-incremental"] = ota_from_target_files.GetBuildProp(
      "ro.build.version.incremental", OPTIONS.target_info_dict)

  source_boot = common.GetBootableImage(
      "boot.img", "boot.img", OPTIONS.source_tmp, "",
      OPTIONS.source_info_dict)
  target_boot = common.GetBootableImage(
      "boot.img", "boot.img", OPTIONS.target_tmp, "BOOT")
  updating_boot = (not OPTIONS.two_step and
                   (source_boot.data != target_boot.data))

  target_recovery = common.GetBootableImage(
      "recovery.img", "recovery.img", OPTIONS.target_tmp, "RECOVERY")

  system_src = GetImage("system", OPTIONS.source_tmp)
  system_tgt = GetImage("system", OPTIONS.target_tmp)

  blockimgdiff_version = 1
  if OPTIONS.info_dict:
    blockimgdiff_version = max(
        int(i) for i in
        OPTIONS.info_dict.get("blockimgdiff_versions", "1").split(","))

  # Check the first block of the source system partition for remount R/W only
  # if the filesystem is ext4.
  system_src_partition = OPTIONS.source_info_dict["fstab"]["/system"]
  check_first_block = system_src_partition.fs_type == "ext4"
  # Disable using imgdiff for squashfs. 'imgdiff -z' expects input files to be
  # in zip formats. However with squashfs, a) all files are compressed in LZ4;
  # b) the blocks listed in block map may not contain all the bytes for a given
  # file (because they're rounded to be 4K-aligned).
  system_tgt_partition = OPTIONS.target_info_dict["fstab"]["/system"]
  disable_imgdiff = (system_src_partition.fs_type == "squashfs" or
                     system_tgt_partition.fs_type == "squashfs")
  system_diff = common.BlockDifference("system", system_tgt, system_src,
                                       check_first_block,
                                       version=blockimgdiff_version,
                                       disable_imgdiff=disable_imgdiff)

  if HasVendorPartition(OPTIONS.target_tmp):
    if not HasVendorPartition(OPTIONS.source_tmp):
      raise RuntimeError("can't generate incremental that adds /vendor")
    vendor_src = GetImage("vendor", OPTIONS.source_tmp)
    vendor_tgt = GetImage("vendor", OPTIONS.target_tmp)

    # Check first block of vendor partition for remount R/W only if
    # disk type is ext4
    vendor_partition = OPTIONS.source_info_dict["fstab"]["/vendor"]
    check_first_block = vendor_partition.fs_type == "ext4"
    disable_imgdiff = vendor_partition.fs_type == "squashfs"
    vendor_diff = common.BlockDifference("vendor", vendor_tgt, vendor_src,
                                         check_first_block,
                                         version=blockimgdiff_version,
                                         disable_imgdiff=disable_imgdiff)
  #else:
  vendor_diff = None

  ota_from_target_files.AppendAssertions(script, OPTIONS.target_info_dict, oem_dicts)
  device_specific.IncrementalOTA_Assertions()

  # Two-step incremental package strategy (in chronological order,
  # which is *not* the order in which the generated script has
  # things):
  #
  # if stage is not "2/3" or "3/3":
  #    do verification on current system
  #    write recovery image to boot partition
  #    set stage to "2/3"
  #    reboot to boot partition and restart recovery
  # else if stage is "2/3":
  #    write recovery image to recovery partition
  #    set stage to "3/3"
  #    reboot to recovery partition and restart recovery
  # else:
  #    (stage must be "3/3")
  #    perform update:
  #       patch system files, etc.
  #       force full install of new boot image
  #       set up system to update recovery partition on first boot
  #    complete script normally
  #    (allow recovery to mark itself finished and reboot)

  if OPTIONS.two_step:
    if not OPTIONS.source_info_dict.get("multistage_support", None):
      assert False, "two-step packages not supported by this build"
    fs = OPTIONS.source_info_dict["fstab"]["/misc"]
    assert fs.fs_type.upper() == "EMMC", \
        "two-step packages only supported on devices with EMMC /misc partitions"
    bcb_dev = {"bcb_dev": fs.device}
    common.ZipWriteStr(output_zip, "recovery.img", target_recovery.data)
    script.AppendExtra("""
if get_stage("%(bcb_dev)s") == "2/3" then
""" % bcb_dev)

    # Stage 2/3: Write recovery image to /recovery (currently running /boot).
    script.Comment("Stage 2/3")
    script.AppendExtra("sleep(20);\n")
    script.WriteRawImage("/recovery", "recovery.img")
    script.AppendExtra("""
set_stage("%(bcb_dev)s", "3/3");
reboot_now("%(bcb_dev)s", "recovery");
else if get_stage("%(bcb_dev)s") != "3/3" then
""" % bcb_dev)

    # Stage 1/3: (a) Verify the current system.
    script.Comment("Stage 1/3")

  # Dump fingerprints
  script.Print("Source: %s" % (source_fp,))
  script.Print("Target: %s" % (target_fp,))

  script.Print("Verifying current system...")

  device_specific.IncrementalOTA_VerifyBegin()

  # When blockimgdiff version is less than 3 (non-resumable block-based OTA),
  # patching on a device that's already on the target build will damage the
  # system. Because operations like move don't check the block state, they
  # always apply the changes unconditionally.
  if blockimgdiff_version <= 2:
    if source_oem_props is None:
      script.AssertSomeFingerprint(source_fp)
    else:
      script.AssertSomeThumbprint(
          ota_from_target_files.GetBuildProp("ro.build.thumbprint", OPTIONS.source_info_dict))

  else: # blockimgdiff_version > 2
    if source_oem_props is None and target_oem_props is None:
      script.AssertSomeFingerprint(source_fp, target_fp)
    elif source_oem_props is not None and target_oem_props is not None:
      script.AssertSomeThumbprint(
          ota_from_target_files.GetBuildProp("ro.build.thumbprint", OPTIONS.target_info_dict),
          ota_from_target_files.GetBuildProp("ro.build.thumbprint", OPTIONS.source_info_dict))
    elif source_oem_props is None and target_oem_props is not None:
      script.AssertFingerprintOrThumbprint(
          source_fp,
          ota_from_target_files.GetBuildProp("ro.build.thumbprint", OPTIONS.target_info_dict))
    else:
      script.AssertFingerprintOrThumbprint(
          target_fp,
          ota_from_target_files.GetBuildProp("ro.build.thumbprint", OPTIONS.source_info_dict))

  # Check the required cache size (i.e. stashed blocks).
  size = []
  if system_diff:
    size.append(system_diff.required_cache)
  if vendor_diff:
    size.append(vendor_diff.required_cache)

  if updating_boot:
    boot_type, boot_device = common.GetTypeAndDevice(
        "/boot", OPTIONS.source_info_dict)
    d = common.Difference(target_boot, source_boot)
    _, _, d = d.ComputePatch()
    if d is None:
      include_full_boot = True
      common.ZipWriteStr(output_zip, "boot.img", target_boot.data)
    else:
      include_full_boot = False

      print("boot      target: %d  source: %d  diff: %d" % (
          target_boot.size, source_boot.size, len(d)))

      common.ZipWriteStr(output_zip, "patch/boot.img.p", d)

      script.PatchCheck("%s:%s:%d:%s:%d:%s" %
                        (boot_type, boot_device,
                         source_boot.size, source_boot.sha1,
                         target_boot.size, target_boot.sha1))
      size.append(target_boot.size)

  if size:
    script.CacheFreeSpaceCheck(max(size))

  #device_specific.IncrementalOTA_VerifyEnd()

  if OPTIONS.two_step:
    # Stage 1/3: (b) Write recovery image to /boot.
    ota_from_target_files._WriteRecoveryImageToBoot(script, output_zip)

    script.AppendExtra("""
set_stage("%(bcb_dev)s", "2/3");
reboot_now("%(bcb_dev)s", "");
else
""" % bcb_dev)

    # Stage 3/3: Make changes.
    script.Comment("Stage 3/3")

  # Verify the existing partitions.
  system_diff.WriteVerifyScript(script, touched_blocks_only=True)
  if vendor_diff:
    vendor_diff.WriteVerifyScript(script, touched_blocks_only=True)

  script.Comment("---- start making changes here ----")

  device_specific.IncrementalOTA_InstallBegin()

  system_diff.WriteScript(script, output_zip,
                          progress=0.8 if vendor_diff else 0.9)

  if vendor_diff:
    vendor_diff.WriteScript(script, output_zip, progress=0.1)

  if OPTIONS.two_step:
    common.ZipWriteStr(output_zip, "boot.img", target_boot.data)
    script.WriteRawImage("/boot", "boot.img")
    print("writing full boot image (forced by two-step mode)")

  if not OPTIONS.two_step:
    if updating_boot:
      if include_full_boot:
        print("boot image changed; including full.")
        script.Print("Installing boot image...")
        script.WriteRawImage("/boot", "boot.img")
      else:
        # Produce the boot image by applying a patch to the current
        # contents of the boot partition, and write it back to the
        # partition.
        print("boot image changed; including patch.")
        script.Print("Patching boot image...")
        script.ShowProgress(0.1, 10)
        script.ApplyPatch("%s:%s:%d:%s:%d:%s"
                          % (boot_type, boot_device,
                             source_boot.size, source_boot.sha1,
                             target_boot.size, target_boot.sha1),
                          "-",
                          target_boot.size, target_boot.sha1,
                          source_boot.sha1, "patch/boot.img.p")
    else:
      print("boot image unchanged; skipping.")

  # Do device-specific installation (eg, write radio image).
  #device_specific.IncrementalOTA_InstallEnd()

  if OPTIONS.extra_script is not None:
    script.AppendExtra(OPTIONS.extra_script)

  if OPTIONS.wipe_user_data:
    script.Print("Erasing user data...")
    script.FormatPartition("/data")
    metadata["ota-wipe"] = "yes"

  if OPTIONS.two_step:
    script.AppendExtra("""
set_stage("%(bcb_dev)s", "");
endif;
endif;
""" % bcb_dev)



  script.SetProgress(1)
  # For downgrade OTAs, we prefer to use the update-binary in the source
  # build that is actually newer than the one in the target build.

  updater_path = OPTIONS.updater_binary
  if updater_path is None:
    if OPTIONS.downgrade:
      updater_path = os.path.join(OPTIONS.source_tmp, "OTA/bin/updater")
    else:
      updater_path = os.path.join(OPTIONS.target_tmp, "OTA/bin/updater")

  script.AddToZip(None, output_zip, input_path=updater_path)
  metadata["ota-required-cache"] = str(script.required_cache)
  ota_from_target_files.WriteMetadata(metadata, output_zip)


def setup_directory(info_dir, out_dir):
  os.mkdir(os.path.join(out_dir, "IMAGES"))
  for item in os.listdir(out_dir):
    src = os.path.join(out_dir, item)
    dst = os.path.join(out_dir, "IMAGES", item)
    if os.path.isfile(src):
      os.rename(src, dst)

  for item in os.listdir(info_dir):
    src = os.path.join(info_dir, item)
    dst = os.path.join(out_dir, item)
    if os.path.isdir(src):
      shutil.copytree(src, dst)
    else:
      shutil.copy(src, dst)


def check_files(prefix):
  assert(os.path.exists(os.path.join(prefix, "IMAGES/system.img")))
  assert(os.path.exists(os.path.join(prefix, "system.map")))
  assert(os.path.exists(os.path.join(prefix, "META/misc_info.txt")))
  assert(os.path.exists(os.path.join(prefix, "RECOVERY/RAMDISK/etc/recovery.fstab")))
  assert(os.path.exists(os.path.join(prefix, "SYSTEM/build.prop")))
  if HasVendorPartition(prefix):
    assert(os.path.exists(os.path.join(prefix, "VENDOR/build.prop")))
    assert(os.path.exists(os.path.join(prefix, "vendor.map")))


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
  OPTIONS.info_dict = common.LoadInfoDict(OPTIONS.tgt_info_dir, OPTIONS.tgt_info_dir)

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

  print("unzipping target img-files...")
  OPTIONS.input_tmp, _ = common.UnzipTemp(args[0], ["*.img"])
  setup_directory(OPTIONS.tgt_info_dir, OPTIONS.input_tmp)
  OPTIONS.target_tmp = OPTIONS.input_tmp
  check_files(OPTIONS.target_tmp)

  if OPTIONS.verbose:
    print("--- target info ---")
    common.DumpInfoDict(OPTIONS.info_dict)

  # If the caller explicitly specified the device-specific extensions
  # path via -s/--device_specific, use that.
  # Otherwise, take the path of the file from 'tool_extensions' in the
  # info dict and look for that in the local filesystem, relative to
  # the current directory.

  if OPTIONS.device_specific is None:
    OPTIONS.device_specific = OPTIONS.info_dict.get("tool_extensions", None)
  else:
    OPTIONS.device_specific = os.path.abspath(OPTIONS.device_specific)

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
    print("--- can't determine the cache partition size ---")
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
    print("unzipping source target-files...")
    # TODO close the zip
    OPTIONS.source_tmp, _ = common.UnzipTemp(OPTIONS.incremental_source, ["*.img"])
    setup_directory(OPTIONS.src_info_dir, OPTIONS.source_tmp)
    check_files(OPTIONS.source_tmp)

    OPTIONS.target_info_dict = OPTIONS.info_dict
    OPTIONS.source_info_dict = common.LoadInfoDict(OPTIONS.src_info_dir,
                                                   OPTIONS.source_tmp)


    if OPTIONS.verbose:
      print("--- source info ---")
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
      print("--- failed to build incremental; falling back to full ---")
      OPTIONS.incremental_source = None
      WriteFullOTAPackage(output_zip)

  common.ZipClose(output_zip)

  # Sign the generated zip package unless no_signing is specified.
  if not OPTIONS.no_signing:
    ota_from_target_files.SignOutput(temp_zip_file.name, args[1])
    temp_zip_file.close()

  print("done.")


if __name__ == '__main__':
  try:
    common.CloseInheritedPipes()
    main(sys.argv[1:])
  except common.ExternalError as e:
    logging.error("\n   ERROR: %s\n" % (e,))
    sys.exit(1)
  #finally:
  #  common.Cleanup()