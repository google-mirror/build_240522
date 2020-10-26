# Test file ops builtins
print('rblf_file_exists("file_ops.star") =', rblf_file_exists("file_ops.star"))
print('rblf_file_exists("no such file") =', rblf_file_exists("no such file"))
print('rblf_wildcard("[fl]*.star", ".") =', rblf_wildcard("[fl]*.star", "."))
print('rblf_wildcard("foo *.star", ".") =', rblf_wildcard("foo *star", "."))