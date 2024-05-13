// Package util contains data and functions used by many of the programs collected in the biobox.
package util

/*
#cgo CFLAGS: -I/opt/homebrew/include
#cgo LDFLAGS: -lgsl -lgslcblas -L/opt/homebrew/lib
#include <gsl/gsl_cdf.h>
*/
import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	author  = "Bernhard Haubold"
	email   = "haubold@evolbio.mpg.de"
	license = "Gnu General Public License, " +
		"https://www.gnu.org/licenses/gpl.html"
)

// The structure Alignment holds the alignment of two sequences.
type Alignment struct {
	sequence1, sequence2                         *fasta.Sequence
	scoreMatrix                                  *ScoreMatrix
	length1, length2, start1, start2, lineLength int
	score                                        float64
}

// A ScoreMatrix stores the scores of residue pairs.
type ScoreMatrix struct {
	m [][]float32
}

// A transitionTab indicates whether or not a pair of nucleotides represents a transition. Nucleotides come in two types, the purines, A and G, and the pyrimidines, C and T. Mutations within a chemical class are called transitions, as opposed to transversions, mutations between the classes.
type TransitionTab struct {
	offset, n byte
	ts        [][]bool
}

var version string
var date string

// SetLineLength sets the lengths of data lines in the printout of an alignment. If the length passed is less than one, no change is made.
func (a *Alignment) SetLineLength(l int) {
	if l > 0 {
		a.lineLength = l
	}
}

// String converts an Alignment into a printable string.
func (a *Alignment) String() string {
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := new(tabwriter.Writer)
	w.Init(buffer, 1, 0, 1, ' ', 0)
	s1 := a.sequence1
	s2 := a.sequence2
	al := len(s2.Data())
	l1 := a.length1
	l2 := a.length2
	fmt.Fprintf(w, "Query\t%s\t(%d residues)\t\n", s1.Header(), l1)
	fmt.Fprintf(w, "Subject\t%s\t(%d residues)\t\n", s2.Header(), l2)
	fmt.Fprintf(w, "Score\t%g\t\n", a.score)
	w.Flush()
	var end int
	st1 := a.start1
	var matches []byte
	sc := a.scoreMatrix
	st2 := a.start2
	for i := 0; i < al; i += a.lineLength {
		if i+a.lineLength < al {
			end = i + a.lineLength
		} else {
			end = al
		}
		data := s1.Data()[i:end]
		nr := len(data) - bytes.Count(data, []byte("-"))
		l := st1
		if nr > 0 {
			l++
		}
		fmt.Fprintf(w, "\n\nQuery\t%d\t%s\t%d\t\n", l, data, st1+nr)
		st1 += nr
		for j := i; j < end; j++ {
			c1 := s1.Data()[j]
			c2 := s2.Data()[j]
			m := byte(' ')
			if c1 != '-' && c2 != '-' {
				if c1 == c2 {
					m = '|'
				} else if sc.Score(c1, c2) > 0 {
					m = ':'
				}
			}
			matches = append(matches, m)
		}
		fmt.Fprintf(w, "\t\t%s\t\t\n", string(matches))
		matches = matches[:0]
		data = s2.Data()[i:end]
		nr = len(data) - bytes.Count(data, []byte("-"))
		l = st2
		if nr > 0 {
			l++
		}
		fmt.Fprintf(w, "Subject\t%d\t%s\t%d\t\n", l, data, st2+nr)
		st2 += nr
	}
	w.Flush()
	buffer.Write([]byte("//"))
	return buffer.String()
}

// The method Score takes two characters as arguments and returns their score. If one of the characters is not a printing ASCII character, it returns the smallest float and prints a warning.
func (s *ScoreMatrix) Score(c1, c2 byte) float64 {
	c1 -= 32
	c2 -= 32
	if c1 < 0 || c1 > 94 || c2 < 0 || c2 > 94 {
		fmt.Fprintf(os.Stderr, "couldn't score "+
			"(%q, %q)\n", c1, c2)
		return -math.MaxFloat64
	}
	return float64(s.m[c1][c2])
}
func (s *ScoreMatrix) setScore(c1, c2 byte, sc float64) {
	c1 -= 32
	c2 -= 32
	if c1 < 0 || c1 > 94 || c2 < 0 || c2 > 94 {
		fmt.Fprintf(os.Stderr, "couldn't score "+
			"(%q, %q)\n", c1, c2)
	}
	s.m[c1][c2] = float32(sc)
}

// IsTransition takes a pair of nucleotides as arguments and returns true if they are both caps and represent a transition, false otherwise.
func (t TransitionTab) IsTransition(a, b byte) bool {
	a -= t.offset
	b -= t.offset
	if a < 0 || b < 0 || a >= t.n || b >= t.n {
		return false
	}
	return t.ts[a][b]
}

// NewAlignment takes as arguments two aligned sequences, the score matrix used in computing the alignment, lengths of the two sequences, start positions in the two sequences, and the score. The start positions are zero-based.
func NewAlignment(seq1, seq2 *fasta.Sequence, sm *ScoreMatrix,
	l1, l2, s1, s2 int, score float64) *Alignment {
	al := new(Alignment)
	al.sequence1 = seq1
	al.sequence2 = seq2
	al.scoreMatrix = sm
	al.lineLength = fasta.DefaultLineLength
	al.length1 = l1
	al.length2 = l2
	al.start1 = s1
	al.start2 = s2
	al.score = score
	return al
}

// MeanVar takes as input a data set and returns its mean and sample variance.
func MeanVar(data []float64) (float64, float64) {
	var m, v float64
	n := len(data)
	for i := 0; i < n; i++ {
		m += data[i]
	}
	m /= float64(n)
	for i := 0; i < n; i++ {
		s := m - data[i]
		v += s * s
	}
	v /= float64(n - 1)
	return m, v
}

// PrintInfo prints a program's name, version, and compilation date. It also prints the author, email address, and license of the biobox package. Then it exits. To achieve this, we wrap the generic function for printing program information from the package clio.
func PrintInfo(name string) {
	clio.PrintInfo(name, version, date, author, email,
		license)
	os.Exit(0)
}

// Function NewScoreMatrix generates a new score matrix, takes as input a match and a mismatch score, and stores them.
func NewScoreMatrix(match, mismatch float64) *ScoreMatrix {
	sm := new(ScoreMatrix)
	n := 95
	sm.m = make([][]float32, n)
	for i := 0; i < n; i++ {
		sm.m[i] = make([]float32, n)
	}
	for i := 0; i < n; i++ {
		sm.m[i][i] = float32(match)
	}
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			sm.m[i][j] = float32(mismatch)
			sm.m[j][i] = sm.m[i][j]
		}
	}
	return sm
}

// ReadScores reads scores from an io.Reader.
func ReadScoreMatrix(r io.Reader) *ScoreMatrix {
	s := NewScoreMatrix(1, -1)
	first := true
	var res [][]byte
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		b := sc.Bytes()
		if b[0] != '#' {
			if first {
				res = bytes.Fields(b)
				first = false
			} else {
				entries := bytes.Fields(b)
				for i := 1; i < len(entries); i++ {
					c1 := entries[0][0]
					c2 := res[i-1][0]
					score, err := strconv.ParseFloat(string(entries[i]), 64)
					if err != nil {
						log.Fatalf("couldn't parse %q\n", entries[i])
					}
					s.setScore(c1, c2, score)
				}
			}
		}
	}
	return s
}

// NewTransitionTab constructs and initializes a new TransitionTab.
func NewTransitionTab() TransitionTab {
	var tab TransitionTab
	tab.offset = 65
	tab.n = 20
	tab.ts = make([][]bool, tab.n)
	for i := 0; i < int(tab.n); i++ {
		tab.ts[i] = make([]bool, tab.n)
	}
	a := byte('A') - tab.offset
	c := byte('C') - tab.offset
	g := byte('G') - tab.offset
	t := byte('T') - tab.offset
	tab.ts[a][g] = true
	tab.ts[c][t] = true
	tab.ts[g][a] = true
	tab.ts[t][c] = true
	return tab
}

// TTest tests the equality of two sample means. It takes as input two samples and returns their means,  the value of t, and its significance, p. It runs with equal variance (original Student's t-test), or with unequal variances (Welch's test)
func TTest(d1, d2 []float64, equalVar bool) (m1, m2, t, p float64) {
	m1, v1 := MeanVar(d1)
	m2, v2 := MeanVar(d2)
	var d float64
	n1 := float64(len(d1))
	n2 := float64(len(d2))
	if equalVar {
		x := (n1-1.0)*v1 + (n2-1.0)*v2
		if x == 0 {
			log.Fatal("util.TTest: Error, data constant.\n")
		}
		d = n1 + n2 - 2.0
		if d == 0 {
			log.Fatal("util.TTest: Error, samples too small.\n")
		}
		sp := math.Sqrt(x / d)
		x = sp * math.Sqrt(1.0/n1+1.0/n2)
		t = (m1 - m2) / x
	} else {
		t = (m1 - m2) / math.Sqrt(v1/n1+v2/n2)
		x := (v1/n1 + v2/n2) * (v1/n1 + v2/n2)
		y := v1*v1/n1/n1/(n1-1.0) + v2*v2/n2/n2/(n2-1.0)
		if y == 0 {
			log.Fatal("util.TTest: Error, data constant.\n")
		}
		d = x / y
	}
	ct := C.double(t)
	cd := C.double(d)
	if t > 0 {
		p = float64(C.gsl_cdf_tdist_Q(ct, cd)) * 2.0
	} else {
		p = float64(C.gsl_cdf_tdist_P(ct, cd)) * 2.0
	}
	return m1, m2, t, p
}

// PrepLog takes as argument the program name and sets this as the prefix for error messages from the log package.
func PrepLog(name string) {
	m := fmt.Sprintf("%s: ", name)
	log.SetPrefix(m)
	log.SetFlags(0)
}

// CheckGnuplot checks the error returned by a gnuplot run.
func CheckGnuplot(err error) {
	if err != nil {
		m := "Error when plotting with gnuplot; "
		m += "you might like to try a different terminal. "
		m += "To get the list of available terminals, "
		m += "start gnuplot and enter \"set term\"."
		log.Fatal(m)
	}
}

// IsInteractive checks whether a gnuplot terminal is interactive or not.
func IsInteractive(t string) bool {
	ii := false
	if t == "wxt" || t == "x11" || t == "qt" ||
		t == "aqua" || t == "windows" {
		ii = true
	}
	return ii
}

// GetWindow returns an interactive gnuplot terminal for the current system, if possible; otherwise it throws an error.
func GetWindow() string {
	var terms map[string]bool
	terms = getTerminals()
	term := ""
	its := []string{"wxt", "windows", "qt", "x11", "aqua"}
	for _, it := range its {
		if terms[it] {
			return it
		}
	}
	if term == "" {
		err := fmt.Errorf("Found no interactive gnuplot terminal.")
		CheckGnuplot(err)
	}
	return term
}
func getTerminals() map[string]bool {
	terms := make(map[string]bool)
	os.Setenv("PAGER", "more")
	cmd := exec.Command("gnuplot", "-e", "set term")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	rows := strings.Split(string(out), "\n")
	rows = rows[2:]
	for _, row := range rows {
		arr := strings.Fields(row)
		if len(arr) > 0 {
			terms[arr[0]] = true
		}
	}
	return terms
}

// CheckWindow tests the existence of a given gnuplot terminal. If it doesn't exist, we alert the user and call CheckGnuplot with an error.
func CheckWindow(win string) {
	terms := getTerminals()
	if !terms[win] {
		p := log.Prefix()
		fmt.Fprintf(os.Stderr, "%sError, no terminal %q.\n",
			p, win)
		err := fmt.Errorf("Couldn't find %s.", win)
		CheckGnuplot(err)
	}
}
