package rbcrun

import (
	"fmt"
	"go.starlark.net/starlark"
	"strings"
)

// Queue is an iterable allowing to remove the front item ("read") and to
// add the items both at the front ("unread") and at the back ("write").
// The iteration traverses the queue from front to back.
// Unlike  the available iterables (lists, etc.) dequeue can be updated
// during the iteration (and this is the reason to have it).

var queueMethods = map[string]*starlark.Builtin{
	"read":   starlark.NewBuiltin("read", queueRead),
	"unread": starlark.NewBuiltin("unread", queueUnread),
	"write":  starlark.NewBuiltin("write", queueWrite),
}

type Queue struct {
	items     []starlark.Value
	itercount uint32
	frozen    bool
}

func (q *Queue) Attr(name string) (starlark.Value, error) {
	if v, found := queueMethods[name]; found {
		return v.BindReceiver(q), nil
	}
	return nil, nil
}

func (q *Queue) AttrNames() []string {
	var names []string
	for k := range queueMethods {
		names = append(names, k)
	}
	return names
}

func (q *Queue) String() string {
	buf := new(strings.Builder)
	buf.WriteString("[")
	const maxQueue = 1000
	for i, item := range q.items {
		if i > 0 {
			buf.WriteString(", ")
		}
		if i > maxQueue {
			buf.WriteString("...")
			break
		}
		buf.WriteString(item.String())
	}
	buf.WriteString("]")
	return buf.String()
}

func (q *Queue) Type() string {
	return "queue"
}

func (q *Queue) Freeze() {
	if q.frozen {
		return
	}
	q.frozen = true
	for _, item := range q.items {
		item.Freeze()
	}
}

func (q *Queue) Truth() starlark.Bool {
	return len(q.items) > 0
}

func (q *Queue) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: queue")
}

func (q *Queue) Iterate() starlark.Iterator {
	if !q.frozen {
		q.itercount++
	}
	return &queueIterator{q}
}

// Adds an item to the front of the dequeue.
func queueUnread(_ *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var v starlark.Value
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &v); err != nil {
		return nil, err
	}
	recv := fn.Receiver().(*Queue)
	if recv.frozen {
		return nil, fmt.Errorf("%s: queue is frozen", fn.Name())
	}
	recv.items = append([]starlark.Value{v}, recv.items...)
	return starlark.None, nil
}

// Adds an item to the back of the dequeue
func queueWrite(_ *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var v starlark.Value
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &v); err != nil {
		return nil, err
	}
	recv := fn.Receiver().(*Queue)
	if recv.frozen {
		return nil, fmt.Errorf("%s: queue is frozen", fn.Name())
	}
	recv.items = append(recv.items, v)
	return starlark.None, nil
}

// Removes the front item
func queueRead(_ *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 0); err != nil {
		return nil, err
	}
	recv := fn.Receiver().(*Queue)
	if len(recv.items) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}
	var v starlark.Value
	v, recv.items = recv.items[0], recv.items[1:]
	return v, nil
}

func MakeQueue(_ *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	q := &Queue{}
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 0); err != nil {
		return nil, nil
	}
	return q, nil
}

type queueIterator struct {
	q *Queue
}

func (iq *queueIterator) Next(p *starlark.Value) bool {
	if len(iq.q.items) == 0 {
		return false
	}
	*p = iq.q.items[0]
	return true
}

func (iq *queueIterator) Done() {
	if !iq.q.frozen {
		iq.q.itercount--
	}
}
