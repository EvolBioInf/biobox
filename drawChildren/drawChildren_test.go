package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestDrawChildren(t *testing.T) {
	test := exec.Command("./drawChildren", "test.txt")
	get, err := test.Output()
	if err != nil {
		t.Error(err)
	}
	want, err := os.ReadFile("r.txt")
	if err != nil {
		t.Errorf("couldn't open %q", "r.txt")
	}
	if !bytes.Equal(get, want) {
		t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
	}
}
