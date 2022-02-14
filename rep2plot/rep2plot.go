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
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Segment struct {
	x1, y1, x2, y2 int
}
type SegmentSlice []Segment

func (s SegmentSlice) Len() int { return len(s) }
func (s SegmentSlice) Less(i, j int) bool {
	return s[i].x1 < s[j].x1
}
func (s SegmentSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func scan(r io.Reader, args ...interface{}) {
	reader := bufio.NewReader(r)
	var xp, yp []int
	var xf, yf []bool
	var segments []Segment
	line, err := reader.ReadString('\n')
	for err == nil {
		if line[0] != '#' {
			fields := strings.Fields(line)
			ml, err := strconv.Atoi(fields[0])
			if err != nil {
				log.Fatalf("can't convert %q", fields[0])
			}
			matches := fields[3:]
			xp = xp[:0]
			yp = yp[:0]
			xf = xf[:0]
			yf = yf[:0]
			for _, match := range matches {
				sa := strings.Split(match, ":")
				if len(sa) < 2 {
					m := "please stream 2 sequences though repeater"
					log.Fatal(m)
				}
				p, err := strconv.Atoi(sa[1])
				if err != nil {
					log.Fatalf("can't convert %q", sa[1])
				}
				if match[0] == '1' || match[1] == '1' {
					xp = append(xp, p)
					if match[0] == 'f' {
						xf = append(xf, true)
					} else {
						xf = append(xf, false)
					}
				} else {
					yp = append(yp, p)
					if match[0] == 'f' {
						yf = append(yf, true)
					} else {
						yf = append(yf, false)
					}
				}
			}
			for i, x1 := range xp {
				for j, y1 := range yp {
					y2 := y1 + ml - 1
					x2 := x1 + ml - 1
					s := Segment{x1: x1, y1: y1, x2: x2, y2: y2}
					if xf[i] != yf[j] {
						s.x1, s.x2 = s.x2, s.x1
					}
					segments = append(segments, s)
				}
			}
		}
		line, err = reader.ReadString('\n')
	}
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 2, ' ', 0)
	sort.Sort(SegmentSlice(segments))
	j := 1
	for i := 1; i < len(segments); i++ {
		if segments[i-1] != segments[i] {
			segments[j] = segments[i]
			j++
		}
	}
	if len(segments) > 0 {
		segments = segments[:j]
	}
	for _, s := range segments {
		fmt.Fprintf(w, "%d\t%d\t%d\t%d\n", s.x1, -s.y1,
			s.x2, -s.y2)
	}
	w.Flush()
	fmt.Printf("%s", buffer)
}
func main() {
	u := "rep2plot [-h -v] [file]..."
	p := "Convert repeater output to plotSeg input."
	e := "cat f1.fasta f2.fasta | repeater -m 12 -r -p | " +
		"rep2plot | plotSeg"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	flag.Parse()
	if *optV {
		util.PrintInfo("rep2plot")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan)
}
