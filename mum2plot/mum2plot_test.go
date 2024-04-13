package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestMum2plot(t *testing.T) {
	test := exec.Command("./mum2plot", "test.mum")
	get, err := test.Output()
	if err != nil {
		t.Error(err.Error())
	}
	want, err := ioutil.ReadFile("r.txt")
	if !bytes.Equal(get, want) {
		t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
	}
}
