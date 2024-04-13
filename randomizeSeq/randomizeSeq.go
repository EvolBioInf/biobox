package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"math/rand"
	"time"
)

var optV = flag.Bool("v", false, "version")
var optS = flag.Int("s", 0, "seed for random number generator; "+
	"default: internal")

func scan(r io.Reader, args ...interface{}) {
	rn := args[0].(*rand.Rand)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		seq.Shuffle(rn)
		seq.AppendToHeader(" - SHUFFLED")
		fmt.Println(seq)
	}
}
func main() {
	util.PrepLog("randomizeSeq")
	u := "randomizeSeq [-h] [options] [files]"
	p := "Shuffle sequences."
	e := "randomizeSeq *.fasta"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("randomizeSeq")
	}
	var rn *rand.Rand
	if *optS != 0 {
		rn = rand.New(rand.NewSource(int64(*optS)))
	} else {
		t := time.Now().UnixNano()
		rn = rand.New(rand.NewSource(t))
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, rn)
}
