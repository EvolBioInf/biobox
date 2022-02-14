package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
)

var optV = flag.Bool("v", false, "print version & "+
	"program information")
var optR = flag.Bool("r", false, "reverse only")

func scan(r io.Reader, args ...interface{}) {
	optR := args[0].(bool)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		seq.AppendToHeader(" - reverse")
		if optR {
			seq.Reverse()
		} else {
			seq.ReverseComplement()
			seq.AppendToHeader("_complement")
		}
		fmt.Println(seq)
	}
}

func main() {
	u := "revComp [-h] [options] [files]"
	d := "Reverse-complement DNA sequences."
	e := "revComp *.fasta"
	clio.Usage(u, d, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("revComp")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optR)
}
