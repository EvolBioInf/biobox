package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"log"
	"math/rand"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type Node struct {
	id             int
	child1, child2 *Node
}

func setTimes(v *nwk.Node, n int) {
	if v == nil {
		return
	}
	setTimes(v.Child, n)
	setTimes(v.Sib, n)
	if v.Parent == nil {
		return
	}
	v.HasLength = true
	parentHeight := v.Parent.Id - n + 1
	nodeHeight := v.Id - n + 1
	if nodeHeight < 0 {
		nodeHeight = 0
	}
	v.Length = float64(parentHeight) - float64(nodeHeight)
}
func main() {
	util.PrepLog("pickChildren")
	u := "pickChildren [-h] [option]"
	p := "Demo the construction of a coalescent topology."
	e := "pickChildren -n 5 -t pc.nwk"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optN := flag.Int("n", 4, "sample size")
	optT := flag.String("t", "", "print tree to file")
	optS := flag.Int("s", 0, "seed for random number generator")
	flag.Parse()
	if *optV {
		util.PrintInfo("pickChildren")
	}
	n := *optN
	if n <= 2 {
		log.Fatal("please use a sample size " +
			"of at least 2")
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	ran := rand.New(rand.NewSource(seed))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "# p\tc_1\tc_2\t\t")
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "\tt_%d", i+1)
	}
	fmt.Fprintf(w, "\t")
	for i := n; i < 2*n-1; i++ {
		fmt.Fprintf(w, "\tt_%d", i+1)
	}
	fmt.Fprintf(w, "\n")
	tree := make([]*nwk.Node, 2*n-1)
	for i := 0; i < 2*n-1; i++ {
		v := new(nwk.Node)
		v.Id = i
		v.Label = strconv.Itoa(i + 1)
		tree[i] = v
	}
	root := tree[2*n-2]
	for i := n; i >= 2; i-- {
		p := 2*n - i
		fmt.Fprintf(w, "%d", p+1)
		c := int(float64(i) * ran.Float64())
		tree[p].AddChild(tree[c])
		fmt.Fprintf(w, "\t%d", c+1)
		tree[c] = tree[i-1]
		tree[c].Parent = tree[p]
		tree[i-1] = nil
		c = int(float64(i-1) * ran.Float64())
		tree[p].AddChild(tree[c])
		fmt.Fprintf(w, "\t%d\t\t", c+1)
		tree[c] = tree[p]
		tree[c].Parent = tree[p]
		tree[p] = nil
		for i := 0; i < n; i++ {
			if tree[i] == nil {
				fmt.Fprintf(w, "\t-")
			} else {
				fmt.Fprintf(w, "\t%d", tree[i].Id+1)
			}
		}
		fmt.Fprintf(w, "\t")
		for i := n; i < 2*n-1; i++ {
			if tree[i] == nil {
				fmt.Fprintf(w, "\t-")
			} else {
				fmt.Fprintf(w, "\t%d", tree[i].Id+1)
			}
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	if *optT != "" {
		root.Parent = nil
		setTimes(root, n)
		f, err := os.Create(*optT)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(f, "%s\n", root)
		f.Close()
	}
}
