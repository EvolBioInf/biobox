package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestRanseq(t *testing.T) {
	cmd := exec.Command("./ranseq", "-s", "13")
	g, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err := ioutil.ReadFile("res1.fasta")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget\n%s\n", w, g)
	}
	cmd = exec.Command("./ranseq")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res1.fasta")
	if bytes.Equal(g, w) {
		t.Errorf("don't want:\n%s\nbut do get\n%s\n", w, g)
	}
	cmd = exec.Command("./ranseq", "-s", "13", "-n", "2")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res2.fasta")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget\n%s\n", w, g)
	}
	cmd = exec.Command("./ranseq", "-s", "13", "-g", "0.3")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res3.fasta")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget\n%s\n", w, g)
	}
}
