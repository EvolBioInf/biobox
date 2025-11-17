package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestPickChildren(t *testing.T) {
	tests := make([]*exec.Cmd, 0)
	s := "1"
	test := exec.Command("./pickChildren", "-s", s)
	tests = append(tests, test)
	test = exec.Command("./pickChildren", "-s", s, "-n", "5",
		"-t", "test.nwk")
	tests = append(tests, test)
	test = exec.Command("cat", "test.nwk")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		file := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("couldn't open %q", file)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
