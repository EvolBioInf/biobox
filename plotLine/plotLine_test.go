package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"testing"
)

func TestPlotLine(t *testing.T) {
	var tests []*exec.Cmd
	gf, err := ioutil.TempFile(".", "tmp_*.gp")
	if err != nil {
		log.Fatal("can't open script file")
	}
	g := gf.Name()
	test := exec.Command("./plotLine", "-s", g, "test2.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g, "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g, "-P", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g, "-L", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-x", "x", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-y", "y", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-x", "x", "-y", "y", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-p", "test.ps", "-d", "340,340", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-l", "x", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-l", "y", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-l", "xy", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-X", "0.1:10", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-Y", "0.2:100", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g, "-X", "0.1:10",
		"-Y", "0.2:100", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-X", "0.1:10", "-l", "x", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-Y", "0.2:100", "-l", "x", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-X", "0.1:10", "-l", "xy", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-X", "0.1:10", "-Y", "0.2:100", "-l", "xy", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-u", "x", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-u", "y", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-u", "xy", "test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-g", "set title \"External Title\"",
		"test3.dat")
	tests = append(tests, test)
	test = exec.Command("./plotLine", "-s", g,
		"-t", "dumb", "test3.dat")
	tests = append(tests, test)
	for i, test := range tests {
		err := test.Run()
		if err != nil {
			log.Fatalf("can't run %q", test)
		}
		get, err := ioutil.ReadFile(g)
		f := "results/r" + strconv.Itoa(i+1)
		if runtime.GOOS == "darwin" {
			f += "d"
		}
		f += ".gp"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("can't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("%s:\nget:\n%s\nwant:\n%s\n",
				test, string(get), string(want))
		}
	}
	err = os.Remove(g)
	if err != nil {
		log.Fatalf("can't delete %q", g)
	}
}
