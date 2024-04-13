package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestNaiveMatcher(t *testing.T) {
	cmd := exec.Command("./naiveMatcher", "ATTA",
		"dmAdhAdhdup.fasta", "dgAdhAdhdup.fasta")
	get, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	want, err := ioutil.ReadFile("r1.txt")
	if err != nil {
		t.Errorf("couldn't open r1.txt\n")
	}
	if !bytes.Equal(get, want) {
		t.Errorf("want:\n%s\nget:\n%s\n", want, get)
	}
	cmd = exec.Command("./naiveMatcher", "-p", "p.fasta",
		"dmAdhAdhdup.fasta", "dgAdhAdhdup.fasta")
	get, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	want, err = ioutil.ReadFile("r2.txt")
	if err != nil {
		t.Errorf("couldn't open r1.txt\n")
	}
	if !bytes.Equal(get, want) {
		t.Errorf("want:\n%s\nget:\n%s\n", want, get)
	}
}
