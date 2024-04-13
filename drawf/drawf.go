package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"math/rand"
	"time"
)

type gene struct {
	a      *gene
	d      []*gene
	i      int
	l      string
	isMrca bool
	p      int
}

func main() {
	util.PrepLog("drawf")
	u := "drawf [-h] [option]..."
	p := "Draw Wright-Fisher population."
	e := "drawf | neato -T x11"
	clio.Usage(u, p, e)
	var optN = flag.Int("n", 10, "number of genes")
	var optG = flag.Int("g", 10, "number of generations")
	var optU = flag.Bool("u", false, "untangled lines of descent")
	var optS = flag.Int64("s", 0, "seed for random number generator")
	var optF = flag.Float64("f", 0.4, "scaling factor for plot")
	var optM = flag.Bool("m", false, "mark most recent common "+
		"ancestor")
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("drawf")
	}
	m := *optG
	n := *optN
	wfp := make([][]*gene, m)
	for i := 0; i < m; i++ {
		wfp[i] = make([]*gene, n)
		for j := 0; j < n; j++ {
			wfp[i][j] = new(gene)
			wfp[i][j].d = make([]*gene, 0)
			wfp[i][j].i = j
			wfp[i][j].l = fmt.Sprintf("i%d_%d", i, j)
		}
	}
	genes := wfp[m-1]
	for _, gene := range genes {
		gene.p = 1
	}
	seed := *optS
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	source := rand.NewSource(seed)
	r := rand.New(source)
	for i := 1; i < m; i++ {
		for j := 0; j < n; j++ {
			p := r.Intn(n)
			a := wfp[i-1][p]
			a.d = append(a.d, wfp[i][j])
			wfp[i][j].a = a
		}
	}
	if *optM {
		found := false
		for i := m - 1; i > 0; i-- {
			for _, gene := range wfp[i] {
				gene.a.p += gene.p
				if gene.a.p == n {
					gene.a.isMrca = true
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	if *optU {
		for i := 0; i < m-1; i++ {
			k := 0
			for j := 0; j < n; j++ {
				for _, d := range wfp[i][j].d {
					wfp[i+1][k] = d
					k++
				}
			}
		}
	}
	fmt.Println("# Wright-Fisher population generated with drawf.")
	fmt.Println("# Render with neato, e.g.")
	fmt.Println("# $ neato -T x11 foo.dot")
	fmt.Println("digraph g {")
	fmt.Println("\tnode [shape=point, penwidth=4.0];")
	f := *optF
	for i, genes := range wfp {
		fmt.Printf("\tg_%d[shape=plaintext,pos=\"%.4g,%.4g!\"];",
			i+1, 0.0, float64(m-i)*f)
		for j, gene := range genes {
			fmt.Printf("%s[pos=\"%.4g,%.4g!\"", gene.l,
				float64(j+1)*f, float64(m-i)*f)
			// if gene.isMrca {
			//    fmt.Printf(",color=\"red\"")
			// }
			fmt.Printf("];")
		}
		fmt.Printf("\n")
	}
	genes = wfp[m-1]
	fmt.Println("\tnode [shape=plaintext]")
	fmt.Printf("\t")
	for i, gene := range genes {
		x := float64(i+1) * f
		fmt.Printf("%d[pos=\"%.4g,%.4g!\"];",
			gene.i+1, x, 0.0)
	}
	fmt.Printf("\n")
	fmt.Println("\tedge [arrowhead=none,penwidth=2.0];")
	for i := 1; i < m; i++ {
		genes := wfp[i]
		fmt.Printf("\t")
		for _, g := range genes {
			fmt.Printf("%s->%s;", g.l, g.a.l)
		}
		fmt.Printf("\n")
	}
	if *optM {
		for i, genes := range wfp {
			for j, gene := range genes {
				if gene.isMrca {
					fmt.Printf("mrca[pos=\"%.4g,%.4g!\"",
						float64(j+1)*f, float64(m-i)*f)
					fmt.Printf("shape=point,color=\"red\"];")
				}
			}
		}
	}
	fmt.Println("}")
}
