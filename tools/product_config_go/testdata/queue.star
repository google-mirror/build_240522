# Tests "queue" data type
load("assert.star", "assert")


def test_enqdeq():
    # Simple test that we can add to the queue while iterating over it.
    q.write(1)
    sum = 0
    for x in q:
        assert.eq(q.read(), x)
        assert.true(x == 1 or x == 2)
        if x == 1:
            q.write(2)
        sum += x
    assert.eq(sum, 3)


def test_traversal1():
    # Test that we can use the queue type to traverse the tree level by level
    children = {
        "1": ["1.1", "1.2"],
        "2": ["2.1"],
        "1.1": ["1.1.1", "1.1.2"],
        "1.2": ["1.2.1", "1.2.2"]
    }
    expected_sequence = ["1", "2", "1.1", "1.2", "2.1", "1.1.1", "1.1.2", "1.2.1", "1.2.2"]
    actual_sequence = []
    myq = queue()
    myq.write("1")
    myq.write("2")
    for x in myq:
        actual_sequence.append(myq.read())
        subnodes = children.get(x)
        if subnodes != None:
            for subnode in subnodes:
                myq.write(subnode)

    assert.eq(expected_sequence, actual_sequence)


def test_traversal2():
    # Test that we can use the queue type to traverse the tree depth first
    children = {
        "1": ["1.1", "1.2"],
        "2": ["2.1"],
        "1.1": ["1.1.1", "1.1.2"],
        "1.2": ["1.2.1", "1.2.2"]
    }
    expected_sequence = ["1", "1.1", "1.1.1", "1.1.2", "1.2", "1.2.1", "1.2.2", "2", "2.1"]
    actual_sequence = []
    myq = queue()
    myq.write("1")
    myq.write("2")
    for x in myq:
        actual_sequence.append(myq.read())
        subnodes = children.get(x)
        if subnodes != None:
            for subnode in reversed(subnodes):
                myq.unread(subnode)

    assert.eq(expected_sequence, actual_sequence)


q = queue()

assert.eq(str(q), "[]")
assert.eq(dir(q), ["read", "unread", "write"])

q.write(1)
q.write(2)  # Now: [1, 2]
assert.eq(q.read(), 1)  # [2]
q.unread(5)  # [5, 2]
assert.eq(q.read(), 5)  # [2]
q.write(4)  # [2, 4]
assert.eq(q.read(), 2)  # [4]
assert.true(q)
assert.eq(q.read(), 4)  # []
assert.true(not q)
assert.fails(lambda: q.read(), "queue is empty")

test_enqdeq()
test_traversal1()
test_traversal2()
