package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"math/rand"
	"strconv"
	"time"
)

var optV = flag.Bool("v", false, "print version & "+
	"program information")
var optL = flag.Int("l", 100, "sequence length")
var optN = flag.Int("n", 1, "number of sequences")
var optG = flag.Float64("g", 0.5, "G/C content")
var optS = flag.Int("s", 0, "seed for random number generator; "+
	"default: internal")

func main() {
	u := "ranseq [-h] [options]"
	d := "Generate random sequence."
	e := "ranseq -l 1000"
	clio.Usage(u, d, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("ranseq")
	}
	var r *rand.Rand
	if *optS != 0 {
		r = rand.New(rand.NewSource(int64(*optS)))
	} else {
		t := time.Now().UnixNano()
		r = rand.New(rand.NewSource(t))
	}
	var s []byte
	var c byte
	for i := 0; i < *optN; i++ {
		s = s[:0]
		h := "Rand" + strconv.Itoa(i+1)
		for j := 0; j < *optL; j++ {
			if r.Float64() < *optG {
				if r.Float64() < 0.5 {
					c = 'G'
				} else {
					c = 'C'
				}
			} else {
				if r.Float64() < 0.5 {
					c = 'A'
				} else {
					c = 'T'
				}
			}
			s = append(s, c)
		}
		seq := fasta.NewSequence(h, s)
		fmt.Println(seq)
	}
}
