package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestBwt(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./bwt", "t1.fasta")
	tests = append(tests, test)
	test = exec.Command("./bwt", "-d", "t2.fasta")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %s", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".fasta"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't read %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s", get, want)
		}
	}
}
