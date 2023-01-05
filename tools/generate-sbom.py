#!/usr/bin/env python3
#
# Copyright (C) 2023 The Android Open Source Project
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
Generate the SBOM of the current target product in SPDX format.
Usage example:
  generate-sbom.py --output_file out/target/product/vsoc_x86_64/sbom.spdx \
                   --metadata out/target/product/vsoc_x86_64/sbom-metadata.csv \
                   --product_out_dir=out/target/product/vsoc_x86_64 \
                   --build_version $(cat out/target/product/vsoc_x86_64/build_fingerprint.txt) \
                   --product_mfr=Google
"""

import argparse
import csv
import datetime
import google.protobuf.text_format as text_format
import hashlib
import os
import project_metadata_pb2

# Common
SPDXID = 'SPDXID'
CREATED = 'Created'
EXTERNAL_DOCUMENT_REF = 'ExternalDocumentRef'

# Package
PACKAGE_NAME = 'PackageName'
PACKAGE_DOWNLOAD_LOCATION = 'PackageDownloadLocation'
PACKAGE_VERSION = 'PackageVersion'
PACKAGE_SUPPLIER = 'PackageSupplier'
FILES_ANALYZED = 'FilesAnalyzed'
PACKAGE_VERIFICATION_CODE = 'PackageVerificationCode'
# Package license
PACKAGE_LICENSE_CONCLUDED = 'PackageLicenseConcluded'
PACKAGE_LICENSE_INFO_FROM_FILES = 'PackageLicenseInfoFromFiles'
PACKAGE_LICENSE_DECLARED = 'PackageLicenseDeclared'
PACKAGE_LICENSE_COMMENTS = 'PackageLicenseComments'

# File
FILE_NAME = 'FileName'
FILE_CHECKSUM = 'FileChecksum'
# File license
FILE_LICENSE_CONCLUDED = 'LicenseConcluded'
FILE_LICENSE_INFO_IN_FILE = 'LicenseInfoInFile'
FILE_LICENSE_COMMENTS = 'LicenseComments'
FILE_COPYRIGHT_TEXT = 'FileCopyrightText'
FILE_NOTICE = 'FileNotice'
FILE_ATTRIBUTION_TEXT = 'FileAttributionText'

# Relationship
REL_DESCRIBES = 'DESCRIBES'
REL_VARIANT_OF = 'VARIANT_OF'
REL_GENERATED_FROM = 'GENERATED_FROM'

# Package type
PKG_SOURCE = 'SOURCE'
PKG_UPSTREAM = 'UPSTREAM'
PKG_PREBUILT = 'PREBUILT'

# Report
ISSUE_NO_METADATA = 'No metadata generated in Make for installed files:'
ISSUE_NO_METADATA_FILE = 'No METADATA file found for installed file:'
ISSUE_METADATA_FILE_INCOMPLETE = 'METADATA file incomplete:'
INFO_METADATA_FOUND_FOR_PACKAGE = 'METADATA file found for packages:'


def get_args():
  parser = argparse.ArgumentParser()
  parser.add_argument('-v', '--verbose', action='store_true', default=False, help='Print more information.')
  parser.add_argument('--output_file', required=True, help='The generated SBOM file in SPDX format.')
  parser.add_argument('--metadata', required=True, help='The SBOM metadata file path.')
  parser.add_argument('--product_out_dir', required=True, help='The parent directory of all the installed files.')
  parser.add_argument('--build_version', required=True, help='The build version.')
  parser.add_argument('--product_mfr', required=True, help='The product manufacturer.')

  return parser.parse_args()


def log(*info):
  if args.verbose:
    for i in info:
      print(i)


def new_doc_header(doc_id):
  return {
      'SPDXVersion': 'SPDX-2.3',
      'DataLicense': 'CC0-1.0',
      SPDXID: doc_id,
      'DocumentName': args.build_version,
      'DocumentNamespace': '<document namespace here>',
      'Creator': 'Organization: Google, LLC',
      'Created': '<timestamp>',
      EXTERNAL_DOCUMENT_REF: [],
  }


def new_package_record(id, name, version, supplier, files_analyzed='false'):
  package = {
      PACKAGE_NAME: name,
      SPDXID: id,
      PACKAGE_DOWNLOAD_LOCATION: 'NONE',
      PACKAGE_VERSION: version,
      PACKAGE_SUPPLIER: 'Organization: ' + supplier,
      FILES_ANALYZED: files_analyzed,
  }
  return package


def new_file_record(id, name, checksum):
  return {
      FILE_NAME: name,
      SPDXID: id,
      FILE_CHECKSUM: checksum
  }


def encode_for_spdxid(s):
  """Simple encode for string values used in SPDXID which uses the charset of A-Za-Z0-9.-"""
  result = ''
  for c in s:
    if c.isalnum() or c in '.-':
      result += c
    elif c in '_@/':
      result += '-'
    else:
      result += '0x' + c.encode('utf-8').hex()

  return result.lstrip('-')


def new_package_id(package_name, type):
  return 'SPDXRef-{}-{}'.format(type, encode_for_spdxid(package_name))


def new_file_id(file_path):
  return 'SPDXRef-' + encode_for_spdxid(file_path)


def new_relationship_record(id1, relationship, id2):
  return 'Relationship {} {} {}'.format(id1, relationship, id2)


def checksum(file_path):
  file_path = args.product_out_dir + '/' + file_path
  h = hashlib.sha1()
  if os.path.islink(file_path):
    h.update(os.readlink(file_path).encode('utf-8'))
  else:
    with open(file_path, "rb") as f:
      h.update(f.read())
  return "SHA1: " + h.hexdigest()


def is_source_package(file_metadata):
  module_path = file_metadata['module_path']
  return module_path.startswith('external/') and module_path.find('chromium-webview') == -1


def is_prebuilt_package(file_metadata):
  module_path = file_metadata['module_path']
  if module_path:
    return module_path.startswith('prebuilts/') or module_path.find('chromium-webview') > -1

  kernel_module_copy_files = file_metadata['kernel_module_copy_files']
  if kernel_module_copy_files and not kernel_module_copy_files.startswith('ANDROID-GEN:'):
    return True

  return False


def get_prebuilt_package_name(file_metadata, name_in_metadata_file, metadata_file_path):
  module_path = file_metadata['module_path']
  kernel_module_copy_files = file_metadata['kernel_module_copy_files']
  name = None
  if name_in_metadata_file:
    name = name_in_metadata_file
  elif metadata_file_path:
    name = metadata_file_path
  elif module_path:
    name = module_path
  elif kernel_module_copy_files:
    src_path = kernel_module_copy_files.split(':')[0]
    name = os.path.dirname(src_path)

  return name.removeprefix('prebuilts/').replace('/', '-')


def get_metadata_file_path(file_metadata):
  metadata_path = ''
  if file_metadata['module_path']:
    metadata_path = file_metadata['module_path']
  elif file_metadata['kernel_module_copy_files']:
    metadata_path = os.path.dirname(file_metadata['kernel_module_copy_files'].split(':')[0])

  while metadata_path and not os.path.exists(metadata_path + '/METADATA'):
    metadata_path = os.path.dirname(metadata_path)

  return metadata_path


def get_sbom_fragments(installed_file_metadata, metadata_file_path):
  external_doc_ref = None
  packages = []
  relationships = []

  name_in_metadata_file = None
  if metadata_file_path:
    name_in_metadata_file = metadata_file_protos[metadata_file_path].name

  if is_source_package(installed_file_metadata):
    # Source fork packges
    # TODO: if module_path != metadata_file_path and module_path uses different license, the name should be derived from module_path so its specific license can be reported
    name = name_in_metadata_file if name_in_metadata_file else os.path.basename(
        metadata_file_path if metadata_file_path else installed_file_metadata['module_path'])
    source_package_id = new_package_id(name, PKG_SOURCE)
    source_package = new_package_record(source_package_id, name, args.build_version, args.product_mfr)
    # TODO: check if upstream package SBOM link exists. The upstream package is created when there is NO upstream package SBOM link.
    upstream_package_id = new_package_id(name, PKG_UPSTREAM)
    upstream_package = new_package_record(upstream_package_id, name, args.build_version, args.product_mfr)
    packages += [source_package, upstream_package]
    relationships.append(new_relationship_record(source_package_id, REL_VARIANT_OF, upstream_package_id))
  elif is_prebuilt_package(installed_file_metadata):
    # TODO: upstream package SBOM link should always exists, and add relationship prebuilt_package_id VARIANT_OF upstream SBOM
    # Prebuilt fork packages
    name = get_prebuilt_package_name(installed_file_metadata, name_in_metadata_file, metadata_file_path)
    prebuilt_package_id = new_package_id(name, PKG_PREBUILT)
    prebuilt_package = new_package_record(prebuilt_package_id, name, args.build_version, args.product_mfr)
    packages.append(prebuilt_package)

  return external_doc_ref, packages, relationships


def generate_package_verification_code(files):
  checksums = [file[FILE_CHECKSUM] for file in files]
  checksums.sort()
  h = hashlib.sha1()
  h.update(''.join(checksums).encode(encoding='utf-8'))
  return h.hexdigest()


def write_record(f, record):
  if record.__class__.__name__ == 'dict':
    for k, v in record.items():
      if k == EXTERNAL_DOCUMENT_REF:
        for ref in v:
          f.write('{}: {}\n'.format(k, ref))
      else:
        f.write('{}: {}\n'.format(k, v))
  elif record.__class__.__name__ == 'str':
    f.write(record + '\n')
  f.write('\n')


def save_report(report):
  prefix, _ = os.path.splitext(args.output_file)
  with open(prefix + '-gen-report.txt', 'w', encoding="utf-8") as report_file:
    for type, issues in report.items():
      report_file.write(type + '\n')
      for issue in issues:
        report_file.write('\t' + issue + '\n')
      report_file.write('\n')


def sort_rels(rel):
  # rel = 'Relationship file_id GENERATED_FROM package_id'
  fields = rel.split(' ')
  return fields[3] + fields[1]


# Validate the metadata generated by Make for installed files and report if there is no metadata.
def installed_file_has_metadata(installed_file_metadata, report):
  installed_file = installed_file_metadata['installed_file']
  module_path = installed_file_metadata['module_path']
  product_copy_files = installed_file_metadata['product_copy_files']
  kernel_module_copy_files = installed_file_metadata['kernel_module_copy_files']

  if (not module_path and
      not product_copy_files and
      not kernel_module_copy_files and
      not installed_file.endswith('.fsv_meta')):
    report[ISSUE_NO_METADATA].append(installed_file)
    return False

  return True


def report_metadata_file(metadata_file_path, installed_file_metadata, report):
  if metadata_file_path:
    report[INFO_METADATA_FOUND_FOR_PACKAGE].append(
        "installed_file: {}, module_path: {}, METADATA file: {}".format(
            installed_file_metadata['installed_file'],
            installed_file_metadata['module_path'],
            metadata_file_path + '/METADATA'))

    package_metadata = project_metadata_pb2.Metadata()
    with open(metadata_file_path + '/METADATA', "rt") as f:
      text_format.Parse(f.read(), package_metadata)

    if not metadata_file_path in metadata_file_protos:
      metadata_file_protos[metadata_file_path] = package_metadata
      if not package_metadata.name:
        report[ISSUE_METADATA_FILE_INCOMPLETE].append('{} does not has "name"'.format(metadata_file_path + '/METADATA'))
      if not package_metadata.version:
        report[ISSUE_METADATA_FILE_INCOMPLETE].append(
          '{} does not has "version"'.format(metadata_file_path + '/METADATA'))
  else:
    report[ISSUE_NO_METADATA_FILE].append(
        "installed_file: {}, module_path: {}".format(
            installed_file_metadata['installed_file'], installed_file_metadata['module_path']))


def main():
  global args
  args = get_args()
  log("Args:", vars(args))

  global metadata_file_protos
  metadata_file_protos = {}

  doc_id = 'SPDXRef-DOCUMENT'
  doc_header = new_doc_header(doc_id)

  product_package_id = 'SPDXRef-PRODUCT'
  product_package = new_package_record(product_package_id, 'PRODUCT', args.build_version, args.product_mfr,
                                       files_analyzed='true')

  platform_package_id = 'SPDXRef-PLATFORM'
  platform_package = new_package_record(platform_package_id, 'PLATFORM', args.build_version, args.product_mfr)

  # Report on some issues and information
  report = {
      ISSUE_NO_METADATA: [],
      ISSUE_NO_METADATA_FILE: [],
      ISSUE_METADATA_FILE_INCOMPLETE: [],
      INFO_METADATA_FOUND_FOR_PACKAGE: []
  }

  # Scan the metadata in CSV file and create the corresponding package and file records in SPDX
  product_files = []
  package_ids = []
  package_records = []
  rels_file_gen_from = []
  with open(args.metadata, newline='') as sbom_metadata_file:
    reader = csv.DictReader(sbom_metadata_file)
    for installed_file_metadata in reader:
      installed_file = installed_file_metadata['installed_file']
      module_path = installed_file_metadata['module_path']
      product_copy_files = installed_file_metadata['product_copy_files']
      kernel_module_copy_files = installed_file_metadata['kernel_module_copy_files']

      if not installed_file_has_metadata(installed_file_metadata, report):
        continue

      file_id = new_file_id(installed_file)
      product_files.append(new_file_record(file_id, installed_file, checksum(installed_file)))

      if is_source_package(installed_file_metadata) or is_prebuilt_package(installed_file_metadata):
        metadata_file_path = get_metadata_file_path(installed_file_metadata)
        report_metadata_file(metadata_file_path, installed_file_metadata, report)

        # File from source fork packages or prebuilt fork packages
        external_doc_ref, pkgs, rels = get_sbom_fragments(installed_file_metadata, metadata_file_path)
        if len(pkgs) > 0:
          if external_doc_ref:
            doc_header[EXTERNAL_DOCUMENT_REF].append(external_doc_ref)
          for p in pkgs:
            if not p[SPDXID] in package_ids:
              package_ids.append(p[SPDXID])
              package_records.append(p)
          for rel in rels:
            if not rel in package_records:
              package_records.append(rel)
          fork_package_id = pkgs[0][SPDXID]  # The first package should be the source/prebuilt fork package
          rels_file_gen_from.append(new_relationship_record(file_id, REL_GENERATED_FROM, fork_package_id))
      elif module_path:
        # File from platform package
        rels_file_gen_from.append(new_relationship_record(file_id, REL_GENERATED_FROM, platform_package_id))
      elif product_copy_files:
        # Format of product_copy_files: <source path>:<dest path>
        src_path = product_copy_files.split(':')[0]
        # So far product_copy_files are copied from directory system, kernel, hardware, frameworks and device,
        # so process them as files from platform package
        rels_file_gen_from.append(new_relationship_record(file_id, REL_GENERATED_FROM, platform_package_id))
      elif installed_file.endswith('.fsv_meta'):
        # See build/make/core/Makefile:2988
        rels_file_gen_from.append(new_relationship_record(file_id, REL_GENERATED_FROM, platform_package_id))
      elif kernel_module_copy_files.startswith('ANDROID-GEN'):
        # For the four files generated for _dlkm, _ramdisk partitions
        # See build/make/core/Makefile:323
        rels_file_gen_from.append(new_relationship_record(file_id, REL_GENERATED_FROM, platform_package_id))

  product_package[PACKAGE_VERIFICATION_CODE] = generate_package_verification_code(product_files)

  all_records = [
      doc_header,
      product_package,
      new_relationship_record(doc_id, REL_DESCRIBES, product_package_id),
  ]
  all_records += product_files
  all_records.append(platform_package)
  all_records += package_records
  rels_file_gen_from.sort(key=sort_rels)
  all_records += rels_file_gen_from

  # Save SBOM records to output file
  doc_header[CREATED] = datetime.datetime.now(tz=datetime.timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')
  with open(args.output_file, 'w', encoding="utf-8") as output_file:
    for rec in all_records:
      write_record(output_file, rec)

  save_report(report)


if __name__ == '__main__':
  main()
