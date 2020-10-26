package rbcrun

import (
	"testing"
)

func init() {
	starlarktestSetup()
}

func TestPropset(t *testing.T) {
	exerciseStarlarkTestFile(t, "testdata/propset.star")
	return
}
