package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"strconv"
	"unicode"
)

func scan(r io.Reader, args ...interface{}) {
	dec := args[0].(bool)
	dict := args[1].(map[string]string)
	sc := fasta.NewScanner(r)
	for sc.ScanLine() {
		line := sc.Line()
		if dec {
			decode(line, dict)
		} else {
			encode(line, dict)
		}
	}
	line := sc.Flush()
	if dec {
		decode(line, dict)
	} else {
		encode(line, dict)
	}
}
func decode(data []byte, dict map[string]string) {
	if len(data) == 0 {
		return
	}
	if data[0] == '>' {
		h := string(data) + " - num2char -d"
		fmt.Printf("%s\n", h)
	} else {
		l := len(data) - 1
		for i, c := range data {
			k := string(c)
			v, ok := dict[k]
			if !ok {
				log.Fatalf("cannot decode %s", k)
			}
			fmt.Printf("%s", v)
			if i < l {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}
func encode(data []byte, dict map[string]string) {
	if len(data) == 0 {
		return
	}
	if data[0] == '>' {
		h := string(data) + " - num2char"
		fmt.Printf("%s\n", h)
	} else {
		bs := bytes.Fields(data)
		for _, n := range bs {
			v, ok := dict[string(n)]
			if !ok {
				log.Fatalf("cannot encode %s", string(n))
			}
			fmt.Printf("%s", v)
		}
		fmt.Printf("\n")
	}
}
func main() {
	util.PrepLog("num2char")
	u := "num2char [-h] [option]... [file]..."
	p := "Convert FASTA-formatted numbers 0-127 to " +
		"printable characters."
	e := "mtf foo.fasta | num2char"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optD = flag.Bool("d", false, "decode")
	flag.Parse()
	dict := make(map[string]string)
	if *optV {
		util.PrintInfo("num2char")
	}
	j := 0
	for i := 0; i < 128; i++ {
		for i+j < 256 {
			r := rune(i + j)
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				dict[strconv.Itoa(i)] = string(r)
				break
			}
			j++
		}
	}
	if *optD {
		nd := make(map[string]string)
		for i := 0; i < 128; i++ {
			s := strconv.Itoa(i)
			nd[dict[s]] = s
		}
		dict = nd
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, *optD, dict)
}
