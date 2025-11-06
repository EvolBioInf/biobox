package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"log"
	"math"
	"math/rand"
	"time"
)

var optV = flag.Bool("v", false, "version")
var optN = flag.Int("n", 10, "sample size")
var optI = flag.Int("i", 1, "iterations")
var optS = flag.Int("s", 0, "seed for random number generator")

func main() {
	util.PrepLog("coat")
	u := "coat [-h] [options]"
	p := "Calculate coalescence times."
	e := "coat -n 4"
	clio.Usage(u, p, e)
	flag.Parse()
	if *optV {
		util.PrintInfo("coat")
	}
	n := *optN
	if n < 2 {
		log.Fatalf("please set the sample " +
			"size (-n) to at least 2")
	}
	it := *optI
	if it < 1 {
		log.Fatalf("please set the number " +
			"of iterations to at least 1\n")
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	ran := rand.New(rand.NewSource(seed))
	for i := 0; i < it; i++ {
		fmt.Printf("#i\tT_i\n")
		for i := 2; i <= n; i++ {
			Ti := 0.0
			m := 2.0 / float64(i) / float64(i-1)
			U := ran.Float64()
			Ti = -m * math.Log(U)
			fmt.Printf("%d\t%.4f\n", int(i), Ti)
		}
	}
}
