package rbcrun

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

// Stack is an iterable.
// Unlike the available iterables (lists, etc.) stack can be updated
// during the iteration (and this is the reason to have it).

var stackMethods = map[string]*starlark.Builtin{
	"pop":  starlark.NewBuiltin("write", stackPop),
	"push": starlark.NewBuiltin("unread", stackPush),
}

type Stack struct {
	items     []starlark.Value
	itercount uint32
	frozen    bool
}

func (st *Stack) Attr(name string) (starlark.Value, error) {
	if v, found := stackMethods[name]; found {
		return v.BindReceiver(st), nil
	}
	return nil, nil
}

func (st *Stack) AttrNames() []string {
	var names []string
	for k := range stackMethods {
		names = append(names, k)
	}
	return names
}

func (st *Stack) String() string {
	buf := new(strings.Builder)
	buf.WriteString("[")
	const maxStack = 1000
	for i, item := range st.items {
		if i > 0 {
			buf.WriteString(", ")
		}
		if i > maxStack {
			buf.WriteString("...")
			break
		}
		buf.WriteString(item.String())
	}
	buf.WriteString("]")
	return buf.String()
}

func (st *Stack) Type() string {
	return "stack"
}

func (st *Stack) Freeze() {
	if st.frozen {
		return
	}
	st.frozen = true
	for _, item := range st.items {
		item.Freeze()
	}
}

func (st *Stack) Truth() starlark.Bool {
	return len(st.items) > 0
}

func (st *Stack) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: stack")
}

func (st *Stack) Iterate() starlark.Iterator {
	if !st.frozen {
		st.itercount++
	}
	return &stackIterator{st}
}

func (st *Stack) pop() (starlark.Value, bool) {
	n := len(st.items) - 1
	if n < 0 {
		return starlark.None, false
	}
	v := st.items[n]
	st.items = st.items[:n]
	return v, true
}

// Pushes an item.
func stackPush(_ *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var v starlark.Value
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &v); err != nil {
		return nil, err
	}
	st := fn.Receiver().(*Stack)
	if st.frozen {
		return nil, fmt.Errorf("%s: stack is frozen", fn.Name())
	}
	st.items = append(st.items, v)
	return starlark.None, nil
}

// Pops the top item
func stackPop(_ *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 0); err != nil {
		return nil, err
	}
	st := fn.Receiver().(*Stack)
	if st.frozen {
		return nil, fmt.Errorf("%s: stack is frozen", fn.Name())
	}
	if v, ok := st.pop(); ok {
		return v, nil
	}
	return nil, fmt.Errorf("stack is empty")
}

func MakeStack(_ *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	q := &Stack{}
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 0); err != nil {
		return nil, nil
	}
	return q, nil
}

type stackIterator struct {
	stack *Stack
}

func (ist *stackIterator) Next(p *starlark.Value) bool {
	var ok bool
	*p, ok = ist.stack.pop()
	return ok
}

func (ist *stackIterator) Done() {
	if !ist.stack.frozen {
		ist.stack.itercount--
	}
}
