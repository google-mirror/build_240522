# Test rblf_wildcard
load("assert.star", "assert")

assert.eq(rblf_wildcard("module*.star"), ["module1.star", "module2.star"])
assert.eq(rblf_wildcard("wildcard.star"), ["wildcard.star"])
assert.eq(rblf_wildcard(["wildcard.star", "shell.star"]), ["wildcard.star", "shell.star"])
