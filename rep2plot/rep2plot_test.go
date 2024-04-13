package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestRep2plot(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./rep2plot", "test1.txt")
	tests = append(tests, test)
	test = exec.Command("./rep2plot", "test2.txt")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("cant' open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
