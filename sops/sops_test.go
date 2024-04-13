package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestSops(t *testing.T) {
	var tests []*exec.Cmd
	f := "test.fasta"
	test := exec.Command("./sops", "-a", "2", f)
	tests = append(tests, test)
	test = exec.Command("./sops", "-i", "-2", f)
	tests = append(tests, test)
	test = exec.Command("./sops", "-g", "-1", f)
	tests = append(tests, test)
	test = exec.Command("./sops", "-m", "sm.txt", f)
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
