package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestWatterson(t *testing.T) {
	cmd := exec.Command("./watterson", "-n", "10", "-t", "20")
	g, err := cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err := ioutil.ReadFile("res1.txt")
	if err != nil {
		t.Errorf("couldnt' open res1.txt")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("get:\n%s\nwant:\n%s\n", g, w)
	}
	cmd = exec.Command("./watterson", "-n", "10", "-t", "20", "-a")
	g, err = cmd.Output()
	if err != nil {
		t.Errorf("couldn't run %q\n", cmd)
	}
	w, err = ioutil.ReadFile("res2.txt")
	if err != nil {
		t.Errorf("couldnt' open res2.txt")
	}
	if !bytes.Equal(g, w) {
		t.Errorf("get:\n%s\nwant:\n%s\n", g, w)
	}
}
