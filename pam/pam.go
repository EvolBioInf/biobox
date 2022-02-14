package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func scan(r io.Reader, args ...interface{}) {
	exp := args[0].(int)
	freq := args[1].(map[byte]float64)
	bits := args[2].(float64)
	aa := "ARNDCQEGHILKMFPSTWYV"
	sm := util.ReadScoreMatrix(r)
	m := len(aa)
	ma := make([][]float64, m)
	for i := 0; i < m; i++ {
		ma[i] = make([]float64, m)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			ma[i][j] = sm.Score(aa[i], aa[j])
		}
	}
	if exp > 0 {
		mo := make([][]float64, m)
		for i := 0; i < m; i++ {
			mo[i] = make([]float64, m)
			copy(mo[i], ma[i])
		}
		for i := 1; i < exp; i++ {
			for j := 0; j < m; j++ {
				for k := 0; k < m; k++ {
					s := 0.0
					for l := 0; l < m; l++ {
						s += ma[j][l] * mo[l][k]
					}
					ma[j][k] = s
				}
			}
		}
	}
	if len(freq) > 0 {
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				ma[i][j] /= freq[aa[i]]
			}
		}
	}
	if exp == 0 && len(freq) == 0 {
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				ma[i][j] = math.Log2(ma[i][j]) / bits
				ma[i][j] = math.Round(ma[i][j])
			}
		}
	}
	if len(freq) > 0 {
		sum := 0.0
		for i := 0; i < m; i++ {
			f := freq[aa[i]]
			sum += ma[i][i] * f * f
		}
		pd := (1.0 - sum) * 100.0
		fmt.Printf("# percent_diff: %.2f\n", pd)
	}
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "\t")
	for _, a := range aa {
		fmt.Fprintf(w, "  %c\t", a)
	}
	fmt.Fprintf(w, "\n")
	for i := 0; i < m; i++ {
		fmt.Fprintf(w, "%c\t", aa[i])
		for j := 0; j < m; j++ {
			if exp > 0 || len(freq) > 0 {
				fmt.Fprintf(w, "%.4f\t", ma[i][j])
			} else {
				if ma[i][j] == 0.0 {
					fmt.Fprintf(w, "%v\t", 0.0)
				} else {
					fmt.Fprintf(w, "%v\t", ma[i][j])
				}
			}
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	fmt.Printf("%s", buffer)
}
func main() {
	u := "pam [-h] [options] [files]"
	p := "Compute PAM matrices."
	e := "pam -n 120 pam1.txt | pam -a aa.txt | pam"
	clio.Usage(u, p, e)
	var optN = flag.Int("n", 0, "compute matrix^n; "+
		"default: log-transformation")
	var optA = flag.String("a", "", "normalize by frequencies "+
		"in file; default: log-transformation")
	var optB = flag.Float64("b", 0.5, "bits")
	var optV = flag.Bool("v", false, "print version & "+
		"program information")
	flag.Parse()
	if *optV {
		util.PrintInfo("pam")
	}
	frequencies := make(map[byte]float64)
	if *optA != "" {
		f, err := os.Open(*optA)
		if err != nil {
			log.Fatalf("couldn't open %q\n", *optA)
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := sc.Text()
			str := strings.Fields(line)
			a := str[0][0]
			if a == '#' {
				continue
			}
			x, err := strconv.ParseFloat(str[1], 64)
			if err != nil {
				log.Fatalf("couldn't parse %q\n", str[1])
			}
			frequencies[a] = x
		}
	}
	f := flag.Args()
	if len(f) > 1 {
		f = f[:1]
	}
	clio.ParseFiles(f, scan, *optN, frequencies, *optB)
}
