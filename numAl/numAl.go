package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type node struct {
	m, n, i int
}

var optV = flag.Bool("v", false, "print version & "+
	"program information")
var optT = flag.Bool("t", false, "top-down (default bottom-up)")
var optP = flag.Bool("p", false, "print data structure (default result)")
var nodeId int

func topDown(m, n int) float64 {
	if m > 0 && n > 0 {
		r := topDown(m-1, n) + topDown(m, n-1) +
			topDown(m-1, n-1)
		return r
	} else {
		return 1.0
	}
}
func bottomUp(m, n int) [][]float64 {
	var mat [][]float64
	mat = make([][]float64, m+1)
	for i := 0; i <= m; i++ {
		mat[i] = make([]float64, n+1)
	}
	for i := 0; i <= m; i++ {
		mat[i][0] = 1
	}
	for i := 1; i <= n; i++ {
		mat[0][i] = 1
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			mat[i][j] = mat[i-1][j] +
				mat[i][j-1] +
				mat[i-1][j-1]
		}
	}
	return mat
}
func newNode(m, n int) *node {
	v := new(node)
	v.m = m
	v.n = n
	v.i = nodeId
	nodeId++
	return v
}
func printTopDown(v *node) {
	if v.m == 0 || v.n == 0 {
		return
	}
	c1 := newNode(v.m-1, v.n)
	c2 := newNode(v.m, v.n-1)
	c3 := newNode(v.m-1, v.n-1)
	fmt.Printf("\tn%d->n%d\n", v.i, c1.i)
	fmt.Printf("\tn%d->n%d\n", v.i, c2.i)
	fmt.Printf("\tn%d->n%d\n", v.i, c3.i)
	printTopDown(c1)
	printTopDown(c2)
	printTopDown(c3)
}
func main() {
	u := "numAl [-h] [options] m n"
	p := "Compute the number of possible global alignments " +
		"between two sequences of lengths m and n"
	e := "numAl 5 10"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("numAl")
	}
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "please provide two "+
			"sequence lengths\n")
		os.Exit(0)
	}
	m, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't convert %q\n", args[0])
		os.Exit(0)
	}
	n, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't convert %q\n", args[1])
		os.Exit(0)
	}
	if !*optP {
		var du time.Duration
		var na float64
		start := time.Now()
		if *optT {
			na = topDown(m, n)
		} else {
			mat := bottomUp(m, n)
			na = mat[m][n]
		}
		end := time.Now()
		du = end.Sub(start)
		fmt.Printf("f(%d, %d) = %g (%g s)\n", m, n, na, du.Seconds())
	} else {
		if *optT {
			fmt.Println("# Recursion tree for computing the number of alignments.")
			fmt.Println("# Generated with numAl, render with")
			fmt.Println("# $ dot -T x11 foo.dot")
			fmt.Println("digraph g {")
			fmt.Println("\tnode[shape=point]")
			fmt.Println("\tedge[arrowhead=none]")
			r := newNode(m, n)
			printTopDown(r)
			fmt.Println("}")
		} else {
			mat := bottomUp(m, n)
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ',
				tabwriter.AlignRight)
			for i := 0; i <= n; i++ {
				fmt.Fprintf(w, "\t%d", i)
			}
			fmt.Fprintf(w, "\t\n")
			for i := 0; i <= m; i++ {
				fmt.Fprintf(w, "%d", i)
				for j := 0; j <= n; j++ {
					fmt.Fprintf(w, "\t%d", int(mat[i][j]))
				}
				fmt.Fprintf(w, "\t\n")
			}

			w.Flush()
		}
	}
}
