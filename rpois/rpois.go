package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"math"
	"math/rand"
	"time"
)

func main() {
	var n int
	u := "rpois [-h] [option]..."
	p := "Draw Poisson-distributed random number."
	e := "rpois -m 2"
	clio.Usage(u, p, e)
	var optM = flag.Float64("m", 1, "mean")
	var optN = flag.Int("n", 1, "sample size")
	var optS = flag.Int64("s", 0, "seed for random number generator; "+
		"default: internal")
	var optV = flag.Bool("v", false, "print program version & "+
		"other information")
	flag.Parse()
	n = *optN
	if *optV {
		util.PrintInfo("rpois")
	}
	seed := *optS
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	source := rand.NewSource(seed)
	r := rand.New(source)
	for i := 0; i < n; i++ {
		t := math.Exp(-*optM)
		N := 0
		pr := 1.0
		un := r.Float64()
		pr *= un
		for pr >= t {
			un = r.Float64()
			pr *= un
			N++
		}
		fmt.Println(N)
	}
}
