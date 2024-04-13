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
	"sort"
	"strconv"
)

type sortableContigs []*fasta.Sequence

func (s sortableContigs) Len() int {
	return len(s)
}
func (s sortableContigs) Less(i, j int) bool {
	return len(s[i].Data()) < len(s[j].Data())
}
func (s sortableContigs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func scan(r io.Reader, args ...interface{}) {
	contigs := args[0].(*([]*fasta.Sequence))
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		s := sc.Sequence()
		(*contigs) = append(*contigs, s)
	}
}
func bestAl(contigs []*fasta.Sequence, sm *pal.ScoreMatrix,
	rev bool, optO, optE float64) (i, j int,
	oal *pal.OverlapAlignment) {
	var mi, mj int
	var mo *pal.OverlapAlignment
	ms := -1.0
	for i := 0; i < len(contigs); i++ {
		for j := i + 1; j < len(contigs); j++ {
			oal := pal.NewOverlapAlignment(contigs[i], contigs[j],
				sm, optO, optE)
			oal.Align()
			if oal.Score() > ms {
				mo = oal
				mi = i
				mj = j
				ms = oal.Score()
			}
			if rev {
				if len(contigs[i].Data()) < len(contigs[j].Data()) {
					contigs[i].ReverseComplement()
				} else {
					contigs[j].ReverseComplement()
				}
				oal = pal.NewOverlapAlignment(contigs[i], contigs[j],
					sm, optO, optE)
				oal.Align()
				if oal.Score() > ms {
					mo = oal
					mi = i
					mj = j
					ms = oal.Score()
				}
			}
		}
	}
	return mi, mj, mo
}
func clean(b []byte) []byte {
	n := 0
	for i := 0; i < len(b); i++ {
		if b[i] != '-' {
			b[n] = b[i]
			n++
		}
	}
	b = b[:n]
	return b
}
func main() {
	util.PrepLog("sass")
	u := "sass [option]... [file]..."
	p := "Calculate assembly using a simple algorithm."
	e := "sass -r reads.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optA = flag.Float64("a", 1, "match")
	var optI = flag.Float64("i", -3, "mismatch")
	var optM = flag.String("m", "", "file containing score matrix")
	var optO = flag.Float64("o", -5, "gap opening")
	var optE = flag.Float64("e", -2, "gap extension")
	var optR = flag.Bool("r", false, "include reverse strand")
	var optMM = flag.Bool("M", false, "print merge steps")
	var optT = flag.Float64("t", 15.0, "score threshold")
	flag.Parse()
	if *optV {
		util.PrintInfo("sass")
	}
	var sm *pal.ScoreMatrix
	if *optM == "" {
		sm = pal.NewByteScoreMatrix(*optA, *optI)
	} else {
		f, err := os.Open(*optM)
		if err != nil {
			log.Fatalf("couldn't open score matrix %q\n",
				(*optM))
		}
		sm = pal.ReadScoreMatrix(f)
		f.Close()
	}
	files := flag.Args()
	contigs := make([]*fasta.Sequence, 0)
	clio.ParseFiles(files, scan, &contigs)
	i, j, bal := bestAl(contigs, sm, *optR, *optO, *optE)
	for len(contigs) > 1 && bal.Score() >= *optT {
		a1, a2 := bal.RawAlignment()
		var m []byte
		for i, c := range a1 {
			if c != '-' {
				m = append(m, c)
			} else {
				m = append(m, a2[i])
			}
		}
		if *optMM {
			s1 := string(clean(a1))
			s2 := string(clean(a2))
			s3 := string(m)
			fmt.Println(s1, s2, s3)
		}
		contig := fasta.NewSequence("", m)
		contigs = append(contigs, contig)
		n := 0
		for k := 0; k < len(contigs); k++ {
			if k != i && k != j {
				contigs[n] = contigs[k]
				n++
			}
		}
		contigs = contigs[:n]
		i, j, bal = bestAl(contigs, sm, *optR, *optO, *optE)
	}
	sc := sortableContigs(contigs)
	sort.Sort(sc)
	nc := 0
	for i := len(sc) - 1; i >= 0; i-- {
		if len(sc[i].Header()) == 0 {
			sc[i].AppendToHeader("Contig_" +
				strconv.Itoa(nc+1))
			nc++
		}
		fmt.Println(sc[i])
	}
}
