package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestGetSeq(t *testing.T) {
	cmd := exec.Command("./getSeq", "Seq1", "test.fasta")
	o, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	e, err := ioutil.ReadFile("res1.txt")
	if !bytes.Equal(o, e) {
		t.Errorf("want:\n%s\nget:\n%s\n", e, o)
	}
	cmd = exec.Command("./getSeq", "1$", "test.fasta")
	o, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	e, err = ioutil.ReadFile("res2.txt")
	if !bytes.Equal(o, e) {
		t.Errorf("want:\n%s\nget:\n%s\n", e, o)
	}
	cmd = exec.Command("./getSeq", "[123]$", "test.fasta")
	o, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	e, err = ioutil.ReadFile("res3.txt")
	if !bytes.Equal(o, e) {
		t.Errorf("want:\n%s\nget:\n%s\n", e, o)
	}
	cmd = exec.Command("./getSeq", "-c", "[123]$",
		"test.fasta")
	o, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	e, err = ioutil.ReadFile("res4.txt")
	if !bytes.Equal(o, e) {
		t.Errorf("want:\n%s\nget:\n%s\n", e, o)
	}
}
