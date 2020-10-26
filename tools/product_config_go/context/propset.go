package context

import (
	"fmt"
	"go.starlark.net/starlark"
	"sort"
	"strconv"
	"strings"
)

// Make is the implementation of a built-in function that instantiates
// a propset.
//
func MakePropset(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("propset: unexpected positional arguments")
	}
	return &PropSet{make(map[string]starlark.Value)}, nil
}

// PropSet is similar to 'struct' or 'module' (see starklarkstruct package in starlark-go)
// but allowing to set the attribute values. Thus, with
//   ps = propset()
// we can later on set its arbitrary attributes
//   ps.x = 1
//   ps.y = [1,2]
// and then reference known attributes
//   print(ps.x)
// Just like for struct or module, propset's currently available attributes can be enumerated
// with dir(), attribute's presence can be checked with hasattr(), and dynamic attribute's
// value can be retrieved with getattr.
// At the same time, PropSet can be manipulated as a dictionary, that is, ps["x"] is equivalent
// to ps.x
// In the parlance of Starlark's Go implementation, PropSet type implements HasAttrs and HasSetKey
// interfaces.
type PropSet struct {
	properties map[string]starlark.Value
}

func (p PropSet) Get(value starlark.Value) (v starlark.Value, found bool, err error) {
	s, err := strconv.Unquote(value.String())
	if err != nil {
		return starlark.None, false, err
	}
	if v, found := p.properties[s]; found {
		return v, true, nil
	}
	return starlark.None, false, nil
}

func (p PropSet) SetKey(k, v starlark.Value) error {
	s, err := strconv.Unquote(k.String())
	if err != nil {
		return err
	}

	p.properties[s] = v
	return nil
}

func (p PropSet) String() string {
	buf := new(strings.Builder)
	buf.WriteString("propset(")
	sep := ""
	for _, n := range p.AttrNames() {
		buf.WriteString(sep)
		buf.WriteString(n)
		buf.WriteString(" = ")
		buf.WriteString(p.properties[n].String())
		sep = ", "
	}
	buf.WriteString(")")
	return buf.String()
}

func (p PropSet) Type() string {
	return "propset"
}

func (p PropSet) Freeze() {
	for _, v := range p.properties {
		v.Freeze()
	}
}

func (p PropSet) Truth() starlark.Bool {
	return len(p.properties) > 0
}

func (p PropSet) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type propset")
}

func (p PropSet) Attr(name string) (starlark.Value, error) {
	if v, found := p.properties[name]; found {
		return v, nil
	}
	return nil, nil
}

func (p PropSet) AttrNames() []string {
	props := make([]string, len(p.properties))
	i := 0
	for n := range p.properties {
		props[i] = n
		i++
	}
	sort.Strings(props)
	return props
}

func (p PropSet) SetField(name string, val starlark.Value) error {
	p.properties[name] = val
	return nil
}
