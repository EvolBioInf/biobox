package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestMaf(t *testing.T) {
	var tests []*exec.Cmd
	f := "t.fasta"
	test := exec.Command("./maf", f)
	tests = append(tests, test)
	test = exec.Command("./maf", "-n", f)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		want = append(want, '\n')
		if err != nil {
			t.Errorf("can't open %q", f)
		}
		if bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s", get, want)
		}
	}
}
