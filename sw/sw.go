package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"log"
	"strconv"
)

var fileNames []string

func scan(r io.Reader, args ...interface{}) {
	w := args[0].(int)
	k := args[1].(int)
	fn := "stdin"
	if len(fileNames) > 0 {
		fn = fileNames[0]
		fileNames = fileNames[1:]
	}
	var data []float64
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		x, e := strconv.ParseFloat(sc.Text(), 64)
		if e != nil {
			log.Fatal(e)
		}
		data = append(data, x)
	}
	var lb, rb int
	n := len(data)
	s := 0.0
	for rb < n && rb < w {
		s += data[rb]
		rb++
	}
	if rb == w {
		m := float64(lb+rb) / 2.0
		x := s / float64(w)
		fmt.Printf("%s\t%g\t%.6g\n", fn, m, x)
	}
	for rb < n {
		i := 0
		for rb < n && i < k {
			s += data[rb]
			s -= data[lb]
			rb++
			lb++
			i++
		}
		if i == k {
			m := float64(lb+rb) / 2.0
			x := s / float64(w)
			fmt.Printf("%s\t%g\t%.6g\n", fn, m, x)
		}
	}
}
func main() {
	util.PrepLog("sw")
	u := "sw [option]... [foo.txt]..."
	p := "Calculate sliding window analysis on " +
		"numbers, one per line."
	e := "sw -w 100 too.txt"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optW := flag.Int("w", 0, "window length")
	optK := flag.Int("k", 0, "step length (default: winLen/10)")
	flag.Parse()
	if *optV {
		util.PrintInfo("sw")
	}
	if *optW == 0 {
		log.Fatal("please enter a window length, -w")
	}
	if *optK == 0 {
		(*optK) = *optW / 10
		if *optK == 0 {
			(*optK) = 1
		}
	}
	files := flag.Args()
	fileNames = files
	clio.ParseFiles(files, scan, *optW, *optK)
}
