# Misc testing
def product(name=None,brand=None, packages=[]):
    if name == None or brand == None:
        fail("name=, brand= are mandatory")
    setFinal("PRODUCT_NAME", name)
    setFinal("PRODUCT_BRAND", brand)
    for p in packages:
        appendTo("PRODUCT_PACKAGES", p)

