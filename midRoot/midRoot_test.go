package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestMidRoot(t *testing.T) {
	cmd := exec.Command("./midRoot", "-p", "test.nwk")
	get, err := cmd.Output()
	if err != nil {
		t.Errorf("can't run %q", cmd)
	}
	want, err := ioutil.ReadFile("r.txt")
	if err != nil {
		t.Errorf("can't optn r.txt")
	}
	if !bytes.Equal(get, want) {
		t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
	}
}
