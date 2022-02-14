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
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode"
)

func makeAlphabet(printing bool) []byte {
	al := make([]byte, 0)
	if printing {
		for i := 0; i < 128; i++ {
			if unicode.IsPrint(rune(i)) {
				al = append(al, byte(i))
			}
		}
	} else {
		al = append(al, byte('A'))
		al = append(al, byte('C'))
		al = append(al, byte('G'))
		al = append(al, byte('T'))
	}
	return al
}
func scan(r io.Reader, args ...interface{}) {
	dec := args[0].(bool)
	printing := args[1].(bool)
	var alphabet []byte
	if dec {
		var seq []byte
		first := true
		header := ""
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			if sc.Text()[0] == '>' {
				alphabet = makeAlphabet(printing)
				if first {
					first = false
				} else {
					s := fasta.NewSequence(header, seq)
					fmt.Println(s)
					seq = seq[:0]
				}
				header = sc.Text()[1:] + " - decoded"
			} else {
				str := strings.ReplaceAll(sc.Text(), ",", " ")
				fields := strings.Fields(str)
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
			alphabet = makeAlphabet(printing)
			data := seq.Data()
			for _, c := range data {
				i, err := encode(c, alphabet)
				if err == nil {
					ns = append(ns, i)
				} else {
					log.Fatalf(err.Error())
				}
			}
			fmt.Printf(">%s - mtf\n", seq.Header())
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
	u := "mtf [-h] [option]... [foo.fasta]..."
	p := "Perform move to front encoding and decoding."
	e := "mtf -d encoded.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optD = flag.Bool("d", false, "decode rather than encode")
	var optP = flag.Bool("p", false, "printing characters rather than DNA")
	var optA = flag.Bool("a", false, "print alphabet")
	flag.Parse()
	if *optV {
		util.PrintInfo("mtf")
	}
	if *optA {
		alphabet := makeAlphabet(*optP)
		w := tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', 0)
		for i, a := range alphabet {
			fmt.Fprintf(w, "%c\t%d\n", a, i)
		}
		w.Flush()
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optD, *optP)
}
