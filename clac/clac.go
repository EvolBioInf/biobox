package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

type clade struct {
	k string
	n int
}
type cladeSlice []clade

func (c cladeSlice) Len() int {
	return len(c)
}
func (c cladeSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c cladeSlice) Less(i, j int) bool {
	if c[i].n != c[j].n {
		return c[i].n < c[j].n
	} else {
		return c[i].k < c[j].k
	}
}
func scan(r io.Reader, args ...interface{}) {
	refTrees := args[0].([]*nwk.Node)
	sc := nwk.NewScanner(r)
	clades := make(map[string]int)
	nt := 0
	for sc.Scan() {
		root := sc.Tree()
		nt++
		countClades(root, clades)
	}
	if len(refTrees) > 0 {
		for _, root := range refTrees {
			annotateTree(root, clades, nt)
			fmt.Println(root)
		}
	} else {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "#ID\tCount\tTaxa\tClade\n")
		cs := make([]clade, 0)
		var c clade
		for k, n := range clades {
			c.k = k
			c.n = n
			cs = append(cs, c)
		}
		sort.Sort(cladeSlice(cs))
		x := 0
		for i := len(cs) - 1; i >= 0; i-- {
			x++
			taxa := strings.Split(cs[i].k, "$")
			t := len(taxa)
			fmt.Fprintf(w, "%d\t%d\t%d\t{", x, cs[i].n, t)
			for j, s := range taxa {
				if j > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, s)
			}
			fmt.Fprintf(w, "}\n")
		}
		w.Flush()
	}
}
func countClades(v *nwk.Node, c map[string]int) {
	if v == nil {
		return
	}
	if v.Parent != nil && v.Child != nil {
		k := v.Key("$")
		c[k]++
	}
	countClades(v.Child, c)
	countClades(v.Sib, c)
}
func annotateTree(v *nwk.Node, clades map[string]int, nc int) {
	if v == nil {
		return
	}
	if v.Parent != nil && v.Child != nil {
		p := float64(clades[v.Key("$")]) /
			float64(nc) * 100.0
		p = math.Round(p)
		v.Label = strconv.Itoa(int(p))
	}
	annotateTree(v.Child, clades, nc)
	annotateTree(v.Sib, clades, nc)
}
func main() {
	util.PrepLog("clac")
	u := "clac [-h] [option]... [trees.nwk]..."
	p := "Count the clades in phylogenies."
	e := "dnaDist -b 1000 foo.fasta | nj | clac"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optR = flag.String("r", "", "file of reference tree(s)")
	flag.Parse()
	if *optV {
		util.PrintInfo("clac")
	}
	var refTrees []*nwk.Node
	if *optR != "" {
		tf, err := os.Open(*optR)
		if err != nil {
			log.Fatalf("couldn't open %q", *optR)
		}
		defer tf.Close()
		sc := nwk.NewScanner(tf)
		for sc.Scan() {
			refTrees = append(refTrees, sc.Tree())
		}
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, refTrees)
}
