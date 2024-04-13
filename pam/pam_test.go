package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestPam(t *testing.T) {
	commands := make([]*exec.Cmd, 0)
	c := exec.Command("./pam", "-n", "170", "pam1.txt")
	commands = append(commands, c)
	c = exec.Command("./pam", "-a", "aa.txt", "p170.txt")
	commands = append(commands, c)
	c = exec.Command("./pam", "p170n.txt")
	commands = append(commands, c)
	c = exec.Command("./pam", "-b", "0.3333", "p170n.txt")
	commands = append(commands, c)
	results := make([]string, len(commands))
	for i, _ := range commands {
		results[i] = "r" + strconv.Itoa(i+1) + ".txt"
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
		if !bytes.Equal(want, get) {
			t.Errorf("%s\nwant:\n%s\nget:\n%s\n", cmd, want, get)
		}

	}
}
