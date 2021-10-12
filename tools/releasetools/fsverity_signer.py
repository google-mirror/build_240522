#!/usr/bin/env python
#
# Copyright 2021 Google Inc. All rights reserved.
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

"""
`fsverity_signer` generates fsverity metadata and signature to a container file

This actually is a simple wrapper around the `fsverity` program. A file is
signed by the program which produces the PKCS#7 signature file, merkle tree file
, and the fsverity_descriptor file. Then the files are packed into a single
output file so that the information about the signing stays together.

Currently, the output of this script is used by `fd_server` which is the host-
side backend of an authfs filesystem. `fd_server` uses this file in case when
the underlying filesystem (ext4, etc.) on the device doesn't support the
fsverity feature natively in which case the information is read directly from
the filesystem using ioctl.
"""

import argparse
import os
import re
import shutil
import subprocess
import sys
import tempfile
from struct import *

class TempDirectory(object):
  def __enter__(self):
    self.name = tempfile.mkdtemp()
    return self.name

  def __exit__(self, *unused):
    shutil.rmtree(self.name)

class FSVeritySigner:
  def __init__(self, fsverity_path):
    self._fsverity_path = fsverity_path

    # Default values for some properties
    self.set_hash_alg("sha256")
    self.set_raw_signature(False)

  def set_key(self, key):
    self._key = key

  def set_cert(self, cert):
    self._cert = cert

  def set_hash_alg(self, hash_alg):
    self._hash_alg = hash_alg

  def set_raw_signature(self, raw_signature):
    self._raw_signature = raw_signature

  def _raw_signature(pkcs7_sig_file):
    """ Extracts raw signature from DER formatted PKCS#7 detached signature file

    Do that by parsing the ASN.1 tree to get the location of the signature
    in the file and then read the portion.
    """

    # Note: there seems to be no public python API (even in 3p modules) that
    # provides direct access to the raw signature at this moment. So, `openssl
    # asn1parse` commandline tool is used instead.
    cmd = ['openssl', 'asn1parse']
    cmd.extend(['-inform', 'DER'])
    cmd.extend(['-in', pkcs7_sig_file])
    out = subprocess.check_output(cmd, universal_newlines=True)

    # The signature is the last element in the tree
    last_line = out.splitlines()[-1]
    m = re.search('(\d+):.*hl=\s*(\d+)\s*l=\s*(\d+)\s*.*OCTET STRING', last_line)
    if not m:
      raise RuntimeError("Failed to parse asn1parse output: " + out)
    offset = int(m.group(1))
    header_len = int(m.group(2))
    size = int(m.group(3))
    with open(pkcs7_sig_file, 'rb') as f:
      f.seek(offset + header_len)
      return f.read(size)

  def sign(self, input_file, output_file=None):
    if not self._key:
      raise RuntimeError("key must be specified.")

    if not self._cert:
      raise RuntimeError("cert must be specified.")

    if not output_file:
      output_file = input_file + '.fsv_meta'

    with TempDirectory() as temp_dir:
      self._do_sign(input_file, output_file, temp_dir)

  def _do_sign(self, input_file, output_file, work_dir):
    # temporary files
    desc_file = os.path.join(work_dir, 'desc')
    merkletree_file = os.path.join(work_dir, 'merkletree')
    sig_file = os.path.join(work_dir, 'signature')

    # convert DER private key to PEM
    pem_key = os.path.join(work_dir, 'key.pem')
    cmd = ['openssl', 'pkcs8']
    cmd.extend(['-inform', 'DER'])
    cmd.extend(['-in', self._key])
    cmd.extend(['-nocrypt'])
    cmd.extend(['-out', pem_key])
    subprocess.check_call(cmd)

    # run the fsverity util to create the temporary files
    cmd = [self._fsverity_path, 'sign']
    cmd.append(input_file)
    cmd.append(sig_file)
    cmd.extend(['--key', pem_key])
    cmd.extend(['--cert', self._cert])
    cmd.extend(['--hash-alg', self._hash_alg])
    cmd.extend(['--block-size', '4096'])
    cmd.extend(['--out-merkle-tree', merkletree_file])
    cmd.extend(['--out-descriptor', desc_file])
    subprocess.check_call(cmd, stdout=open(os.devnull, 'w'))

    with open(output_file, 'wb') as out:
      # 1. version
      out.write(pack('<I', 1))

      # 2. fsverity_descriptor
      with open(desc_file, 'rb') as f:
        out.write(f.read())

      # 3. signature
      SIG_TYPE_PKCS7 = 1
      SIG_TYPE_RAW = 2
      if self._raw_signature:
        out.write(pack('<I', SIG_TYPE_RAW))
        sig = self.raw_signature(sig_file)
        out.write(pack('<I', len(sig)))
        out.write(sig)
      else:
        with open(sig_file, 'rb') as f:
          out.write(pack('<I', SIG_TYPE_PKCS7))
          sig = f.read()
          out.write(pack('<I', len(sig)))
          out.write(sig)

      # 4. merkle tree
      with open(merkletree_file, 'rb') as f:
        # merkle tree is placed at the next nearest page boundary to make
        # mmapping possible
        out.seek(next_page(out.tell()))
        out.write(f.read())

def next_page(n):
  """ Returns the next nearest page boundary from `n` """
  PAGE_SIZE = 4096
  return (n + PAGE_SIZE - 1) // PAGE_SIZE * PAGE_SIZE

if __name__ == '__main__':
  p = argparse.ArgumentParser()
  p.add_argument(
      '--output',
      help='output file. If omitted, print to <INPUT>.fsv_meta',
      metavar='output',
      default=None)
  p.add_argument(
      'input',
      help='input file to be signed')
  p.add_argument(
      '--key',
      help='PKCS#8 private key file in DER format',
      required=True)
  p.add_argument(
      '--cert',
      help='x509 certificate file in PEM format',
      required=True)
  p.add_argument(
      '--hash-alg',
      help='hash algorithm to use to build the merkle tree',
      choices=['sha256', 'sha512'],
      default='sha256')
  p.add_argument(
      '--raw-signature',
      help='raw signature is stored instead of PKCS#7 format',
      action='store_true')
  p.add_argument(
      '--fsverity-path',
      help='path to the fsverity program',
      required=True)
  args = p.parse_args(sys.argv[1:])

  signer = FSVeritySigner(args.fsverity_path)
  signer.set_key(args.key)
  signer.set_cert(args.cert)
  signer.set_hash_alg(args.hash_alg)
  signer.set_raw_signature(args.raw_signature)

  signer.sign(args.input, args.output)
