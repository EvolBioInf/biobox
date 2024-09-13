package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestSw(t *testing.T) {
	var tests []*exec.Cmd
	ds := "test1.dat"
	test := exec.Command("./sw", "-w", "5", ds)
	tests = append(tests, test)
	test = exec.Command("./sw", "-w", "5", "-k", "2", ds)
	tests = append(tests, test)
	ds = "test2.dat"
	test = exec.Command("./sw", "-w", "100", ds)
	tests = append(tests, test)
	test = exec.Command("./sw", "-w", "100", "-k", "1", ds)
	tests = append(tests, test)
	test = exec.Command("./sw", "-w", "100", "-k", "20", ds)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := os.ReadFile(f)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
