<<<<<<< HEAD   (f7b9b7 Merge "Merge empty history for sparse-8547496-L6510000095455)
=======
# Module loaded my load.star
load("assert.star", "assert")

# Make sure that builtins are defined for the loaded module, too
assert.true(rblf_wildcard("module1.star"))
assert.true(not rblf_wildcard("no_such file"))
test = "module1"
>>>>>>> BRANCH (c458fa Merge "Version bump to TKB1.220517.001.A1 [core/build_id.mk])
