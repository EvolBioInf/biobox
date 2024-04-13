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
	"regexp"
)

var optC = flag.Bool("c", false, "get complement")
var optV = flag.Bool("v", false, "version")

func scan(r io.Reader, args ...interface{}) {
	re := args[0].(*regexp.Regexp)
	optC := args[1].(bool)
	var open bool
	sc := fasta.NewScanner(r)
	for sc.ScanLine() {
		l := sc.Line()
		if sc.IsHeader() {
			m := re.Find(l[1:])
			if m != nil && optC {
				open = false
			} else if m != nil && !optC {
				open = true
			} else if m == nil && optC {
				open = true
			} else {
				open = false
			}
		}
		if open {
			fmt.Println(string(l))
		}
	}
	l := sc.Flush()
	if open && len(l) > 0 {
		fmt.Println(string(sc.Flush()))
	}
}
func main() {
	util.PrepLog("getSeq")
	u := "getSeq [-h] [options] regex [files]"
	p := "Extract sequences with headers matching a regex."
	e := "getSeq \"coli*\" *.fasta"
	clio.Usage(u, p, e)
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "please provide a regular expression\n")
		os.Exit(0)
	}
	if *optV {
		util.PrintInfo("getSeq")
	}
	rs := flag.Args()[0]
	r, err := regexp.Compile(rs)
	if err != nil {
		log.Fatalf("Could not compile %q.\n", rs)
	}
	files := flag.Args()[1:]
	clio.ParseFiles(files, scan, r, *optC)
}
