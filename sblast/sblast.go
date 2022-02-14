package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/kt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type Opts struct {
	a, i, t float64
	w, s    int
	n, l    bool
}
type Alignment struct {
	qs, qe, ss, se int
	score          float64
	forward        bool
}
type AlSliceStart []Alignment
type AlSliceEnd []Alignment
type AlSliceScore []Alignment

func (a AlSliceStart) Len() int {
	return len(a)
}
func (a AlSliceStart) Less(i, j int) bool {
	if a[i].ss == a[j].ss {
		return a[i].score > a[j].score
	} else {
		return a[i].ss < a[j].ss
	}
}
func (a AlSliceStart) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AlSliceEnd) Len() int {
	return len(a)
}
func (a AlSliceEnd) Less(i, j int) bool {
	if a[i].se == a[j].se {
		return a[i].score > a[j].score
	} else {
		return a[i].se < a[j].se
	}
}
func (a AlSliceEnd) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AlSliceScore) Len() int {
	return len(a)
}
func (a AlSliceScore) Less(i, j int) bool {
	return a[i].score > a[j].score
}
func (a AlSliceScore) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func scan(r io.Reader, args ...interface{}) {
	opts := args[0].(*Opts)
	qName := args[1].(string)
	out := args[2].(*tabwriter.Writer)
	sScanner := fasta.NewScanner(r)
	for sScanner.ScanSequence() {
		subject := sScanner.Sequence()
		qFile, err := os.Open(qName)
		if err != nil {
			log.Fatalf("couldn't open %s\n", qName)
		}
		defer qFile.Close()
		qScanner := fasta.NewScanner(qFile)
		for qScanner.ScanSequence() {
			query := qScanner.Sequence()
			if opts.l {
				words := getWords(query, opts.w)
				qa := strings.Fields(query.Header())[0]
				for i, word := range words {
					fmt.Fprintf(out, "%s\t%d\t%s\n", qa, i+1, word)
				}
			} else {
				forward := true
				alignments := align(query, subject, opts, forward)
				query.ReverseComplement()
				forward = false
				a := align(query, subject, opts, forward)
				alignments = append(alignments, a...)
				qa := strings.Fields(query.Header())[0]
				sa := strings.Fields(subject.Header())[0]
				for _, a := range alignments {
					if !a.forward {
						a.ss, a.se = a.se, a.ss
					}
					fmt.Fprintf(out, "%s\t%s\t%d\t%d\t%d\t%d\t%.1f\n",
						qa, sa, a.qs+1, a.qe+1, a.ss+1, a.se+1, a.score)
				}
			}
		}
	}
}
func getWords(seq *fasta.Sequence, w int) []string {
	var words []string
	d := seq.Data()
	l := len(d)
	for i := 0; i <= l-w; i++ {
		word := string(d[i : i+w])
		words = append(words, word)
	}
	return words
}
func align(query, subject *fasta.Sequence,
	opts *Opts, forward bool) []Alignment {
	var alignments []Alignment
	if opts.n {
		q := query.Data()
		m := len(q)
		s := subject.Data()
		n := len(s)
		w := opts.w
		for i := 0; i < m-w; i++ {
			p := q[i : i+w]
			for j := 0; j < n-w; j++ {
				var k int
				for k = 0; k < w; k++ {
					if s[j+k] != p[k] {
						break
					}
				}
				if k == opts.w {
					a := Alignment{qs: i, qe: i + w - 1, ss: j, se: j + w - 1,
						score: float64(w) * opts.a, forward: forward}
					alignments = append(alignments, a)
				}
			}
		}
	} else {
		var patterns []string
		q := query.Data()
		m := len(q)
		w := opts.w
		for i := 0; i <= m-w; i++ {
			p := string(q[i : i+w])
			patterns = append(patterns, p)
		}
		tree := kt.NewKeywordTree(patterns)
		matches := tree.Search(subject.Data(), patterns)
		for _, m := range matches {
			qs := m.Pattern
			ss := m.Position
			qe := qs + w - 1
			se := ss + w - 1
			sc := float64(w) * opts.a
			a := Alignment{qs: qs, ss: ss, qe: qe,
				se: se, score: sc, forward: forward}
			alignments = append(alignments, a)
		}
	}
	q := query.Data()
	m := len(q)
	s := subject.Data()
	n := len(s)
	for i, _ := range alignments {
		cq := alignments[i].qs - 1
		cs := alignments[i].ss - 1
		score := alignments[i].score
		is := 0
		for cq >= 0 && cs >= 0 && is <= opts.s {
			if q[cq] == s[cs] {
				score += opts.a
			} else {
				score += opts.i
			}
			if score > alignments[i].score {
				alignments[i].score = score
				alignments[i].qs = cq
				alignments[i].ss = cs
				is = 0
			} else {
				is++
			}
			cq--
			cs--
		}
		cq = alignments[i].qe + 1
		cs = alignments[i].se + 1
		score = alignments[i].score
		is = 0
		for cq < m && cs < n && is <= opts.s {
			if q[cq] == s[cs] {
				score += opts.a
			} else {
				score += opts.i
			}
			if score > alignments[i].score {
				alignments[i].score = score
				alignments[i].qe = cq
				alignments[i].se = cs
				is = 0
			} else {
				is++
			}
			cq++
			cs++
		}
	}
	i := 0
	max := -1.0
	for _, al := range alignments {
		if al.score >= opts.t {
			if max < al.score {
				max = al.score
			}
			alignments[i] = al
			i++
		}
	}
	alignments = alignments[:i]
	sort.Sort(AlSliceStart(alignments))
	j := 0
	if len(alignments) > 0 {
		j = 1
	}
	for i := 1; i < len(alignments); i++ {
		if alignments[i].ss != alignments[i-1].ss {
			alignments[j] = alignments[i]
			j++
		}
	}
	alignments = alignments[:j]
	sort.Sort(AlSliceEnd(alignments))
	j = 0
	if len(alignments) > 0 {
		j = 1
	}
	for i := 1; i < len(alignments); i++ {
		if alignments[i].se != alignments[i-1].se {
			alignments[j] = alignments[i]
			j++
		}
	}
	alignments = alignments[:j]
	sort.Sort(AlSliceScore(alignments))
	return alignments
}
func main() {
	u := "sblast [-h] [option]... query.fasta [subject.fasta]..."
	p := "Carry out a simple version of BLAST."
	e := "sblast query.fasta subject.fasta"
	clio.Usage(u, p, e)
	var optA = flag.Float64("a", 1.0, "match")
	var optI = flag.Float64("i", -3.0, "mismatch")
	var optW = flag.Int("w", 11, "word length")
	var optS = flag.Int("s", 30, "maximum number "+
		"of idle extension steps")
	var optT = flag.Float64("t", 50.0, "threshold score")
	var optN = flag.Bool("n", false, "naive matching")
	var optL = flag.Bool("l", false, "print word list")
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("sblast")
	}
	opts := new(Opts)
	opts.a = *optA
	opts.i = *optI
	opts.w = *optW
	opts.s = *optS
	opts.t = *optT
	opts.n = *optN
	opts.l = *optL
	files := flag.Args()
	if len(files) == 0 {
		log.Fatal("please provide a query")
	}
	out := tabwriter.NewWriter(os.Stdout, 2, 1, 2, ' ', 0)
	if !opts.l {
		fmt.Fprintf(out, "#qa\tsa\tqs\tqe\tsa\tse\tscore\n")
	} else {
		fmt.Fprintf(out, "#qa\tn\tword\n")
	}
	clio.ParseFiles(files[1:], scan, opts, files[0], out)
	out.Flush()
}
