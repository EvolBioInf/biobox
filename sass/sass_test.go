package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestSass(t *testing.T) {
	var tests []*exec.Cmd
	f := "f.fasta"
	test := exec.Command("./sass", f)
	tests = append(tests, test)
	test = exec.Command("./sass", "-r", f)
	tests = append(tests, test)
	test = exec.Command("./sass", "-r", "-M", f)
	tests = append(tests, test)
	test = exec.Command("./sass", "-r", "-t", "20", f)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("couldn't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("couldn't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
