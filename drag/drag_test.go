package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestDrag(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./drag", "-s", "1")
	tests = append(tests, test)
	test = exec.Command("./drag", "-s", "1", "-g", "5")
	tests = append(tests, test)
	test = exec.Command("./drag", "-s", "1", "-n", "5")
	tests = append(tests, test)
	test = exec.Command("./drag", "-s", "1", "-t", "3,4")
	tests = append(tests, test)
	test = exec.Command("./drag", "-s", "1", "-t", "-1")
	tests = append(tests, test)
	test = exec.Command("./drag", "-s", "1", "-a")
	tests = append(tests, test)
	test = exec.Command("./drag", "-s", "1", "-G", "-t", "5")
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
			t.Errorf("get:\n%s\nwant:\n%s", get, want)
		}
	}
}
