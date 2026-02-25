package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"io"
	"math/rand"
	"time"
)

func parse(r io.Reader, args ...interface{}) {
	n := args[0].(int)
	sc := nwk.NewScanner(r)
	for sc.Scan() {
		tree := sc.Tree()
		labels := []string{}
		labels = extractLabels(tree, labels)
		for i := 0; i < n; i++ {
			rand.Shuffle(len(labels), func(i, j int) {
				labels[i], labels[j] = labels[j], labels[i]
			})
			l := 0
			l = relabelLeaves(tree, labels, l)
			fmt.Println(tree)
		}
	}
}
func extractLabels(v *nwk.Node, labels []string) []string {
	if v != nil {
		if v.Child == nil {
			labels = append(labels, v.Label)
		}
		labels = extractLabels(v.Child, labels)
		labels = extractLabels(v.Sib, labels)
	}
	return labels
}
func relabelLeaves(v *nwk.Node, labels []string, l int) int {
	if v != nil {
		l = relabelLeaves(v.Child, labels, l)
		l = relabelLeaves(v.Sib, labels, l)
		if v.Child == nil {
			v.Label = labels[l]
			l++
		}
	}
	return l
}
func main() {
	util.PrepLog("shuphyl")
	u := "shuphyl [-h] [options] [trees]"
	p := "The phogram shuphyl shuffles the leaf labels of phylogenies"
	e := "shuphyl -n 10 foo.nwk"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optN := flag.Int("n", 1, "number of iterations")
	optS := flag.Int("s", 0, "seed of random number generator (default internal)")
	flag.Parse()
	if *optV {
		util.PrintInfo("shuphyl")
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	files := flag.Args()
	clio.ParseFiles(files, parse, *optN)
}
