package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestAl(t *testing.T) {
	var tests []*exec.Cmd
	test := exec.Command("./al", "-m", "BLOSUM62", "s1.fasta",
		"s2.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "dmAdhAdhdup.fasta",
		"dgAdhAdhdup.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-o", "o1.fasta", "o2.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-l", "dmAdhAdhdup.fasta",
		"dgAdhAdhdup.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-l", "-n", "3", "dmAdhAdhdup.fasta",
		"dgAdhAdhdup.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-P", "v", "s3.fasta", "s4.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-P", "e", "s3.fasta", "s4.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-P", "f", "s3.fasta", "s4.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-P", "g", "s3.fasta", "s4.fasta")
	tests = append(tests, test)
	test = exec.Command("./al", "-P", "t", "s3.fasta", "s4.fasta")
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
