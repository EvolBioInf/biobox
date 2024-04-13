package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestKerror(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./kerror", "-l", "q.fasta", "s.fasta")
	tests = append(tests, test)
	test = exec.Command("./kerror", "-k", "6", "q.fasta", "s.fasta")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("want:\n%s\nget:\n%s", get, want)
		}
	}
}
