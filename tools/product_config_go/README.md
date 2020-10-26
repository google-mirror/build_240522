# Product Config Toy

Product config toy reads and executes Starlark files, resulting in creation
of string/list "variables". It then generates the sequence of makefile-style
assignments for them, e.g.:
```
PRODUCT_NAME:=coral
PRODUCT_PACKAGES:=foo bar baz
...
```

## Usage
`product_config_go` *file*

## Example
`example1.star`:
```python
load("product.star", "product")
product(name = "Coral",
        brand = "google",
        packages = ["foo", "bar"]
)
```

`product.star`:
```python
def product(name=None,brand=None, manufacturer=None):
    if name == None or brand == None:
        fail("name=, brand= are mandatory")
    setFinal("PRODUCT_NAME", name)
    setFinal("PRODUCT_BRAND", brand)
    for p in packages:
        appendTo("PRODUCT_PACKAGES", p)
```

We get
```
$ product_config_go example1.star
PRODUCT_NAME:=Coral
PRODUCT_BRAND:=google
PRODUCT_PACKAGES:=foo bar
```


## Predefined Functions
Starlark code can use the following functions:

##### set

`set("X", "value")` creates variable `X` and sets its value to `"value"`. It will fail if
`X` already exists and has been frozen (see `setFinal`)

##### setFinal

`setFinal("X", "value")` works as `set`, only that `X` value cannot be changed after the call.

##### appendTo
`appendTo("X", "value")` creates a list variable if it does not exist and appends `"value"` to
 its value list.

##### loadGenerated

`loadGenerated("cmd", ["arg1", ...])` runs command which generates Starlark script on stdout, which
 is then executed.