package rbcrun

import (
	"testing"
)

func init() {
	starlarktestSetup()
}

func TestQueue(t *testing.T) {
	exerciseStarlarkTestFile(t, "testdata/queue.star")
}
