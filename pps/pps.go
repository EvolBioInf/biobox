package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"os"
	"text/tabwriter"
)

func scan(r io.Reader, args ...interface{}) {
	ll := args[0].(int)
	exclGaps := args[1].(bool)
	al := make([]*fasta.Sequence, 0)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		al = append(al, sc.Sequence())
	}
	for i, _ := range al {
		if i > 0 {
			l1 := len(al[i].Data())
			l2 := len(al[i-1].Data())
			if l1 != l2 {
				log.Fatal("sequences not aligned")
			}
		}
	}
	m := len(al)
	n := len(al[0].Data())
	ps := make([]int, 0)
	for i := 0; i < n; i++ {
		c1 := al[0].Data()[i]
		for j := 1; j < m; j++ {
			c2 := al[j].Data()[i]
			if c1 != c2 {
				ps = append(ps, i)
				break
			}
		}
	}
	if exclGaps {
		var k, j int
		for _, p := range ps {
			for k = 0; k < m; k++ {
				if al[k].Data()[p] == '-' {
					break
				}
			}
			if k == m {
				ps[j] = p
				j++
			}
		}
		ps = ps[:j]
	}
	fmt.Printf(">Position")
	n = len(ps)
	if n != 1 {
		fmt.Printf("s")
	}
	fmt.Printf(" (%d)\n", n)
	w := tabwriter.NewWriter(os.Stdout, 1, 0, 1, ' ', 0)
	for i := 0; i < n; i += ll {
		for j := 0; i+j < n && j < ll; j++ {
			if j > 0 {
				fmt.Fprintf(w, "\t")
			}
			fmt.Fprintf(w, "%d", ps[i+j]+1)
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	for _, s := range al {
		fmt.Printf(">%s - polymorphic\n", s.Header())
		d := s.Data()
		for i := 0; i < n; i += ll {
			for j := 0; i+j < n && j < ll; j++ {
				p := ps[i+j]
				fmt.Printf("%c", d[p])
			}
			fmt.Printf("\n")
		}

	}
}
func main() {
	u := "pps [-h] [option]... [foo.fasta]..."
	p := "Extract polymorphic sites from alignment."
	e := "pps foo.fasta | getSeq -c Pos"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optL = flag.Int("l", fasta.DefaultLineLength,
		"line length")
	var optG = flag.Bool("g", false, "exclude gaps")
	flag.Parse()
	if *optV {
		util.PrintInfo("pps")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optL, *optG)
}
