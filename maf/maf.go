package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/esa"
	"github.com/evolbioinf/fasta"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

func scan(r io.Reader, args ...interface{}) {
	printNum := args[0].(bool)
	w := args[1].(*tabwriter.Writer)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		t := seq.Data()
		sa := esa.Sa(t)
		isa := make([]int, len(sa))
		lcp := esa.Lcp(t, sa)
		lcp = append(lcp, 0)
		for i, s := range sa {
			isa[s] = i
		}
		factors := make([][]byte, 0)
		i := 0
		for i < len(sa) {
			l1 := lcp[isa[i]]
			l2 := lcp[isa[i]+1]
			j := i + max(max(l1, l2), 1)
			factors = append(factors, t[i:j])
			i = j
		}
		if printNum {
			n := len(factors)
			m := len(t)
			a := strings.Fields(seq.Header())[0]
			fmt.Fprintf(w, "%s\t%d\t%d\t%.3g\n", a, n,
				m, float64(n)/float64(m))
		} else {
			var fs *fasta.Sequence
			fd := make([]byte, 0)
			n := len(factors)
			for i := 0; i < n-1; i++ {
				fd = append(fd, factors[i]...)
				fd = append(fd, '.')
			}
			fd = append(fd, factors[n-1]...)
			h := seq.Header() + " - match factors"
			fs = fasta.NewSequence(h, fd)
			fmt.Println(fs)
		}
	}
	if printNum {
		w.Flush()
	}
}
func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
func main() {
	util.PrepLog("maf")
	u := "maf [-h] [option]... [foo.fasta]..."
	p := "Compute the match factors of a sequence."
	e := "maf foo.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optN = flag.Bool("n", false, "print number of factors "+
		"instead of factors")
	flag.Parse()
	if *optV {
		util.PrintInfo("maf")
	}
	var w *tabwriter.Writer
	if *optN {
		w = tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', 0)
		fmt.Fprintf(w, "#acc\tfactors\tresidues\t"+
			"factors/residues\n")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optN, w)
}
