package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/nwk"
	"io"
	"math"
	"sort"
)

type leafSlice []*nwk.Node

func (l leafSlice) Len() int { return len(l) }
func (l leafSlice) Less(i, j int) bool {
	return l[i].Length < l[j].Length
}
func (l leafSlice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func scan(r io.Reader, args ...interface{}) {
	bits := args[0].(bool)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		counts := make([]int, 256)
		for _, c := range seq.Data() {
			counts[c]++
		}
		n := 0
		for _, count := range counts {
			if count > 0 {
				n++
			}
		}
		if n == 1 {
			if counts['A'] == 0 {
				counts['A'] = 1
			} else {
				counts['C'] = 1
			}
		}
		leaves := make([]*nwk.Node, 0)
		labels := make(map[int]byte)
		for i, c := range counts {
			if c > 0 {
				n := nwk.NewNode()
				n.Length = float64(c) / float64(len(seq.Data()))
				n.HasLength = true
				labels[n.Id] = byte(i)
				leaves = append(leaves, n)
			}
		}
		for len(leaves) > 1 {
			sort.Sort(leafSlice(leaves))
			leaves[0].Label = "0"
			leaves[1].Label = "1"
			n := nwk.NewNode()
			n.AddChild(leaves[0])
			n.AddChild(leaves[1])
			n.Length = leaves[0].Length + leaves[1].Length
			n.HasLength = true
			leaves = append(leaves, n)
			leaves = leaves[2:]
		}
		root := leaves[0]
		traverse(root, labels)
		if bits {
			sl := len(seq.Data())
			lw := 0.0
			lw = sumLeafWeights(root, lw)
			nb := float64(sl) * lw
			nb = math.Ceil(nb)
			fmt.Printf(">%s\n", seq.Header())
			fmt.Printf("Bits: %d\n", int(nb))
		} else {
			fmt.Println(root)
		}
	}
}
func traverse(n *nwk.Node, labels map[int]byte) {
	if n == nil {
		return
	}
	if n.Child == nil {
		code := n.Label
		v := n.Parent
		for v != nil {
			code += v.Label
			v = v.Parent
		}
		b := []byte(code)
		for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
			b[i], b[j] = b[j], b[i]
		}
		code = string(b)
		n.Label = fmt.Sprintf("\"%s-%c/%s\"", n.Label, labels[n.Id], code)
	}
	traverse(n.Child, labels)
	traverse(n.Sib, labels)
}
func sumLeafWeights(v *nwk.Node, w float64) float64 {
	if v == nil {
		return w
	}
	x := 0.0
	if v.Child == nil {
		cl := float64(len(v.Label) - 6)
		w += v.Length * cl
	}
	w = sumLeafWeights(v.Child, w+x)
	w = sumLeafWeights(v.Sib, w+x)
	return w
}
func main() {
	util.PrepLog("hut")
	m := "hut [-h] [option]... [file]..."
	p := "Convert sequences into their Huffman trees."
	e := "hut foo.fasta"
	clio.Usage(m, p, e)
	var optV = flag.Bool("v", false, "version")
	var optB = flag.Bool("b", false, "bits")
	flag.Parse()
	if *optV {
		util.PrintInfo("hut")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optB)
}
