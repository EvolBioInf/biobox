package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestShuphyl(t *testing.T) {
	tests := []*exec.Cmd{}
	f := "test.nwk"
	test := exec.Command("./shuphyl", "-s", "1", f)
	tests = append(tests, test)
	test = exec.Command("./shuphyl", "-s", "1", "-n", "2", f)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		f = "r" + strconv.Itoa(i+1) + ".txt"
		want, err := os.ReadFile(f)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
