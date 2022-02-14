package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"math/rand"
	"time"
)

func main() {
	u := "ranDot [-h] [option]..."
	p := "Draw random graph in dot notation."
	e := "ranDot -c lightsalmon -C lightgray -d"
	clio.Usage(u, p, e)
	var optN = flag.Int("n", 10, "number of nodes")
	var optP = flag.Float64("p", 0.05, "edge probability")
	var optD = flag.Bool("d", false, "directed edges")
	var optSS = flag.Bool("S", false, "allow edge to self")
	var optC = flag.String("c", "", "color of connected nodes")
	var optCC = flag.String("C", "", "color of singleton nodes; "+
		"color names: www.graphviz.org/doc/info/colors.html")
	var optS = flag.Int64("s", 0, "seed for random number generator "+
		"(default internal)")
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("ranDot")
	}
	seed := *optS
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	source := rand.NewSource(seed)
	r := rand.New(source)
	n := *optN
	nodes := make([]int, n)
	for i := 0; i < n; i++ {
		nodes[i] = i + 1
	}
	edges := make([][]bool, n)
	for i := 0; i < n; i++ {
		edges[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j && !*optSS {
				continue
			}
			if r.Float64() <= *optP {
				edges[i][j] = true
			}
		}
	}
	fmt.Println("# Graph written by ranDot.")
	fmt.Println("# Render: dot|neato|circo foo.dot")
	fmt.Println("graph G {")
	if *optCC != "" {
		fmt.Printf("node [style=filled, color=%s]\n", *optCC)
	}
	singletons := make([]bool, n)
	for i := 0; i < n; i++ {
		singletons[i] = true
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if edges[i][j] {
				singletons[i] = false
				singletons[j] = false
			}
		}
	}
	for i := 0; i < n; i++ {
		if singletons[i] {
			fmt.Printf("\t%d\n", nodes[i])
		}
	}
	if *optCC != "" && *optC == "" {
		fmt.Printf("node [style=\"\", color=\"\"]\n")
	}
	if *optC != "" {
		fmt.Printf("node [style=filled, color=%s]\n", *optC)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if edges[i][j] {
				fmt.Printf("\t%d -- %d", nodes[i], nodes[j])
				if *optD {
					fmt.Printf("[dir=")
					if edges[j][i] {
						fmt.Printf("both]")
						edges[j][i] = false
					} else {
						fmt.Printf("forward]")
					}
				}
				fmt.Printf("\n")
				edges[i][j] = false
			}
		}
	}
	fmt.Println("}")
}
