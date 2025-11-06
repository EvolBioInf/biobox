package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestCoa(t *testing.T) {
	tests := []*exec.Cmd{}
	s := "3"
	test := exec.Command("./coat", "-s", s)
	tests = append(tests, test)
	test = exec.Command("./coat", "-s", s, "-n", "4")
	tests = append(tests, test)
	test = exec.Command("./coat", "-s", s, "-i", "2")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := os.ReadFile(f)
		if err != nil {
			t.Errorf("couldn't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
