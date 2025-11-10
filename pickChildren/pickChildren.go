package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"log"
	"math/rand"
	"time"
)

type Node struct {
	child1, child2 *Node
}

func main() {
	util.PrepLog("pickChildren")
	u := "pickChildren [-h] [option]"
	p := "Demo the construction of a coalescent topology."
	e := "pickChildren -n 4 | tee pc.txt | drawChildren |" +
		"dot -T x11"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optN := flag.Int("n", 10, "sample size")
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
	fmt.Printf("# p\tc_1\tc_2\n")
	tree := make([]*Node, 2*n-1)
	for i := 0; i < 2*n-1; i++ {
		tree[i] = new(Node)
	}
	for i := n; i >= 2; i-- {
		p := 2*n - i
		fmt.Printf("%d\t", p+1)
		c := int(float64(i) * ran.Float64())
		tree[p].child1 = tree[c]
		fmt.Printf("%d\t", c+1)
		tree[c] = tree[i]
		c = int(float64(i-1) * ran.Float64())
		tree[p].child2 = tree[c]
		fmt.Printf("%d\n", c+1)
		tree[c] = tree[p]
	}
}
