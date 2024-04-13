package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestMutator(t *testing.T) {
	commands := make([]*exec.Cmd, 0)
	c := exec.Command("./mutator", "-s", "3", "dna.fa")
	commands = append(commands, c)
	c = exec.Command("./mutator", "-s", "3", "-p", "0,1,3,100,101",
		"dna.fa")
	commands = append(commands, c)
	c = exec.Command("./mutator", "-s", "3", "-m", "0.2", "dna.fa")
	commands = append(commands, c)
	c = exec.Command("./mutator", "-s", "3", "-P", "pro.fa")
	commands = append(commands, c)
	c = exec.Command("./mutator", "-s", "3", "-n", "2", "dna.fa")
	commands = append(commands, c)
	results := make([]string, len(commands))
	for i, _ := range commands {
		results[i] = "r" + strconv.Itoa(i+1) + ".fa"
	}
	for i, cmd := range commands {
		get, err := cmd.Output()
		if err != nil {
			t.Errorf("couldn't run %q\n", cmd)
		}
		want, err := ioutil.ReadFile(results[i])
		if err != nil {
			t.Errorf("couldn't open %q\n", results[i])
		}
		if !bytes.Equal(get, want) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}
}
