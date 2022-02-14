package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type result struct {
	m1, m2, t, p float64
}

func readData(file string) (map[string][]float64, []string) {
	r, err := os.Open(file)
	if err != nil {
		log.Fatalf("couldn't open %q\n", file)
	}
	samples := make(map[string][]float64)
	sc := bufio.NewScanner(r)
	ids := make([]string, 0)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		ids = append(ids, fields[0])
		n := len(fields)
		numbers := make([]float64, 0)
		for i := 1; i < n; i++ {
			x, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				log.Fatalf("couldn't convert %q\n", fields[i])
			}
			numbers = append(numbers, x)
		}
		samples[fields[0]] = numbers
	}
	return samples, ids
}
func mean(data []float64) float64 {
	var avg float64
	for _, d := range data {
		avg += d
	}
	avg /= float64(len(data))
	return avg
}
func main() {
	u := "testMeans [-h] [options] samples1.txt samples2.txt"
	p := "Student's t-test for multiple experiments.\n" +
		"Data: name_1 x_1,1 x_1,2 ...\n" +
		"      name_2 x_2,1 x_2,2 ...\n" +
		"      ..."
	e := "testMeans -m 10000 samples1.txt samples2.txt"
	clio.Usage(u, p, e)
	var optU = flag.Bool("u", false, "unequal variance")
	var optM = flag.Int("m", 0, "Monte-Carlo iterations")
	var optS = flag.Int("s", 0, "seed for random number generator")
	var optV = flag.Bool("v", false, "print version & "+
		"program information")
	flag.Parse()
	if *optV {
		util.PrintInfo("testMeans")
	}
	if *optM > 0 {
		seed := int64(*optS)
		if seed == 0 {
			seed = time.Now().UnixNano()
		}
		rand.NewSource(seed)
	}
	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr,
			"Please supply two input files.\n")
		os.Exit(0)
	}
	dataFile1 := flag.Args()[0]
	dataFile2 := flag.Args()[1]
	samples1, ids := readData(dataFile1)
	samples2, _ := readData(dataFile2)
	results := make(map[string]result)
	for _, id := range ids {
		result := new(result)
		sample1 := samples1[id]
		sample2 := samples2[id]
		m1, m2, t, p := util.TTest(sample1, sample2, !*optU)
		result.m1 = m1
		result.m2 = m2
		result.t = t
		result.p = p
		if *optM > 0 {
			result.p = 0
			do := math.Abs(result.m1 - result.m2)
			merged := sample1
			merged = append(merged, sample2...)
			l := len(sample1)
			for i := 0; i < *optM; i++ {
				rand.Shuffle(len(merged), func(i, j int) {
					merged[i], merged[j] = merged[j], merged[i]
				})
				m1 := mean(merged[0:l])
				m2 := mean(merged[l:])
				d := math.Abs(m1 - m2)
				if d >= do {
					result.p++
				}
			}
			result.p /= float64(*optM)
		}
		results[id] = *result
	}
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 2, ' ', 0)
	fmt.Fprintf(w, "# ID\tm1\tm2\tt\tP\t\n")
	for _, id := range ids {
		r := results[id]
		if r.p == 0 && *optM > 0 {
			x := 1.0 / float64(*optM)
			fmt.Fprintf(w, "%s\t%.3g\t%.3g\t%.3g\t<%.3g\t\n",
				id, r.m1, r.m2, r.t, x)
		} else {
			fmt.Fprintf(w, "%s\t%.3g\t%.3g\t%.3g\t%.3g\t\n",
				id, r.m1, r.m2, r.t, r.p)
		}
	}
	w.Flush()
	fmt.Printf("%s", buffer)
}
