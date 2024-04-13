package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestMtf(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./mtf", "t1.fasta")
	tests = append(tests, test)
	test = exec.Command("./mtf", "-d", "r1.fasta")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".fasta"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s", get, want)
		}
	}
}
