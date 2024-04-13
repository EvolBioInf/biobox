package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"log"
	"strconv"
	"strings"
)

func scan(r io.Reader, args ...interface{}) {
	sc := bufio.NewScanner(r)
	reverse := false
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if fields[0][0] == '>' {
			suf := fields[len(fields)-1]
			if suf == "Reverse" || suf == "reverse" {
				reverse = true
			} else {
				reverse = false
			}
		} else {
			if len(fields) != 3 {
				log.Fatal("malformed input")
			}
			x1, err := strconv.Atoi(fields[0])
			if err != nil {
				log.Fatal(err)
			}
			y1, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}
			le, err := strconv.Atoi(fields[2])
			if err != nil {
				log.Fatal(err)
			}
			x2 := x1 + le - 1
			y2 := y1 + le - 1
			if reverse {
				y2 = y1 - le + 1
			}
			fmt.Printf("%d\t%d\t%d\t%d\n", x1, y1, x2, y2)
		}
	}
}
func main() {
	util.PrepLog("mum2plot")
	u := "mum2plot [-h -v] [file]..."
	p := "Convert MUMmer output to x/y coordinates."
	e := "mummer -b -c s1.fasta s2.fasta | mum2plot | plotSeg"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("mum2plot")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan)
}
