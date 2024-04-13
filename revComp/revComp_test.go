package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestRevComp(t *testing.T) {
	cmd := exec.Command("./revComp", "test.fasta")
	g, err := cmd.Output()
	if err != nil {
		t.Errorf("coutdn't run %q\n", cmd)
	}
	w, err := ioutil.ReadFile("res1.fasta")
	if err != nil {
		t.Error("couldn't open res1.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./revComp", "-r", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res2.fasta")
	if err != nil {
		t.Error("couldnt' open res2.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
}
