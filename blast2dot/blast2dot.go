package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"strings"
)

func scan(r io.Reader, args ...interface{}) {
	optS := args[0].(bool)
	optC := args[1].(string)
	optCC := args[2].(string)
	accessions := make(map[string]int)
	families := make(map[string]map[string]bool)
	n := 1
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Fields(line)
		query := fields[0]
		sbjct := fields[1]
		if accessions[query] == 0 {
			accessions[query] = n
			n++
		}
		if query != sbjct {
			if families[query] == nil {
				families[query] = make(map[string]bool)
			}
			qm := families[query]
			qm[sbjct] = true
		}
	}
	mm := make([][]bool, n)
	for i := 0; i < n; i++ {
		mm[i] = make([]bool, n)
	}
	for k, v := range accessions {
		accessions[k] = v - 1
	}
	for q, m := range families {
		i := accessions[q]
		for s, _ := range m {
			j := accessions[s]
			mm[i][j] = true
		}
	}
	names := make([]string, n)
	for k, v := range accessions {
		names[v] = k
	}
	fmt.Println("# Graph written by blast2dot.")
	fmt.Println("# Render: dot|neato|circo foo.dot")
	fmt.Println("graph G {")
	if optC != "" {
		fmt.Printf("node [style=filled, color=%s]\n", optC)
	}
	for i, v := range mm {
		for j, _ := range v {
			if i != j && mm[i][j] {
				fmt.Printf("\t%s -- %s[dir=", names[i], names[j])
				if mm[j][i] {
					fmt.Printf("both]\n")
					mm[j][i] = false
				} else {
					fmt.Printf("forward]\n")
				}
				mm[i][j] = false
			}
		}
	}
	if optS {
		if optCC != "" {
			fmt.Printf("node [style=filled, color=%s]\n", optCC)
		} else if optC != "" {
			fmt.Println("nod [style=\"\", color=\"\"]")
		}
		for k, _ := range accessions {
			if families[k] == nil {
				fmt.Printf("\t%s\n", k)
			}
		}
	}
	fmt.Println("}")
}
func main() {
	u := "blast2dot [-h] [option]... [file]..."
	p := "Convert BLAST output to dot code " +
		"for plotting with GraphViz programs " +
		"like dot, neato, or circo."
	e := "blast2dot -C lightgray -c lightsalmon foo.bl | neato -T x11"
	clio.Usage(u, p, e)
	var optS = flag.Bool("s", false, "include singletons")
	var optC = flag.String("c", "", "color of gene families")
	var optCC = flag.String("C", "", "color of singletons; color names: "+
		"www.graphviz.org/doc/info/colors.html")
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("blast2dot")
	}
	if *optCC != "" {
		*optS = true
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optS, *optC, *optCC)
}
