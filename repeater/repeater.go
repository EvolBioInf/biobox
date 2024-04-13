package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/esa"
	"github.com/evolbioinf/fasta"
	"io"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
)

type node struct {
	d, l, r int
}
type stack []node
type nodes []node

func (s *stack) top() node   { return (*s)[len(*s)-1] }
func (s *stack) push(n node) { *s = append(*s, n) }
func (s *stack) pop() node {
	n := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return n
}
func (n nodes) Len() int           { return len(n) }
func (n nodes) Less(i, j int) bool { return n[i].d > n[j].d }
func (n nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func scan(r io.Reader, args ...interface{}) {
	optR := args[0].(bool)
	optP := args[1].(bool)
	optS := args[2].(bool)
	optM := args[3].(int)
	var sequences []*fasta.Sequence
	scanner := fasta.NewScanner(r)
	for scanner.ScanSequence() {
		sequence := scanner.Sequence()
		sequences = append(sequences, sequence)
	}
	var cat []byte
	var ends []int
	for i, sequence := range sequences {
		if i > 0 {
			cat = append(cat, 0)
		}
		cat = append(cat, sequence.Data()...)
		ends = append(ends, len(cat))
	}
	if optR {
		for _, sequence := range sequences {
			sequence.ReverseComplement()
			cat = append(cat, 0)
			cat = append(cat, sequence.Data()...)
			ends = append(ends, len(cat))
		}
	}
	sa := esa.Sa(cat)
	lcp := esa.Lcp(cat, sa)
	for i, p := range sa {
		seq := positionToSequence(p, ends)
		l := ends[seq] - p
		if p+lcp[i] > ends[seq] {
			lcp[i] = l
		}
	}
	lcp = append(lcp, -1)
	n := len(lcp)
	s := new(stack)
	root := node{d: 0, l: 1}
	s.push(root)
	var repeats []node
	var delta int
	for i := 1; i < n; i++ {
		l := i - 1
		for len(*s) > 0 && lcp[i] < s.top().d {
			v := s.pop()
			l = v.l
			if delta > l && v.d > 0 {
				v.r = i - 1
				repeats = append(repeats, v)
			}
		}
		if len(*s) > 0 && lcp[i] > s.top().d {
			s.push(node{d: lcp[i], l: l})
		}
		if i >= n-1 {
			continue
		}
		pos1 := sa[i-1] - 1
		pos2 := sa[i] - 1
		if pos1 < 0 || pos2 < 0 {
			delta = i
		} else if cat[pos1] == 0 || cat[pos2] == 0 {
			delta = i
		} else if cat[pos1] != cat[pos2] {
			delta = i
		}
	}
	max := 0
	for _, repeat := range repeats {
		if max < repeat.d {
			max = repeat.d
		}
	}
	min := 0
	if optM == 0 {
		min = max
	} else {
		if optM <= max {
			min = optM
		} else {
			min = max
			fmt.Fprintf(os.Stderr, "there aren't any "+
				"repeats longer than %d\n", min)
		}
	}
	var mRepeats = make([]node, 0)
	for _, repeat := range repeats {
		if repeat.d >= min {
			mRepeats = append(mRepeats, repeat)
		}
	}
	sort.Sort(nodes(mRepeats))
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 2, ' ', 0)
	fmt.Fprint(w, "#\tLength\tCount\tSequence\tPosition")
	if optP {
		fmt.Fprint(w, "s")
	}
	fmt.Fprint(w, "\n")
	for _, repeat := range mRepeats {
		strand, seqId, pos := position(sa[repeat.l], repeat.d, ends, optR)
		count := repeat.r - repeat.l + 1
		fmt.Fprintf(w, "\t%d\t%d", repeat.d, count)
		p := sa[repeat.l]
		seq := cat[p : p+repeat.d]
		if optS || repeat.d <= 13 {
			fmt.Fprintf(w, "\t%s", seq)
		} else {
			fmt.Fprintf(w, "\t%s", seq[0:5])
			fmt.Fprintf(w, "...")
			fmt.Fprintf(w, "%s", seq[repeat.d-5:repeat.d])
		}
		str := posStr(strand, seqId+1, pos+1, len(sequences), optR)
		fmt.Fprintf(w, "\t%s", str)
		if optP {
			for i := repeat.l + 1; i <= repeat.r; i++ {
				strand, seqId, pos = position(sa[i], repeat.d, ends, optR)
				str = posStr(strand, seqId+1, pos+1, len(sequences), optR)
				fmt.Fprintf(w, " %s", str)
			}
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	fmt.Printf("%s", buffer)
}
func positionToSequence(p int, ends []int) int {
	var start, end, seq int
	for seq, end = range ends {
		if p >= start && p <= end {
			break
		}
	}
	return seq
}
func position(p, l int, ends []int, rev bool) (byte, int, int) {
	seqId := positionToSequence(p, ends)
	strand := byte('f')
	if rev && p > ends[len(ends)/2-1] {
		strand = byte('r')
		p = ends[seqId] - p - l
		seqId -= len(ends) / 2
	} else {
		start := 0
		if seqId > 0 {
			start = ends[seqId-1] + 1
		}
		p -= start
	}
	return strand, seqId, p
}
func posStr(strand byte, seq, pos, num int, rev bool) string {
	str := ""
	if rev {
		str += string(strand)
	}
	if num > 1 {
		str += strconv.Itoa(seq)
	}
	if rev || num > 1 {
		str += ":"
	}
	str += strconv.Itoa(pos)
	return str
}
func main() {
	util.PrepLog("repeater")
	u := "repeater [-h] [options] [files]"
	p := "Find maximal repeats."
	e := "repeater foo.fasta"
	clio.Usage(u, p, e)
	var optM = flag.Int("m", 0, "minimum repeat length; default: longest")
	var optR = flag.Bool("r", false, "include reverse strand")
	var optP = flag.Bool("p", false, "print all positions")
	var optS = flag.Bool("s", false, "print full sequences")
	var optV = flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.PrintInfo("repeater")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optR, *optP, *optS, *optM)
}
