package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/dist"
	"github.com/evolbioinf/nwk"
	"io"
)

func scan(r io.Reader, args ...interface{}) {
	printMat := args[0].(bool)
	sc := dist.NewScanner(r)
	for sc.Scan() {
		dm := sc.DistanceMatrix()
		dm.MakeSymmetrical()
		n := len(dm.Names)
		var root *nwk.Node
		t := make([]*nwk.Node, n)
		for i := 0; i < n; i++ {
			t[i] = nwk.NewNode()
			t[i].Label = dm.Names[i]
		}
		for i := n; i > 1; i-- {
			if printMat {
				fmt.Printf("%s", dm)
			}
			md, mj, mk := dm.Min()
			c1 := t[mj]
			c2 := t[mk]
			root = nwk.NewNode()
			l := fmt.Sprintf("(%s,%s)", c1.Label, c2.Label)
			root.Label = l
			root.Length = md / 2
			root.AddChild(c1)
			root.AddChild(c2)
			data := make([]float64, i-2)
			k := 0
			for j := 0; j < i; j++ {
				if j == mj || j == mk {
					continue
				}
				data[k] = (dm.Matrix[j][mj] + dm.Matrix[j][mk]) / 2.0
				k++
			}
			dm.DeletePair(mj, mk)
			dm.Append(root.Label, data)
			j := 0
			for k := 0; k < i; k++ {
				if k == mj || k == mk {
					continue
				}
				t[j] = t[k]
				j++
			}
			t = t[:j]
			t = append(t, root)
		}
		branchLengths(root)
		fmt.Println(root)
	}
}
func branchLengths(v *nwk.Node) {
	if v == nil {
		return
	}
	branchLengths(v.Child)
	branchLengths(v.Sib)
	if v.Child != nil {
		v.Label = ""
	}
	if v.Parent != nil {
		v.Length = v.Parent.Length - v.Length
		v.HasLength = true
	}
}
func main() {
	util.PrepLog("upgma")
	u := "upgma [-h] [option]... [foo.dist]..."
	p := "Cluster a distance matrix into a tree using UPGMA."
	e := "upgma foo.dist"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optM = flag.Bool("m", false, "print intermediate "+
		"matrices")
	flag.Parse()
	if *optV {
		util.PrintInfo("upgma")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optM)
}
