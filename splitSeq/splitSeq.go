package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"strconv"
)

func parse(r io.Reader, args ...interface{}) {
	optL := args[0].(*int)
	optS := args[1].(*int)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		header := seq.Header()
		data := seq.Data()
		n := len(seq.Data())
		m := n - *optL
		s := 0
		for ; s < m; s += *optS {
			e := s + *optL
			h := header + "|"
			h += strconv.Itoa(s+1) +
				".." +
				strconv.Itoa(e)
			d := data[s:e]
			ns := fasta.NewSequence(h, d)
			fmt.Println(ns)
		}
		e := n
		if s < e {
			h := header + "|"
			h += strconv.Itoa(s+1) +
				".." +
				strconv.Itoa(e)
			d := data[s:e]
			ns := fasta.NewSequence(h, d)
			fmt.Println(ns)
		}
	}
}
func main() {
	util.PrepLog("splitSeq")
	u := "splitSeq [-h] [options] [files]"
	p := "Split sequences into fragments."
	e := "splitSeq -l 300 -s 30 eco.fasta"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optL := flag.Int("l", 1000, "fragment length")
	optS := flag.Int("s", 0, "step length (default -l)")
	flag.Parse()
	if *optV {
		util.PrintInfo("splitSeq")
	}
	if *optL < 1 {
		log.Fatal("please use fragment length > 0")
	}
	if *optS == 0 {
		(*optS) = *optL
	} else if *optS < 1 {
		log.Fatal("please use step length > 0")
	}
	files := flag.Args()
	clio.ParseFiles(files, parse, optL, optS)
}
