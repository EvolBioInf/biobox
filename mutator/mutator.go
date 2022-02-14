package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func scan(r io.Reader, args ...interface{}) {
	alphabet := args[0].(string)
	mu := args[1].(float64)
	pos := args[2].([]int)
	ran := args[3].(*rand.Rand)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		res := seq.Data()
		if len(pos) > 0 {
			for _, p := range pos {
				l := len(res)
				if p < l {
					res[p] = mutate(res[p], ran, alphabet)
				} else {
					fmt.Fprintf(os.Stderr, "trying to mutate "+
						"position %d, but sequence only "+
						"contains %d residues\n", p+1, l)
				}
			}
		} else {
			l := len(res)
			for i := 0; i < l; i++ {
				if ran.Float64() < mu {
					r := ran.Intn(l)
					res[r] = mutate(res[r], ran, alphabet)
				}
			}
		}
		h := seq.Header() + " - mutated"
		ns := fasta.NewSequence(h, res)
		fmt.Println(ns)
	}
}
func mutate(res byte, ran *rand.Rand, alphabet string) byte {
	n := len(alphabet)
	new := res
	for new == res {
		p := ran.Intn(n)
		new = alphabet[p]
	}
	return new
}
func main() {
	u := "mutator [-h] [options] [fasta file(s)]"
	p := "Mutate input sequences."
	e := "mutator -p 1,10,100 foo.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print version & "+
		"program information")
	var optM = flag.Float64("m", 0.01, "mutation rate")
	var optP = flag.String("p", "", "positions to be mutated; "+
		"comma-separated, one-based")
	var optPP = flag.Bool("P", false, "protein instead of DNA")
	var optS = flag.Int("s", 0, "seed for random number genrator; "+
		"default: internal")
	flag.Parse()
	if *optV {
		util.PrintInfo("mutator")
	}
	var positions []int
	if *optP != "" {
		str := strings.Split(*optP, ",")
		for _, ps := range str {
			position, err := strconv.Atoi(ps)
			if err != nil {
				log.Fatalf("couldn't convert %q\n", ps)
			}
			position--
			if position < 0 {
				fmt.Fprintf(os.Stderr, "position %d cannot be mutated\n", position+1)
			} else {
				positions = append(positions, position)
			}
		}
	}
	alphabet := "ACGT"
	if *optPP {
		alphabet = "ACDEFGHIKLMNPQRSTVWY"
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	ran := rand.New(rand.NewSource(seed))
	f := flag.Args()
	clio.ParseFiles(f, scan, alphabet, *optM, positions, ran)
}
