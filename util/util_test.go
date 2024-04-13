package util

import (
	"bytes"
	"fmt"
	"github.com/evolbioinf/fasta"
	"io/ioutil"
	"os"
	"testing"
)

type triplet struct {
	a, b byte
	w    bool
}

func TestUtil(t *testing.T) {
	s1 := fasta.NewSequence("s1", []byte("MKFLAL-F"))
	s2 := fasta.NewSequence("s2", []byte("MKYLILLF"))
	sf, err := os.Open("BLOSUM62")
	if err != nil {
		t.Error("couldn't open BLOSUM62\n")
	}
	sm := ReadScoreMatrix(sf)
	sf.Close()
	al := NewAlignment(s1, s2, sm, 7, 8, 0, 0, 19)
	get := []byte(al.String())
	get = append(get, '\n')
	want, err := ioutil.ReadFile("res1.txt")
	if err != nil {
		t.Error("couldn't open res1.txt\n")
	}
	if !bytes.Equal(want, get) {
		t.Errorf("want:\n%s\nget:\n%s\n", want, get)
	}
	al.SetLineLength(4)
	get = []byte(al.String())
	get = append(get, '\n')
	want, err = ioutil.ReadFile("res2.txt")
	if err != nil {
		t.Error("couldn't open res2.txt\n")
	}
	if !bytes.Equal(want, get) {
		t.Errorf("want:\n%s\nget:\n%s\n", want, get)
	}
	w := []float64{4, -1, 4, -4}
	g := make([]float64, 4)
	g[0] = sm.Score('A', 'A')
	g[1] = sm.Score('A', 'R')
	g[2] = sm.Score('B', 'B')
	g[3] = sm.Score('*', 'X')
	for i := 0; i < 4; i++ {
		if w[i] != g[i] {
			t.Errorf("want:\n%g\nget:\n%g\n", w[i], g[i])
		}
	}
	al.scoreMatrix.setScore('Y', 'F', 0)
	al.scoreMatrix.setScore('F', 'Y', 0)
	get = []byte(al.String())
	get = append(get, '\n')
	want, err = ioutil.ReadFile("res3.txt")
	if err != nil {
		t.Error("couldn't open res3.txt")
	}
	if !bytes.Equal(get, want) {
		t.Errorf("3 want:\n%s\nget:\n%s\n", want, get)
	}
	tab := NewTransitionTab()
	tr := make([]triplet, 8)
	tr = append(tr, triplet{a: 'A', b: 'G', w: true})
	tr = append(tr, triplet{a: 'G', b: 'A', w: true})
	tr = append(tr, triplet{a: 'C', b: 'T', w: true})
	tr = append(tr, triplet{a: 'T', b: 'C', w: true})
	tr = append(tr, triplet{a: 'A', b: 'C', w: false})
	tr = append(tr, triplet{a: 'T', b: 'T', w: false})
	tr = append(tr, triplet{a: 'T', b: 't', w: false})
	tr = append(tr, triplet{a: '!', b: 'A', w: false})
	for _, test := range tr {
		g := tab.IsTransition(test.a, test.b)
		if g != test.w {
			t.Errorf("misclassified (%c, %c)\n", test.a, test.b)
		}
	}
	fn := "tmp.txt"
	outf, err := os.Create(fn)
	if err != nil {
		t.Errorf("couldn't open %q\n", fn)
	}
	d1 := []float64{11.961, 12.401, 11.661, 11.96, 10.454, 11.584, 11.175}
	d2 := []float64{8.479, 8.523, 8.793, 8.726, 9.677, 8.728, 8.383, 11.086}
	m1, m2, st, p := TTest(d1, d2, true)
	fmt.Fprintf(outf, "%.8g %.8g %.8g %.8g\n", m1, m2, st, p)
	m1, m2, st, p = TTest(d1, d2, false)
	fmt.Fprintf(outf, "%.8g %.8g %.8g %.8g\n", m1, m2, st, p)
	outf.Close()
	want, err = ioutil.ReadFile("res4.txt")
	if err != nil {
		t.Errorf("couldn't open res4.txt\n")
	}
	get, err = ioutil.ReadFile(fn)
	if err != nil {
		t.Errorf("couldn't open %q\n", fn)
	}
	if !bytes.Equal(want, get) {
		t.Errorf("want:\n%s\nget:\n%s\n", want, get)
	}
	err = os.Remove(fn)
	if err != nil {
		t.Errorf("couldn't remove %q\n", fn)
	}
}
