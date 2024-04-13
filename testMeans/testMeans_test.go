package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestTestMeans(t *testing.T) {
	tests := make([]*exec.Cmd, 0)
	test := exec.Command("./testMeans", "d1.txt", "d2.txt")
	tests = append(tests, test)
	test = exec.Command("./testMeans", "-u", "d1.txt", "d2.txt")
	tests = append(tests, test)
	test = exec.Command("./testMeans", "-s", "3", "-m", "1000",
		"d1.txt", "d2.txt")
	tests = append(tests, test)
	results := make([]string, 0)
	for i, _ := range tests {
		r := "r" + strconv.Itoa(i+1) + ".txt"
		results = append(results, r)
	}
	for i, test := range tests {
		want, err := ioutil.ReadFile(results[i])
		if err != nil {
			t.Errorf("couldn't open %q\n", results[i])
		}
		get, err := test.Output()
		if err != nil {
			t.Errorf("couldn't run %q\n", test)
		}
		if !bytes.Equal(want, get) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}
}
