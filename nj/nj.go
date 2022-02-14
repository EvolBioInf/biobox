package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/dist"
	"github.com/evolbioinf/nwk"
	"io"
	"os"
	"text/tabwriter"
)

func scan(r io.Reader, args ...interface{}) {
	printMat := args[0].(bool)
	sc := dist.NewScanner(r)
	for sc.Scan() {
		dm := sc.DistanceMatrix()
		dm.MakeSymmetrical()
		r := rowSums(dm)
		sm := smat(dm, r)
		var root *nwk.Node
		n := len(dm.Names)
		t := make([]*nwk.Node, n)
		for i := 0; i < n; i++ {
			t[i] = nwk.NewNode()
			t[i].Label = dm.Names[i]
		}
		for i := n; i > 3; i-- {
			if printMat {
				printMatrices(dm, sm, r)
			}
			_, mj, mk := sm.Min()
			c1 := t[mj]
			c2 := t[mk]
			root = nwk.NewNode()
			l := fmt.Sprintf("(%s,%s)", c1.Label, c2.Label)
			root.Label = l
			x := float64(i-2) * dm.Matrix[mj][mk]
			denom := float64(2 * (i - 2))
			c1.Length = (x + r[mj] - r[mk]) / denom
			c2.Length = (x + r[mk] - r[mj]) / denom
			c1.HasLength = true
			c2.HasLength = true
			root.AddChild(c1)
			root.AddChild(c2)
			data := make([]float64, i-2)
			k := 0
			for j := 0; j < i; j++ {
				if j == mj || j == mk {
					continue
				}
				data[k] = (dm.Matrix[j][mj] +
					dm.Matrix[j][mk] - dm.Matrix[mj][mk]) / 2.0
				k++
			}
			dm.DeletePair(mj, mk)
			dm.Append(root.Label, data)
			k = 0
			for j := 0; j < i; j++ {
				if j == mj || j == mk {
					continue
				}
				t[k] = t[j]
				k++
			}
			t = t[:k]
			t = append(t, root)
			r = rowSums(dm)
			sm = smat(dm, r)
		}
		if printMat {
			printMatrices(dm, sm, r)
		}
		c1 := t[0]
		c2 := t[1]
		c3 := t[2]
		c1.Length = (dm.Matrix[0][1] + dm.Matrix[0][2] -
			dm.Matrix[1][2]) / 2.0
		c2.Length = (dm.Matrix[1][0] + dm.Matrix[1][2] -
			dm.Matrix[0][2]) / 2.0
		c3.Length = (dm.Matrix[2][0] + dm.Matrix[2][1] -
			dm.Matrix[0][1]) / 2.0
		c1.HasLength = true
		c2.HasLength = true
		c3.HasLength = true
		root = nwk.NewNode()
		root.AddChild(c1)
		root.AddChild(c2)
		root.AddChild(c3)
		resetLabels(root)
		fmt.Println(root)
	}
}
func rowSums(dm *dist.DistMat) []float64 {
	n := len(dm.Names)
	r := make([]float64, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			r[i] += dm.Matrix[i][j]
		}
	}
	return r
}
func smat(dm *dist.DistMat, r []float64) *dist.DistMat {
	n := len(dm.Names)
	sm := dist.NewDistMat(n)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			sm.Matrix[i][j] = dm.Matrix[i][j] -
				(r[i]+r[j])/float64(n-2)
			sm.Matrix[j][i] = sm.Matrix[i][j]
		}
	}
	return sm
}
func printMatrices(dm, sm *dist.DistMat, r []float64) {
	w := tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', 0)
	n := len(dm.Names)
	fmt.Fprintf(w, "%d\n", n)
	w.Flush()
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "%s", dm.Names[i])
		for j := 0; j < n; j++ {
			x := sm.Matrix[i][j]
			if i < j {
				x = dm.Matrix[i][j]
			}
			fmt.Fprintf(w, "\t%.3g", x)
		}
		fmt.Fprintf(w, "\t%.3g\n", r[i])
	}
	w.Flush()
}
func resetLabels(v *nwk.Node) {
	if v == nil {
		return
	}
	resetLabels(v.Child)
	resetLabels(v.Sib)
	if v.Child != nil {
		v.Label = ""
	}
}
func main() {
	u := "nj [-h] [option]... [foo.dist]..."
	p := "Calculate neighbor-joining tree."
	e := "nj foo.dist"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optM = flag.Bool("m", false, "print intermediate "+
		"matrices")
	flag.Parse()
	if *optV {
		util.PrintInfo("nj")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optM)
}
