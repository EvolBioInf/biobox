package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestCutSeq(t *testing.T) {
	cmd := exec.Command("./cutSeq", "-r", "10-20", "test.fasta")
	g, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err := ioutil.ReadFile("res1.fasta")
	if err != nil {
		t.Errorf("couldn't open res1.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}

	cmd = exec.Command("./cutSeq", "-r", "10-20,25-50", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res2.fasta")
	if err != nil {
		t.Errorf("couldn't open res2.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./cutSeq", "-j", "-r", "10-20,25-50", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res3.fasta")
	if err != nil {
		t.Errorf("couldn't open res3.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./cutSeq", "-f", "coord1.txt", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res1.fasta")
	if err != nil {
		t.Errorf("couldn't open res1.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./cutSeq", "-f", "coord2.txt", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res2.fasta")
	if err != nil {
		t.Errorf("couldn't open res2.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./cutSeq", "-j", "-f", "coord2.txt", "test.fasta")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res3.fasta")
	if err != nil {
		t.Errorf("couldn't open res3.fasta")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
}
