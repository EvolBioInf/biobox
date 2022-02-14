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
	"strings"
)

func scan(r io.Reader, args ...interface{}) {
	tree := args[0].(*kt.Node)
	pseq := args[1].([]*fasta.Sequence)
	pstr := args[2].([]string)
	optR := args[3].(bool)
	scanner := fasta.NewScanner(r)
	for scanner.ScanSequence() {
		seq := scanner.Sequence()
		matches := tree.Search(seq.Data(), pstr)
		fmt.Printf("# %s\n", seq.Header())
		printMatches(matches, pseq)
		if optR {
			seq.ReverseComplement()
			matches = tree.Search(seq.Data(), pstr)
			fmt.Printf("# %s - Reverse\n", seq.Header())
			printMatches(matches, pseq)
		}
	}
}
func printMatches(matches []kt.Match, patterns []*fasta.Sequence) {
	for _, m := range matches {
		s := patterns[m.Pattern]
		fmt.Printf("%d\t%s\n", m.Position+1, s.Header())
	}
}
func main() {
	u := "keyMat [-h] [options] [file(s)]"
	p := "Set matching in sequence data"
	e := "keyMat -p ATTC,ATTG foo.fasta"
	clio.Usage(u, p, e)
	var optP = flag.String("p", "", "comma-separated patterns")
	var optF = flag.String("f", "", "file with FASTA-formatted patterns")
	var optR = flag.Bool("r", false, "include reverse strand")
	var optV = flag.Bool("v", false, "print version & "+
		"other program information")
	flag.Parse()
	if *optV {
		util.PrintInfo("keyMat")
	}
	var patterns []*fasta.Sequence
	if *optP != "" {
		p := strings.Split(*optP, ",")
		for _, s := range p {
			seq := fasta.NewSequence(s, []byte(s))
			patterns = append(patterns, seq)
		}
	}
	if *optF != "" {
		file, err := os.Open(*optF)
		if err != nil {
			log.Fatalf("couldn't open %q\n", *optF)
		}
		scanner := fasta.NewScanner(file)
		for scanner.ScanSequence() {
			patterns = append(patterns, scanner.Sequence())
		}
		file.Close()
	}
	if len(patterns) == 0 {
		fmt.Fprintf(os.Stderr, "please enter at least one pattern\n")
		os.Exit(-1)
	}
	var sp []string
	for _, s := range patterns {
		sp = append(sp, string(s.Data()))
	}
	tree := kt.NewKeywordTree(sp)
	files := flag.Args()
	clio.ParseFiles(files, scan, tree, patterns, sp, *optR)
}
