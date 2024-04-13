package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestTranslate(t *testing.T) {
	var tests []*exec.Cmd
	f := "test.fasta"
	test := exec.Command("./translate", f)
	tests = append(tests, test)
	test = exec.Command("./translate", "-f", "1", f)
	tests = append(tests, test)
	test = exec.Command("./translate", "-f", "2", f)
	tests = append(tests, test)
	test = exec.Command("./translate", "-f", "3", f)
	tests = append(tests, test)
	test = exec.Command("./translate", "-f", "-1", f)
	tests = append(tests, test)
	test = exec.Command("./translate", "-f", "-2", f)
	tests = append(tests, test)
	test = exec.Command("./translate", "-f", "-3", f)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("couldn't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".fasta"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("couldn't read %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
