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

func scan(r io.Reader, args ...interface{}) {
	mat := args[0].(*pal.ScoreMatrix)
	g := args[1].(float64)
	sc := fasta.NewScanner(r)
	var msa [][]byte
	for sc.ScanSequence() {
		msa = append(msa, sc.Sequence().Data())
	}
	for i := 1; i < len(msa); i++ {
		l1 := len(msa[i-1])
		l2 := len(msa[i])
		if l1 != l2 {
			m := "sequence %d has length %d, " +
				"but sequence %d has length %d; " +
				"this doesn't look like an alignment"
			log.Fatalf(m, i, l1, i+1, l2)
		}
	}
	m := len(msa)
	n := len(msa[0])
	s := 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < m-1; j++ {
			for k := j + 1; k < m; k++ {
				r1 := msa[j][i]
				r2 := msa[k][i]
				if r1 == '-' && r2 == '-' {
					continue
				}
				if r1 == '-' || r2 == '-' {
					s += g
				} else {
					s += mat.Score(r1, r2)
				}
			}
		}
	}
	fmt.Printf("sum-of-pairs_score\t%g\n", s)
}

func main() {
	util.PrepLog("sops")
	u := "sops [-h] [option]... [foo.fasta]..."
	p := "Calculate the sum-of-pairs score of a multiple sequence alignment."
	e := "sops msa.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optA = flag.Float64("a", 1, "match")
	var optI = flag.Float64("i", -3, "mismatch")
	var optM = flag.String("m", "", "score matrix")
	var optG = flag.Float64("g", -2, "gap")
	flag.Parse()
	if *optV {
		util.PrintInfo("sops")
	}
	var mat *pal.ScoreMatrix
	if *optM == "" {
		mat = pal.NewScoreMatrix(*optA, *optI)
	} else {
		f, err := os.Open(*optM)
		if err != nil {
			log.Fatalf("couldn't open score matrix %q",
				(*optM))
		}
		defer f.Close()
		mat = pal.ReadScoreMatrix(f)
	}
	f := flag.Args()
	clio.ParseFiles(f, scan, mat, *optG)
}
