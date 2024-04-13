package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestNumAl(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./numAl", "10", "10")
	tests = append(tests, test)
	test = exec.Command("./numAl", "-t", "10", "10")
	tests = append(tests, test)
	test = exec.Command("./numAl", "-p", "3", "3")
	tests = append(tests, test)
	test = exec.Command("./numAl", "-p", "-t", "3", "3")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if i < 2 {
			get = get[:24]
		}
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
