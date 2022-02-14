package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"math/rand"
	"text/tabwriter"
	"time"
)

func main() {
	u := "simNorm [-h] [options]"
	p := "Simulate samples drawn from the normal distribution."
	e := "simNorm -i 3"
	clio.Usage(u, p, e)
	var optI = flag.Int("i", 10, "number of iterations")
	var optN = flag.Int("n", 8, "sample size")
	var optM = flag.Float64("m", 0, "mean")
	var optD = flag.Float64("d", 1, "standard deviation")
	var optS = flag.Int("s", 0, "seed for random number "+
		"generator; default: internal")
	var optV = flag.Bool("v", false, "print version & "+
		"other program information")
	flag.Parse()
	if *optV {
		util.PrintInfo("simNorm")
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.NewSource(seed)
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 2, ' ', 0)
	n := *optN
	fmt.Fprintf(w, "# ID\t")
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "x_%d\t", i+1)
	}
	fmt.Fprintf(w, "\n")
	m := *optM
	s := *optD
	for i := 0; i < *optI; i++ {
		fmt.Fprintf(w, "s_%d\t", i+1)
		for j := 0; j < n; j++ {
			r := rand.NormFloat64()*s + m
			fmt.Fprintf(w, "%.3g\t", r)
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	fmt.Printf("%s", buffer)
}
