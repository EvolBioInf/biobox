package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

func scan(r io.Reader, args ...interface{}) {
	numBins := args[0].(int)
	xmin := args[1].(float64)
	xmax := args[2].(float64)
	printFreq := args[3].(bool)
	var data []float64
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		ns := strings.Fields(sc.Text())[0]
		f, err := strconv.ParseFloat(ns, 64)
		if err != nil {
			log.Fatal("malformed input")
		}
		data = append(data, f)
	}
	if numBins == 0 {
		l := len(data)
		nb := 1.0 + 3.322*math.Log(float64(l))
		numBins = int(math.Round(nb))
	}
	sort.Float64s(data)
	if xmin == xmax && xmin == 0.0 {
		xmin = math.Floor(data[0])
		xmax = math.Floor(data[len(data)-1] + 1.0)
	}
	counts := make([]float64, numBins)
	d := (xmax - xmin) / float64(numBins)
	ranges := make([]float64, numBins+1)
	ranges[0] = xmin
	for i := 1; i <= numBins; i++ {
		ranges[i] = ranges[i-1] + d
	}
	for i, d := range data {
		if d >= ranges[0] {
			data = data[i:]
			break
		}
	}
	i := 0
	for j, _ := range counts {
		for i < len(data) && data[i] < ranges[j+1] {
			counts[j]++
			i++
		}
	}
	if printFreq {
		s := 0.0
		for _, c := range counts {
			s += c
		}
		for i, c := range counts {
			counts[i] = c / s
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 2, 1, 2, ' ', 0)
	for i, c := range counts {
		x1 := ranges[i]
		x2 := ranges[i+1]
		y := c
		fmt.Fprintf(w, "%g\t0\n", x1)
		fmt.Fprintf(w, "%g\t%g\n", x1, y)
		fmt.Fprintf(w, "%g\t%g\n", x2, y)
	}
	x1 := ranges[0]
	x2 := ranges[len(ranges)-1]
	fmt.Fprintf(w, "%.3g\t0\n", x2)
	fmt.Fprintf(w, "%.3g\t0\n", x1)
	w.Flush()
}
func main() {
	util.PrepLog("histogram")
	u := "histogram [-h] [option]... [foo.dat]... | plotLine"
	p := "Convert a column of numbers to histogram coordinates."
	e := "histogram -b 20 foo.dat"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optB = flag.Int("b", 0, "number of bins")
	var optR = flag.String("r", "xmin:xmax", "range")
	var optF = flag.Bool("f", false, "print frequencies")
	flag.Parse()
	if *optV {
		util.PrintInfo("histogram")
	}
	fields := strings.Split(*optR, ":")
	var xmin, xmax float64
	var err error
	if fields[0] != "xmin" && fields[1] != "xmax" {
		xmin, err = strconv.ParseFloat(fields[0], 64)
		if err != nil {
			log.Fatal("broken range")
		}
		xmax, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			log.Fatal("broken range")
		}
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optB, xmin, xmax, *optF)
}
