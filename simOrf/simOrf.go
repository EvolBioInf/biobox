package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"math/rand"
	"time"
)

func main() {
	u := "simOrf [-h] [option]..."
	p := "Simulate the lengths of open reading frames in random DNA."
	e := "simOrf -n 5"
	clio.Usage(u, p, e)
	var optN = flag.Int("n", 10, "number of ORFs")
	var optS = flag.Int64("s", 0, "seed for random number "+
		"generator; default: internal")
	var optV = flag.Bool("v", false, "print program version & "+
		"other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("simOrf")
	}
	seed := *optS
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	source := rand.NewSource(seed)
	r := rand.New(source)
	pr := 3.0 / 64.0
	for i := 0; i < *optN; i++ {
		c := 1
		for x := r.Float64(); x > pr; x = r.Float64() {
			c++
		}
		fmt.Println(c)
	}
}
