"""Tool to find static libraries that maybe should be shared libraries and shared libraries that maybe should be static libraries.

This tool only looks at the module-info.json for the current target.
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
  parser = argparse.ArgumentParser(description="TODO(devinmoore) do things.")
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
  return parser.parse_args()


"""Loads the module-info.json file and analyzes the modules

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
"""


def main() -> None:
  module_info = json.load(open(ANDROID_PRODUCT_OUT + "/module-info.json"))

  global args
  args = parse_args()

  if args.module:
    if args.module not in module_info:
      print("Module {} does not exist".format(args.module))
      exit(1)

  includedStatically = defaultdict(int)
  includedSharedly = defaultdict(int)
  for module, value in module_info.items():
    if value["class"][0] == "NATIVE_TESTS":
      # We don't care about how tests are including libraries
      # Unfortunately not all tests fall under this class.
      # cc_fuzz is EXECUTABLES instead of NATIVE_TESTS.
      continue

    # count all of the shared and static libs included in this module
    for lib in value["shared_libs"]:
      includedSharedly[lib] += 1
    for lib in value["static_libs"]:
      includedStatically[lib] += 1

  if args.print_shared:
    # list of libraries included sharedly and used a little
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
    # list of libraries included statically and used a lot in this target
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
    # print info about the given library
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


if __name__ == "__main__":
  main()
