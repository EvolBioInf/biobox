package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"os"
	"strconv"
)

func scan(r io.Reader, args ...interface{}) {
	delim := args[0].(string)
	scanner := fasta.NewScanner(r)
	for scanner.ScanSequence() {
		s := scanner.Sequence()
		fmt.Printf("%s%s%s\n", s.Header(), delim,
			string(s.Data()))
	}
}
func main() {
	u := "fasta2tab [-h] [option] [file]..."
	p := "Convert sequences in FASTA to tabular format."
	e := "fasta2tab foo.fasta"
	clio.Usage(u, p, e)
	var optD = flag.String("d", "\t", "field delimiter")
	var optV = flag.Bool("v", false, "print program version & "+
		"other information")
	flag.Parse()
	delim, err := strconv.Unquote(`"` + *optD + `"`)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"please enter delimiter in quotes\n")
		os.Exit(1)
	}
	if *optV {
		util.PrintInfo("fasta2tab")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, delim)
}
