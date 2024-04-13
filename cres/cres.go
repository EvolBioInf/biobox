package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"os"
	"text/tabwriter"
)

var optS = flag.Bool("s", false, "count sequences separately")
var optV = flag.Bool("v", false, "version")
var version, date string

func scan(r io.Reader, args ...interface{}) {
	counts := args[0].([]int64)
	separate := args[1].(bool)
	isFirstSequence := args[2].(*bool)
	scanner := fasta.NewScanner(r)
	for scanner.ScanLine() {
		if scanner.IsHeader() {
			if separate {
				if *isFirstSequence {
					*isFirstSequence = false
				} else {
					write(counts, *optS)
					reset(counts)
				}
				fmt.Printf("%s: ", scanner.Line())
			}
		} else {
			count(counts, scanner.Line())
		}
	}
	count(counts, scanner.Flush())
}
func write(counts []int64, separate bool) {
	var s int64
	for _, v := range counts {
		s += v
	}
	if !separate {
		fmt.Printf("Total: ")
	}
	fmt.Printf("%d\n", s)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 0, 1, ' ', 0)
	if s > 0 {
		fmt.Fprintf(w, "Residue\tCount\tFraction\t\n")
	}
	for i, v := range counts {
		if v > 0 {
			fmt.Fprintf(w, "%c\t%d\t%.3g\t\n", i, v,
				float64(v)/float64(s))
		}
	}
	w.Flush()
}
func reset(counts []int64) {
	for i, _ := range counts {
		counts[i] = 0
	}
}
func count(counts []int64, data []byte) {
	for _, c := range data {
		counts[c]++
	}
}
func main() {
	util.PrepLog("cres")
	u := "cres [-h] [options] [files]"
	p := "Count residues in input."
	e := "cres -s *.fasta"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("cres")
	}
	files := flag.Args()
	counts := make([]int64, 256)
	isFirstSequence := true
	clio.ParseFiles(files, scan, counts, *optS, &isFirstSequence)
	write(counts, *optS)
}
