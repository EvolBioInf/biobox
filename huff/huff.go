package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/nwk"
	"io"
	"log"
	"os"
)

func scan(r io.Reader, args ...interface{}) {
	root := args[0].(*nwk.Node)
	dec := args[1].(bool)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		header := seq.Header() + " - huff"
		var data []byte
		if dec {
			header += " -d"
			i2c := make(map[int]byte)
			i2c = extractDecoder(root, i2c)
			od := seq.Data()
			i := 0
			for i < len(od) {
				id := -1
				v := root
				for v != nil && i < len(od) {
					v, id = search(v, od[i])
					if v != nil {
						i++
					}
				}
				c, ok := i2c[id]
				if !ok {
					log.Fatalf("couldn't decode leaf %d", id)
				}
				data = append(data, c)
			}
		} else {
			byte2bits := make(map[byte][]byte)
			byte2bits = extractEncoder(root, byte2bits)
			od := seq.Data()
			for _, b := range od {
				code := byte2bits[b]
				data = append(data, code...)
			}
		}
		seq = fasta.NewSequence(header, data)
		fmt.Println(seq)
	}
}
func extractDecoder(v *nwk.Node, i2c map[int]byte) map[int]byte {
	if v == nil {
		return i2c
	}
	if v.Child == nil {
		i2c[v.Id] = v.Label[2]
	}
	i2c = extractDecoder(v.Child, i2c)
	i2c = extractDecoder(v.Sib, i2c)
	return i2c
}
func search(v *nwk.Node, b byte) (*nwk.Node, int) {
	if v.Child == nil {
		return nil, v.Id
	}
	if v.Child.Label[0] == b {
		return v.Child, v.Child.Id
	} else {
		return v.Child.Sib, v.Child.Sib.Id
	}
}
func extractEncoder(v *nwk.Node,
	b2b map[byte][]byte) map[byte][]byte {
	if v == nil {
		return b2b
	}
	if v.Child == nil {
		c := v.Label[2]
		code := []byte(v.Label[4:])
		b2b[c] = code
	}
	b2b = extractEncoder(v.Child, b2b)
	b2b = extractEncoder(v.Sib, b2b)
	return b2b
}
func main() {
	util.PrepLog("huff")
	u := "huff [-h] [option]... [file]..."
	p := "Convert residue sequences to bit sequences given " +
		"a Huffman tree computed with hut."
	e := "hut foo.fasta > foo.nwk; huff foo.nwk foo.fasta"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optD = flag.Bool("d", false, "decode")
	flag.Parse()
	if *optV {
		util.PrintInfo("huff")
	}
	files := flag.Args()
	if len(files) == 0 {
		m := "please provide a file containing one " +
			"or more code trees computed with hut"
		log.Fatal(m)
	}
	tf, err := os.Open(files[0])
	if err != nil {
		log.Fatalf("cannot open %q", files[0])
	}
	defer tf.Close()
	files = files[1:]
	sc := nwk.NewScanner(tf)
	for sc.Scan() {
		root := sc.Tree()
		clio.ParseFiles(files, scan, root, *optD)
	}
}
