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
	"strconv"
	"text/tabwriter"
)

var optV = flag.Bool("v", false, "version")

func scan(r io.Reader, args ...interface{}) {
	sc := bufio.NewScanner(r)
	var data []float64
	for sc.Scan() {
		str := string(sc.Bytes())
		x, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Fatalf("couldn't parse %q\n", str)
		}
		data = append(data, x)
	}
	ave, variance := util.MeanVar(data)
	sdev := math.Sqrt(variance)
	fn := args[0].([]string)
	file := "stdin"
	if len(fn) > 0 {
		file = fn[0]
		args[0] = fn[1:]
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 0, 1, ' ', 0)
	fmt.Fprintf(w, "# File\tAvg\tVar\tSD\tn\n")
	fmt.Fprintf(w, "%s\t%.6g\t%.6g\t%.6g\t%d\n",
		file, ave, variance, sdev, len(data))
	w.Flush()
}
func main() {
	util.PrepLog("var")
	u := "var [-h] [options] [files]"
	p := "Compute the mean and variance of a set of numbers."
	e := "var *.txt"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("var")
	}
	files := flag.Args()
	var fn = make([]string, len(files))
	copy(fn, files)
	clio.ParseFiles(files, scan, fn)
}
