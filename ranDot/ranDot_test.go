package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestRanDot(t *testing.T) {
	var tests []*exec.Cmd
	var test *exec.Cmd
	args := []string{"-s", "13"}
	test = exec.Command("./ranDot", args...)
	tests = append(tests, test)
	args = append(args, "-C", "lightgray")
	test = exec.Command("./ranDot", args...)
	tests = append(tests, test)
	args = append(args, "-c", "lightsalmon")
	test = exec.Command("./ranDot", args...)
	tests = append(tests, test)
	args = append(args, "-n", "11")
	test = exec.Command("./ranDot", args...)
	tests = append(tests, test)
	args = append(args, "-p", "0.5")
	test = exec.Command("./ranDot", args...)
	tests = append(tests, test)
	args = append(args, "-S")
	test = exec.Command("./ranDot", args...)
	tests = append(tests, test)
	args = append(args, "-d")
	test = exec.Command("./ranDot", args...)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err.Error())
		}
		f := "r" + strconv.Itoa(i+1) + ".dot"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Error(err.Error())
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
