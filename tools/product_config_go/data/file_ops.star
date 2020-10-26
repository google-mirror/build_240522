# Test file ops builtins
def run():
    for f in ("file_ops.star", "no such file"):
        print('rblf_file_exists("%s") =' % f, rblf_file_exists(f))
    for pattern, top in [("[fl]*.star", "."), ("[fl]*.star", ".."), ("foo *star", ".")]:
        print('rblf_wildcard("%s", "%s") =' % (pattern, top), rblf_wildcard(pattern, top))

run()