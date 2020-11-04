"""Runtime functions."""

def _prodconf(module_name, parts, **kwargs):
    """Defines the module with given name having a set of attributes.

    E.g., with
       part1 = cfg("part1", [], VAR11="value11", VAR12="value12")
       part2 = cfg("part2", [], VAR2=["value2"])
       top = cfg("top", [part1, part2], "VAR12"="top12", VARX="topx", VAR2=["top2"])
    defines:
         Ref          Value
       part1.VAR11    "value11"
       part1.VAR12    "value12"
       part2.VAR2     ["value2"]
       top.VAR11      "value11"     i.e., single value from part is visible unless overridden
       top.VAR12      "top12"       i.e., aggregate's single value overrides part's one
       top.VAR2       ["value2", "top2"]  i.e., list values are appended
       top.VARX       "topx"

    Args:
        module_name: module name
        parts: parts configurations
        kwargs:

    Return:
        module
    """
    variables = dict()
    for part in parts:
        for nm in dir(part):
            variables[nm] = __updated_value(variables.get(nm), getattr(part, nm))
    for nm in kwargs:
        variables[nm] = __updated_value(variables.get(nm), kwargs[nm])
    return module(module_name, **variables)

def _merge(vars, mod):
    """Merges variables defined in the given module into the dictionary.

    For each member X in module:
    * if there is no "X" in dictionary, it is added with value module.X
    * if "X" is in dictionary but is not a list, module.X replaces it
    * finally, if the value of "X" in dictionary is a list, module.X
      is appended to it unless it is already present

    Args:
      vars: dictionary
      mod: module
    """
    for v in dir(mod):
        vars[v] = __updated_value(vars.get(v), getattr(mod, v))

def __updated_value(old_value, new_value):
    """Utility function used by prodconf and merge."""
    if old_value == None or type(old_value) != 'list':
        return new_value
    elif type(new_value) == 'list':
        ret = old_value
        for v in new_value:
            if v not in old_value:
                ret += [v]
        return ret
    else:
        return old_value + [new_value] if new_value not in old_value else old_value


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


def _soong_var(name):
    """Returns Soong variable value."""
    return []


def _copy_if_exists(path_pair):
    """If from file exists, returns [from:to] pair."""
    l = path_pair.split(":", 2)
    # Check that l[0] exists
    return [":".join(l)] if rblf_file_exists(l[0]) else []


def _find_and_copy(pattern, from_dir, to_dir):
    """Return a copy list for the files matching the pattern."""
    return ["%s/%s:%s/%s" % (from_dir, f, to_dir, f) for f in rblf_wildcard(pattern, from_dir)]


def _is_defined(name):
    """Returns true if variable is defined"""
    return False


def _is_board_platform_in(boards):
    """Returns True if TARGET_BOARD_PLATFORM is in the list."""
    return rblf.soong_var("TARGET_BOARD_PLATFORM") in boards


def _is_board_platform(board):
    """Returns True is $TARGET_BOARD_PLATFORM is equal ot given value."""
    return rblf.soong_var("TARGET_BOARD_PLATFORM") == board


def _is_vendor_board_platform(vendor):
    """Returns True is $TARGET_BOARD_PLATFORM is in $vendor_BOARD_PLATFORMS list"""
    return rblf.soong_var("TARGET_BOARD_PLATFORM") in rblf.soong_var(vendor + "_BOARD_PLATFORMS")


def _is_product_in(products):
    """Returns true if $TARGET_PRODUCT is in the list"""
    return rblf.soong_var("TARGET_PRODUCT") in products


def _require_artifacts_in_path(paths, allowed_paths):
    """TODO."""
    print("require_artifacts_in_path(", __words(paths), ",", __words(allowed_paths), ")")


def _warning(file, message=""):
    """Prints warning."""
    print("%s: warning: %s" % (file, message))


rblf = module("rblf",
              addprefix = _addprefix,
              addsuffix = _addsuffix,
              copy_if_exists = _copy_if_exists,
              find_and_copy = _find_and_copy,
              is_defined = _is_defined,
              is_board_platform = _is_board_platform,
              is_board_platform_in = _is_board_platform_in,
              is_product_in = _is_product_in,
              is_vendor_board_platform = _is_vendor_board_platform,
              merge = _merge,
              printvars = _printvars,
              prodconf = _prodconf,
              soong_var = _soong_var,
              require_artifacts_in_path =_require_artifacts_in_path,
              warning = _warning)
