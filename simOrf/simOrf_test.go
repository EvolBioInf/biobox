package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestSimOrf(t *testing.T) {
	var tests []*exec.Cmd
	cmd := exec.Command("./simOrf", "-s", "23")
	tests = append(tests, cmd)
	cmd = exec.Command("./simOrf", "-s", "23", "-n", "20")
	tests = append(tests, cmd)
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
