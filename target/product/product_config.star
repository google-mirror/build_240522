"""Runtime functions."""

def _merge(dest, src):
    """Merges source propset into the destination one

    The rules for merging propset src into propset dest are as follows;
    * if there is no "X" in dest, dest.X is set to src.X
    * if "X" is in dest but is not a list, src.X replaces dest.X
    * finally, if the value of "X" in dest is a list, src.X
      is appended to it unless it is already present

    Args:
      dest: configuration propset
      src: module's config to be merged
    """
    for p in dir(src):
        __update_property(dest, p, getattr(src, p))


def __update_property(ps, attr, new_value):
    old_value = getattr(ps, attr) if hasattr(ps, attr) else None
    if old_value == None or type(old_value) != 'list':
        ps[attr] = new_value
    elif type(new_value) == 'list':
        ret = old_value
        for v in new_value:
            if v not in old_value:
                ret += [v]
        ps[attr] = ret
    else:
        ps[attr] = old_value + [new_value] if new_value not in old_value else old_value


def _printvars(mod):
    """Prints known configuration variables"""
    for attr in sorted(dir(mod)):
        print(attr, "=", getattr(mod, attr))


def _addprefix(prefix, string_or_list):
    """Adds prefix and returns a list.

    If string_or_list is a list, prepends prefix to each element.
    Otherwise, string_or_list is considered to be a string which
    is split into words and then prefix is prepended to each one.

    Args:
        prefix
        string_or_list

    """
    return [ prefix + x for x in __words(string_or_list)]


def _addsuffix(suffix, string_or_list):
    """Adds suffix and returns a list.

    If string_or_list is a list, appends suffix to each element.
    Otherwise, string_or_list is considered to be a string which
    is split into words and then suffix is appended to each one.

    Args:
      suffix
      string_or_list
    """
    return [ x + suffix for x in __words(string_or_list)]


def __words(string_or_list):
    if type(string_or_list) == "list":
        return string_or_list
    return string_or_list.split()


def _copy_if_exists(path_pair):
    """If from file exists, returns [from:to] pair."""
    l = path_pair.split(":", 2)
    # Check that l[0] exists
    return [":".join(l)] if rblf_file_exists(l[0]) else []


def _find_and_copy(pattern, from_dir, to_dir):
    """Return a copy list for the files matching the pattern."""
    return ["%s/%s:%s/%s" % (from_dir, f, to_dir, f) for f in rblf_wildcard(pattern, from_dir)]


def _global_init():
    """Returns PropSet created from the runtime environment."""
    ps = propset()
    # From build/make/core/envsetup.mk:
    ps.ART_APEX_JARS = [
        "com.android.art:core-oj",
        "com.android.art:core-libart",
        "com.android.art:okhttp",
        "com.android.art:bouncycastle",
        "com.android.art:apache-xml"
    ]
    #TODO(asmundak): which other variables set in envsetup.mk are to be prepopulated

    # Environment variables
    for k in dir(rblf_env):
        ps[k] = getattr(rblf_env, k)

    # Variables set as KEY_var command line arguments
    for k in dir(rblf_cli):
        ps[k] = getattr(rblf_cli, k)

    # Variables that should be defined.
    build_vars = [
        "BUILD_ID",
        "HOST_ARCH", "HOST_2ND_ARCH", "HOST_OS", "HOST_OS_EXTRA", "HOST_CROSS_OS", "HOST_CROSS_ARCH",
        "HOST_CROSS_2ND_ARCH", "HOST_BUILD_TYPE", "OUT_DIR",
        "PLATFORM_VERSION_CODENAME", "PLATFORM_VERSION", "PRODUCT_SOONG_NAMESPACES",
        "TARGET_PRODUCT", "TARGET_BUILD_VARIANT", "TARGET_BUILD_TYPE", "TARGET_ARCH", "TARGET_ARCH_VARIANT",
    ]
    for bv in build_vars:
        if not hasattr(ps, bv):
            rblf_error(bv + " is not defined")

    return ps


def _require_artifacts_in_path(paths, allowed_paths):
    """TODO."""
    print("require_artifacts_in_path(", __words(paths), ",", __words(allowed_paths), ")")


def _mkerror(file, message=""):
    """Prints error and stops."""
    rblf_error("%s: %s. Stop" % (file, message))


def _mkwarning(file, message=""):
    """Prints warning."""
    print("%s: warning: %s" % (file, message))


def _mkinfo(file, message=""):
    """Prints info."""
    print(message)


rblf = module("rblf",
              addprefix = _addprefix,
              addsuffix = _addsuffix,
              copy_if_exists = _copy_if_exists,
              find_and_copy = _find_and_copy,
              global_init = _global_init,
              merge = _merge,
              mkinfo = _mkinfo,
              mkerror = _mkerror,
              mkwarning = _mkwarning,
              printvars = _printvars,
              require_artifacts_in_path =_require_artifacts_in_path,
              warning = _mkwarning)
