# Tests "queue" data type
load("assert.star", "assert")


def test_enqdeq():
    # Simple test that we can add to the queue while iterating over it.
    q.enqueue(1)
    sum = 0
    for x in q:
        assert.eq(q.dequeue(), x)
        assert.true(x == 1 or x == 2)
        if x == 1:
            q.enqueue(2)
        sum += x
    assert.eq(sum, 3)


def test_traversal():
    # Test that we can use the queue type to traverse the tree level by level
    tr = dict([
        ("1", ["1.1", "1.2"]),
        ("2", ["2.1"]),
        ("1.1", ["1.1.1", "1.1.2"]),
        ("1.2", ["1.2.1", "1.2.2"])
    ])
    expected_sequence = ["1", "2", "1.1", "1.2", "2.1", "1.1.1", "1.1.2", "1.2.1", "1.2.2"]
    actual_sequence = []
    myq = queue()
    myq.enqueue("1")
    myq.enqueue("2")
    for x in myq:
        myq.dequeue()
        actual_sequence.append(x)
        l = tr.get(x)
        if l == None:
            continue
        for c in l:
            myq.enqueue(c)

    assert.eq(expected_sequence, actual_sequence)


q = queue()

assert.eq(str(q), "[]")
assert.eq(dir(q), ["dequeue", "enqueue"])

q.enqueue(1)
q.enqueue(2)
assert.eq(q.dequeue(), 1)
q.enqueue(4)
assert.eq(q.dequeue(), 2)
assert.true(q)
assert.eq(q.dequeue(), 4)
assert.true(not q)
assert.fails(lambda: q.dequeue(), "queue is empty")

test_enqdeq()
test_traversal()
