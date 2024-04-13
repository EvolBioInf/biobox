package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"testing"
)

func TestPlotTree(t *testing.T) {
	var tests []*exec.Cmd
	gf, err := ioutil.TempFile(".", "tmp_*.gp")
	if err != nil {
		t.Error("can't open temp file")
	}
	g := gf.Name()
	f := "newick.nwk"
	test := exec.Command("./plotTree", "-r", "-s", g, f)
	tests = append(tests, test)
	test = exec.Command("./plotTree", "-u", "-s", g, f)
	tests = append(tests, test)
	test = exec.Command("./plotTree", "-r", "-s", g, "-n", f)
	tests = append(tests, test)
	test = exec.Command("./plotTree", "-u", "-s", g, "-n", f)
	tests = append(tests, test)
	test = exec.Command("./plotTree", "-t", "dumb", "-s", g, f)
	tests = append(tests, test)
	for i, test := range tests {
		err := test.Run()
		if err != nil {
			t.Errorf("couldn't run %q", test)
		}
		get, err := ioutil.ReadFile(g)
		if err != nil {
			t.Errorf("couldn't open %q", g)
		}
		f := "results/r" + strconv.Itoa(i+1)
		if runtime.GOOS == "darwin" {
			f += "d"
		}
		f += ".gp"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("couldn't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
	err = os.Remove(g)
	if err != nil {
		t.Errorf("can't remove %q", g)
	}
}
