package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestRepeater(t *testing.T) {
	var commands []*exec.Cmd
	c := exec.Command("./repeater", "test.fasta")
	commands = append(commands, c)
	c = exec.Command("./repeater", "-m", "13", "test.fasta")
	commands = append(commands, c)
	c = exec.Command("./repeater", "-r", "test.fasta")
	commands = append(commands, c)
	c = exec.Command("./repeater", "-p", "test.fasta")
	commands = append(commands, c)
	c = exec.Command("./repeater", "-s", "test.fasta")
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
			t.Errorf("couldnt' open %q\n", results[i])
		}
		if !bytes.Equal(want, get) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}
}
