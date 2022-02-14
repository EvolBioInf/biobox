package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/esa"
	"github.com/evolbioinf/fasta"
	"io"
	"strings"
)

var nodeId int

type node struct {
	d, l, r, id, level int
	child, sib, parent *node
}
type stack []*node
type nodeAction func(*node, ...interface{})

func (v *node) String() string {
	if v == nil {
		return "!"
	}
	s := fmt.Sprintf("%d-[%d..%d]", v.d, v.l, v.r)
	return s
}
func (s *stack) push(n *node) { *s = append(*s, n) }
func (s *stack) pop() *node {
	n := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return n
}
func (s *stack) top() *node { return (*s)[len(*s)-1] }
func (p *node) addChild(c *node) {
	c.parent = p
	if p.child == nil {
		p.child = c
	} else {
		w := p.child
		if c.l < w.l {
			p.child = c
			c.sib = w
			return
		}
		for w.sib != nil {
			if c.l > w.r && c.l < w.sib.l {
				c.sib = w.sib
				w.sib = c
				return
			}
			w = w.sib
		}
		w.sib = c
	}
}
func scan(r io.Reader, args ...interface{}) {
	optI := args[0].(bool)
	optN := args[1].(bool)
	optD := args[2].(bool)
	optL := args[3].(bool)
	optX := args[4].(float64)
	optY := args[5].(float64)
	optS := args[6].(bool)
	scanner := fasta.NewScanner(r)
	for scanner.ScanSequence() {
		sequence := scanner.Sequence()
		data := sequence.Data()
		if optS {
			data = append(data, '$')
		}
		sa := esa.Sa(data)
		lcp := esa.Lcp(data, sa)
		lcp = append(lcp, -1)
		n := len(lcp)
		var v *node
		root := newNode(0, 0, -1, nil)
		stack := new(stack)
		stack.push(root)
		for i := 1; i < n; i++ {
			l := i - 1
			for len(*stack) > 0 && lcp[i] < stack.top().d {
				stack.top().r = i - 1
				v = stack.pop()
				l = v.l
				if len(*stack) > 0 && lcp[i] <= stack.top().d {
					p := stack.top()
					p.addChild(v)
					v = nil
				}
			}
			if len(*stack) > 0 && lcp[i] > stack.top().d {
				w := newNode(lcp[i], l, -1, v)
				stack.push(w)
				v = nil
			}
		}
		if !optI {
			traverse(root, addLeaves, sa)
			preorder(root)
		}
		if optI {
			fmt.Printf("\\psset{nodesep=2pt, levelsep=1cm}\n")
			printIntervals(root)
		} else if optN {
			printNewick(root, sa)
		} else {
			l := len(data)
			x := float64(l) * optX
			m := maxNodeLevel(root, 0)
			y := float64(m) * optY
			fmt.Printf("\\begin{pspicture}(%.2g,%.2g)(%.2g,%.2g)\n",
				0.0, -y, x, 0.0)
			fmt.Printf("\\psset{xunit=%.3g, yunit=%.3g}\n", optX, optY)
			traverse(root, drawCedge, sa, data)
			traverse(root, drawCnode, sa, optL, optD)
			fmt.Printf("\\end{pspicture}\n")
		}
	}
}
func newNode(d, l, r int, child *node) *node {
	n := new(node)
	n.d = d
	n.l = l
	n.r = r
	n.id = nodeId
	nodeId++
	if child != nil {
		n.child = child
		child.parent = n
	}
	return n
}
func preorder(v *node) {
	if v != nil {
		if v.parent != nil {
			v.level = v.parent.level + 1
		}
		preorder(v.child)
		preorder(v.sib)
	}
}
func traverse(v *node, fn nodeAction, args ...interface{}) {
	if v != nil {
		traverse(v.child, fn, args...)
		traverse(v.sib, fn, args...)
		fn(v, args...)
	}
}
func addLeaves(p *node, args ...interface{}) {
	sa := args[0].([]int)
	l := len(sa)
	if p.child == nil {
		for i := p.l; i <= p.r; i++ {
			c := newNode(l-sa[i], i, i, nil)
			p.addChild(c)
		}
	} else {
		for i := p.l; i < p.child.l; i++ {
			c := newNode(l-sa[i], i, i, nil)
			p.addChild(c)
		}
		v := p.child
		for v.sib != nil {
			x := v.sib.l
			for i := v.r + 1; i < x; i++ {
				c := newNode(l-sa[i], i, i, nil)
				p.addChild(c)
			}
			v = v.sib
		}
		for i := v.r + 1; i <= p.r; i++ {
			c := newNode(l-sa[i], i, i, nil)
			p.addChild(c)
		}
	}
}
func maxNodeLevel(v *node, m int) int {
	if v != nil {
		if v.level > m {
			m = v.level
		}
		m = maxNodeLevel(v.child, m)
		m = maxNodeLevel(v.sib, m)
	}
	return m
}
func drawCnode(v *node, args ...interface{}) {
	sa := args[0].([]int)
	nodeLabel := args[1].(bool)
	depth := args[2].(bool)
	x := float64(v.l+v.r) / 2.0
	if nodeLabel {
		fmt.Printf("\\rput(%.3g,%d){\\rnode{%d}{"+
			"\\psframebox[linecolor=lightgray]{%d}}}",
			x, -v.level, v.id, v.id)
	} else {
		fmt.Printf("\\dotnode(%.3g,%d){%d}\n",
			x, -v.level, v.id)
	}
	if v.child == nil {
		fmt.Printf("\\nput{-90}{%d}{%d}\n",
			v.id, sa[v.l]+1)
	} else if depth {
		fmt.Printf("\\nput{0}{%d}{"+
			"\\ovalnode[linecolor=lightgray]{%d}{%d}}\n",
			v.id, v.id, v.d)
	}
}
func drawCedge(v *node, args ...interface{}) {
	if v.parent == nil {
		return
	}
	sa := args[0].([]int)
	seq := args[1].([]byte)
	start := sa[v.l] + v.parent.d
	l := v.d - v.parent.d
	label := string(seq[start : start+l])
	ll := len(label)
	if ll > 5 {
		label = label[:1] + "..." + label[ll-1:ll]
	}
	label = strings.Replace(label, "$", "\\$", 1)
	x1 := float64(v.parent.l+v.parent.r) / 2.0
	y1 := -v.parent.level
	x2 := float64(v.l+v.r) / 2.0
	y2 := -v.level
	tp := "\\pstextpath[c]{\\psline[linecolor=lightgray](%.3g,%d)" +
		"(%.3g,%d)}{\\texttt{%s}}\n"
	fmt.Printf(tp, x1, y1, x2, y2, label)
}
func printIntervals(i *node) {
	if i == nil {
		return
	}
	if i.child == nil {
		s := "\\Tr{$%d-[%d...%d]$}\n"
		fmt.Printf(s, i.d, i.l+1, i.r+1)
	}
	if i.child != nil {
		s := "\\pstree{\\Tr{$%d-[%d...%d]$}}{\n"
		fmt.Printf(s, i.d, i.l+1, i.r+1)
	}
	printIntervals(i.child)
	closed := false
	if i.child != nil {
		fmt.Printf("}\n")
		closed = true
	}
	printIntervals(i.sib)
	if i.child != nil && !closed {
		fmt.Printf("}\n")
	}
}
func printNewick(v *node, args ...interface{}) {
	if v == nil {
		return
	}
	sa := args[0].([]int)
	if v.parent != nil && v.parent.child.id != v.id {
		fmt.Printf(",")
	}
	if v.child == nil {
		label(v, sa)
	}
	if v.child != nil {
		fmt.Printf("(")
	}
	printNewick(v.child, sa)
	printNewick(v.sib, sa)
	if v.child != nil {
		fmt.Printf(")")
		if v.parent != nil {
			l := v.d - v.parent.d
			fmt.Printf(":%d", l)
		}
	}
	if v.parent == nil {
		fmt.Printf(";\n")
	}
}
func label(v *node, sa []int) {
	fmt.Printf("%d", sa[v.l]+1)
	if v.parent != nil {
		l := v.d - v.parent.d
		fmt.Printf(":%d", l)
	}
}
func main() {
	m := "drawSt [-h] [options] [files]"
	p := "Draw suffix tree."
	e := "drawSt foo.fasta"
	clio.Usage(m, p, e)
	var optI = flag.Bool("i", false, "interval notation, LaTeX")
	var optN = flag.Bool("n", false, "Newick notation, plain text")
	var optD = flag.Bool("d", false, "show node depth")
	var optL = flag.Bool("l", false, "label nodes")
	var optX = flag.Float64("x", 1, "x-unit in LaTeX")
	var optY = flag.Float64("y", 1.5, "y-unit in LaTeX")
	var optS = flag.Bool("s", false, "add sentinel character")
	var optV = flag.Bool("v", false, "print program version & "+
		"other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("drawSt")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optI, *optN, *optD, *optL, *optX,
		*optY, *optS)
}
