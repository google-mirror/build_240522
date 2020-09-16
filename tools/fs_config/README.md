# FS Config Generator

The `fs_config_generator.py` tool uses the platform `android_filesystem_config.h` and the
`TARGET_FS_CONFIG_GEN` files to generate the following:
* `fs_config_dirs` and `fs_config_files` files for each partition
* `passwd` and `group` files for each partition
* The `generated_oem_aid.h` header

## Outputs

### `fs_config_dirs` and `fs_config_files`

The `fs_config_dirs` and `fs_config_files` binary files are interpreted by the libcutils
`fs_config()` function, along with the built-in defaults, to serve as overrides to complete the
results. The Target files are used by filesystem and adb tools to ensure that the file and directory
properties are preserved during runtime operations. The host files in the `$OUT` directory are used
in the final stages when building the filesystem images to set the file and directory properties.

See `./fs_config_generator.py fsconfig --help` for how these files are generated.

### `passwd` and `group` files

The `passwd` and `group` files are formatted as documented in man pages passwd(5) and group(5) and
used by bionic for implementing `getpwnam()` and related functions.

See `./fs_config_generator.py passwd --help` and `./fs_config_generator.py group --help` for how
these files are generated.

### The `generated_oem_aid.h` header

The `generated_oem_aid.h` creates identifiers for non-platform AIDs for developers wishing to use
them in their native code.  To do so, include the `oemaids_headers` header library in the
corresponding makefile and `#include "generated_oem_aid.h"` in the code wishing to use these
identifiers.

See `./fs_config_generator.py oemaid --help` for how this file is generated.

## Parsing

The parsing of the `TARGET_FS_CONFIG_GEN` files follows the Python `ConfigParser` specification,
with the sections and fields as defined below. There are two types of sections, both require all
options to be specified.

### Filesystem capabilities

File system capabilities follow the below syntax:

    [path]
    mode: Octal file mode
    user: AID_<user>
    group: AID_<group>
    caps: cap*

Where:

`[path]` \
  The filesystem path to configure. A path ending in / is considered a dir, else its a file.

`mode:` \
  A valid octal file mode of at least 3 digits. If 3 is specified, it is prefixed with a 0, else
  mode is used as is.

`user:` \
  Either the C define for a valid AID or the friendly name. For instance both `AID_RADIO` and
  `radio` are acceptable. Note custom AIDs can be defined in the AID section documented below.

`group:` \
  Same as user.

`caps:` \
  The name as declared in
  [bionic/libc/kernel/uapi/linux/capability.h](../../../../bionic/libc/kernel/uapi/linux/capability.h) without the leading `CAP_`.
  Mixed case is allowed. Caps can also be the raw:
   * binary (0b0101)
   * octal (0455)
   * int (42)
   * hex (0xFF)

  For multiple caps, just separate by whitespace.

It is an error to specify multiple sections with the same `[path]` in different
files. Note that the same file may contain sections that override the previous
section in Python versions <= 3.2. In Python 3.2 it's set to strict mode.

### OEM AIDs

The AID section allows specifying OEM specific AIDs and follows the below syntax:

    [AID_<name>]
    value: <number>

Where:

`[AID_<name>]` \
The `<name>` can contain characters in the set uppercase, numbers, and underscores. The lowercase
version is used as the friendly name. The generated header file for code inclusion uses the exact
`AID_<name>`.

It is an error to specify multiple sections with the same `AID_<name>` (case insensitive with the
same constraints as [path]).

`<name>` must begin with a partition name to ensure that it does not conflict with different
sources.


`value:` \
A valid C style number string (hex, octal, binary and decimal).

It is an error to specify multiple sections with the same value option.

Value options must be specified in the range corresponding to the partition used in `<name>`. The
list of valid partitions and their corresponding ranges is defined in
[system/core/include/private/android_filesystem_config.h](../../../../system/core/include/private/android_filesystem_config.h).
The options are:
* Vendor Partition
  * AID_OEM_RESERVED_START(2900) - AID_OEM_RESERVED_END(2999)
  * AID_OEM_RESERVED_2_START(5000) - AID_OEM_RESERVED_2_END(5999)
* System Partition
  * AID_SYSTEM_RESERVED_START(6000) - AID_SYSTEM_RESERVED_END(6499)
* ODM Partition
  * AID_ODM_RESERVED_START(6500) - AID_ODM_RESERVED_END(6999)
* Product Partition
  * AID_PRODUCT_RESERVED_START(7000) - AID_PRODUCT_RESERVED_END(7499)
* System_ext Partition
  * AID_SYSTEM_EXT_RESERVED_START(7500) - AID_SYSTEM_EXT_RESERVED_END(7999)

### Ordering

Ordering within the `TARGET_FS_CONFIG_GEN` files is not relevant. The paths for files are sorted
like so within their respective array definition:
 * specified path before prefix match
   * for example: foo before f*
 * lexicographical less than before other
   * for example: boo before foo

Given these paths:

    paths=['ac', 'a', 'acd', 'an', 'a*', 'aa', 'ac*']

The sort order would be:

    paths=['a', 'aa', 'ac', 'acd', 'an', 'ac*', 'a*']

Thus the `fs_config` tools will match on specified paths before attempting prefix, and match on the
longest matching prefix.

The declared AIDs are sorted in ascending numerical order based on the option "value". The string
representation of value is preserved. Both choices were made for maximum readability of the
generated file and to line up files. Sync lines are placed with the source file as comments in the
generated header file.

## Unit Tests

From within the `fs_config` directory, unit tests can be executed like so:

    $ python -m unittest test_fs_config_generator.Tests
    .............
    ----------------------------------------------------------------------
    Ran 13 tests in 0.004s

    OK

One could also use nose if they would like:

    $ nose2

To add new tests, simply add a `test_<xxx>` method to the test class. It will automatically
get picked up and added to the test suite.
