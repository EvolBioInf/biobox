package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestDrawSt(t *testing.T) {
	var tests []*exec.Cmd
	file := "test.fasta"
	cmd := exec.Command("./drawSt", "-s", file)
	tests = append(tests, cmd)
	cmd = exec.Command("./drawSt", "-s", "-i", file)
	tests = append(tests, cmd)
	cmd = exec.Command("./drawSt", "-s", "-n", file)
	tests = append(tests, cmd)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err.Error())
		}
		file = "res" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(file)
		if err != nil {
			t.Error(err.Error())
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n",
				string(get), string(want))
		}
	}
}
