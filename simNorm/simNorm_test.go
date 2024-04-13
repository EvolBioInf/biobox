package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestSimNorm(t *testing.T) {
	tests := make([]*exec.Cmd, 0)
	test := exec.Command("./simNorm", "-s", "3")
	tests = append(tests, test)
	test = exec.Command("./simNorm", "-s", "3", "-i", "3")
	tests = append(tests, test)
	test = exec.Command("./simNorm", "-s", "3", "-m", "10.1")
	tests = append(tests, test)
	test = exec.Command("./simNorm", "-s", "3", "-d", "2.5")
	tests = append(tests, test)
	results := make([]string, 0)
	for i, _ := range tests {
		r := "r" + strconv.Itoa(i+1) + ".txt"
		results = append(results, r)
	}
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("couldn't run %q\n", test)
		}
		want, err := ioutil.ReadFile(results[i])
		if err != nil {
			t.Errorf("couldn't open %q\n", results[i])
		}
		if !bytes.Equal(want, get) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}
}
