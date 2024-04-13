package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestFasta2tab(t *testing.T) {
	var tests []*exec.Cmd
	f := "test.fasta"
	cmd := exec.Command("./fasta2tab", f)
	tests = append(tests, cmd)
	cmd = exec.Command("./fasta2tab", "-d", "\\t", f)
	tests = append(tests, cmd)
	cmd = exec.Command("./fasta2tab", "-d", "\\n", f)
	tests = append(tests, cmd)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err.Error())
		}
		f = "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Error(err.Error())
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
