package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestDrawGenes(t *testing.T) {
	test := exec.Command("./drawGenes", "t.txt")
	get, err := test.Output()
	if err != nil {
		t.Errorf("can't run %q", test)
	}
	want, err := ioutil.ReadFile("r.txt")
	if err != nil {
		t.Errorf("can't open %q", "r.txt")
	}
	if !bytes.Equal(get, want) {
		t.Errorf("get:\n%s\nwant:\n%s", get, want)
	}
}
