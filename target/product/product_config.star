"""Runtime functions."""

def _prodconf(name, parts, **kwargs):
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
       top.VAR11      "value11"     i.e., single value from part is visible unless overidden
       top.VAR12      "top12"       i.e., aggregate's single value overrides part's one
       top.VAR2       ["value2", "top2"]  i.e., list values are appended
       top.VARX       "topx"

    Args:
        name: module name
        parts: subconfigurations
        kwargs:

    Return:
        module
    """
    newattrs = dict()
    for part in parts:
        for attr in dir(part):
            newattrs[attr] = __updated_value(newattrs.get(attr), getattr(part, attr))
    for attr in kwargs:
        newattrs[attr] = __updated_value(newattrs.get(attr), kwargs[attr])
    return module(name, **newattrs)


def __updated_value(old_value, new_value):
    """Utility function used by prodconf."""
    if old_value == None or type(old_value) != 'list':
        return new_value
    elif type(new_value) == 'list':
        return old_value + new_value
    else:
        return old_value + [new_value]


def _printvars(mod):
    """Prints known configuration variables"""
    for attr in dir(mod):
        print(attr, "=", getattr(mod, attr))

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

def _warning(file, message=""):
    """Prints warning."""
    print("%s: warning: %s" % (file, message))

rblf = module("rblf",
              copy_if_exists = _copy_if_exists,
              find_and_copy = _find_and_copy,
              is_defined = _is_defined,
              printvars = _printvars,
              prodconf = _prodconf,
              soong_var = _soong_var,
              warning = _warning)
