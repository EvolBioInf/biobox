package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/esa"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"math"
	"regexp"
	"text/tabwriter"
)

func scan(r io.Reader, args ...interface{}) {
	local := args[0].(bool)
	reverse := args[1].(bool)
	quiet := args[2].(bool)
	seqReg := args[3].(*regexp.Regexp)
	scanner := fasta.NewScanner(r)
	var sequences []*fasta.Sequence
	for scanner.ScanSequence() {
		sequence := scanner.Sequence()
		sequences = append(sequences, sequence)
	}
	var cat []byte
	var start, end []int
	start = append(start, 0)
	for i, sequence := range sequences {
		if i > 0 {
			cat = append(cat, 0)
			start = append(start, end[i-1]+1)
		}
		cat = append(cat, sequence.Data()...)
		end = append(end, start[i]+len(sequence.Data()))
	}
	if reverse {
		for _, sequence := range sequences {
			sequence.ReverseComplement()
			cat = append(cat, 0)
			cat = append(cat, sequence.Data()...)
		}
	}
	sa := esa.Sa(cat)
	lcp := esa.Lcp(cat, sa)
	isa := make([]int, len(sa))
	for i, _ := range sa {
		isa[sa[i]] = i
	}
	shu := make([]int, len(sa))
	lcp = append(lcp, -1)
	for i, _ := range sequences {
		for j := start[i]; j < end[i]; j++ {
			is := isa[j]
			shu[is] = lcp[is]
			if lcp[is+1] > shu[is] {
				shu[is] = lcp[is+1]
			}
			shu[is]++
			if sa[is]+shu[is] > end[i] {
				shu[is] = math.MaxInt64
			}
		}
	}
	var maxima []int
	for i, _ := range sequences {
		maxima = append(maxima, math.MaxInt64-1)
		if local {
			continue
		}
		for j := start[i]; j < end[i]; j++ {
			l := shu[isa[j]]
			if l < maxima[i] {
				maxima[i] = l
			}
		}
	}
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 2, ' ', 0)
	for i, sequence := range sequences {
		header := []byte(sequence.Header())
		match := seqReg.Find(header)
		if match == nil {
			continue
		}
		fmt.Printf(">%s\n", sequence.Header())
		buffer.Reset()
		fmt.Fprintf(w, "#\t")
		if !local {
			fmt.Fprint(w, "Count\t")
		}
		fmt.Fprint(w, "Position\tLength")
		if !quiet {
			fmt.Fprintf(w, "\tShustring")
		}
		fmt.Fprint(w, "\n")
		count := 0
		for j := start[i]; j < end[i]; j++ {
			is := isa[j]
			if shu[is] <= maxima[i] {
				count++
				s := sa[is] - start[i]
				l := shu[is]
				if !local {
					fmt.Fprintf(w, "\t%d", count)
				}
				fmt.Fprintf(w, "\t%d\t%d", s+1, l)
				if !quiet {
					s = sa[is]
					str := string(cat[s : s+l])
					fmt.Fprintf(w, "\t%s", str)
				}
				fmt.Fprintf(w, "\n")
			}
		}
		w.Flush()
		fmt.Printf("%s", buffer)
	}
}
func main() {
	util.PrepLog("shustring")
	u := "shustring [-h] [options] [files]"
	p := "Compute shortest unique substrings."
	e := "shustring foo.fasta"
	clio.Usage(u, p, e)
	var optL = flag.Bool("l", false, "local")
	var optS = flag.String("s", ".", "restrict output to sequences "+
		"described by regex")
	var optR = flag.Bool("r", false, "include reverse strand")
	var optQ = flag.Bool("q", false, "quiet, don't print shustrings; saves memory")
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	seqReg, err := regexp.Compile(*optS)
	if err != nil {
		log.Fatalf("couldn't compile %q.\n", *optS)
	}
	if *optV {
		util.PrintInfo("shustring")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optL, *optR, *optQ, seqReg)
}
