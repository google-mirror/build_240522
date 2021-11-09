#!/usr/bin/env python
#
# Copyright (C) 2014 The Android Open Source Project
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
Given a target-files zipfile that does not contain images (ie, does
not have an IMAGES/ top-level subdirectory), produce the images and
add them to the zipfile.

Usage:  add_img_to_target_files [flag] target_files

  -a  (--add_missing)
      Build and add missing images to "IMAGES/". If this option is
      not specified, this script will simply exit when "IMAGES/"
      directory exists in the target file.

  -r  (--rebuild_recovery)
      Rebuild the recovery patch and write it to the system image. Only
      meaningful when system image needs to be rebuilt.

  --replace_verity_private_key
      Replace the private key used for verity signing. (same as the option
      in sign_target_files_apks)

  --replace_verity_public_key
       Replace the certificate (public key) used for verity verification. (same
       as the option in sign_target_files_apks)

  --is_signing
      Skip building & adding the images for "userdata" and "cache" if we
      are signing the target files.
"""

from __future__ import print_function

import datetime
import os
import shlex
import shutil
import subprocess
import sys
import uuid
import zipfile

import build_image
import common
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
import rangelib
import sparse_img
=======
import verity_utils
import ota_metadata_pb2

from apex_utils import GetApexInfoFromTargetFiles
from common import AddCareMapForAbOta
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

if sys.hexversion < 0x02070000:
  print("Python 2.7 or newer is required.", file=sys.stderr)
  sys.exit(1)

OPTIONS = common.OPTIONS

OPTIONS.add_missing = False
OPTIONS.rebuild_recovery = False
OPTIONS.replace_updated_files_list = []
OPTIONS.replace_verity_public_key = False
OPTIONS.replace_verity_private_key = False
OPTIONS.is_signing = False


# Partitions that should have their care_map added to META/care_map.txt.
PARTITIONS_WITH_CARE_MAP = ('system', 'vendor', 'product')


class OutputFile(object):
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  def __init__(self, output_zip, input_dir, prefix, name):
=======
  """A helper class to write a generated file to the given dir or zip.

  When generating images, we want the outputs to go into the given zip file, or
  the given dir.

  Attributes:
    name: The name of the output file, regardless of the final destination.
  """

  def __init__(self, output_zip, input_dir, *args):
    # We write the intermediate output file under the given input_dir, even if
    # the final destination is a zip archive.
    self.name = os.path.join(input_dir, *args)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
    self._output_zip = output_zip
    self.input_name = os.path.join(input_dir, prefix, name)

    if self._output_zip:
      self._zip_name = os.path.join(*args)

      root, suffix = os.path.splitext(name)
      self.name = common.MakeTempFile(prefix=root + '-', suffix=suffix)
    else:
      self.name = self.input_name

  def Write(self):
    if self._output_zip:
      common.ZipWrite(self._output_zip, self.name, self._zip_name)


<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
def GetCareMap(which, imgname):
  """Returns the care_map string for the given partition.

  Args:
    which: The partition name, must be listed in PARTITIONS_WITH_CARE_MAP.
    imgname: The filename of the image.

  Returns:
    (which, care_map_ranges): care_map_ranges is the raw string of the care_map
    RangeSet.
  """
  assert which in PARTITIONS_WITH_CARE_MAP

  simg = sparse_img.SparseImage(imgname)
  care_map_ranges = simg.care_map
  key = which + "_adjusted_partition_size"
  adjusted_blocks = OPTIONS.info_dict.get(key)
  if adjusted_blocks:
    assert adjusted_blocks > 0, "blocks should be positive for " + which
    care_map_ranges = care_map_ranges.intersect(rangelib.RangeSet(
        "0-%d" % (adjusted_blocks,)))

  return [which, care_map_ranges.to_string_raw()]


=======
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
def AddSystem(output_zip, recovery_img=None, boot_img=None):
  """Turn the contents of SYSTEM into a system image and store it in
  output_zip. Returns the name of the system image file."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "system.img")
  if os.path.exists(img.input_name):
    print("system.img already exists; no need to rebuild...")
    return img.input_name

  def output_sink(fn, data):
    ofile = open(os.path.join(OPTIONS.input_tmp, "SYSTEM", fn), "w")
    ofile.write(data)
    ofile.close()

    arc_name = "SYSTEM/" + fn
    if arc_name in output_zip.namelist():
      OPTIONS.replace_updated_files_list.append(arc_name)
    else:
      common.ZipWrite(output_zip, ofile.name, arc_name)

  if OPTIONS.rebuild_recovery:
    print("Building new recovery patch")
    common.MakeRecoveryPatch(OPTIONS.input_tmp, output_sink, recovery_img,
                             boot_img, info_dict=OPTIONS.info_dict)

  block_list = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "system.map")
  CreateImage(OPTIONS.input_tmp, OPTIONS.info_dict, "system", img,
              block_list=block_list)
  return img.name


def AddSystemOther(output_zip):
  """Turn the contents of SYSTEM_OTHER into a system_other image
  and store it in output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "system_other.img")
  if os.path.exists(img.input_name):
    print("system_other.img already exists; no need to rebuild...")
    return

  CreateImage(OPTIONS.input_tmp, OPTIONS.info_dict, "system_other", img)


def AddVendor(output_zip):
  """Turn the contents of VENDOR into a vendor image and store in it
  output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "vendor.img")
  if os.path.exists(img.input_name):
    print("vendor.img already exists; no need to rebuild...")
    return img.input_name

  block_list = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "vendor.map")
  CreateImage(OPTIONS.input_tmp, OPTIONS.info_dict, "vendor", img,
              block_list=block_list)
  return img.name


def AddProduct(output_zip):
  """Turn the contents of PRODUCT into a product image and store it in
  output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "product.img")
  if os.path.exists(img.input_name):
    print("product.img already exists; no need to rebuild...")
    return img.input_name

  block_list = OutputFile(
      output_zip, OPTIONS.input_tmp, "IMAGES", "product.map")
  CreateImage(
      OPTIONS.input_tmp, OPTIONS.info_dict, "product", img,
      block_list=block_list)
  return img.name


<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
def AddSystemExt(output_zip):
  """Turn the contents of SYSTEM_EXT into a system_ext image and store it in
  output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES",
                   "system_ext.img")
  if os.path.exists(img.name):
    logger.info("system_ext.img already exists; no need to rebuild...")
    return img.name

  block_list = OutputFile(
      output_zip, OPTIONS.input_tmp, "IMAGES", "system_ext.map")
  CreateImage(
      OPTIONS.input_tmp, OPTIONS.info_dict, "system_ext", img,
      block_list=block_list)
  return img.name


def AddOdm(output_zip):
  """Turn the contents of ODM into an odm image and store it in output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "odm.img")
  if os.path.exists(img.name):
    logger.info("odm.img already exists; no need to rebuild...")
    return img.name

  block_list = OutputFile(
      output_zip, OPTIONS.input_tmp, "IMAGES", "odm.map")
  CreateImage(
      OPTIONS.input_tmp, OPTIONS.info_dict, "odm", img,
      block_list=block_list)
  return img.name


def AddVendorDlkm(output_zip):
  """Turn the contents of VENDOR_DLKM into an vendor_dlkm image and store it in output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "vendor_dlkm.img")
  if os.path.exists(img.name):
    logger.info("vendor_dlkm.img already exists; no need to rebuild...")
    return img.name

  block_list = OutputFile(
      output_zip, OPTIONS.input_tmp, "IMAGES", "vendor_dlkm.map")
  CreateImage(
      OPTIONS.input_tmp, OPTIONS.info_dict, "vendor_dlkm", img,
      block_list=block_list)
  return img.name


def AddOdmDlkm(output_zip):
  """Turn the contents of OdmDlkm into an odm_dlkm image and store it in output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "odm_dlkm.img")
  if os.path.exists(img.name):
    logger.info("odm_dlkm.img already exists; no need to rebuild...")
    return img.name

  block_list = OutputFile(
      output_zip, OPTIONS.input_tmp, "IMAGES", "odm_dlkm.map")
  CreateImage(
      OPTIONS.input_tmp, OPTIONS.info_dict, "odm_dlkm", img,
      block_list=block_list)
  return img.name


>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
def AddDtbo(output_zip):
  """Adds the DTBO image.

  Uses the image under IMAGES/ if it already exists. Otherwise looks for the
  image under PREBUILT_IMAGES/, signs it as needed, and returns the image name.
  """
  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "dtbo.img")
  if os.path.exists(img.input_name):
    print("dtbo.img already exists; no need to rebuild...")
    return img.input_name

  dtbo_prebuilt_path = os.path.join(
      OPTIONS.input_tmp, "PREBUILT_IMAGES", "dtbo.img")
  assert os.path.exists(dtbo_prebuilt_path)
  shutil.copy(dtbo_prebuilt_path, img.name)

  # AVB-sign the image as needed.
  if OPTIONS.info_dict.get("avb_enable") == "true":
    avbtool = os.getenv('AVBTOOL') or OPTIONS.info_dict["avb_avbtool"]
    part_size = OPTIONS.info_dict["dtbo_size"]
    # The AVB hash footer will be replaced if already present.
    cmd = [avbtool, "add_hash_footer", "--image", img.name,
           "--partition_size", str(part_size), "--partition_name", "dtbo"]
    common.AppendAVBSigningArgs(cmd, "dtbo")
    args = OPTIONS.info_dict.get("avb_dtbo_add_hash_footer_args")
    if args and args.strip():
      cmd.extend(shlex.split(args))
    p = common.Run(cmd, stdout=subprocess.PIPE)
    p.communicate()
    assert p.returncode == 0, \
        "avbtool add_hash_footer of %s failed" % (img.name,)

  img.Write()
  return img.name

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======

def AddPvmfw(output_zip):
  """Adds the pvmfw image.

  Uses the image under IMAGES/ if it already exists. Otherwise looks for the
  image under PREBUILT_IMAGES/, signs it as needed, and returns the image name.
  """
  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "pvmfw.img")
  if os.path.exists(img.name):
    logger.info("pvmfw.img already exists; no need to rebuild...")
    return img.name

  pvmfw_prebuilt_path = os.path.join(
      OPTIONS.input_tmp, "PREBUILT_IMAGES", "pvmfw.img")
  assert os.path.exists(pvmfw_prebuilt_path)
  shutil.copy(pvmfw_prebuilt_path, img.name)

  # AVB-sign the image as needed.
  if OPTIONS.info_dict.get("avb_enable") == "true":
    # Signing requires +w
    os.chmod(img.name, os.stat(img.name).st_mode | stat.S_IWUSR)

    avbtool = OPTIONS.info_dict["avb_avbtool"]
    part_size = OPTIONS.info_dict["pvmfw_size"]
    # The AVB hash footer will be replaced if already present.
    cmd = [avbtool, "add_hash_footer", "--image", img.name,
           "--partition_size", str(part_size), "--partition_name", "pvmfw"]
    common.AppendAVBSigningArgs(cmd, "pvmfw")
    args = OPTIONS.info_dict.get("avb_pvmfw_add_hash_footer_args")
    if args and args.strip():
      cmd.extend(shlex.split(args))
    common.RunAndCheckOutput(cmd)

  img.Write()
  return img.name


def AddCustomImages(output_zip, partition_name):
  """Adds and signs custom images in IMAGES/.

  Args:
    output_zip: The output zip file (needs to be already open), or None to
        write images to OPTIONS.input_tmp/.

  Uses the image under IMAGES/ if it already exists. Otherwise looks for the
  image under PREBUILT_IMAGES/, signs it as needed, and returns the image name.

  Raises:
    AssertionError: If image can't be found.
  """

  key_path = OPTIONS.info_dict.get("avb_{}_key_path".format(partition_name))
  algorithm = OPTIONS.info_dict.get("avb_{}_algorithm".format(partition_name))
  extra_args = OPTIONS.info_dict.get(
      "avb_{}_add_hashtree_footer_args".format(partition_name))
  partition_size = OPTIONS.info_dict.get(
      "avb_{}_partition_size".format(partition_name))

  builder = verity_utils.CreateCustomImageBuilder(
      OPTIONS.info_dict, partition_name, partition_size,
      key_path, algorithm, extra_args)

  for img_name in OPTIONS.info_dict.get(
      "avb_{}_image_list".format(partition_name)).split():
    custom_image = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", img_name)
    if os.path.exists(custom_image.name):
      continue

    custom_image_prebuilt_path = os.path.join(
        OPTIONS.input_tmp, "PREBUILT_IMAGES", img_name)
    assert os.path.exists(custom_image_prebuilt_path), \
      "Failed to find %s at %s" % (img_name, custom_image_prebuilt_path)

    shutil.copy(custom_image_prebuilt_path, custom_image.name)

    if builder is not None:
      builder.Build(custom_image.name)

    custom_image.Write()

  default = os.path.join(OPTIONS.input_tmp, "IMAGES", partition_name + ".img")
  assert os.path.exists(default), \
      "There should be one %s.img" % (partition_name)
  return default

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

def CreateImage(input_dir, info_dict, what, output_file, block_list=None):
  print("creating " + what + ".img...")

  image_props = build_image.ImagePropFromGlobalDict(info_dict, what)
  fstab = info_dict["fstab"]
  mount_point = "/" + what
  if fstab and mount_point in fstab:
    image_props["fs_type"] = fstab[mount_point].fs_type

  # Use a fixed timestamp (01/01/2009) when packaging the image.
  # Bug: 24377993
  epoch = datetime.datetime.fromtimestamp(0)
  timestamp = (datetime.datetime(2009, 1, 1) - epoch).total_seconds()
  image_props["timestamp"] = int(timestamp)

  if what == "system":
    fs_config_prefix = ""
  else:
    fs_config_prefix = what + "_"

  fs_config = os.path.join(
      input_dir, "META/" + fs_config_prefix + "filesystem_config.txt")
  if not os.path.exists(fs_config):
    fs_config = None

  # Override values loaded from info_dict.
  if fs_config:
    image_props["fs_config"] = fs_config
  if block_list:
    image_props["block_list"] = block_list.name

  # Use repeatable ext4 FS UUID and hash_seed UUID (based on partition name and
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  # build fingerprint).
  uuid_seed = what + "-"
  if "build.prop" in info_dict:
    build_prop = info_dict["build.prop"]
    if "ro.build.fingerprint" in build_prop:
      uuid_seed += build_prop["ro.build.fingerprint"]
    elif "ro.build.thumbprint" in build_prop:
      uuid_seed += build_prop["ro.build.thumbprint"]
=======
  # build fingerprint). Also use the legacy build id, because the vbmeta digest
  # isn't available at this point.
  build_info = common.BuildInfo(info_dict, use_legacy_id=True)
  uuid_seed = what + "-" + build_info.GetPartitionFingerprint(what)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
  image_props["uuid"] = str(uuid.uuid5(uuid.NAMESPACE_URL, uuid_seed))
  hash_seed = "hash_seed-" + uuid_seed
  image_props["hash_seed"] = str(uuid.uuid5(uuid.NAMESPACE_URL, hash_seed))

  succ = build_image.BuildImage(os.path.join(input_dir, what.upper()),
                                image_props, output_file.name)
  assert succ, "build " + what + ".img image failed"

  output_file.Write()
  if block_list:
    block_list.Write()

  # Set the 'adjusted_partition_size' that excludes the verity blocks of the
  # given image. When avb is enabled, this size is the max image size returned
  # by the avb tool.
  is_verity_partition = "verity_block_device" in image_props
  verity_supported = (image_props.get("verity") == "true" or
                      image_props.get("avb_enable") == "true")
  is_avb_enable = image_props.get("avb_hashtree_enable") == "true"
  if verity_supported and (is_verity_partition or is_avb_enable):
    adjusted_blocks_value = image_props.get("partition_size")
    if adjusted_blocks_value:
      adjusted_blocks_key = what + "_adjusted_partition_size"
      info_dict[adjusted_blocks_key] = int(adjusted_blocks_value)/4096 - 1


def AddUserdata(output_zip):
  """Create a userdata image and store it in output_zip.

  In most case we just create and store an empty userdata.img;
  But the invoker can also request to create userdata.img with real
  data from the target files, by setting "userdata_img_with_data=true"
  in OPTIONS.info_dict.
  """

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "userdata.img")
  if os.path.exists(img.input_name):
    print("userdata.img already exists; no need to rebuild...")
    return

  # Skip userdata.img if no size.
  image_props = build_image.ImagePropFromGlobalDict(OPTIONS.info_dict, "data")
  if not image_props.get("partition_size"):
    return

  print("creating userdata.img...")

  # Use a fixed timestamp (01/01/2009) when packaging the image.
  # Bug: 24377993
  epoch = datetime.datetime.fromtimestamp(0)
  timestamp = (datetime.datetime(2009, 1, 1) - epoch).total_seconds()
  image_props["timestamp"] = int(timestamp)

  if OPTIONS.info_dict.get("userdata_img_with_data") == "true":
    user_dir = os.path.join(OPTIONS.input_tmp, "DATA")
  else:
    user_dir = common.MakeTempDir()

  fstab = OPTIONS.info_dict["fstab"]
  if fstab:
    image_props["fs_type"] = fstab["/data"].fs_type
  succ = build_image.BuildImage(user_dir, image_props, img.name)
  assert succ, "build userdata.img image failed"

  common.CheckSize(img.name, "userdata.img", OPTIONS.info_dict)
  img.Write()


def AppendVBMetaArgsForPartition(cmd, partition, img_path, public_key_dir):
  if not img_path:
    return

  # Check if chain partition is used.
  key_path = OPTIONS.info_dict.get("avb_" + partition + "_key_path")
  if key_path:
    # extract public key in AVB format to be included in vbmeta.img
    avbtool = os.getenv('AVBTOOL') or OPTIONS.info_dict["avb_avbtool"]
    public_key_path = os.path.join(public_key_dir, "%s.avbpubkey" % partition)
    p = common.Run([avbtool, "extract_public_key", "--key", key_path,
                    "--output", public_key_path],
                   stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    p.communicate()
    assert p.returncode == 0, \
        "avbtool extract_public_key fail for partition: %r" % partition

    rollback_index_location = OPTIONS.info_dict[
        "avb_" + partition + "_rollback_index_location"]
    cmd.extend(["--chain_partition", "%s:%s:%s" % (
        partition, rollback_index_location, public_key_path)])
  else:
    cmd.extend(["--include_descriptors_from_image", img_path])


def AddVBMeta(output_zip, partitions):
  """Creates a VBMeta image and store it in output_zip.

  Args:
    output_zip: The output zip file, which needs to be already open.
    partitions: A dict that's keyed by partition names with image paths as
        values. Only valid partition names are accepted, which include 'boot',
        'recovery', 'system', 'vendor', 'dtbo'.
  """
  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "vbmeta.img")
  if os.path.exists(img.input_name):
    print("vbmeta.img already exists; not rebuilding...")
    return img.input_name

  avbtool = os.getenv('AVBTOOL') or OPTIONS.info_dict["avb_avbtool"]
  cmd = [avbtool, "make_vbmeta_image", "--output", img.name]
  common.AppendAVBSigningArgs(cmd, "vbmeta")

  public_key_dir = common.MakeTempDir(prefix="avbpubkey-")
  for partition, path in partitions.items():
    assert partition in common.AVB_PARTITIONS, 'Unknown partition: %s' % (
        partition,)
    assert os.path.exists(path), 'Failed to find %s for partition %s' % (
        path, partition)
    AppendVBMetaArgsForPartition(cmd, partition, path, public_key_dir)

  args = OPTIONS.info_dict.get("avb_vbmeta_args")
  if args and args.strip():
    split_args = shlex.split(args)
    for index, arg in enumerate(split_args[:-1]):
      # Sanity check that the image file exists. Some images might be defined
      # as a path relative to source tree, which may not be available at the
      # same location when running this script (we have the input target_files
      # zip only). For such cases, we additionally scan other locations (e.g.
      # IMAGES/, RADIO/, etc) before bailing out.
      if arg == '--include_descriptors_from_image':
        image_path = split_args[index + 1]
        if os.path.exists(image_path):
          continue
        found = False
        for dir_name in ['IMAGES', 'RADIO', 'VENDOR_IMAGES', 'PREBUILT_IMAGES']:
          alt_path = os.path.join(
              OPTIONS.input_tmp, dir_name, os.path.basename(image_path))
          if os.path.exists(alt_path):
            split_args[index + 1] = alt_path
            found = True
            break
        assert found, 'failed to find %s' % (image_path,)
    cmd.extend(split_args)

  p = common.Run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
  p.communicate()
  assert p.returncode == 0, "avbtool make_vbmeta_image failed"
  img.Write()


def AddPartitionTable(output_zip):
  """Create a partition table image and store it in output_zip."""

  img = OutputFile(
      output_zip, OPTIONS.input_tmp, "IMAGES", "partition-table.img")
  bpt = OutputFile(
      output_zip, OPTIONS.input_tmp, "IMAGES", "partition-table.bpt")

  # use BPTTOOL from environ, or "bpttool" if empty or not set.
  bpttool = os.getenv("BPTTOOL") or "bpttool"
  cmd = [bpttool, "make_table", "--output_json", bpt.name,
         "--output_gpt", img.name]
  input_files_str = OPTIONS.info_dict["board_bpt_input_files"]
  input_files = input_files_str.split(" ")
  for i in input_files:
    cmd.extend(["--input", i])
  disk_size = OPTIONS.info_dict.get("board_bpt_disk_size")
  if disk_size:
    cmd.extend(["--disk_size", disk_size])
  args = OPTIONS.info_dict.get("board_bpt_make_table_args")
  if args:
    cmd.extend(shlex.split(args))

  p = common.Run(cmd, stdout=subprocess.PIPE)
  p.communicate()
  assert p.returncode == 0, "bpttool make_table failed"

  img.Write()
  bpt.Write()


def AddCache(output_zip):
  """Create an empty cache image and store it in output_zip."""

  img = OutputFile(output_zip, OPTIONS.input_tmp, "IMAGES", "cache.img")
  if os.path.exists(img.input_name):
    print("cache.img already exists; no need to rebuild...")
    return

  image_props = build_image.ImagePropFromGlobalDict(OPTIONS.info_dict, "cache")
  # The build system has to explicitly request for cache.img.
  if "fs_type" not in image_props:
    return

  print("creating cache.img...")

  # Use a fixed timestamp (01/01/2009) when packaging the image.
  # Bug: 24377993
  epoch = datetime.datetime.fromtimestamp(0)
  timestamp = (datetime.datetime(2009, 1, 1) - epoch).total_seconds()
  image_props["timestamp"] = int(timestamp)

  user_dir = common.MakeTempDir()

  fstab = OPTIONS.info_dict["fstab"]
  if fstab:
    image_props["fs_type"] = fstab["/cache"].fs_type
  succ = build_image.BuildImage(user_dir, image_props, img.name)
  assert succ, "build cache.img image failed"

  common.CheckSize(img.name, "cache.img", OPTIONS.info_dict)
  img.Write()


def AddRadioImagesForAbOta(output_zip, ab_partitions):
  """Adds the radio images needed for A/B OTA to the output file.

  It parses the list of A/B partitions, looks for the missing ones from RADIO/
  or VENDOR_IMAGES/ dirs, and copies them to IMAGES/ of the output file (or
  dir).

  It also ensures that on returning from the function all the listed A/B
  partitions must have their images available under IMAGES/.

  Args:
    output_zip: The output zip file (needs to be already open), or None to
        write images to OPTIONS.input_tmp/.
    ab_partitions: The list of A/B partitions.

  Raises:
    AssertionError: If it can't find an image.
  """
  for partition in ab_partitions:
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    img_name = partition.strip() + ".img"
    prebuilt_path = os.path.join(OPTIONS.input_tmp, "IMAGES", img_name)
    if os.path.exists(prebuilt_path):
      print("%s already exists, no need to overwrite..." % (img_name,))
      continue

    img_radio_path = os.path.join(OPTIONS.input_tmp, "RADIO", img_name)
    if os.path.exists(img_radio_path):
      if output_zip:
        common.ZipWrite(output_zip, img_radio_path, "IMAGES/" + img_name)
      else:
        shutil.copy(img_radio_path, prebuilt_path)
      continue

    # Walk through VENDOR_IMAGES/ since files could be under subdirs.
    img_vendor_dir = os.path.join(OPTIONS.input_tmp, "VENDOR_IMAGES")
    for root, _, files in os.walk(img_vendor_dir):
      if img_name in files:
        if output_zip:
          common.ZipWrite(output_zip, os.path.join(root, img_name),
                          "IMAGES/" + img_name)
        else:
          shutil.copy(os.path.join(root, img_name), prebuilt_path)
        break
=======
    img_name = partition + ".img"
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

    # Assert that the image is present under IMAGES/ now.
    if output_zip:
      # Zip spec says: All slashes MUST be forward slashes.
      img_path = 'IMAGES/' + img_name
      assert img_path in output_zip.namelist(), "cannot find " + img_name
    else:
      img_path = os.path.join(OPTIONS.input_tmp, "IMAGES", img_name)
      assert os.path.exists(img_path), "cannot find " + img_name


<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
def AddCareMapTxtForAbOta(output_zip, ab_partitions, image_paths):
  """Generates and adds care_map.txt for system and vendor partitions.

  Args:
    output_zip: The output zip file (needs to be already open), or None to
        write images to OPTIONS.input_tmp/.
    ab_partitions: The list of A/B partitions.
    image_paths: A map from the partition name to the image path.
  """
  care_map_list = []
  for partition in ab_partitions:
    partition = partition.strip()
    if partition not in PARTITIONS_WITH_CARE_MAP:
      continue

    verity_block_device = "{}_verity_block_device".format(partition)
    avb_hashtree_enable = "avb_{}_hashtree_enable".format(partition)
    if (verity_block_device in OPTIONS.info_dict or
        OPTIONS.info_dict.get(avb_hashtree_enable) == "true"):
      image_path = image_paths[partition]
      assert os.path.exists(image_path)
      care_map_list += GetCareMap(partition, image_path)

  if care_map_list:
    care_map_path = "META/care_map.txt"
    if output_zip and care_map_path not in output_zip.namelist():
      common.ZipWriteStr(output_zip, care_map_path, '\n'.join(care_map_list))
    else:
      with open(os.path.join(OPTIONS.input_tmp, care_map_path), 'w') as fp:
        fp.write('\n'.join(care_map_list))
      if output_zip:
        OPTIONS.replace_updated_files_list.append(care_map_path)


=======
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
def AddPackRadioImages(output_zip, images):
  """Copies images listed in META/pack_radioimages.txt from RADIO/ to IMAGES/.

  Args:
    output_zip: The output zip file (needs to be already open), or None to
        write images to OPTIONS.input_tmp/.
    images: A list of image names.

  Raises:
    AssertionError: If a listed image can't be found.
  """
  for image in images:
    img_name = image.strip()
    _, ext = os.path.splitext(img_name)
    if not ext:
      img_name += ".img"

    prebuilt_path = os.path.join(OPTIONS.input_tmp, "IMAGES", img_name)
    if os.path.exists(prebuilt_path):
      print("%s already exists, no need to overwrite..." % (img_name,))
      continue

    img_radio_path = os.path.join(OPTIONS.input_tmp, "RADIO", img_name)
    assert os.path.exists(img_radio_path), \
        "Failed to find %s at %s" % (img_name, img_radio_path)

    if output_zip:
      common.ZipWrite(output_zip, img_radio_path, "IMAGES/" + img_name)
    else:
      shutil.copy(img_radio_path, prebuilt_path)


def ReplaceUpdatedFiles(zip_filename, files_list):
  """Updates all the ZIP entries listed in files_list.

  For now the list includes META/care_map.txt, and the related files under
  SYSTEM/ after rebuilding recovery.
  """
  common.ZipDelete(zip_filename, files_list)
  output_zip = zipfile.ZipFile(zip_filename, "a",
                               compression=zipfile.ZIP_DEFLATED,
                               allowZip64=True)
  for item in files_list:
    file_path = os.path.join(OPTIONS.input_tmp, item)
    assert os.path.exists(file_path)
    common.ZipWrite(output_zip, file_path, arcname=item)
  common.ZipClose(output_zip)


<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
def HasPartition(partition_name):
  """Determines if the target files archive should build a given partition."""

  return ((os.path.isdir(
      os.path.join(OPTIONS.input_tmp, partition_name.upper())) and
           OPTIONS.info_dict.get(
               "building_{}_image".format(partition_name)) == "true") or
          os.path.exists(
              os.path.join(OPTIONS.input_tmp, "IMAGES",
                           "{}.img".format(partition_name))))


def AddApexInfo(output_zip):
  apex_infos = GetApexInfoFromTargetFiles(OPTIONS.input_tmp, 'system',
                                          compressed_only=False)
  apex_metadata_proto = ota_metadata_pb2.ApexMetadata()
  apex_metadata_proto.apex_info.extend(apex_infos)
  apex_info_bytes = apex_metadata_proto.SerializeToString()

  output_file = os.path.join(OPTIONS.input_tmp, "META", "apex_info.pb")
  with open(output_file, "wb") as ofile:
    ofile.write(apex_info_bytes)
  if output_zip:
    arc_name = "META/apex_info.pb"
    if arc_name in output_zip.namelist():
      OPTIONS.replace_updated_files_list.append(arc_name)
    else:
      common.ZipWrite(output_zip, output_file, arc_name)


def AddVbmetaDigest(output_zip):
  """Write the vbmeta digest to the output dir and zipfile."""

  # Calculate the vbmeta digest and put the result in to META/
  boot_images = OPTIONS.info_dict.get("boot_images")
  # Disable the digest calculation if the target_file is used as a container
  # for boot images. A boot container might contain boot-5.4.img, boot-5.10.img
  # etc., instead of just a boot.img and will fail in vbmeta digest calculation.
  boot_container = boot_images and (
      len(boot_images.split()) >= 2 or boot_images.split()[0] != 'boot.img')
  if (OPTIONS.info_dict.get("avb_enable") == "true" and not boot_container and
      OPTIONS.info_dict.get("avb_building_vbmeta_image") == "true"):
    avbtool = OPTIONS.info_dict["avb_avbtool"]
    digest = verity_utils.CalculateVbmetaDigest(OPTIONS.input_tmp, avbtool)
    vbmeta_digest_txt = os.path.join(OPTIONS.input_tmp, "META",
                                     "vbmeta_digest.txt")
    with open(vbmeta_digest_txt, 'w') as f:
      f.write(digest)
    # writes to the output zipfile
    if output_zip:
      arc_name = "META/vbmeta_digest.txt"
      if arc_name in output_zip.namelist():
        OPTIONS.replace_updated_files_list.append(arc_name)
      else:
        common.ZipWriteStr(output_zip, arc_name, digest)


>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
def AddImagesToTargetFiles(filename):
  """Creates and adds images (boot/recovery/system/...) to a target_files.zip.

  It works with either a zip file (zip mode), or a directory that contains the
  files to be packed into a target_files.zip (dir mode). The latter is used when
  being called from build/make/core/Makefile.

  The images will be created under IMAGES/ in the input target_files.zip.

  Args:
    filename: the target_files.zip, or the zip root directory.
  """
  if os.path.isdir(filename):
    OPTIONS.input_tmp = os.path.abspath(filename)
  else:
    OPTIONS.input_tmp = common.UnzipTemp(filename)

  if not OPTIONS.add_missing:
    if os.path.isdir(os.path.join(OPTIONS.input_tmp, "IMAGES")):
      print("target_files appears to already contain images.")
      sys.exit(1)

  OPTIONS.info_dict = common.LoadInfoDict(OPTIONS.input_tmp, OPTIONS.input_tmp)

  has_recovery = OPTIONS.info_dict.get("no_recovery") != "true"

  # {vendor,product}.img is unlike system.img or system_other.img. Because it
  # could be built from source, or dropped into target_files.zip as a prebuilt
  # blob. We consider either of them as {vendor,product}.img being available,
  # which could be used when generating vbmeta.img for AVB.
  has_vendor = (os.path.isdir(os.path.join(OPTIONS.input_tmp, "VENDOR")) or
                os.path.exists(os.path.join(OPTIONS.input_tmp, "IMAGES",
                                            "vendor.img")))
  has_product = (os.path.isdir(os.path.join(OPTIONS.input_tmp, "PRODUCT")) or
                 os.path.exists(os.path.join(OPTIONS.input_tmp, "IMAGES",
                                             "product.img")))
  has_system_other = os.path.isdir(os.path.join(OPTIONS.input_tmp,
                                                "SYSTEM_OTHER"))

  # Set up the output destination. It writes to the given directory for dir
  # mode; otherwise appends to the given ZIP.
  if os.path.isdir(filename):
    output_zip = None
  else:
    output_zip = zipfile.ZipFile(filename, "a",
                                 compression=zipfile.ZIP_DEFLATED,
                                 allowZip64=True)

  # Always make input_tmp/IMAGES available, since we may stage boot / recovery
  # images there even under zip mode. The directory will be cleaned up as part
  # of OPTIONS.input_tmp.
  images_dir = os.path.join(OPTIONS.input_tmp, "IMAGES")
  if not os.path.isdir(images_dir):
    os.makedirs(images_dir)

  # A map between partition names and their paths, which could be used when
  # generating AVB vbmeta image.
  partitions = dict()

  def banner(s):
    print("\n\n++++ " + s + " ++++\n\n")

  banner("boot")
  # common.GetBootableImage() returns the image directly if present.
  boot_image = common.GetBootableImage(
      "IMAGES/boot.img", "boot.img", OPTIONS.input_tmp, "BOOT")
  # boot.img may be unavailable in some targets (e.g. aosp_arm64).
  if boot_image:
    partitions['boot'] = os.path.join(OPTIONS.input_tmp, "IMAGES", "boot.img")
    if not os.path.exists(partitions['boot']):
      boot_image.WriteToDir(OPTIONS.input_tmp)
      if output_zip:
        boot_image.AddToZip(output_zip)

  recovery_image = None
  if has_recovery:
    banner("recovery")
    recovery_image = common.GetBootableImage(
        "IMAGES/recovery.img", "recovery.img", OPTIONS.input_tmp, "RECOVERY")
    assert recovery_image, "Failed to create recovery.img."
    partitions['recovery'] = os.path.join(
        OPTIONS.input_tmp, "IMAGES", "recovery.img")
    if not os.path.exists(partitions['recovery']):
      recovery_image.WriteToDir(OPTIONS.input_tmp)
      if output_zip:
        recovery_image.AddToZip(output_zip)

      banner("recovery (two-step image)")
      # The special recovery.img for two-step package use.
      recovery_two_step_image = common.GetBootableImage(
          "IMAGES/recovery-two-step.img", "recovery-two-step.img",
          OPTIONS.input_tmp, "RECOVERY", two_step_image=True)
      assert recovery_two_step_image, "Failed to create recovery-two-step.img."
      recovery_two_step_image_path = os.path.join(
          OPTIONS.input_tmp, "IMAGES", "recovery-two-step.img")
      if not os.path.exists(recovery_two_step_image_path):
        recovery_two_step_image.WriteToDir(OPTIONS.input_tmp)
        if output_zip:
          recovery_two_step_image.AddToZip(output_zip)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  banner("system")
  partitions['system'] = AddSystem(
      output_zip, recovery_img=recovery_image, boot_img=boot_image)
=======
  def add_partition(partition, has_partition, add_func, add_args):
    if has_partition:
      banner(partition)
      partitions[partition] = add_func(output_zip, *add_args)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  if has_vendor:
    banner("vendor")
    partitions['vendor'] = AddVendor(output_zip)
=======
  add_partition_calls = (
      ("system", has_system, AddSystem, [recovery_image, boot_image]),
      ("vendor", has_vendor, AddVendor, [recovery_image, boot_image]),
      ("product", has_product, AddProduct, []),
      ("system_ext", has_system_ext, AddSystemExt, []),
      ("odm", has_odm, AddOdm, []),
      ("vendor_dlkm", has_vendor_dlkm, AddVendorDlkm, []),
      ("odm_dlkm", has_odm_dlkm, AddOdmDlkm, []),
      ("system_other", has_system_other, AddSystemOther, []),
  )
  for call in add_partition_calls:
    add_partition(*call)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  if has_product:
    banner("product")
    partitions['product'] = AddProduct(output_zip)

  if has_system_other:
    banner("system_other")
    AddSystemOther(output_zip)
=======
  AddApexInfo(output_zip)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

  if not OPTIONS.is_signing:
    banner("userdata")
    AddUserdata(output_zip)
    banner("cache")
    AddCache(output_zip)

  if OPTIONS.info_dict.get("board_bpt_enable") == "true":
    banner("partition-table")
    AddPartitionTable(output_zip)

  add_partition("dtbo",
                OPTIONS.info_dict.get("has_dtbo") == "true", AddDtbo, [])
  add_partition("pvmfw",
                OPTIONS.info_dict.get("has_pvmfw") == "true", AddPvmfw, [])

  if OPTIONS.info_dict.get("avb_enable") == "true":
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    banner("vbmeta")
    AddVBMeta(output_zip, partitions)
=======
    # vbmeta_partitions includes the partitions that should be included into
    # top-level vbmeta.img, which are the ones that are not included in any
    # chained VBMeta image plus the chained VBMeta images themselves.
    # Currently custom_partitions are all chained to VBMeta image.
    vbmeta_partitions = common.AVB_PARTITIONS[:] + tuple(custom_partitions)

    vbmeta_system = OPTIONS.info_dict.get("avb_vbmeta_system", "").strip()
    if vbmeta_system:
      banner("vbmeta_system")
      partitions["vbmeta_system"] = AddVBMeta(
          output_zip, partitions, "vbmeta_system", vbmeta_system.split())
      vbmeta_partitions = [
          item for item in vbmeta_partitions
          if item not in vbmeta_system.split()]
      vbmeta_partitions.append("vbmeta_system")

    vbmeta_vendor = OPTIONS.info_dict.get("avb_vbmeta_vendor", "").strip()
    if vbmeta_vendor:
      banner("vbmeta_vendor")
      partitions["vbmeta_vendor"] = AddVBMeta(
          output_zip, partitions, "vbmeta_vendor", vbmeta_vendor.split())
      vbmeta_partitions = [
          item for item in vbmeta_partitions
          if item not in vbmeta_vendor.split()]
      vbmeta_partitions.append("vbmeta_vendor")

    if OPTIONS.info_dict.get("avb_building_vbmeta_image") == "true":
      banner("vbmeta")
      AddVBMeta(output_zip, partitions, "vbmeta", vbmeta_partitions)

  if OPTIONS.info_dict.get("use_dynamic_partitions") == "true":
    if OPTIONS.info_dict.get("build_super_empty_partition") == "true":
      banner("super_empty")
      AddSuperEmpty(output_zip)

  if OPTIONS.info_dict.get("build_super_partition") == "true":
    if OPTIONS.info_dict.get(
        "build_retrofit_dynamic_partitions_ota_package") == "true":
      banner("super split images")
      AddSuperSplit(output_zip)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

  banner("radio")
  ab_partitions_txt = os.path.join(OPTIONS.input_tmp, "META",
                                   "ab_partitions.txt")
  if os.path.exists(ab_partitions_txt):
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    with open(ab_partitions_txt, 'r') as f:
      ab_partitions = f.readlines()
=======
    with open(ab_partitions_txt) as f:
      ab_partitions = f.read().splitlines()
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

    # For devices using A/B update, copy over images from RADIO/ and/or
    # VENDOR_IMAGES/ to IMAGES/ and make sure we have all the needed
    # images ready under IMAGES/. All images should have '.img' as extension.
    AddRadioImagesForAbOta(output_zip, ab_partitions)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
    # Generate care_map.txt for system and vendor partitions (if present), then
    # write this file to target_files package.
    AddCareMapTxtForAbOta(output_zip, ab_partitions, partitions)
=======
    # Generate care_map.pb for ab_partitions, then write this file to
    # target_files package.
    output_care_map = os.path.join(OPTIONS.input_tmp, "META", "care_map.pb")
    AddCareMapForAbOta(output_zip if output_zip else output_care_map,
                       ab_partitions, partitions)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

  # Radio images that need to be packed into IMAGES/, and product-img.zip.
  pack_radioimages_txt = os.path.join(
      OPTIONS.input_tmp, "META", "pack_radioimages.txt")
  if os.path.exists(pack_radioimages_txt):
    with open(pack_radioimages_txt, 'r') as f:
      AddPackRadioImages(output_zip, f.readlines())

  AddVbmetaDigest(output_zip)

  if output_zip:
    common.ZipClose(output_zip)
    if OPTIONS.replace_updated_files_list:
      ReplaceUpdatedFiles(output_zip.filename,
                          OPTIONS.replace_updated_files_list)


def main(argv):
  def option_handler(o, a):
    if o in ("-a", "--add_missing"):
      OPTIONS.add_missing = True
    elif o in ("-r", "--rebuild_recovery",):
      OPTIONS.rebuild_recovery = True
    elif o == "--replace_verity_private_key":
      OPTIONS.replace_verity_private_key = (True, a)
    elif o == "--replace_verity_public_key":
      OPTIONS.replace_verity_public_key = (True, a)
    elif o == "--is_signing":
      OPTIONS.is_signing = True
    else:
      return False
    return True

  args = common.ParseOptions(
      argv, __doc__, extra_opts="ar",
      extra_long_opts=["add_missing", "rebuild_recovery",
                       "replace_verity_public_key=",
                       "replace_verity_private_key=",
                       "is_signing"],
      extra_option_handler=option_handler)


  if len(args) != 1:
    common.Usage(__doc__)
    sys.exit(1)

  AddImagesToTargetFiles(args[0])
  print("done.")

if __name__ == '__main__':
  try:
    common.CloseInheritedPipes()
    main(sys.argv[1:])
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  except common.ExternalError as e:
    print("\n   ERROR: %s\n" % (e,))
    sys.exit(1)
=======
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
  finally:
    common.Cleanup()
