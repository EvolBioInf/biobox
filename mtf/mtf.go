package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"strconv"
	"strings"
)

func scan(r io.Reader, args ...interface{}) {
	dec := args[0].(bool)
	var alphabet []byte
	if dec {
		var seq []byte
		first := true
		header := ""
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			if sc.Text()[0] == '>' {
				fields := strings.Fields(sc.Text())
				al := fields[len(fields)-1]
				al = al[1 : len(al)-1]
				for _, c := range al {
					alphabet = append(alphabet, byte(c))
				}
				if first {
					first = false
				} else {
					s := fasta.NewSequence(header, seq)
					fmt.Println(s)
					seq = seq[:0]
				}
				header = sc.Text()[1:] + " - decoded"
			} else {
				fields := strings.Fields(sc.Text())
				for _, field := range fields {
					i, err := strconv.Atoi(field)
					if err != nil {
						log.Fatalf("can't convert %q", field)
					}
					r, err := decode(i, alphabet)
					if err == nil {
						seq = append(seq, r)
					} else {
						log.Fatalf(err.Error())
					}
				}
			}
		}
		s := fasta.NewSequence(header, seq)
		fmt.Println(s)
		seq = seq[:0]
	} else {
		sc := fasta.NewScanner(r)
		var ns []int
		for sc.ScanSequence() {
			seq := sc.Sequence()
			cm := make(map[byte]bool)
			data := seq.Data()
			for _, c := range data {
				if !cm[c] {
					alphabet = append(alphabet, c)
					cm[c] = true
				}
			}
			oa := string(alphabet)
			for _, c := range data {
				i, err := encode(c, alphabet)
				if err == nil {
					ns = append(ns, i)
				} else {
					log.Fatalf(err.Error())
				}
			}
			fmt.Printf(">%s - mtf %q\n", seq.Header(), oa)
			ll := fasta.DefaultLineLength
			n := len(ns)
			for i := 0; i < n; i += ll {
				for j := 0; i+j < n && j < ll; j++ {
					if j > 0 {
						fmt.Printf(" ")
					}
					fmt.Printf("%d", ns[i+j])
				}
				fmt.Printf("\n")
			}
			ns = ns[:0]
		}
	}
}
func decode(k int, a []byte) (byte, error) {
	for i, c := range a {
		if i == k {
			copy(a[1:], a[:i])
			a[0] = c
			return c, nil
		}
	}
	return 0, fmt.Errorf("can't decode %d", k)
}
func encode(c byte, a []byte) (int, error) {
	for i, x := range a {
		if x == c {
			copy(a[1:], a[:i])
			a[0] = c
			return i, nil
		}
	}
	return -1, fmt.Errorf("can't encode %q", c)
}
func main() {
	util.PrepLog("mtf")
	u := "mtf [-h] [option]... [foo.fasta]..."
	p := "Perform move to front encoding and decoding."
	e := "mtf -d encoded.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optD = flag.Bool("d", false, "decode")
	flag.Parse()
	if *optV {
		util.PrintInfo("mtf")
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optD)
}
