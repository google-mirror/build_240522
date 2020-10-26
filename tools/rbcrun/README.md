# Roboleaf configuration files interpreter

Reads and executes Roboleaf product configuration files.

## Usage

`rbcrun` *options* *VAR=value*... [ *file* ]

A Roboleaf configuration file is a Starlark script. Usually it is read from *file*, but it can be specified inline as a
value of `-c` option. If file's name contains `=`, use `-f` option to specify file's name
(my=file.rbc sets the value of `my` to `file.rbc`).

### Options

`-d` *dir*\
Root directory for load("//path",...)

`-c` *text*\
Read script from *text*

`--perf` *file*\
Gather performance statistics and save it to *file*. Use \
`       go tool prof -top`*file*\
to show top CPU users

`-f` *file*\
File to run.

## Extensions

The runner allows Starlark scripts to use the following features that Bazel's Starlark interpreter does not support:

### Propset Data Type

`propset` is similar to `struct` or `module` (see `starklarkstruct` package in starlark-go) but allows to set the
attribute values for arbitrary attribute names. Thus, with

```
   ps = propset()
```

we can set its arbitrary attributes

```
   ps.x = 1
   ps.y = [1,2]
```

can then reference known attributes

```
   print(ps.x)
```

Just like for a struct or module, propset's currently available attributes can
be enumerated with `dir()`, attribute's presence can be checked with
`hasattr()`, and dynamic attribute's value can be retrieved with `getattr()`.

At the same time, a propset can be manipulated as a dictionary, that is,
`ps["x"]` is equivalent to `ps.x`.

### Stack

The only loop construct in Starlark is the iteration over an iterable data type,
and the language prohibits recursive calls. This makes it very difficult to
traverse product's configuration hierarchy the same way as `inherit-product`
macro does for the makefiles. The runner provides an iterable data type called
`stack` which can be modified during the iteration and thus provides the
controlled workaround for the restrictions above.

A stack has `push` and `pop` methods.

Example:

```python
def foo():
    s = stack()
    s.push(1)
    for item in q:
        if item == 1:
            q.push(2)
```

the loop body will be executed twice, because during the first iteration `2`
will be pushed to the `s`.

### Load statement URI

Starlark does not define the format of the load statement's first argument.
The Roboleaf configuration interpreter supports the format that Bazel uses
(`":file"` or `"//path:file"`). In addition, it allows the URI to end with
`"|symbol"` which defines a single variable `symbol` with `None` value if a
module does not exist. Thus,

```
load(":mymodule.rbc|init", mymodule_init="init")
```

will load the module `mymodule.rbc` and export a symbol `init` in it as
`mymodule_init` if `mymodule.rbc` exists. If `mymodule.rbc` is missing,
`mymodule_init` will be set to `None`

### Predefined Symbols

#### rblf_env

A propset containing environment variables. E.g., `rblf_env.USER` is the
username when running on Unix.

#### rblf_cli

A propset containing the variable set by the interpreter's command line. That
is, running

```
rbcrun FOO=bar myfile.rbc
```

will have the value of `rblf_cli.FOO` be `"bar"`

### Predefined Functions

#### rblf_file_exists(*file*)

Returns `True`  if *file* exists

#### rblf_wildcard(*glob*, *top* = None)

Expands *glob*. If *top* is supplied, expands "*top*/*glob*", then removes
"*top*/" prefix from the matching file names.

#### rblf_regex(*pattern*, *text*)

Returns *True* if *text* matches *pattern*.

#### rblf_shell(*command*)

Runs `sh -c "`*command*`"`, reads its output, converts all newlines into spaces,
chops trailing newline returns this string. This is equivalent to Make's
`shell` builtin function

##### loadGenerated(*command*, [*arg1*, ...])

Runs command which generates Starlark script on stdout, which is then executed.
*This function is disabled at the moment.*

##### psetdefault(*propset*, "attribute", *value*)

If *propset* has no *attribute*, sets *propset*.*attribute* to *value*. Otherwise does nothing. Most of a propset
attributes are incrementally built lists, and this function allows the code converter generate attribute initialization
code.
