import argparse
import struct
import sys
import update_payload
from update_payload import Payload
import ota_metadata_pb2
import tempfile
import zipfile
import os
import care_map_pb2

import common
from typing import BinaryIO, List
from update_metadata_pb2 import DeltaArchiveManifest, DynamicPartitionMetadata, DynamicPartitionGroup
from ota_metadata_pb2 import OtaMetadata

from payload_signer import PayloadSigner
from ota_utils import PayloadGenerator, METADATA_PROTO_NAME, FinalizeMetadata

CARE_MAP_ENTRY = "care_map.pb"

def WriteDataBlob(payload: Payload, outfp: BinaryIO, read_size=1024*64):
  for i in range(0, payload.total_data_length, read_size):
    blob = payload.ReadDataBlob(
        i, min(i+read_size, payload.total_data_length)-i)
    outfp.write(blob)


def ConcatBlobs(payloads: List[Payload], outfp: BinaryIO):
  for payload in payloads:
    WriteDataBlob(payload, outfp)


def TotalDataLength(partitions):
  for partition in reversed(partitions):
    for op in reversed(partition.operations):
      if op.data_offset > 0:
        return op.data_offset + op.data_length
  return 0


def ExtendPartitionUpdates(partitions, new_partitions):
  prefix_blob_length = TotalDataLength(partitions)
  partitions.extend(new_partitions)
  for part in partitions[-len(new_partitions):]:
    for op in part.operations:
      if op.HasField("data_length") and op.data_length != 0:
        op.data_offset += prefix_blob_length


def MergeDynamicPartitionGroups(groups: List[DynamicPartitionGroup], new_groups: List[DynamicPartitionGroup]):
  new_groups = {new_group.name: new_group for new_group in new_groups}
  for group in groups:
    if group.name not in new_groups:
      continue
    new_group = new_groups[group.name]
    assert set(group.partition_names).intersection(set(new_group.partition_names)) == set(
    ), "Old group and new group should not have any intersections"
    group.partition_names.extend(new_group.partition_names)
    group.size = max(new_group.size, group.size)
    del new_groups[group.name]
  for new_group in new_groups.values():
    groups.append(new_group)


def MergeDynamicPartitionMetadata(metadata: DynamicPartitionMetadata, new_metadata: DynamicPartitionMetadata):
  MergeDynamicPartitionGroups(metadata.groups, new_metadata.groups)
  metadata.snapshot_enabled &= new_metadata.snapshot_enabled
  metadata.vabc_enabled &= new_metadata.vabc_enabled
  assert metadata.vabc_compression_param == new_metadata.vabc_compression_param, f"{metadata.vabc_compression_param} vs. {new_metadata.vabc_compression_param}"
  metadata.cow_version = max(metadata.cow_version, new_metadata.cow_version)


def MergeManifests(payloads: List[Payload]) -> DeltaArchiveManifest:
  if len(payloads) == 0:
    return None
  if len(payloads) == 1:
    return payloads[0].manifest

  output_manifest = DeltaArchiveManifest()
  output_manifest.block_size = payloads[0].manifest.block_size
  output_manifest.partial_update = True
  output_manifest.dynamic_partition_metadata.snapshot_enabled = payloads[
      0].manifest.dynamic_partition_metadata.snapshot_enabled
  output_manifest.dynamic_partition_metadata.vabc_enabled = payloads[
      0].manifest.dynamic_partition_metadata.vabc_enabled
  output_manifest.dynamic_partition_metadata.vabc_compression_param = payloads[
      0].manifest.dynamic_partition_metadata.vabc_compression_param
  apex_info = {}
  for payload in payloads:
    manifest = payload.manifest
    assert manifest.block_size == output_manifest.block_size
    output_manifest.minor_version = max(
        output_manifest.minor_version, manifest.minor_version)
    output_manifest.max_timestamp = max(
        output_manifest.max_timestamp, manifest.max_timestamp)
    output_manifest.apex_info.extend(manifest.apex_info)
    for apex in manifest.apex_info:
      apex_info[apex.package_name] = apex
    ExtendPartitionUpdates(output_manifest.partitions, manifest.partitions)
    MergeDynamicPartitionMetadata(
        output_manifest.dynamic_partition_metadata, manifest.dynamic_partition_metadata)

  for apex_name in sorted(apex_info.keys()):
    output_manifest.apex_info.extend(apex_info[apex_name])

  return output_manifest


def MergePayloads(payloads: List[Payload]):
  with tempfile.NamedTemporaryFile(prefix="payload_blob") as tmpfile:
    ConcatBlobs(payloads, tmpfile)


def MergeCareMap(paths: List[str]):
  care_map = care_map_pb2.CareMap()
  for path in paths:
    with zipfile.ZipFile(path, "r", allowZip64=True) as zfp:
      if CARE_MAP_ENTRY in zfp.namelist():
        care_map_bytes = zfp.read(CARE_MAP_ENTRY)
        partial_care_map = care_map_pb2.CareMap()
        partial_care_map.ParseFromString(care_map_bytes)
        care_map.partitions.extend(partial_care_map.partitions)
  return care_map.SerializeToString()


def WriteHeaderAndManifest(manifest: DeltaArchiveManifest, fp: BinaryIO):
  __MAGIC = b"CrAU"
  __MAJOR_VERSION = 2
  manifest_bytes = manifest.SerializeToString()
  fp.write(struct.pack(f">4sQQL", __MAGIC,
           __MAJOR_VERSION, len(manifest_bytes), 0))
  fp.write(manifest_bytes)


def AddOtaMetadata(input_ota, metadata_ota, output_ota, package_key, pw):
  with zipfile.ZipFile(metadata_ota, 'r') as zfp:
    metadata = OtaMetadata()
    metadata.ParseFromString(zfp.read(METADATA_PROTO_NAME))
    FinalizeMetadata(metadata, input_ota, output_ota,
                     package_key=package_key, pw=pw, no_signing=package_key is None)
    return output_ota


def CheckOutput(output_ota):
  payload = update_payload.Payload(output_ota)
  payload.CheckOpDataHash()

def main():
  parser = argparse.ArgumentParser(description='Merge multiple partial OTAs')
  parser.add_argument('packages', type=str, nargs='+',
                      help='Paths to OTA packages to merge')
  parser.add_argument('--package_key', type=str,
                      help='Paths to private key for signing payload')
  parser.add_argument('--search_path', type=str,
                      help='Search path for framework/signapk.jar')
  parser.add_argument('--output', type=str,
                      help='Paths to output merged ota', required=True)
  parser.add_argument('--metadata_ota', type=str,
                      help='Output zip will use build metadata from this OTA package, if unspecified, use the last OTA package in merge list')
  parser.add_argument('--private_key_suffix', type=str,
                      help='Suffix to be appended to package_key path', default=".pk8")
  args = parser.parse_args()
  file_paths = args.packages
  print(args)

  common.OPTIONS.search_path = args.search_path

  metadata_ota = args.packages[-1]
  if args.metadata_ota is not None:
    metadata_ota = args.metadata_ota
    assert os.path.exists(metadata_ota)

  payloads = [Payload(path) for path in file_paths]
  merged_manifest = MergeManifests(payloads)

  with tempfile.NamedTemporaryFile() as unsigned_payload:
    WriteHeaderAndManifest(merged_manifest, unsigned_payload)
    ConcatBlobs(payloads, unsigned_payload)
    unsigned_payload.flush()

    generator = PayloadGenerator()
    generator.payload_file = unsigned_payload.name
    print("Payload size:", os.path.getsize(generator.payload_file))

    if args.package_key:
      print("Signing payload...")
      signer = PayloadSigner(args.package_key, args.private_key_suffix)
      signed_payload = signer.SignPayload(unsigned_payload.name)
      generator.payload_file = signed_payload

    print("Payload size:", os.path.getsize(generator.payload_file))

    print("Writing to", args.output)
    key_passwords = common.GetKeyPasswords([args.package_key])
    with tempfile.NamedTemporaryFile(prefix="signed_ota", suffix=".zip") as signed_ota:
      with zipfile.ZipFile(signed_ota, "w") as zfp:
        generator.WriteToZip(zfp)
        zfp.writestr(CARE_MAP_ENTRY, MergeCareMap(args.packages))
      AddOtaMetadata(signed_ota.name, metadata_ota,
                     args.output, args.package_key, key_passwords[args.package_key])




if __name__ == '__main__':
  sys.exit(main())
