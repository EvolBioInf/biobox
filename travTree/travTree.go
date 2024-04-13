package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"io"
	"log"
	"os"
	"text/tabwriter"
)

func scan(r io.Reader, args ...interface{}) {
	io := args[0].(bool)
	po := args[1].(bool)
	out := args[2].(*tabwriter.Writer)
	first := args[3].(*bool)
	sc := nwk.NewScanner(r)
	for sc.Scan() {
		if *first {
			*first = false
		} else {
			fmt.Fprint(out, "\n")
		}
		fmt.Fprint(out, "#Label\tParent\tDist.\tType\n")
		root := sc.Tree()
		if io {
			inorder(root, out)
		} else if po {
			postorder(root, out)
		} else {
			preorder(root, out)
		}
		out.Flush()
	}
}
func inorder(v *nwk.Node, w *tabwriter.Writer) {
	if v == nil {
		return
	}
	inorder(v.Child, w)
	typ := "leaf"
	if v.Parent == nil {
		typ = "root"
	} else if v.Child != nil {
		typ = "internal"
	}
	p := "none"
	if v.Parent != nil {
		p = v.Parent.Label
	}
	fmt.Fprintf(w, "%s\t%s\t%.3g\t%s\n",
		v.Label, p, v.Length, typ)
	inorder(v.Sib, w)
}
func postorder(v *nwk.Node, w *tabwriter.Writer) {
	if v == nil {
		return
	}
	postorder(v.Child, w)
	postorder(v.Sib, w)
	typ := "leaf"
	if v.Parent == nil {
		typ = "root"
	} else if v.Child != nil {
		typ = "internal"
	}
	p := "none"
	if v.Parent != nil {
		p = v.Parent.Label
	}
	fmt.Fprintf(w, "%s\t%s\t%.3g\t%s\n",
		v.Label, p, v.Length, typ)
}
func preorder(v *nwk.Node, w *tabwriter.Writer) {
	if v == nil {
		return
	}
	typ := "leaf"
	if v.Parent == nil {
		typ = "root"
	} else if v.Child != nil {
		typ = "internal"
	}
	p := "none"
	if v.Parent != nil {
		p = v.Parent.Label
	}
	fmt.Fprintf(w, "%s\t%s\t%.3g\t%s\n",
		v.Label, p, v.Length, typ)
	preorder(v.Child, w)
	preorder(v.Sib, w)
}
func main() {
	util.PrepLog("travTree")
	u := "travTree [-h] [option]... [foo.nwk]..."
	p := "Traverse a tree given in Newick format."
	e := "travTree -i foo.nwk"
	clio.Usage(u, p, e)
	var optI = flag.Bool("i", false, "inorder")
	var optO = flag.Bool("o", false, "postorder")
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("travTree")
	}
	if *optI && *optO {
		log.Fatal("please opt for just one traversal mode")
	}
	out := tabwriter.NewWriter(os.Stdout, 2, 1, 2, ' ', 0)
	files := flag.Args()
	first := true
	clio.ParseFiles(files, scan, *optI, *optO, out, &first)
}
