package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
)

func scan(r io.Reader, args ...interface{}) {
	reads := args[0].(*([]*fasta.Sequence))
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		s := sc.Sequence()
		(*reads) = append(*reads, s)
	}
}
func overlap(a, b []byte, k int) int {
	p := 0
	if k > len(b) {
		log.Fatal("can't have longer overlaps than reads")
	}
	for true {
		s := bytes.Index(a[p:], b[0:k])
		if s == -1 {
			return 0
		}
		if bytes.HasPrefix(b, a[p+s:]) {
			return len(a) - s - p
		}
		p += s + 1
	}
	return p
}
func main() {
	util.PrepLog("olga")
	u := "olga [-v|-h] [file]..."
	p := "Calculate overlap graph from input strings."
	e := "olga foo.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optK = flag.Int("k", 1, "overlap")
	flag.Parse()
	if *optV {
		util.PrintInfo("olga")
	}
	if *optK < 0 {
		log.Fatal("please enter positive minimum overlap")
	}
	files := flag.Args()
	reads := make([]*fasta.Sequence, 0)
	clio.ParseFiles(files, scan, &reads)
	n := len(reads)
	graph := make([][]int, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			a := reads[i].Data()
			b := reads[j].Data()
			graph[i][j] = overlap(a, b, *optK)
		}
	}
	fmt.Println("# Overlap graph generated with olga.")
	fmt.Println("# Render: dot foo.dot")
	fmt.Println("digraph G {")
	for i := 0; i < n; i++ {
		fmt.Printf("\tn%d [label=\"%s\"]\n", i,
			string(reads[i].Data()))
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if graph[i][j] == 0 {
				continue
			}
			fmt.Printf("\tn%d -> n%d [label=\" %d\"]\n",
				i, j, graph[i][j])
		}
	}
	fmt.Println("}")
}
