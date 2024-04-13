package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestBlast2dot(t *testing.T) {
	var tests []*exec.Cmd
	f := "test.bl"
	cmd := exec.Command("./blast2dot", "-c", "lightsalmon", f)
	tests = append(tests, cmd)
	cmd = exec.Command("./blast2dot", "-c", "lightsalmon",
		"-C", "lightgray", f)
	tests = append(tests, cmd)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err.Error())
		}
		f = "r" + strconv.Itoa(i+1) + ".dot"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Error(err.Error())
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n",
				string(get), string(want))
		}
	}
}
