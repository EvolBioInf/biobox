package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/pal"
	"io"
	"log"
	"os"
)

var optV = flag.Bool("v", false, "print version & "+
	"program information")
var optL = flag.Bool("l", false, "local (default global)")
var optO = flag.Bool("o", false, "overlap (default global)")
var optI = flag.Float64("i", -3, "mismatch")
var optA = flag.Float64("a", 1, "match")
var optM = flag.String("m", "", "file containing score matrix")
var optP = flag.Float64("p", -5, "gap opening")
var optE = flag.Float64("e", -2, "gap extension")
var optN = flag.Int("n", 1, "number of local alignments")
var optLL = flag.Int("L", fasta.DefaultLineLength, "line length")
var optPP = flag.String("P", "", "print programming matrix (v|e|f|g|t)")

func scan(r io.Reader, args ...interface{}) {
	q := args[0].(*fasta.Sequence)
	mat := args[1].(*pal.ScoreMatrix)
	isLocal := *optL
	isOverlap := *optO
	gapO := *optP
	gapE := *optE
	numAl := *optN
	var printMat byte
	if *optPP != "" {
		printMat = []byte(*optPP)[0]
	}
	ll := *optLL
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		s := sc.Sequence()
		if isLocal {
			al := pal.NewLocalAlignment(q, s, mat, gapO, gapE)
			al.SetLineLength(ll)
			for i := 0; i < numAl && al.Align(); i++ {
				if printMat != 0 {
					s := al.PrintMatrix(printMat)
					fmt.Printf(s)
				} else {
					fmt.Println(al)
				}
			}
		} else if isOverlap {
			al := pal.NewOverlapAlignment(q, s, mat, gapO, gapE)
			al.SetLineLength(ll)
			al.Align()
			if printMat != 0 {
				s := al.PrintMatrix(printMat)
				fmt.Printf(s)
			} else {
				fmt.Println(al)
			}
		} else {
			al := pal.NewGlobalAlignment(q, s, mat, gapO, gapE)
			al.SetLineLength(ll)
			al.Align()
			if printMat != 0 {
				s := al.PrintMatrix(printMat)
				fmt.Printf(s)
			} else {
				fmt.Println(al)
			}
		}
	}
}
func main() {
	u := "al [-h] [options] query.fasta [subject files]"
	p := "Align two sequences."
	e := "al query.fasta subject.fasta"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("al")
	}
	m := "-P should be v, e, f, g for the cell element " +
		"or t for the traceback"
	if *optPP != "" {
		if *optPP != "v" && *optPP != "e" &&
			(*optPP) != "f" && *optPP != "g" &&
			(*optPP) != "t" {
			fmt.Println(m)
			os.Exit(-1)
		}
	}
	files := flag.Args()
	if len(files) < 1 {
		fmt.Fprintf(os.Stderr, "please give the name "+
			"of a query file\n")
		os.Exit(0)
	}
	query := files[0]
	subject := files[1:]
	var mat *pal.ScoreMatrix
	if *optM == "" {
		mat = pal.NewScoreMatrix(*optA, *optI)
	} else {
		f, err := os.Open(*optM)
		if err != nil {
			log.Fatalf("couldn't open score matrix %q\n",
				*optM)
		}
		mat = pal.ReadScoreMatrix(f)
		f.Close()
	}
	qf, err := os.Open(query)
	if err != nil {
		log.Fatalf("couldn't open %q\n", query)
	}
	sc := fasta.NewScanner(qf)
	for sc.ScanSequence() {
		q := sc.Sequence()
		clio.ParseFiles(subject, scan, q, mat)
	}
}
