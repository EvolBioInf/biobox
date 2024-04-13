package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
)

var optV = flag.Bool("v", false, "version")
var optL = flag.Int("l", fasta.DefaultLineLength, "line length, "+
	"< 1 for unbroken lines")

func scan(r io.Reader, args ...interface{}) {
	l := args[0].(int)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		se := sc.Sequence()
		se.SetLineLength(l)
		fmt.Println(se)
	}
}
func main() {
	util.PrepLog("wrapSeq")
	u := "wrapSeq [-h] [options] [files]"
	p := "Wrap lines of sequence data."
	e := "wrapSeq -l 50 *.fasta"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("wrapSeq")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optL)
}
