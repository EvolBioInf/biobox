package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestDnaDist(t *testing.T) {
	tests := make([]*exec.Cmd, 0)
	c := exec.Command("./dnaDist", "test.fa")
	tests = append(tests, c)
	c = exec.Command("./dnaDist", "-k", "test.fa")
	tests = append(tests, c)
	c = exec.Command("./dnaDist", "pr.fa")
	tests = append(tests, c)
	c = exec.Command("./dnaDist", "-k", "pr.fa")
	tests = append(tests, c)
	c = exec.Command("./dnaDist", "-b", "5", "-s", "3", "pr.fa")
	tests = append(tests, c)
	results := make([]string, len(tests))
	for i, _ := range tests {
		results[i] = "r" + strconv.Itoa(i+1) + ".txt"
	}
	for i, c := range tests {
		get, err := c.Output()
		if err != nil {
			t.Errorf("couldn't run %q\n", c)
		}
		want, err := ioutil.ReadFile(results[i])
		if err != nil {
			t.Errorf("couldn't open %q\n", results[i])
		}
		if !bytes.Equal(get, want) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}

}
