package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/kt"
	"io"
	"log"
	"os"
	"strings"
)

func scan(r io.Reader, args ...interface{}) {
	tree := args[0].(*kt.Node)
	pseq := args[1].([]*fasta.Sequence)
	pstr := args[2].([]string)
	optR := args[3].(bool)
	optI := args[4].(bool)
	scanner := fasta.NewScanner(r)
	for scanner.ScanSequence() {
		seq := scanner.Sequence()
		if optI {
			d := seq.Data()
			d = bytes.ToUpper(d)
			seq = fasta.NewSequence(seq.Header(), d)
		}
		matches := tree.Search(seq.Data(), pstr)
		fmt.Printf("# %s\n", seq.Header())
		printMatches(matches, pseq)
		if optR {
			seq.ReverseComplement()
			matches = tree.Search(seq.Data(), pstr)
			ls := len(seq.Data())
			for i, match := range matches {
				lm := len(pseq[match.Pattern].Data())
				pr := match.Position
				pf := ls - pr - lm
				match.Position = pf
				matches[i] = match
			}
			fmt.Printf("# %s - Reverse\n", seq.Header())
			printMatches(matches, pseq)
		}
	}
}
func printMatches(matches []kt.Match,
	patterns []*fasta.Sequence) {
	for _, m := range matches {
		s := patterns[m.Pattern]
		fmt.Printf("%d\t%s\n", m.Position+1,
			s.Header())
	}
}
func main() {
	util.PrepLog("keyMat")
	u := "keyMat [-h] [options] [patterns] [file(s)]"
	p := "Match one or more patterns in sequence data."
	e := "keyMat -r ATTC,ATTG foo.fasta"
	clio.Usage(u, p, e)
	m := "file with FASTA-formatted patterns"
	var optP = flag.String("p", "", m)
	var optR = flag.Bool("r", false, "include reverse strand")
	var optI = flag.Bool("i", false, "ignore case")
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("keyMat")
	}
	var patterns []*fasta.Sequence
	files := flag.Args()
	if *optP != "" {
		file, err := os.Open(*optP)
		if err != nil {
			log.Fatalf("couldn't open %q\n", *optP)
		}
		scanner := fasta.NewScanner(file)
		for scanner.ScanSequence() {
			patterns = append(patterns, scanner.Sequence())
		}
		file.Close()
	} else if len(files) > 0 {
		p := strings.Split(files[0], ",")
		for _, s := range p {
			seq := fasta.NewSequence(s, []byte(s))
			patterns = append(patterns, seq)
		}
		files = files[1:]
	}
	if len(patterns) == 0 {
		m := "please enter at least one pattern\n"
		fmt.Fprintf(os.Stderr, m)
		os.Exit(-1)
	}
	var sp []string
	for _, s := range patterns {
		seq := string(s.Data())
		if *optI {
			seq = strings.ToUpper(seq)
		}
		sp = append(sp, seq)
	}
	tree := kt.NewKeywordTree(sp)
	clio.ParseFiles(files, scan, tree, patterns, sp, *optR, *optI)
}
