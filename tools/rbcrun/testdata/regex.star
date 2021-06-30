# Tests rblf_regex
load("assert.star", "assert")


def test():
    pattern = "^(foo.*bar|abc.*d|1.*)$"
    for w in ("foobar", "fooxbar", "abcxd", "123"):
        assert.true(rblf_regex(pattern, w), "%s should match %s" % (w, pattern))
    for w in ("afoobar", "abcde"):
        assert.true(not rblf_regex(pattern, w), "%s should not match %s" % (w, pattern))
    for (pattern, subst, text, expected) in [
        ["abc", "def", "abcd", "defd"],
        ["abc(.*)", "def$1", "abcfoo", "deffoo"],
        ["a(.*)x", "b-${1}y", "adfhxx", "b-dfhxy"],
        ["a(.*)b(.*)c", "x${1}y${2}z", "a12b34c", "x12y34z"],
        ["a(.)b", "c${1}d", "a1bxa2ba3b", "c1dxc2dc3d"]]:
        actual = rblf_regex_subst(pattern, subst, text)
        assert.true(actual == expected, "subst(%s,%s,%s) should be '%s', got '%s'" % (
            pattern, subst, text, expected, actual))


test()
