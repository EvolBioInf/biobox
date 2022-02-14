package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/kt"
	"github.com/evolbioinf/pal"
	"io"
	"log"
	"os"
	"text/tabwriter"
)

type opts struct {
	o, e float64
	k    int
	l    bool
}

func scan(r io.Reader, args ...interface{}) {
	queries := args[0].([]*fasta.Sequence)
	sm := args[1].(*pal.ScoreMatrix)
	op := args[2].(*opts)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		subject := sc.Sequence()
		for _, query := range queries {
			fragments := make([]string, 0)
			starts := make([]int, 0)
			q := query.Data()
			m := len(q)
			r := m / (op.k + 1)
			for i := 0; i <= m-r; i += r {
				f := q[i : i+r]
				fragments = append(fragments, string(f))
				starts = append(starts, i)
			}
			nf := len(fragments)
			fragments[nf-1] = string(q[starts[nf-1]:])
			if op.l {
				w := tabwriter.NewWriter(os.Stdout, 1, 0, 1, ' ', 0)
				fmt.Fprintf(w, "#Id\tStart\tFragment\n")
				for i, f := range fragments {
					fmt.Fprintf(w, "%d\t%d\t%s\n", i+1, starts[i], f)
				}
				w.Flush()
			} else {
				var matches []kt.Match
				ktree := kt.NewKeywordTree(fragments)
				matches = ktree.Search(subject.Data(), fragments)
				var l, r int
				for _, match := range matches {
					if match.Position < r && match.Position > l {
						continue
					}
					i := starts[match.Pattern]
					j := match.Position
					l = j - i - op.k
					if l < 0 {
						l = 0
					}
					r = j + m - i + op.k
					if r > len(subject.Data()) {
						r = len(subject.Data())
					}
					sbjctFrag := subject.Data()[l:r]
					sf := fasta.NewSequence(subject.Header(), sbjctFrag)
					oal := pal.NewOverlapAlignment(query, sf, sm, op.o, op.e)
					oal.Align()
					oal.SetSubjectStart(l)
					oal.TrimQuery()
					e := oal.Mismatches() + oal.Gaps()
					if e <= op.k {
						oal.SetSubjectLength(len(subject.Data()))
						fmt.Printf("%s\n", oal)
					}
				}
			}
		}
	}
}
func main() {
	u := "kerror [-h] [option]... query.fasta [subject.fasta]..."
	p := "Calculate k-error alignments between a short " +
		"query and a long subject."
	e := "kerror -k 3 query.fasta subject.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optK = flag.Int("k", 1, "number of errors")
	var optA = flag.Float64("a", 1, "match")
	var optI = flag.Float64("i", -3, "mismatch")
	var optM = flag.String("m", "", "file containing score matrix")
	var optO = flag.Float64("o", -5, "gap opening")
	var optE = flag.Float64("e", -2, "gap extension")
	var optL = flag.Bool("l", false, "print fragment list")
	flag.Parse()
	if *optV {
		util.PrintInfo("kerror")
	}
	var sm *pal.ScoreMatrix
	if *optM == "" {
		sm = pal.NewScoreMatrix(*optA, *optI)
	} else {
		f, err := os.Open(*optM)
		if err != nil {
			log.Fatalf("can't open %q", *optM)
		}
		sm = pal.ReadScoreMatrix(f)
		f.Close()
	}
	op := new(opts)
	op.k = *optK
	op.o = *optO
	op.e = *optE
	op.l = *optL
	files := flag.Args()
	var queries []*fasta.Sequence
	if len(files) < 1 {
		log.Fatal("please enter query file")
	} else {
		f, err := os.Open(files[0])
		if err != nil {
			log.Fatalf("can't open %q", files[0])
		}
		sc := fasta.NewScanner(f)
		for sc.ScanSequence() {
			q := sc.Sequence()
			queries = append(queries, q)
		}
		f.Close()
	}
	clio.ParseFiles(files[1:], scan, queries, sm, op)
}
