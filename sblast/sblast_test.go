package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestSblast(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./sblast", "-a", "2",
		"test.fasta", "test.fasta")
	tests = append(tests, test)
	test = exec.Command("./sblast", "-i", "-2",
		"test.fasta", "test.fasta")
	tests = append(tests, test)
	test = exec.Command("./sblast", "-w", "20",
		"test.fasta", "test.fasta")
	tests = append(tests, test)
	test = exec.Command("./sblast", "-s", "20",
		"test.fasta", "test.fasta")
	tests = append(tests, test)
	test = exec.Command("./sblast", "-t", "40",
		"test.fasta", "test.fasta")
	tests = append(tests, test)
	test = exec.Command("./sblast", "-n",
		"test.fasta", "test.fasta")
	tests = append(tests, test)
	test = exec.Command("./sblast", "-l",
		"test.fasta", "test.fasta")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("couldn't run %s\n", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("couldn't open %s\n", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n",
				string(get), string(want))
		}
	}
}
