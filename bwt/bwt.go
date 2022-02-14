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
)

func scan(r io.Reader, args ...interface{}) {
	decode := args[0].(bool)
	var in, out *fasta.Sequence
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		in = sc.Sequence()
		if decode {
			var count, first [256]int
			var prior []int
			transform := in.Data()
			prior = make([]int, len(transform))
			for i, t := range transform {
				prior[i] = count[t]
				count[t]++
			}
			s := 0
			for i, c := range count {
				first[i] = s
				s += c
			}
			o := make([]byte, len(transform))
			if c := bytes.Count(transform, []byte("$")); c != 1 {
				m := "sentinel, $, appears %d times rather than once"
				log.Fatalf(m, c)
			}
			i := bytes.IndexByte(transform, '$')
			for j := len(transform) - 1; j > -1; j-- {
				o[j] = transform[i]
				i = prior[i] + first[transform[i]]
			}
			h := in.Header() + " - decoded"
			out = fasta.NewSequence(h, o[:len(o)-1])
		} else {
			data := in.Data()
			sent := byte('$')
			data = append(data, sent)
			sa := esa.Sa(data)
			o := make([]byte, len(data))
			for i, s := range sa {
				o[i] = sent
				if s > 0 {
					o[i] = data[s-1]
				}
			}
			h := in.Header() + " - bwt"
			out = fasta.NewSequence(h, o)
		}
		fmt.Printf("%s\n", out)
	}
}
func main() {
	u := "bwt [-h] [option]... [foo.fasta]..."
	p := "Compute the Burrows-Wheeler transform."
	e := "bwt foo.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optD = flag.Bool("d", false, "decode")
	flag.Parse()
	if *optV {
		util.PrintInfo("bwt")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optD)
}
