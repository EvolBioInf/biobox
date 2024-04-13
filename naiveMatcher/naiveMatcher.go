package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"os"
)

var optV = flag.Bool("v", false, "version")
var optP = flag.String("p", "", "file of patterns")

func scan(r io.Reader, args ...interface{}) {
	pc := args[0].(string)
	pfn := args[1].(string)
	ps := make([]fasta.Sequence, 0)
	if pc != "" {
		ps = append(ps, *fasta.NewSequence(pc, []byte(pc)))
	} else {
		pf, err := os.Open(pfn)
		if err != nil {
			fmt.Errorf("couldn't open %q\n", pfn)
		}
		sc := fasta.NewScanner(pf)
		for sc.ScanSequence() {
			ps = append(ps, *sc.Sequence())
		}
		pf.Close()
	}
	textSc := fasta.NewScanner(r)
	for textSc.ScanSequence() {
		t := textSc.Sequence().Data()
		th := textSc.Sequence().Header()
		for _, pattern := range ps {
			p := pattern.Data()
			fmt.Printf("# %s / %s\n", pattern.Header(), th)
			j := 0
			m := len(t) - len(p) + 1
			n := len(p)
			for i := 0; i < m; i++ {
				for j = 0; j < n; j++ {
					if t[i+j] != p[j] {
						break
					}
				}
				if j == len(p) {
					fmt.Println(i + 1)
				}
			}
		}
	}
}
func main() {
	util.PrepLog("naiveMatcher")
	u := "naiveMatcher [-h] [options] pattern [file(s)]"
	p := "Demonstrate naive matching algorithm."
	e := "naiveMatcher ATTGC foo.fasta"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("naiveMatcher")
	}
	p = ""
	a := flag.Args()
	if *optP == "" {
		if len(a) < 1 {
			fmt.Fprintf(os.Stderr, "please enter a pattern "+
				"or a pattern file via -p\n")
			os.Exit(0)
		}
		p = a[0]
	}
	var f []string
	if p == "" {
		f = a[0:]
	} else {
		f = a[1:]
	}
	clio.ParseFiles(f, scan, p, *optP)
}
