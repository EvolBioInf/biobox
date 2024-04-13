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

func TestPlotSeg(t *testing.T) {
	gf, err := ioutil.TempFile(".", "tmp_*.gp")
	if err != nil {
		log.Fatal("cant open output file")
	}
	g := gf.Name()
	var tests []*exec.Cmd
	f := "test.dat"
	te := exec.Command("./plotSeg", "-s", g, f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-x", "x", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-y", "y", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-x", "x",
		"-y", "y", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-X", "100:500", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-Y", "100:500", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-X", "100:500",
		"-Y", "100:500", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-d", "300,300", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-g",
		"set title \"External Title\"", f)
	tests = append(tests, te)
	te = exec.Command("./plotSeg", "-s", g, "-t", "dumb", f)
	tests = append(tests, te)
	for i, test := range tests {
		err = test.Run()
		if err != nil {
			log.Fatalf("can't run %q", test)
		}
		get, err := ioutil.ReadFile(g)
		if err != nil {
			log.Fatalf("can't read %q", g)
		}
		f = "results/r" + strconv.Itoa(i+1)
		if runtime.GOOS == "darwin" {
			f += "d"
		}
		f += ".gp"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("can't read %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n",
				string(get), string(want))
		}
	}
	err = os.Remove(g)
	if err != nil {
		log.Fatalf("can't delete %q", g)
	}
}
