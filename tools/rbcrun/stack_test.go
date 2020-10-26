package rbcrun

import (
	"testing"
)

func init() {
	starlarktestSetup()
}

func TestStack(t *testing.T) {
	exerciseStarlarkTestFile(t, "testdata/stack.star")
}
