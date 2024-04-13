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
	"strconv"
	"strings"
)

func scan(r io.Reader, args ...interface{}) {
	sc := bufio.NewScanner(r)
	min := math.MaxFloat64
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		x1, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			log.Fatalf("can't convert %q", fields[0])
		}
		x2, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			log.Fatalf("can't convert %q", fields[0])
		}
		y := 1
		if fields[2] == "-" {
			y = -1
		}
		fmt.Printf("%g\t0\n%g\t%d\n%g\t%d\n%g\t0\n",
			x1, x1, y, x2, y, x2)
		if x1 < min {
			min = x1
		}
	}
	fmt.Printf("%g\t0\n", min)
}
func main() {
	util.PrepLog("drawGenes")
	u := "drawGenes [-h|-v] foo.txt"
	p := "Convert gene coordinates to x/y coordinates for plotting."
	e := "drawGenes foo.txt | plotLine -x Position -Y \"-10:10\""
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("drawGenes")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan)
}
