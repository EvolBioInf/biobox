package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestPps(t *testing.T) {
	var tests []*exec.Cmd
	f := "hom.fasta"
	test := exec.Command("./pps", f)
	tests = append(tests, test)
	test = exec.Command("./pps", "-g", f)
	tests = append(tests, test)
	test = exec.Command("./pps", "-l", "20", f)
	tests = append(tests, test)
	test = exec.Command("./pps", "-d", f)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %s", test)
		}
		f = "r" + strconv.Itoa(i+1) + ".fasta"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s", get, want)
		}
	}
}
