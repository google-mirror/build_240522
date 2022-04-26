import argparse
from curses import meta
import sys
import update_payload
from update_payload import Payload
import ota_metadata_pb2
import tempfile
from typing import BinaryIO, List
from update_metadata_pb2 import DeltaArchiveManifest, DynamicPartitionMetadata, DynamicPartitionGroup


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
  for op in partitions[-len(new_partitions):].operations:
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
    group.size += new_group.size
    del new_groups[group.name]
  for new_group in new_groups.values():
    groups.append(new_group)


def MergeDynamicPartitionMetadata(metadata: DynamicPartitionMetadata, new_metadata: DynamicPartitionMetadata):
  MergeDynamicPartitionGroups(metadata.groups, new_metadata.groups)
  metadata.snapshot_enabled &= new_metadata.snapshot_enabled
  metadata.vabc_enabled &= new_metadata.vabc_enabled
  assert metadata.vabc_compression_param == new_metadata.vabc_compression_param
  metadata.cow_version = max(metadata.cow_version, new_metadata.cow_version)


def MergeManifests(payloads: List[Payload]) -> DeltaArchiveManifest:
  # TODO(zhangkelvin) Implement manifest merging
  output_manifest = DeltaArchiveManifest()
  output_manifest.block_size = payloads[0].manifest.block_size
  output_manifest.partial_update = True
  for payload in payloads:
    manifest = payload.manifest
    assert manifest.block_size == output_manifest.block_size
    output_manifest.minor_version = max(
        output_manifest.minor_version, manifest.minor_version)
    output_manifest.max_timestamp = max(
        output_manifest.max_timestamp, manifest.max_timestamp)
    output_manifest.apex_info.extend(manifest.apex_info)
    ExtendPartitionUpdates(output_manifest.partitions, manifest.partitions)
    MergeDynamicPartitionMetadata(
        output_manifest.dynamic_partition_metadata, manifest.dynamic_partition_metadata)

  return output_manifest


def MergePayloads(payloads: List[Payload]):
  with tempfile.NamedTemporaryFile(prefix="payload_blob") as tmpfile:
    ConcatBlobs(payloads, tmpfile)


def main():
  parser = argparse.ArgumentParser(description='Merge multiple partial OTAs')
  parser.add_argument('packages', type=str, nargs='+',
                      help='Paths to OTA packages to merge')
  args = parser.parse_args()
  print(args)


if __name__ == '__main__':
  sys.exit(main())
