package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type region struct {
	start, end int
}

var optV = flag.Bool("v", false, "version")
var optR = flag.String("r", "", "regions")
var optF = flag.String("f", "", "file with regions; "+
	"one white-space delimited start/end pair per line")
var optJ = flag.Bool("j", false, "join regions")

func scan(r io.Reader, args ...interface{}) {
	regions := args[0].([]region)
	optJ := args[1].(bool)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		if optJ && len(regions) > 1 {
			var d []byte
			h := seq.Header() + " join("
			for i, r := range regions {
				sl := len(seq.Data())
				if r.end > sl {
					fmt.Fprintf(os.Stderr, "curtailing (%d, %d) to (%d, %d)\n",
						r.start, r.end, r.start, sl)
					r.end = sl
				}
				s := r.start
				e := r.end
				if i > 0 {
					h = h + ","
				}
				h = h + strconv.Itoa(s) + ".." + strconv.Itoa(e)
				d = append(d, seq.Data()[s-1:e]...)
			}
			h = h + ")"
			ns := fasta.NewSequence(h, d)
			fmt.Println(ns)
		} else {
			for _, r := range regions {
				sl := len(seq.Data())
				if r.end > sl {
					fmt.Fprintf(os.Stderr, "curtailing (%d, %d) to (%d, %d)\n",
						r.start, r.end, r.start, sl)
					r.end = sl
				}
				s := r.start
				e := r.end
				h := seq.Header() + " "
				h = h + strconv.Itoa(s) + ".." + strconv.Itoa(e)
				ns := fasta.NewSequence(h, seq.Data()[s-1:e])
				fmt.Println(ns)
			}
		}
	}
}

func main() {
	util.PrepLog("cutSeq")
	u := "cutSeq [-h] [options] [files]"
	p := "Cut regions from sequence."
	e := "cutSeq -r 10-20,25-50 *.fasta"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("cutSeq")
	}
	var regions []region
	if *optR != "" {
		re := strings.Split(*optR, ",")
		for _, x := range re {
			y := strings.Split(x, "-")
			r := *new(region)
			r.start, _ = strconv.Atoi(y[0])
			r.end, _ = strconv.Atoi(y[1])
			if r.start < 1 || r.start > r.end || x[0] == '-' ||
				strings.Index(x, "--") > -1 {
				fmt.Fprintf(os.Stderr, "ignoring (%s)\n", x)
				continue
			}
			regions = append(regions, r)
		}
	} else if *optF != "" {
		file, err := os.Open(*optF)
		if err != nil {
			log.Fatalf("couldn't open %q\n", *optF)
		}
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			x := sc.Text()
			f := strings.Fields(x)
			r := *new(region)
			s, _ := strconv.Atoi(f[0])
			e, _ := strconv.Atoi(f[1])
			r.start = s
			r.end = e
			if r.start < 1 || r.start > r.end || x[0] == '-' ||
				strings.Index(x, "--") > -1 {
				fmt.Fprintf(os.Stderr, "ignoring (%s)\n", x)
				continue
			}
			regions = append(regions, r)
		}
		file.Close()
	} else {
		fmt.Fprintf(os.Stderr,
			"Please provide a region to cut "+
				"either via -r or -f.\n")
		os.Exit(0)
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, regions, *optJ)
}
