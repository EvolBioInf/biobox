package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type indiv struct {
	a             [2]*indiv
	p             int
	isMale        bool
	n             string
	ag            [2]string
	g             [2]bool
	isOnPath      bool
	isUa, isNonUa bool
}

func main() {
	util.PrepLog("drag")
	u := "drag [-h] [option]..."
	p := "Draw genealogy of diploid individuals."
	e := "drag -t 4,6 | neato -T x11"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optG = flag.Int("g", 10, "number of generations")
	var optN = flag.Int("n", 10, "number of individuals")
	var optT = flag.String("t", "", "trace genealogy of "+
		"individuals, e.g. 3,4,5; -1 for all")
	var optGG = flag.Bool("G", false, "trace genes")
	var optF = flag.Float64("f", 1.0, "scale factor for plot")
	var optA = flag.Bool("a", false, "ancestor statistics")
	var optS = flag.Int64("s", 0, "seed for random number generator")
	flag.Parse()
	if *optV {
		util.PrintInfo("drag")
	}
	var tr []int
	if *optT != "" {
		s := *optT
		if s[0] == '-' {
			for i := 0; i < *optN; i++ {
				tr = append(tr, i)
			}
		} else {
			fields := strings.Split(*optT, ",")
			for _, field := range fields {
				i, err := strconv.Atoi(field)
				if err != nil {
					log.Fatalf("can't convert %q", field)
				}
				tr = append(tr, i-1)
			}
		}
	}
	seed := *optS
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	m := *optG
	n := *optN
	pop := make([][]*indiv, m)
	for i := 0; i < m; i++ {
		pop[i] = make([]*indiv, n)
		nf := n
		for j := 0; j < n; j++ {
			pop[i][j] = new(indiv)
			pop[i][j].n = fmt.Sprintf("i_%d_%d", i, j)
			if rand.Float64() < 0.5 {
				pop[i][j].isMale = true
				nf--
			}
		}
		if nf == 0 {
			r := rand.Intn(n)
			pop[i][r].isMale = false
		} else if nf == n {
			r := rand.Intn(n)
			pop[i][r].isMale = true
		}
	}
	for i := m - 1; i > 0; i-- {
		for j := 0; j < n; j++ {
			pop[i][j].a[0] = pop[i-1][rand.Intn(n)]
			pop[i][j].a[1] = pop[i-1][rand.Intn(n)]
			for pop[i][j].a[0].isMale == pop[i][j].a[1].isMale {
				pop[i][j].a[1] = pop[i-1][rand.Intn(n)]
			}
			for k := 0; k < 2; k++ {
				r := rand.Intn(2)
				name := fmt.Sprintf("%s_%d", pop[i][j].a[k].n, r)
				pop[i][j].ag[k] = name
			}
		}
	}
	for i := 0; i < n; i++ {
		pop[m-1][i].p = 1
	}
	for i := m - 1; i > 0; i-- {
		for j := 0; j < n; j++ {
			for k := 0; k < 2; k++ {
				if pop[i][j].a[k].p < n {
					pop[i][j].a[k].p += pop[i][j].p
				}
			}
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if pop[i][j].p >= n {
				pop[i][j].isUa = true
			} else if pop[i][j].p == 0 {
				pop[i][j].isNonUa = true
			}
		}
	}
	for _, t := range tr {
		pop[m-1][t].isOnPath = true
		pop[m-1][t].a[0].isOnPath = true
		pop[m-1][t].a[1].isOnPath = true
		if *optGG {
			pop[m-1][t].g[0] = true
			pop[m-1][t].g[1] = true
		}
	}
	for i := m - 2; i > 0; i-- {
		for j := 0; j < n; j++ {
			if pop[i][j].isOnPath {
				pop[i][j].a[0].isOnPath = true
				pop[i][j].a[1].isOnPath = true
			}
		}
	}
	if *optGG {
		for _, t := range tr {
			for i := 0; i < 2; i++ {
				pop[m-1][t].g[i] = true
				l := len(pop[m-1][t].ag[i])
				if pop[m-1][t].ag[i][l-1] == '0' {
					pop[m-1][t].a[i].g[0] = true
				} else {
					pop[m-1][t].a[i].g[1] = true
				}
			}
		}
		for i := m - 2; i > 0; i-- {
			for j := 0; j < n; j++ {
				for k := 0; k < 2; k++ {
					if pop[i][j].g[k] {
						l := len(pop[i][j].ag[k])
						if pop[i][j].ag[k][l-1] == '0' {
							pop[i][j].a[k].g[0] = true
						} else {
							pop[i][j].a[k].g[1] = true
						}
					}
				}
			}
		}
	}
	if *optA {
		tua := 0
		foundUa := false
		for i := m - 1; i > -1; i-- {
			for j := 0; j < n; j++ {
				if pop[i][j].isUa {
					foundUa = true
					break
				}
			}
			if foundUa {
				break
			}
			tua++
		}
		tpa := 0
		foundPa := false
		for i := m - 1; i > -1; i-- {
			npa := 0
			for j := 0; j < n; j++ {
				if !pop[i][j].isUa && !pop[i][j].isNonUa {
					npa++
				}
			}
			if npa == 0 {
				foundPa = true
				break
			}
			tpa++
		}
		if !foundUa {
			tua = 0
		}
		if !foundPa {
			tpa = 0
		}
		m1 := "Generations_to_first_universal_ancestor\t%d\n"
		m2 := "Generations_to_no_partial_ancestor\t%d\n"
		fmt.Printf(m1, tua)
		fmt.Printf(m2, tpa)
	} else {
		fmt.Println("# Genealogy generated with drag.")
		fmt.Println("# Render with neato.")
		fmt.Println("graph g {")
		f := *optF
		t := "%c_%d[shape=plaintext,pos=\"%.4g,%.4g!\"];"
		for i := 0; i < m; i++ {
			y := float64(m-i) * f
			fmt.Printf("\t"+t, 'g', i+1, 0.0, y)
			for j := 0; j < n; j++ {
				in := pop[i][j]
				var c, s string
				var x float64
				c = "lightgreen"
				if in.isUa {
					c = "salmon"
				} else if in.isNonUa {
					c = "lightblue"
				}
				s = "ellipse"
				if in.isMale {
					s = "box"
				}
				x = float64(j+1) * f
				tmpl := "%s[label=\"\",color=%s,shape=%s," +
					"style=filled,pos=\"%.4g,%.4g!\"];"
				fmt.Printf(tmpl, in.n, c, s, x, y)
			}
			fmt.Printf(t+"\n", 'b', m-1-i, float64(n+1)*f, y)
		}
		y := -0.0 * f
		for i := 0; i < n; i++ {
			x := float64(i+1) * f
			fmt.Printf(t, 'i', i+1, x, y)
		}
		fmt.Printf("\tnode[shape=point,penwidth=4];\n")
		for i := 0; i < m; i++ {
			fmt.Printf("\t")
			for j := 0; j < n; j++ {
				in := pop[i][j]
				x := float64(j+1) * f
				y := float64(m-i) * f
				name := in.n + "_0"
				tmpl := "%s[pos=\"%.4g,%.4g!\"];"
				fmt.Printf(tmpl, name, x-0.15, y)
				name = in.n + "_1"
				fmt.Printf(tmpl, name, x+0.15, y)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("edge[color=black]")
		for i := 1; i < m; i++ {
			for j := 0; j < n; j++ {
				in := pop[i][j]
				if in.isOnPath {
					if *optGG {
						for k := 0; k < 2; k++ {
							if in.g[k] {
								so := fmt.Sprintf("%s_%d", in.n, k)
								de := fmt.Sprintf("%s", in.ag[k])
								fmt.Printf("\t%s--%s\n", so, de)
							}
						}
					} else {
						fmt.Printf("\t%s_0--%s;", in.n, in.ag[0])
						fmt.Printf("%s_1--%s;\n", in.n, in.ag[1])
					}
				}
			}
		}
		fmt.Printf("}\n")
	}
}
