package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestGenTree(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./genTree", "-s", "13")
	tests = append(tests, test)
	test = exec.Command("./genTree", "-s", "13", "-a")
	tests = append(tests, test)
	test = exec.Command("./genTree", "-s", "13", "-c")
	tests = append(tests, test)
	test = exec.Command("./genTree", "-s", "13", "-i", "2")
	tests = append(tests, test)
	test = exec.Command("./genTree", "-s", "13", "-l")
	tests = append(tests, test)
	test = exec.Command("./genTree", "-s", "13", "-n", "9")
	tests = append(tests, test)
	test = exec.Command("./genTree", "-s", "13", "-t", "500")
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
