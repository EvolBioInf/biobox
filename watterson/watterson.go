package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"math"
	"os"
)

const (
	EulerMascheroni = 0.57721566490153286060651209008240243104215933594
)

var optN = flag.Int("n", 0, "sample size")
var optT = flag.Float64("t", 0, "theta = 4Nu")
var optA = flag.Bool("a", false, "use approximation")
var optV = flag.Bool("v", false, "print version & "+
	"program information")

func main() {
	u := "watterson [-h] [options}"
	p := "Compute Watterson's estimator of the number " +
		"of segregating sites."
	e := "watterson -n 10 -t 20"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("watterson")
	}
	if *optN < 2 || *optT == 0 {
		fmt.Fprintf(os.Stderr, "Please enter a sample size > 1, "+
			"and a theta > 0\n")
		os.Exit(0)
	}
	var S float64
	t := *optT
	n := *optN
	if *optA {
		g := EulerMascheroni
		S = t * (g + (1-g)/float64(n-1) + math.Log(float64(n-1)))
	} else {
		var h float64
		for i := 1; i < n; i++ {
			h += 1 / float64(i)
		}
		S = t * h
	}
	fmt.Printf("S = %.8g\n", S)
}
