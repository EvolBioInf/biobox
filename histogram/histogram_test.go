package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestHistogram(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./histogram", "-r", "0:16", "-b", "16",
		"test.dat")
	tests = append(tests, test)
	test = exec.Command("./histogram", "-r", "0:16", "-b", "16",
		"-f", "test.dat")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err.Error())
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Error(err.Error())
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
