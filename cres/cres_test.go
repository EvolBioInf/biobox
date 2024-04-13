package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestCres(t *testing.T) {
	cmd := exec.Command("./cres", "test.fasta")
	o, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	e, err := ioutil.ReadFile("res1.txt")
	if err != nil {
		t.Error("couldn't open res1.txt")
	}
	if !bytes.Equal(o, e) {
		t.Errorf("wanted:\n%s\ngot:\n%s\n", string(e), string(o))
	}
	cmd = exec.Command("./cres", "-s", "test.fasta")
	o, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	e, err = ioutil.ReadFile("res2.txt")
	if err != nil {
		t.Error("couldn't open res2.txt")
	}
	if !bytes.Equal(o, e) {
		t.Errorf("wanted:\n%s\ngot:\n%s\n", string(e), string(o))
	}
}
