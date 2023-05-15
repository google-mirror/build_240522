#!/usr/bin/env python3

"""Tool to find static libraries that maybe should be shared libraries and shared libraries that maybe should be static libraries.

This tool only looks at the module-info.json for the current target.

Example of "class" types for each of the modules in module-info.json
  "EXECUTABLES": 2307,
  "ETC": 9094,
  "NATIVE_TESTS": 10461,
  "APPS": 2885,
  "JAVA_LIBRARIES": 5205,
  "EXECUTABLES/JAVA_LIBRARIES": 119,
  "FAKE": 553,
  "SHARED_LIBRARIES/STATIC_LIBRARIES": 7591,
  "STATIC_LIBRARIES": 11535,
  "SHARED_LIBRARIES": 10852,
  "HEADER_LIBRARIES": 1897,
  "DYLIB_LIBRARIES": 1262,
  "RLIB_LIBRARIES": 3413,
  "ROBOLECTRIC": 39,
  "PACKAGING": 5,
  "PROC_MACRO_LIBRARIES": 36,
  "RENDERSCRIPT_BITCODE": 17,
  "DYLIB_LIBRARIES/RLIB_LIBRARIES": 8,
  "ETC/FAKE": 1

None of the "SHARED_LIBRARIES/STATIC_LIBRARIES" are double counted in the
modules with one class
RLIB/

All of these classes have shared_libs and/or static_libs
    "EXECUTABLES",
    "SHARED_LIBRARIES",
    "STATIC_LIBRARIES",
    "SHARED_LIBRARIES/STATIC_LIBRARIES", # cc_library
    "HEADER_LIBRARIES",
    "NATIVE_TESTS", # test modules
    "DYLIB_LIBRARIES", # rust
    "RLIB_LIBRARIES", # rust
    "ETC", # rust_bindgen
"""

from collections import defaultdict

import json, os, argparse

ANDROID_PRODUCT_OUT = os.environ.get("ANDROID_PRODUCT_OUT")
# If a shared library is used less than MAX_SHARED_INCLUSIONS times in a target,
# then it will likely save memory by changing it to a static library
# This move will also use less storage
MAX_SHARED_INCLUSIONS = 2
# If a static library is used more than MAX_STATIC_INCLUSIONS times in a target,
# then it will likely save memory by changing it to a shared library
# This move will also likely use less storage
MIN_STATIC_INCLUSIONS = 3


def parse_args():
  parser = argparse.ArgumentParser(
      description=(
          "Parse module-info.jso and display information about static and"
          " shared library dependencies."
      )
  )
  parser.add_argument(
      "--module", dest="module", help="Print the info for the module."
  )
  parser.add_argument(
      "--shared",
      dest="print_shared",
      action=argparse.BooleanOptionalAction,
      help=(
          "Print the list of libraries that are shared_libs for fewer than {}"
          " modules.".format(MAX_SHARED_INCLUSIONS)
      ),
  )
  parser.add_argument(
      "--static",
      dest="print_static",
      action=argparse.BooleanOptionalAction,
      help=(
          "Print the list of libraries that are static_libs for more than {}"
          " modules.".format(MIN_STATIC_INCLUSIONS)
      ),
  )
  parser.add_argument(
      "--recursive",
      dest="recursive",
      action=argparse.BooleanOptionalAction,
      default=True,
      help=(
          "Gather all dependencies of EXECUTABLES recursvily before calculating"
          " the stats. This eliminates duplicates from multiple libraries"
          " including the same dependencies in a single binary."
      ),
  )
  parser.add_argument(
      "--both",
      dest="both",
      action=argparse.BooleanOptionalAction,
      default=False,
      help=(
          "Print a list of libraries that are including libraries as both"
          " static and shared"
      ),
  )
  return parser.parse_args()


class TransitiveHelper:

  class Libs:

    def __init__(self):
      self.shared = set()
      self.static = set()

  def __init__(self):
    # keep a list of already expanded libraries so we don't end up in a cycle
    self.visited = defaultdict(lambda: self.Libs())

  def extendUnique(self, module_libs, new_libs):
    module_libs.extend(filter(lambda x: x not in module_libs, new_libs))

  # module is an object from the module-info dictionary
  # module_info is the dictionary from module-info.json
  # modify the module's shared_libs and static_libs with all of the transient
  # dependencies required from all of the explicit dependencies
  def flattenDeps(self, module, module_info):
    shared = module["shared_libs"]
    static = module["static_libs"]

    for lib in shared:
      if not lib or lib not in module_info:
        continue
      if lib in self.visited:
        self.extendUnique(module["shared_libs"], self.visited[lib].shared)
      else:
        res = self.flattenDeps(module_info[lib], module_info)
        self.extendUnique(module["shared_libs"], res["shared_libs"])
        self.visited[lib].shared.update(res["shared_libs"])

    for lib in static:
      if not lib or lib not in module_info:
        continue
      if lib in self.visited:
        self.extendUnique(module["static_libs"], self.visited[lib].shared)
      else:
        res = self.flattenDeps(module_info[lib], module_info)
        self.extendUnique(module["static_libs"], res["static_libs"])
        self.visited[lib].static.update(res["static_libs"])

    return module


def main():
  module_info = json.load(open(ANDROID_PRODUCT_OUT + "/module-info.json"))

  global args
  args = parse_args()

  if args.module:
    if args.module not in module_info:
      print("Module {} does not exist".format(args.module))
      exit(1)

  includedStatically = defaultdict(int)
  includedSharedly = defaultdict(int)
  includedBothly = defaultdict(set)
  transitive = TransitiveHelper()
  for name, module in module_info.items():
    if args.recursive:
      # in this recursive mode we only want to see what is included by the executables
      if module["class"][0] != "EXECUTABLES":
        continue
      module = transitive.flattenDeps(module, module_info)
      # filter out fuzzers by their dependency on clang
      if "libclang_rt.fuzzer" in module["static_libs"]:
        continue
    else:
      if module["class"][0] == "NATIVE_TESTS":
        # We don't care about how tests are including libraries
        continue

    # count all of the shared and static libs included in this module
    for lib in module["shared_libs"]:
      includedSharedly[lib] += 1
    for lib in module["static_libs"]:
      includedStatically[lib] += 1

    intersection = set(module["shared_libs"]).intersection(
        module["static_libs"]
    )
    if intersection:
      includedBothly[name] = intersection

  if args.print_shared:
    print(
        "Shared libraries that are included by fewer than {} modules on a"
        " device:".format(MAX_SHARED_INCLUSIONS)
    )
    print(
        json.dumps(
            {
                name: includedSharedly[name]
                for name in includedSharedly
                if includedSharedly[name] < MAX_SHARED_INCLUSIONS
            },
            indent=2,
        )
    )

  if args.print_static:
    print(
        "Libraries that are included statically by more than {} modules on a"
        " device:".format(MIN_STATIC_INCLUSIONS)
    )
    print(
        json.dumps(
            {
                name: includedStatically[name]
                for name in includedStatically
                if includedStatically[name] > MIN_STATIC_INCLUSIONS
            },
            indent=2,
        )
    )

  if args.module:
    print(json.dumps(module_info[args.module], indent=2))
    print(
        "{} is included in shared_libs {} times".format(
            args.module, includedSharedly[args.module]
        )
    )
    print(
        "{} is included in static_libs {} times".format(
            args.module, includedStatically[args.module]
        )
    )

  if args.both:
    allIncludedBothly = set()
    for name, libs in includedBothly.items():
      print(
          "{} includes the following libraries as both shared and static: {}"
          .format(name, libs)
      )
      allIncludedBothly.update(libs)
    print(
        "List of libraries used both statically and shared in the same"
        " processes: {}".format(allIncludedBothly)
    )


if __name__ == "__main__":
  main()
