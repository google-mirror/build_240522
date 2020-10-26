# Test conditional loading
load(":mod.star|sym", "sym")
print(sym.VAR)
load(":nosuchfile.star|empty", "empty")
print(dir(empty))