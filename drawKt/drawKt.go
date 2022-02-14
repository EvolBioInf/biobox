package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/kt"
	"log"
	"os"
	"strings"
)

func writeLatex(root *kt.Node, optL bool) string {
	w := new(bytes.Buffer)
	var maxDepth int
	kt.BreadthFirst(root, findMaxDepth, &maxDepth)
	var nl int
	kt.BreadthFirst(root, countLeaves, &nl)
	dist := float64(nl) / (float64(nl) - 1.0)
	xcoords := make([]float64, kt.NodeCount())
	var curX float64
	postorder(root, setX, &curX, dist, xcoords)
	var x1 float64
	y1 := -0.8
	x2 := float64(nl)
	y2 := float64(maxDepth) + 0.7
	if optL {
		x1 -= 0.3
		y1 -= 0.2
		x2 += 0.3
		y2 += 0.5
	}
	fmt.Fprintf(w, "\\begin{pspicture}(%.2g,%.2g)(%.2g,%.2g)\n",
		x1, y1, x2, y2)
	fmt.Fprint(w, "%% Nodes\n")
	kt.BreadthFirst(root, writeLatexNode, w, xcoords, maxDepth, optL)
	fmt.Fprint(w, "%% Match links\n")
	fmt.Fprintf(w, "\\psset{linecolor=lightgray}")
	kt.BreadthFirst(root, writeMatchLink, w)
	fmt.Fprint(w, "%% Failure links\n")
	fmt.Fprint(w, "\\psset{linecolor=red,linewidth=0.5pt,nodesep=2pt}\n")
	kt.BreadthFirst(root, writeFailureLink, w, xcoords)
	fmt.Fprint(w, "%% Output sets\n")
	kt.BreadthFirst(root, writeOutputSet, w, optL)
	fmt.Fprintf(w, "\\end{pspicture}")
	return (w.String())
}
func findMaxDepth(v *kt.Node, args ...interface{}) {
	maxDepth := args[0].(*int)
	if v.Depth > *maxDepth {
		*maxDepth = v.Depth
	}
}
func countLeaves(n *kt.Node, args ...interface{}) {
	nl := args[0].(*int)
	if n.Child == nil {
		*nl = *nl + 1
	}
}
func postorder(v *kt.Node, fn kt.NodeAction, args ...interface{}) {
	if v != nil {
		postorder(v.Child, fn, args...)
		fn(v, args...)
		postorder(v.Sib, fn, args...)
	}
}
func setX(v *kt.Node, args ...interface{}) {
	curX := args[0].(*float64)
	dist := args[1].(float64)
	xcoords := args[2].([]float64)
	if v.Child == nil {
		xcoords[v.Id] = *curX
		*curX = *curX + dist
	} else {
		x1 := xcoords[v.Child.Id]
		cp := v.Child
		for cp.Sib != nil {
			cp = cp.Sib
		}
		x2 := xcoords[cp.Id]
		xcoords[v.Id] = (x1 + x2) / 2.0
	}
}
func writeLatexNode(v *kt.Node, args ...interface{}) {
	w := args[0].(*bytes.Buffer)
	xcoords := args[1].([]float64)
	maxDepth := args[2].(int)
	optL := args[3].(bool)
	y := maxDepth - v.Depth
	x := xcoords[v.Id]
	if optL {
		fmt.Fprintf(w, "\\cnodeput(%.3g,%d){%d}{%d}\n",
			x, y, v.Id, v.Id+1)
	} else {
		fmt.Fprintf(w, "\\dotnode(%.3g,%d){%d}\n",
			x, y, v.Id)
	}
}
func writeMatchLink(v *kt.Node, args ...interface{}) {
	w := args[0].(*bytes.Buffer)
	if v.Parent != nil {
		p := v.Parent.Id
		n := v.Id
		c := v.In
		fmt.Fprintf(w, "\\ncline{%d}{%d}"+
			"\\ncput[nrot=:U]{\\texttt{%c}}\n",
			p, n, c)
	}
}
func writeFailureLink(v *kt.Node, args ...interface{}) {
	w := args[0].(*bytes.Buffer)
	xcoords := args[1].([]float64)
	if v.Parent == nil {
		fmt.Fprint(w, "\\nccurve[angleA=130,angleB=50,ncurv=6]"+
			"{->}{0}{0}\n")
	} else {
		angle := 50
		x1 := xcoords[v.Id]
		x2 := xcoords[v.Fail.Id]
		if x1 > x2 {
			angle = -50
		}
		fmt.Fprintf(w, "\\ncarc[arcangle=%d]{->}{%d}{%d}\n",
			angle, v.Id, v.Fail.Id)
	}
}
func writeOutputSet(v *kt.Node, args ...interface{}) {
	if len(v.Output) == 0 {
		return
	}
	w := args[0].(*bytes.Buffer)
	optL := args[1].(bool)
	angle := 0
	if v.Child == nil {
		angle = -90
	}
	fmt.Fprintf(w, "\\nput{%d}{%d}{$\\{", angle, v.Id)
	if optL {
		fmt.Fprintf(w, "p_")
	}
	fmt.Fprintf(w, "%d", v.Output[0]+1)
	for i := 1; i < len(v.Output); i++ {
		fmt.Fprintf(w, ",")
		if optL {
			fmt.Fprintf(w, "p_")
		}
		fmt.Fprintf(w, "%d", v.Output[i]+1)
	}
	fmt.Fprint(w, "\\}$}\n")
}
func main() {
	u := "drawKt [-h] [options] [patterns]"
	p := "Draw the keyword tree of a set of patterns"
	e := "drawKt ATTT ATTC AT TG TT > kt.tex"
	clio.Usage(u, p, e)
	var optW = flag.String("w", "", "LaTeX wrapper file")
	var optL = flag.Bool("l", false, "labeled nodes; default: plain")
	var optT = flag.Bool("t", false, "plain text; default: LaTeX")
	var optV = flag.Bool("v", false, "print program version & "+
		"other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("drawKt")
	}
	var patterns []string
	if len(flag.Args()) > 0 {
		patterns = flag.Args()
	} else {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			patterns = append(patterns, sc.Text())
		}
	}
	tree := kt.NewKeywordTree(patterns)
	if *optT {
		fmt.Println(tree)
	} else {
		fmt.Println(writeLatex(tree, *optL))
		if *optW != "" {
			f, err := os.Create(*optW)
			if err != nil {
				log.Fatalf("couldn't open %q\n", *optW)
			}
			fmt.Fprintf(f, "\\documentclass{article}\n")
			fmt.Fprintf(f, "\\usepackage{pst-all}\n")
			fmt.Fprintf(f, "\\begin{document}\n")
			fmt.Fprintf(f, "\\begin{center}\n\\input{kt.tex}\n\\end{center}\n")
			fmt.Fprintf(f, "\\end{document}\n")
			f.Close()
			old := *optW
			new := strings.TrimSuffix(old, ".tex")
			fmt.Fprintf(os.Stderr, "# Wrote wrapper %s; if the keyword tree is in "+
				"kt.tex, run \n# latex %s\n# dvips %s -o -q\n# "+
				"ps2pdf %s.ps\n", old, new, new, new)
		}
	}
}
