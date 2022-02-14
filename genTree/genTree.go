package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func setBranchLen(v *nwk.Node) {
	if v == nil {
		return
	}
	setBranchLen(v.Child)
	l := 0.0
	if v.Parent != nil {
		l = v.Parent.Length - v.Length
	}
	v.Length = l
	setBranchLen(v.Sib)
}
func addMut(v *nwk.Node, t float64, r *rand.Rand) {
	if v == nil {
		return
	}
	lambda := t * v.Length / 2.0
	x := math.Exp(-lambda)
	p := 1.0
	c := 0.0
	for p > x {
		p *= r.Float64()
		c++
	}
	v.Length = c
	addMut(v.Child, t, r)
	addMut(v.Sib, t, r)
}
func labelLeaves(v *nwk.Node, nc int) int {
	if v == nil {
		return nc
	}
	nc = labelLeaves(v.Child, nc)
	if v.Child == nil {
		nc++
		v.Label = "T" + strconv.Itoa(nc)
	}
	nc = labelLeaves(v.Sib, nc)
	return nc
}
func labelInternalNodes(v *nwk.Node, nc int) int {
	if v == nil {
		return nc
	}
	if v.Child != nil {
		nc++
		v.Label = "N" + strconv.Itoa(nc)
	}
	nc = labelInternalNodes(v.Child, nc)
	nc = labelInternalNodes(v.Sib, nc)
	return nc
}
func main() {
	u := "genTree [-h] [option]..."
	p := "Generate random trees."
	e := "genTree -n 15"
	clio.Usage(u, p, e)
	var optN = flag.Int("n", 10, "sample size")
	var optI = flag.Int("i", 1, "iterations")
	var optT = flag.Float64("t", 1000, "theta=2Nu")
	var optC = flag.Bool("c", false, "coalescent")
	var optA = flag.Bool("a", false, "absolute branch lengths")
	var optL = flag.Bool("l", false, "label internal branches")
	var optS = flag.Int("s", 0, "seed for random number generator")
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("genTree")
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	ran := rand.New(rand.NewSource(seed))
	n := *optN
	tree := make([]*nwk.Node, 2*n-1)
	for ii := 0; ii < *optI; ii++ {
		for i := 0; i < 2*n-1; i++ {
			tree[i] = nwk.NewNode()
		}
		t := 0.0
		for i := 0; i < n; i++ {
			tree[i].HasLength = true
		}
		for i := n; i > 1; i-- {
			lambda := float64(n * (n - 1) / 2)
			if *optC {
				lambda = float64(i * (i - 1) / 2)
			}
			t += rand.ExpFloat64() / lambda
			j := 2*n - i
			tree[j].Length = t
			tree[j].HasLength = true
		}
		for i := n; i > 1; i-- {
			p := tree[2*n-i]
			r := ran.Intn(i)
			c := tree[r]
			p.AddChild(c)
			tree[r] = tree[i-1]
			r = ran.Intn(i - 1)
			c = tree[r]
			p.AddChild(c)
			tree[r] = p
		}
		root := tree[len(tree)-1]
		setBranchLen(root)
		if !*optA {
			addMut(root, *optT, ran)
		}
		nc := 0
		nc = labelLeaves(root, nc)
		if *optL {
			nc = 0
			labelInternalNodes(root, nc)
		}
		fmt.Println(root)
	}
}
