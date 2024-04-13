package main

import (
	"os/exec"
	"strconv"
	"testing"
)

func TestRpois(t *testing.T) {
	var tests []*exec.Cmd
	seeds := []int{1, 2, 3, 4, 5}
	for i := 0; i < len(seeds); i++ {
		s := strconv.Itoa(seeds[i])
		cmd := exec.Command("./rpois", "-s", s)
		tests = append(tests, cmd)
	}
	want := []string{"3", "0", "3", "0", "3"}
	for i, test := range tests {
		get, err := test.Output()
		get = get[0 : len(get)-1]
		if err != nil {
			t.Error(err.Error())
		}
		if string(get) != want[i] {
			t.Errorf("get: %s\nwant: %s\n",
				get, want[i])
		}
	}
}
