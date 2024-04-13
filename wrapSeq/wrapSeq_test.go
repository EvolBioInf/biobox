package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestWrapSeq(t *testing.T) {
	cmd := exec.Command("./wrapSeq", "test.fasta")
	g, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err := ioutil.ReadFile("test.fasta")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./wrapSeq", "-l", "100", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res1.fasta")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./wrapSeq", "-l", "0", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res1.fasta")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./wrapSeq", "-l", "50", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res2.fasta")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
}
