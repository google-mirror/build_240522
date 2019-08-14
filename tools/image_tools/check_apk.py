#!/usr/bin/python
# -*- coding: utf-8 -*-

from __future__ import print_function, unicode_literals
import os
import zipfile
from apkverify import ApkSignature

def is_apk_file(file_path, validate=False):
    result = False
    errors = ''
    algs_for_zip_dict = {}

    if not os.path.exists(file_path):
        errors += 'apk check: file not found. file path - {0}\n'.format(file_path)
        return result, algs_for_zip_dict,errors

    try:
        if zipfile.is_zipfile(file_path):
            try:
                with zipfile.ZipFile(file_path) as arch:
                    path_files = arch.namelist()
            except Exception as error_message:
                errors += 'apk check: file can\'t check. file: {0} | ' \
                          'Error message: {1}\n'.format(file_path, str(error_message))
                path_files = []
            if 'AndroidManifest.xml' in path_files:
                result = True
                #if 'META-INF/MANIFEST.MF' in path_files:
                #    result = True
                #else:
                #    for filename in path_files:
                #        if filename.startswith('META-INF/') and filename.endswith('.SF'):
                #            result = True
                #            break
            if result and validate and ApkSignature:
                try:
                    result, algs_for_zip_dict = ApkSignature(file_path).verify()
                except Exception as error_message:
                    result = False
                    errors += 'apk check: file can\'t check with ' \
                              'ApkSignature. file: {0} | ' \
                              'Error message: {1}\n'.format(file_path, str(error_message))
    except Exception as error_message:
        errors += 'apk check: file can\'t check. file: {0} | ' \
                  'Error message: {1}\n'.format(file_path, str(error_message))
    return result, algs_for_zip_dict, errors
