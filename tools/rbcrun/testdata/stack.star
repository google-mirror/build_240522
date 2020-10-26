# Tests "stack" data type
load("assert.star", "assert")

def test_iter():
    # Simple test that we can add to the stack while iterating over it.
    s1 = stack()
    s1.push(1)
    s1.push(2)
    assert.eq(2, s1.pop())
    sum = 0
    for x in s1:
        assert.true(x == 1 or x == 2)
        if x == 1:
            s1.push(2)
        sum += x
    assert.eq(sum, 3)


def test_tree():
    # Test that we can use the stack type to traverse the tree in prefix order
    children = {
        "1": ["1.1", "1.2"],
        "2": ["2.1"],
        "1.1": ["1.1.1", "1.1.2"],
        "1.2": ["1.2.1", "1.2.2"]
    }
    expected_sequence = ["1", "1.1", "1.1.1", "1.1.2", "1.2", "1.2.1", "1.2.2", "2", "2.1"]
    actual_sequence = []
    myst = stack()
    myst.push("2")
    myst.push("1")
    for node in myst:
        actual_sequence.append(node)
        subnodes = children.get(node)
        if subnodes != None:
            for subnode in reversed(subnodes):
                myst.push(subnode)

    assert.eq(expected_sequence, actual_sequence)


st = stack()

assert.eq(str(st), "[]")
assert.eq(dir(st), ["pop", "push"])

st.push(1)
st.push(2)  # Now: [1, 2]
assert.eq(st.pop(), 2)  # [2]
assert.true(st)
assert.eq(st.pop(), 1)
assert.true(not st)
assert.fails(lambda: st.pop(), "stack is empty")

test_iter()
test_tree()
