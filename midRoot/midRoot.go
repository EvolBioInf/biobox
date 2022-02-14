package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"io"
	"math"
)

func scan(r io.Reader, args ...interface{}) {
	printPair := args[0].(bool)
	sc := nwk.NewScanner(r)
	for sc.Scan() {
		root := sc.Tree()
		var leaves []*nwk.Node
		leaves = collectLeaves(root, leaves)
		n := len(leaves)
		max := -math.MaxFloat64
		var mi, mj int
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				l1 := leaves[i]
				l2 := leaves[j]
				a := l1.LCA(l2)
				d := l1.UpDistance(a) + l2.UpDistance(a)
				if max < d {
					max = d
					mi = i
					mj = j
				}
			}
		}
		if printPair {
			fmt.Printf("# d(%s, %s): %.3g\n",
				leaves[mi].Label, leaves[mj].Label, max)
		}
		l1 := leaves[mi]
		l2 := leaves[mj]
		a := l1.LCA(l2)
		v := l1
		if l1.UpDistance(a) < l2.UpDistance(a) {
			v = l2
		}
		s := v.Length
		for s < max/2.0 {
			v = v.Parent
			s += v.Length
		}
		r := nwk.NewNode()
		p := v.Parent
		p.AddChild(r)
		p.RemoveChild(v)
		r.AddChild(v)
		x2 := s - max/2.0
		x1 := v.Length - x2
		v.Length = x1
		r.Length = x2
		parentToChild(r)
		root = r
		fmt.Println(root)
	}
}
func collectLeaves(v *nwk.Node, l []*nwk.Node) []*nwk.Node {
	if v == nil {
		return l
	}
	l = collectLeaves(v.Child, l)
	l = collectLeaves(v.Sib, l)
	if v.Child == nil {
		l = append(l, v)
	}
	return l
}
func parentToChild(v *nwk.Node) {
	if v.Parent.Parent != nil {
		parentToChild(v.Parent)
	}
	p := v.Parent
	p.RemoveChild(v)
	v.AddChild(p)
	p.Length = v.Length
	p.HasLength = true
}
func main() {
	u := "midRoot [-h] [option]... [foo.nwk]..."
	p := "Add midpoint root to a tree."
	e := "midRoot foo.nwk"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optP = flag.Bool("p", false, "print most distant pair")
	flag.Parse()
	if *optV {
		util.PrintInfo("midRoot")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optP)
}
