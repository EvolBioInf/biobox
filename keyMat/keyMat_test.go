package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestKeyMat(t *testing.T) {
	var commands []*exec.Cmd
	r := "./keyMat"
	p := "ATTT,ATTC,AT,TG,TT"
	f := "patterns.fasta"
	i := "test.fasta"
	c := exec.Command(r, p, i)
	commands = append(commands, c)
	c = exec.Command(r, "-r", p, i)
	commands = append(commands, c)
	c = exec.Command(r, "-p", f, i)
	commands = append(commands, c)
	c = exec.Command(r, "-p", f, "-r", i)
	commands = append(commands, c)
	var files []string
	for i, _ := range commands {
		f := "r" + strconv.Itoa(i+1) + ".txt"
		files = append(files, f)
	}
	for i, command := range commands {
		get, err := command.Output()
		if err != nil {
			t.Errorf("couldn't run %q\n", command)
		}
		want, err := ioutil.ReadFile(files[i])
		if err != nil {
			t.Errorf("couldn't open %q\n", files[i])
		}
		if !bytes.Equal(want, get) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}
}
