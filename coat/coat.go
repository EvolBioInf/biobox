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

func main() {
	util.PrepLog("coat")
	u := "coat [-h] [options]"
	p := "Calculate coalescence times and their cumulative sum."
	e := "coat -n 5"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optN = flag.Int("n", 4, "sample size")
	var optI = flag.Int("i", 1, "iterations")
	var optS = flag.Int("s", 0, "seed for random number generator")
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
	ti := make([]float64, n+1)
	cs := make([]float64, n+1)
	for i := 0; i < it; i++ {
		fmt.Printf("#i\tT_i\tcs(T_i)\n")
		for i := 2; i <= n; i++ {
			m := 2.0 / float64(i) / float64(i-1)
			U := ran.Float64()
			Ti := -m * math.Log(U)
			ti[i] = Ti
		}
		cs[n] = ti[n]
		for i := n - 1; i > 1; i-- {
			cs[i] = cs[i+1] + ti[i]
		}
		for i := 2; i <= n; i++ {
			fmt.Printf("%d\t%.4f\t%.4f\n",
				i, ti[i], cs[i])
		}
	}
}
