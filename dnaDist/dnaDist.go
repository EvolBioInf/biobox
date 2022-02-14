package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"math"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
)

type cell struct {
	a, b int
	d    float64
}

func scan(r io.Reader, args ...interface{}) {
	optB := args[0].(int)
	optR := args[1].(bool)
	optU := args[2].(bool)
	optK := args[3].(bool)
	ran := args[4].(*rand.Rand)
	ts := args[5].(util.TransitionTab)
	sc := fasta.NewScanner(r)
	var sa []*fasta.Sequence
	for sc.ScanSequence() {
		sa = append(sa, sc.Sequence())
		if len(sa) > 1 {
			i := len(sa) - 1
			if len(sa[i].Data()) != len(sa[i-1].Data()) {
				fmt.Fprintf(os.Stderr, "this doesn't look "+
					"like an alignment\n")
				os.Exit(-1)
			}
		}
	}
	m := len(sa)
	n := len(sa[0].Data())
	msa := make([][]byte, m)
	for i, s := range sa {
		msa[i] = bytes.ToUpper(s.Data())
	}
	pol := make([]bool, n)
	gap := byte('-')
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			if msa[i][j] != gap && msa[i][j] != msa[0][j] {
				pol[j] = true
				break
			}
		}
	}
	ind := make([]int, n)
	dm := make([][]cell, m)
	for i := 0; i < m; i++ {
		dm[i] = make([]cell, m)
	}
	if optB > 0 {
		for i := 0; i < optB; i++ {
			for j := 0; j < n; j++ {
				ind[j] = ran.Intn(n)
			}
			distMat(dm, msa, pol, ind, optR, optU, optK, ts)
			printDist(dm, sa)
			for i := 0; i < m-1; i++ {
				for j := i + 1; j < m; j++ {
					dm[i][j].a = 0
					dm[i][j].b = 0
				}
			}
		}
	} else {
		for i, _ := range ind {
			ind[i] = i
		}
		distMat(dm, msa, pol, ind, optR, optU, optK, ts)
		printDist(dm, sa)
	}
}
func distMat(dm [][]cell, msa [][]byte, pol []bool,
	ind []int, optR, optU, optK bool,
	ts util.TransitionTab) {
	m := len(msa)
	n := len(msa[0])
	for i := 0; i < n; i++ {
		if pol[ind[i]] {
			for j := 0; j < m-1; j++ {
				c1 := msa[j][i]
				for k := j + 1; k < m; k++ {
					c2 := msa[k][i]
					if c1 != c2 {
						if ts.IsTransition(c1, c2) {
							dm[j][k].a++
						} else {
							dm[j][k].b++
						}
					}
				}
			}
		}
	}
	for i := 0; i < m-1; i++ {
		for j := i + 1; j < m; j++ {
			a := float64(dm[i][j].a) / float64(n)
			b := float64(dm[i][j].b) / float64(n)
			if optR {
				dm[i][j].d = float64(dm[i][j].a + dm[i][j].b)
			} else if optU {
				dm[i][j].d = a + b
			} else if optK {
				dm[i][j].d = -math.Log((1-2*a-b)*math.Sqrt(1-2*b)) / 2
			} else {
				p := a + b
				dm[i][j].d = -0.75 * math.Log(1-4./3.*p)
			}
			dm[j][i].d = dm[i][j].d
		}
	}
}
func printDist(dm [][]cell, sa []*fasta.Sequence) {
	n := len(dm)
	fmt.Printf("%d\n", n)
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 1, ' ', 0)
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "%s\t", sa[i].Header())
		for j := 0; j < n; j++ {
			if math.Signbit(dm[i][j].d) {
				fmt.Fprintf(w, "%.6g\t", 0.0)
			} else {
				fmt.Fprintf(w, "%.6g\t", dm[i][j].d)
			}
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	fmt.Printf("%s", buffer)
}
func main() {
	u := "dnaDist [-h] [options] [file(s)]"
	p := "Calculate distances between DNA sequences."
	e := "dnaDist foo.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print version "+
		"program information")
	var optR = flag.Bool("r", false, "raw mismatches")
	var optU = flag.Bool("u", false, "uncorrected mismatches")
	var optK = flag.Bool("k", false, "Kimura distances (default: Jukes-Cantor)")
	var optB = flag.Int("b", 0, "number of bootstrap replicates")
	var optS = flag.Int("s", 0, "seed for random number generator "+
		"(default: internal)")
	flag.Parse()
	if *optV {
		util.PrintInfo("dnaDist")
	}
	if *optB < 0 {
		fmt.Fprintf(os.Stderr, "resetting %d bootstrap "+
			"replicates to zero", *optB)
		*optB = 0
	}
	var ran *rand.Rand
	if *optB > 0 {
		if *optS != 0 {
			ran = rand.New(rand.NewSource(int64(*optS)))
		} else {
			t := time.Now().UnixNano()
			ran = rand.New(rand.NewSource(t))
		}
	}
	files := flag.Args()
	ts := util.NewTransitionTab()
	clio.ParseFiles(files, scan, *optB, *optR, *optU, *optK, ran, ts)
}
