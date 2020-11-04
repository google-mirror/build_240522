# Defines the module with given name having a set of attributes.
# E.g., with 
#   part1 = cfg("part1", [], VAR11="value11", VAR12="value12")
#   part2 = cfg("part2", [], VAR2=["value2"])
#   top = cfg("top", [part1, part2], "VAR12"="top12", VARX="topx", VAR2=["top2"])
# defines:
#     Ref          Value
#   part1.VAR11    "value11"
#   part1.VAR12    "value12"
#   part2.VAR2     ["value2"]
#   top.VAR11      "value11"     i.e., single value from part is visible unless overidden
#   top.VAR12      "top12"       i.e., aggregate's single value overrides part's one
#   top.VAR2       ["value2", "top2"]  i.e., list values are appended
#   top.VARX       "topx"

def prodconf(name, parts, **kwargs):
    newattrs = dict()
    for part in parts:
        for attr in dir(part):
            new_value = getattr(part, attr)
            old_value = newattrs.get(attr)
            if old_value == None or type(old_value) != 'list':
                # Add/replace item
                newattrs[attr] = new_value
            else:
                if type(new_value) == "list":
                    newattrs[attr] = old_value + new_value
                else:
                    newattrs[attr] = old_value + [new_value]
    for attr in kwargs:
        part_value = newattrs.get(attr)
        composite_value = kwargs[attr]
        if part_value == None or type(part_value) != list:
            # Composite's item value overrides part's
            newattrs[attr] = composite_value
        else:
            if type(composite_value) == list:
                newattrs[attr] = part_value + composite_value
            else:
                newattrs[attr] = part_value + [composite_value]
    return module(name, **newattrs)


def printvars(mod):
    for attr in dir(mod):
        print(attr, "=", getattr(mod, attr))
