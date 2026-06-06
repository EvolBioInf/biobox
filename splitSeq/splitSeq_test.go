package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestSplitSeq(t *testing.T) {
	var tests []*exec.Cmd
	f := "test.fasta"
	p := "./splitSeq"
	test := exec.Command(p, f)
	tests = append(tests, test)
	test = exec.Command(p, "-l", "3", f)
	tests = append(tests, test)
	test = exec.Command(p, "-l", "3", "-s", "2", f)
	tests = append(tests, test)
	for i, test := range tests {
		g, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		w, err := os.ReadFile(f)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(g, w) {
			t.Errorf("%s - get:\n%s\nwant:\n%s\n", f, g, w)
		}
	}
}
