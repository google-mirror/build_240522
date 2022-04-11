#!/usr/bin/python3
#
# Copyright (C) 2022 The Android Open Source Project
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

sys.dont_write_bytecode = True
import lunch

EXIT_STATUS_OK = 0
EXIT_STATUS_ERROR = 1

API_DOMAIN_SYSTEM = "system"
API_DOMAIN_VENDOR = "vendor"
API_DOMAIN_MODULE = "module"

class ApiDomain(object):
    def __init__(self, name, tree, product):
        # Product will be null for modules
        self.name = name
        self.tree = tree
        self.product = product

        
class InnerTree(object):
    def __init__(self, root):
        """Initialize with the inner tree root (relative to the workspace root)"""
        self.root = root
        self.domains = {}

    def Invoke(self, args):
        """Call the inner tree command for this inner tree."""
        pass


def ProcessConfig(config):
    """Returns the InnerTree and ApiDomain dicts for the config."""
    def Add(name, tree, product):
        if tree in trees:
            tree = trees[tree]
        else:
            tree = InnerTree(tree)
            trees[tree.root] = tree
        domain = ApiDomain(name, tree, product)
        domains[name] = domain

    trees = {}
    domains = {}

    system_entry = config.get("system")
    if system_entry:
        Add(API_DOMAIN_SYSTEM, system_entry["tree"], system_entry["product"])

    vendor_entry = config.get("vendor")
    if vendor_entry:
        Add(API_DOMAIN_VENDOR, vendor_entry["tree"], vendor_entry["product"])

    for module_name, module_entry in config.get("modules", []).items():
        Add(module_name, module_entry["tree"], None)

    return trees, domains


def Build():
    #
    # Load lunch combo
    #

    # Read the config file
    try:
        config = lunch.LoadCurrentConfig()
    except lunch.ConfigException as ex:
        sys.stderr.write("%s\n" % ex)
        return EXIT_STATUS_ERROR
    print("Config: %s" % config)

    # Construct the trees and domains dicts
    trees, domains = ProcessConfig(config)
    print("trees=%s" % trees)
    print("domains=%s" % domains)

    #
    # Interrogate the trees
    #

    #
    # API Export
    #
    

    #
    # API Surface Assembly
    #

    #
    # API Domain Analysis
    #

    #
    # Final Packaging Rules
    #

    #
    # Build Execution
    #


    #
    # Success!
    #
    return EXIT_STATUS_OK

def main(argv):
    return Build()

if __name__ == "__main__":
    sys.exit(main(sys.argv))

