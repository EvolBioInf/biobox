package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestVar(t *testing.T) {
	cmd := exec.Command("./var", "data1.txt")
	g, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err := ioutil.ReadFile("res1.txt")
	if !bytes.Equal(g, w) {
		t.Errorf("want:\n%s\nget:\n%s\n", w, g)
	}
	cmd = exec.Command("./var", "data1.txt", "data2.txt")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res2.txt")
	if !bytes.Equal(g, w) {
		t.Errorf("want\n%ss\nget:\n%s\n", w, g)
	}
}
