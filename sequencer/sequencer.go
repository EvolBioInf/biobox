package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

type opts struct {
	c, r, R, i, I, e float64
	p, o, S          bool
}

func scan(r io.Reader, args ...interface{}) {
	op := args[0].(*opts)
	rn := args[1].(*rand.Rand)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		n := len(seq.Data())
		se := make([][]byte, 2)
		se[0] = make([]byte, n)
		copy(se[0], seq.Data())
		se[1] = make([]byte, n)
		seq.ReverseComplement()
		copy(se[1], seq.Data())
		cov := int(math.Round(float64(n) * op.c))
		var ns, rc int
		w := bufio.NewWriter(os.Stdout)
		for ns < cov {
			if op.p {
				pos := rn.Intn(n)
				il := int(math.Round(rn.NormFloat64()*op.I + op.i))
				if pos+il < n || op.o {
					rc++
					fmt.Fprintf(w, ">Read%d mate=1\n", rc)
					rl := int(math.Round(rn.NormFloat64()*op.R + op.r))
					if rl < 0 {
						rl *= -1
					}
					for i := pos; i < pos+rl; i++ {
						c := se[0][i%n]
						c = mutate(c, rn, op.e)
						fmt.Fprintf(w, "%c", c)
					}
					fmt.Fprintf(w, "\n")
					ns += rl
					pos = n - (pos + il - 1)
					fmt.Fprintf(w, ">Read%d mate=2\n", rc)
					rl = int(math.Round(rn.NormFloat64()*op.R + op.r))
					if rl < 0 {
						rl *= -1
					}
					for i := pos; i < pos+rl; i++ {
						c := se[1][i%n]
						c = mutate(c, rn, op.e)
						fmt.Fprintf(w, "%c", c)
					}
					fmt.Fprintf(w, "\n")
					ns += rl
				}
			} else {
				pos := rn.Intn(n)
				rl := int(math.Round(rn.NormFloat64()*op.R + op.r))
				if rl < 0 {
					rl *= -1
				}
				strand := 0
				if rn.Float64() < 0.5 && !op.S {
					strand = 1
				}
				if pos+rl <= n || op.o {
					rc++
					fmt.Fprintf(w, ">Read%d\n", rc)
					for i := pos; i < pos+rl; i++ {
						c := se[strand][i%n]
						c = mutate(c, rn, op.e)
						fmt.Fprintf(w, "%c", c)
					}
					fmt.Fprintf(w, "\n")
					ns += rl
				}
			}
		}
		w.Flush()
	}
}

const dna = "ACGT"

func mutate(c byte, r *rand.Rand, e float64) byte {
	if r.Float64() >= e {
		return c
	}
	m := dna[r.Intn(4)]
	for m == c {
		m = dna[r.Intn(4)]
	}
	return m
}
func main() {
	util.PrepLog("sequencer")
	u := "sequencer [-h] [option]... [foo.fasta]..."
	p := "Simulate a DNA sequencing machine."
	e := "sequencer -c 20 foo.fasta"
	clio.Usage(u, p, e)
	var optC = flag.Float64("c", 1.0, "coverage")
	var optR = flag.Float64("r", 100.0, "mean read length")
	var optRR = flag.Float64("R", 0.0, "standard deviation of "+
		"read length")
	var optP = flag.Bool("p", false, "paired end")
	var optI = flag.Float64("i", 500.0, "mean insert length")
	var optII = flag.Float64("I", 0.0, "standard deviation of "+
		"insert length")
	var optE = flag.Float64("e", 0.001, "error rate")
	var optS = flag.Int("s", 0, "seed for random number generator")
	var optO = flag.Bool("o", false, "circular template")
	var optSS = flag.Bool("S", false, "shredder - forward strand only")
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("sequencer")
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rn := rand.New(rand.NewSource(int64(seed)))
	op := new(opts)
	op.c = *optC
	op.r = *optR
	op.R = *optRR
	op.i = *optI
	op.I = *optII
	op.e = *optE
	op.p = *optP
	op.o = *optO
	op.S = *optSS
	files := flag.Args()
	clio.ParseFiles(files, scan, op, rn)
}
