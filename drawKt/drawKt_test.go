package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestDrawKt(t *testing.T) {
	var commands []*exec.Cmd
	c := exec.Command("./drawKt", "ATTT", "ATTC", "AT", "TG", "TT")
	commands = append(commands, c)
	c = exec.Command("./drawKt", "-t", "ATTT", "ATTC", "AT", "TG", "TT")
	commands = append(commands, c)
	c = exec.Command("./drawKt", "-l", "ATTT", "ATTC", "AT", "TG", "TT")
	commands = append(commands, c)

	var names []string
	for i, _ := range commands {
		s := "r" + strconv.Itoa(i+1) + ".txt"
		names = append(names, s)
	}
	for i, command := range commands {
		get, err := command.Output()
		if err != nil {
			t.Errorf("couldn't run %q\n", command)
		}
		want, err := ioutil.ReadFile(names[i])
		if err != nil {
			t.Errorf("couldnt' open %q\n", names[i])
		}
		if !bytes.Equal(want, get) {
			t.Errorf("want:\n%s\nget:\n%s\n", want, get)
		}
	}
}
