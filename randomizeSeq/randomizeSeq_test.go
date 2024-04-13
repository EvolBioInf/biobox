package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestRandomizeSeq(t *testing.T) {
	cmd := exec.Command("./randomizeSeq", "-s", "13", "test.fasta")
	g, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err := ioutil.ReadFile("shuf.fasta")
	if err != nil {
		t.Errorf("couldn't open file %q\n", "shuf.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
}
