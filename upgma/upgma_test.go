package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestUpgma(t *testing.T) {
	cmd := exec.Command("./upgma", "-m", "test.phy")
	get, err := cmd.Output()
	if err != nil {
		t.Errorf("can't run %q", cmd)
	}
	want, err := ioutil.ReadFile("r.txt")
	if err != nil {
		t.Errorf("can't open r.txt")
	}
	if !bytes.Equal(get, want) {
		t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
	}
}
