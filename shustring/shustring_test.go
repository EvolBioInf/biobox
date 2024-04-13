package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestShustring(t *testing.T) {
	var commands []*exec.Cmd
	p := "./shustring"
	f := "test.fasta"
	c := exec.Command(p, f)
	commands = append(commands, c)
	c = exec.Command(p, "-l", f)
	commands = append(commands, c)
	c = exec.Command(p, "-s", "1", f)
	commands = append(commands, c)
	c = exec.Command(p, "-r", f)
	commands = append(commands, c)
	c = exec.Command(p, "-q", f)
	commands = append(commands, c)
	var results []string
	for i, _ := range commands {
		name := "r" + strconv.Itoa(i+1) + ".txt"
		results = append(results, name)
	}
	for i, command := range commands {
		get, err := command.Output()
		if err != nil {
			t.Errorf("couldn't run %q\n", command)
		}
		want, err := ioutil.ReadFile(results[i])
		if err != nil {
			t.Errorf("couldn't open %q\n", results[i])
		}
		if !bytes.Equal(want, get) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}
}
