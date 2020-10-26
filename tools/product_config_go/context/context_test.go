package context

import (
	"bytes"
	"testing"
)

func TestBuiltins(t *testing.T) {
	tests := []struct {
		src, want string
	}{
		{
			"set('VAR', 'value')",
			"VAR:=value\n",
		},
		{
			"setFinal('VAR', 'value')",
			"VAR:=value\n",
		},
		{
			"appendTo('LVAR', 'item2')\nappendTo('LVAR', 'item1')",
			"LVAR:=item1 item2\n",
		},
		{
			"setFinal('VAR', 'value1')\nset('VAR', 'value2')",
			"VAR has been already set and cannot be changed",
		},
		{
			"set('V2', 'value2')\nappendTo('V1', 'value1')",
			"V1:=value1\nV2:=value2\n",
		},
		{
			`loadGenerated('/bin/echo',['set("VAR","value")'])`,
			"VAR:=value\n",
		},
	}
	for _, test := range tests {
		err := Run("string", test.src, nil)
		if err == nil {
			buf := bytes.NewBufferString("")
			PrintConfig(buf)
			out := buf.String()
			if out != test.want {
				t.Errorf("got:\n%swant:\n%s", out, test.want)
			}
		} else {
			out := err.Error()
			if out != test.want {
				t.Errorf("got error: %s\nwant: %s", out, test.want)
			}
		}
	}
}
