# Roboleaf configuration files interpreter

Reads and executes Roboleaf product configuration files.

## Usage

`rbcrun` *options* *VAR=value*... *file*

### Options

` -d` *dir*
Root directory for load("//path",...)
` -c` *text*
Read script from *text*

## Extensions

The runner allows Starlark scripts to use the following features that Bazel's Starlark interpreter does not support:

### Propset Data Type

`propset` is similar to `struct` or `module` (see `starklarkstruct` package in starlark-go)
but allows to set the attribute values for arbitrary attribute names. Thus, with

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

Just like for a struct or module, propset's currently available attributes can be enumerated with `dir()`, attribute's
presence can be checked with `hasattr()`, and dynamic attribute's value can be retrieved with `getattr()`.

At the same time, a propset can be manipulated as a dictionary, that is, `ps["x"]` is equivalent to `ps.x`.

### Load statement URI

Starlark does not define the format of the load statement's first argument. The Roboleaf configuration interpreter
supports the format that Bazel uses (`":file"` or `"//path:file"`). In addition, it allows the URI to end
with `"|symbol"` which defines a single variable
`symbol` with `None` value if a module does not exist. Thus,

```
load(":mymodule.rbc|init", mymodule_init="init")
```

will load the module `mymodule.rbc` and export a symbol `init` in it as `mymodule_init` if
`mymodule.rbc` exists. If `mymodule.rbc` is missing, `mymodule_init` will be set to `None`

### Predefined Symbols

#### rblf_env

A propset containing environment variables. E.g., `rblf_env.USER` is the username when running on Unix.

#### rblf_cli

A propset containing the variable set by the interpreter's command line. That is, running

```
rbcrun FOO=bar myfile.rbc
```

will have the value of `rblf_cli.FOO` be `"bar"`

### Predefined Functions

#### rblf_file_exists(*file*)

Returns `True`  if *file* exists

#### rblf_wildcard(*pattern*, *root-path*)

Returns the list of the files in *root-path* that match *pattern*.

#### rblf_regex(*pattern*, *text*)

Returns *True* if *text matches *pattern*.

##### loadGenerated

`loadGenerated("cmd", ["arg1", ...])` runs command which generates Starlark script on stdout, which is then executed.