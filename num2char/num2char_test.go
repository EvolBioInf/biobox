package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestNum2char(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./num2char", "mtf.fasta")
	tests = append(tests, test)
	test = exec.Command("./num2char", "-d", "r1.txt")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("couldn't run %s", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("couldn't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:%s\n", get, want)
		}
	}
}
