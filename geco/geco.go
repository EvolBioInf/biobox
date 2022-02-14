package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type geneticCode struct {
	codons    []string
	mutants   map[string][]string
	codon2int map[string]int
	int2int   []int
	int2aa    []string
}

func (g *geneticCode) aa(codon string) string {
	ai1 := g.codon2int[codon]
	ai2 := g.int2int[ai1]
	return g.int2aa[ai2]
}
func (g *geneticCode) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 1, 0, 2, ' ', 0)
	dna := "TCAG"
	for i := 0; i < 4; i++ {
		fmt.Fprintf(w, "\t %c", dna[i])
	}
	fmt.Fprint(w, "\n")
	for i := 0; i < 4; i++ {
		fmt.Fprintf(w, "%c", dna[i])
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				c := dna[i:i+1] + dna[k:k+1] +
					dna[j:j+1]
				fmt.Fprintf(w, "\t%s", g.aa(c))
			}
			fmt.Fprintf(w, "\t%c\n", dna[j])
		}
	}
	w.Flush()
	return buf.String()
}
func newGeneticCode() *geneticCode {
	gc := new(geneticCode)
	gc.codons = make([]string, 0)
	gc.mutants = make(map[string][]string)
	gc.codon2int = make(map[string]int)
	gc.int2int = make([]int, 0)
	gc.int2aa = make([]string, 21)
	dna := "TCAG"
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				codon := dna[i : i+1]
				codon += dna[j : j+1]
				codon += dna[k : k+1]
				gc.codons = append(gc.codons, codon)
			}
		}
	}
	b := make([]byte, 3)
	for _, codon := range gc.codons {
		mutants := make([]string, 0)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				b[j] = codon[j]
			}
			for j := 0; j < 4; j++ {
				b[i] = dna[j]
				if b[i] != codon[i] {
					mutants = append(mutants, string(b))
				}
			}
		}
		gc.mutants[codon] = mutants
	}
	aa := "FLSYCWPHQRIMTNKVADEG*"
	ai := make(map[byte]int)
	for i, a := range aa {
		ai[byte(a)] = i
	}
	aaTab := "FFLLSSSSYY**CC*W" +
		"LLLLPPPPHHQQRRRR" +
		"IIIMTTTTNNKKSSRR" +
		"VVVVAAAADDEEGGGG"
	for i, codon := range gc.codons {
		gc.codon2int[codon] = ai[aaTab[i]]
	}
	for _, a := range aa {
		gc.int2int = append(gc.int2int, ai[byte(a)])
	}
	names := []string{
		"Phe", "Leu", "Ser", "Tyr", "Cys",
		"Trp", "Pro", "His", "Gln", "Arg",
		"Ile", "Met", "Thr", "Asn", "Lys",
		"Val", "Ala", "Asp", "Glu", "Gly",
		"Ter"}
	for i := 0; i < 21; i++ {
		gc.int2aa[i] = names[i]
	}
	return gc
}
func scan(r io.Reader, args ...interface{}) {
	n := args[0].(int)
	gc := args[1].(*geneticCode)
	aap := make(map[string]float64)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if fields[0][0] == '#' {
			continue
		}
		aa := fields[0]
		x, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			log.Fatalf("can't convert %q", fields[1])
		}
		aap[aa] = x
	}
	if n == 0 {
		d := meanDiff(gc, aap)
		fmt.Printf("%sd: %.4g\n", gc, d)
	} else {
		for i := 0; i < n; i++ {
			rand.Shuffle(20, func(i, j int) {
				gc.int2int[i], gc.int2int[j] =
					gc.int2int[j], gc.int2int[i]
			})
			d := meanDiff(gc, aap)
			fmt.Printf("%sd: %.4g\n", gc, d)
		}
	}
}
func meanDiff(gc *geneticCode, aap map[string]float64) float64 {
	var d, c float64
	for _, codon := range gc.codons {
		if gc.codon2int[codon] == 20 {
			continue
		}
		aa := gc.aa(codon)
		x := aap[aa]
		mutants := gc.mutants[codon]
		for _, mutant := range mutants {
			if gc.codon2int[mutant] == 20 {
				continue
			}
			aa := gc.aa(mutant)
			y := aap[aa]
			d += (x - y) * (x - y)
			c++
		}
	}
	return d / c
}
func main() {
	u := "geco [-h] [option]... property.dat"
	p := "Explore the genetic code."
	e := "geco -n 10000 polarity.dat | grep '^d'"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optN = flag.Int("n", 0, "number of iterations")
	var optS = flag.Int("s", 0, "seed of random number generator")
	flag.Parse()
	if *optV {
		util.PrintInfo("geco")
	}
	if *optN > 0 {
		seed := int64(*optS)
		if seed == 0 {
			seed = time.Now().UnixNano()
		}
		rand.Seed(seed)
	}
	files := flag.Args()
	gc := newGeneticCode()
	clio.ParseFiles(files, scan, *optN, gc)
}
