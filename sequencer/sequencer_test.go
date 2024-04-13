package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestSequencer(t *testing.T) {
	tests := make([]*exec.Cmd, 0)
	f := "test.fasta"
	test := exec.Command("./sequencer", "-s", "3", f)
	tests = append(tests, test)
	test = exec.Command("./sequencer", "-s", "3", "-p", f)
	tests = append(tests, test)
	test = exec.Command("./sequencer", "-s", "3", "-c", "2", f)
	tests = append(tests, test)
	test = exec.Command("./sequencer", "-s", "3", "-r", "50", f)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %q", test)
		}
		f = "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s", get, want)
		}
	}
}
