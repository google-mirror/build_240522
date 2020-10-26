# Tests 'propset'' extension.
load("assert.star", "assert")

s = propset()
assert.eq(str(s), "propset()")
assert.eq(dir(s), [])
s.host = "localhost"
s.port = 80
assert.eq(str(s), 'propset(host = "localhost", port = 80)')
assert.eq(s.host, "localhost")
assert.eq(getattr(s, "host"), "localhost")
assert.eq(s.port, 80)

s.port = 81
assert.eq(s.port, 81)
assert.eq(getattr(s, "port"), 81)
assert.eq(dir(s), ["host", "port"])

s.protocol = "http"
assert.eq(str(s), 'propset(host = "localhost", port = 81, protocol = "http")')

assert.eq(hasattr(s, "protocol"), True)

assert.eq(hasattr(s, "bad"), False)
assert.eq(getattr(s, "bad", None), None)
assert.fails(lambda: s.bad, "propset has no .bad field or method")

s["protocol"] = "https"
assert.eq(s.protocol, "https")

assert.eq(s["host"], "localhost")

s1 = propset(s)
assert.eq(str(s1), 'propset(host = "localhost", port = 81, protocol = "https")')
