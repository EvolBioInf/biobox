package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"log"
	"strconv"
	"strings"
)

type Pair struct {
	c1, c2 int
}
type Node struct {
	id             int
	child1, child2 *Node
}

func parse(r io.Reader, args ...interface{}) {
	pairs := []Pair{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if line[0] == '#' {
			continue
		}
		fields := strings.Fields(line)
		c1, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal(err)
		}
		c2, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Fatal(err)
		}
		pair := Pair{c1 - 1, c2 - 1}
		pairs = append(pairs, pair)
	}
	n := len(pairs) + 1
	tree := make([]*Node, 2*n-1)
	for i := 0; i < 2*n-1; i++ {
		tree[i] = new(Node)
		tree[i].id = i
	}
	for i := n; i > 1; i = i - 1 {
		p := 2*n - i
		c := pairs[p-n].c1
		tree[p].child1 = tree[c]
		tree[c] = tree[i-1]
		c = pairs[p-n].c2
		tree[p].child2 = tree[c]
		tree[c] = tree[p]
	}
	fmt.Println("digraph g {")
	fmt.Println("\tedge [arrowhead=\"none\"]")
	for i := n; i < 2*n-1; i++ {
		p := tree[i].id
		c1 := tree[i].child1.id
		c2 := tree[i].child2.id
		fmt.Printf("\t%d -> %d\n", p+1, c1+1)
		fmt.Printf("\t%d -> %d\n", p+1, c2+1)
	}
	fmt.Println("}")
}
func main() {
	util.PrepLog("drawChildren")
	u := "drawChildren [-h] [option] [foo.txt]..."
	p := "Draw the coalescent implied by the output of pickChildren."
	e := "pickChildren | tee pc.out | drawChildren | dot -T x11"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("drawChildren")
	}
	files := flag.Args()
	clio.ParseFiles(files, parse)
}
