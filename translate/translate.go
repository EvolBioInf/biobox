package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
)

func scan(r io.Reader, args ...interface{}) {
	gc := args[0].(map[string]byte)
	frame := args[1].(int)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		if frame < 0 {
			seq.ReverseComplement()
			frame *= -1
		}
		d := seq.Data()
		var aa []byte
		for i := frame - 1; i < len(seq.Data())-2; i += 3 {
			codon := string(d[i : i+3])
			aa = append(aa, gc[codon])
		}
		h := seq.Header() + " - translated"
		aaSeq := fasta.NewSequence(h, aa)
		fmt.Println(aaSeq)
	}
}
func main() {
	util.PrepLog("translate")
	u := "translate [-h] [option]... [foo.fasta]..."
	p := "Translate DNA sequences."
	e := "translate -f 2 foo.fasta"
	clio.Usage(u, p, e)
	var optF = flag.Int("f", 1, "reading frame -3|-2|-1|1|2|3")
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("tranlate")
	}
	if *optF < -3 || *optF > 3 {
		m := "please use a reading frame " +
			"between -3 and 3"
		log.Fatal(m)
	}
	gc := make(map[string]byte)
	dna := "TCAG"
	aa := "FFLLSSSSYY**CC*W" +
		"LLLLPPPPHHQQRRRR" +
		"IIIMTTTTNNKKSSRR" +
		"VVVVAAAADDEEGGGG"
	codon := make([]byte, 3)
	n := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				codon[0] = dna[i]
				codon[1] = dna[j]
				codon[2] = dna[k]
				gc[string(codon)] = aa[n]
				n++
			}
		}
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, gc, *optF)
}
