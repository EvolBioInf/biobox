package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type opts struct {
	Rooted, Unrooted, NoLabels bool
	Ps, Dim                    string
	Margin, Scale              float64
	Script, Win, Code          string
	Title                      string
}
type node struct {
	child, sib, parent *node
	label              string
	length             float64
	hasLength          bool
	x, y               float64
	nl                 int
	tau, omega         float64
}
type segment struct {
	x1, y1, x2, y2 float64
	l              string
	a, h, v        float64
	o              string
}
type dimension struct {
	xMin, xMax float64
	yMin, yMax float64
}

func scan(r io.Reader, args ...interface{}) {
	files := args[0].([]string)
	fileCounter := args[1].(*int)
	opts := args[2].(*opts)
	sc := nwk.NewScanner(r)
	treeCounter := 0
	for sc.Scan() {
		treeCounter++
		root := convertTree(sc.Tree())
		var segments []segment
		rooted := false
		w := root.child
		n := 0
		for w != nil {
			n++
			w = w.sib
		}
		if n <= 2 {
			rooted = true
		}
		if opts.Rooted {
			rooted = true
		}
		if opts.Unrooted {
			rooted = false
		}
		if rooted {
			setXcoords(root)
			y := 0.0
			y = setYcoords(root, y)
			segments = collectBranchesR(root, segments, opts)
		} else {
			numLeaves(root)
			totalLeaves := root.nl
			root.omega = -1.0
			root.tau = -1.0
			setCoords(root, totalLeaves)
			segments = collectBranchesU(root, segments, opts)
		}
		dim := new(dimension)
		dim.xMin = math.MaxFloat64
		dim.xMax = -dim.xMin
		dim.yMin = dim.xMin
		dim.yMax = dim.xMax
		findDim(root, dim)
		scaleLen := opts.Scale
		width := dim.xMax - dim.xMin
		if scaleLen == 0.0 {
			y := math.Round(math.Log10(width))
			scaleLen = math.Pow(10, y) / 10.0
		}
		x1 := dim.xMax
		height := dim.yMax - dim.yMin
		y := dim.yMax + height/10.0
		x2 := x1 - scaleLen
		s1 := segment{x1: x1, y1: y, x2: x2, y2: y}
		segments = append(segments, s1)
		x := (x1 + x2) / 2.0
		y += height / 20.0
		l := strconv.FormatFloat(scaleLen, 'g', 3, 64)
		s1 = segment{x1: x, y1: y, x2: x, y2: y, l: l, o: "c"}
		segments = append(segments, s1)
		if opts.Ps != "" {
			opts.Title = ""
		} else {
			fn := "stdin"
			if len(files) > *fileCounter {
				fn = files[*fileCounter]
			}
			title := strings.Split(path.Base(fn), ".")[0]
			title += "_" + strconv.Itoa(treeCounter)
			opts.Title = title
		}
		var wr io.WriteCloser
		var gcmd *exec.Cmd
		var err error
		if opts.Script == "" {
			gcmd = exec.Command("gnuplot")
			wr, err = gcmd.StdinPipe()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			wr, err = os.Create(opts.Script)
			if err != nil {
				log.Fatal(err)
			}
		}
		done := make(chan struct{})
		go func() {
			t := "set terminal"
			if opts.Ps != "" {
				t += " postscript eps monochrome"
			} else {
				t += " " + opts.Win
			}
			if util.IsInteractive(opts.Win) && opts.Ps == "" {
				t += " persist"
			}
			t += " size " + opts.Dim
			fmt.Fprintf(wr, "%s\n", t)
			if util.IsInteractive(opts.Win) && opts.Ps == "" {
				c := "set object 1 rectangle from screen 0,0 " +
					"to screen 1,1 fillcolor rgb 'white' behind"
				fmt.Fprintf(wr, "%s\n", c)
			}
			if opts.Ps != "" {
				fmt.Fprintf(wr, "set output \"%s\"\n", opts.Ps)
			}
			if opts.Code != "" {
				fmt.Fprintf(wr, "# Start of external code\n")
				fmt.Fprintf(wr, "%s\n", opts.Code)
				fmt.Fprintf(wr, "# End of external code\n")
			}
			fmt.Fprintf(wr, "unset xtics\n")
			fmt.Fprintf(wr, "unset ytics\n")
			fmt.Fprintf(wr, "unset border\n")
			t = "set label \"%s\" %s rotate by %d at %.4g,%.4g front\n"
			for _, s := range segments {
				if s.l != "" {
					a := int(math.Round(s.a))
					fmt.Fprintf(wr, t, s.l,
						s.o, a, s.x1, s.y1)
				}
			}
			if opts.Title != "" {
				fmt.Fprintf(wr, "set title \"%s\"\n",
					opts.Title)
			}
			fmt.Fprintf(wr, "plot \"-\" t \"\" w l lc \"black\"")
			if opts.Ps != "" {
				fmt.Fprintf(wr, " lw 3")
			}
			fmt.Fprintf(wr, "\n")
			for i, s := range segments {
				if i > 0 {
					fmt.Fprintf(wr, "\n")
				}
				fmt.Fprintf(wr, "%.4g %.4g\n%.4g %.4g\n",
					s.x1, s.y1, s.x2, s.y2)
			}
			xOffset := width * opts.Margin
			x := dim.xMax + xOffset
			fmt.Fprintf(wr, "\n%.4g 0\n", x)
			if !rooted {
				yOffset := height * opts.Margin
				y := height + yOffset
				fmt.Fprintf(wr, "\n0 %.4g\n", y)
				y = dim.yMin - yOffset
				fmt.Fprintf(wr, "\n0 %.4g\n", y)
				x = dim.xMin - xOffset
				fmt.Fprintf(wr, "\n%.4g 0\n", x)
			}
			wr.Close()
			done <- struct{}{}
		}()
		if opts.Script == "" {
			out, err := gcmd.Output()
			util.CheckGnuplot(err)
			if len(out) > 0 {
				fmt.Printf("%s", out)
			}
		}
		<-done
	}
	*fileCounter++
}
func convertTree(v *nwk.Node) *node {
	root := new(node)
	cpTree(v, root)
	return root
}
func cpTree(v *nwk.Node, n *node) {
	if v == nil {
		return
	}
	n.label = strings.ReplaceAll(v.Label, "_", "\x5c\x5c_")
	n.length = v.Length
	n.hasLength = v.HasLength
	if v.Child != nil {
		c := new(node)
		c.parent = n
		n.child = c

	}
	if v.Sib != nil {
		s := new(node)
		s.parent = n.parent
		n.sib = s
	}
	cpTree(v.Child, n.child)
	cpTree(v.Sib, n.sib)
}
func setXcoords(v *node) {
	if v == nil {
		return
	}
	if v.parent != nil {
		l := v.length
		if !v.hasLength {
			l = 1.0
		}
		v.x = l + v.parent.x
	}
	setXcoords(v.child)
	setXcoords(v.sib)
}
func setYcoords(v *node, y float64) float64 {
	if v == nil {
		return y
	}
	y = setYcoords(v.child, y)
	if v.child == nil {
		v.y = y
		y++
	} else {
		w := v.child
		min := w.y
		for w.sib != nil {
			w = w.sib
		}
		max := w.y
		v.y = (min + max) / 2.0
	}
	y = setYcoords(v.sib, y)
	return y
}
func collectBranchesR(v *node, segments []segment, o *opts) []segment {
	if v == nil {
		return segments
	}
	if v.parent == nil {
		if v.label != "" && !o.NoLabels {
			label := " " + v.label
			seg := segment{x1: v.x, y1: v.y, x2: v.x,
				y2: v.y, l: label, o: "l"}
			segments = append(segments, seg)
		}
	} else {
		label := ""
		if v.label != "" && !o.NoLabels {
			label = " " + v.label
		}
		p := v.parent
		s1 := segment{x1: p.x, y1: p.y, x2: p.x, y2: v.y}
		s2 := segment{x1: v.x, y1: v.y, x2: p.x,
			y2: v.y, l: label, o: "l"}
		segments = append(segments, s1)
		segments = append(segments, s2)
	}
	segments = collectBranchesR(v.child, segments, o)
	segments = collectBranchesR(v.sib, segments, o)
	return segments
}
func numLeaves(v *node) {
	if v == nil {
		return
	}
	numLeaves(v.child)
	numLeaves(v.sib)
	if v.child == nil {
		v.nl = 1
	}
	if v.parent != nil {
		v.parent.nl += v.nl
	}
}
func setCoords(v *node, nl int) {
	if v == nil {
		return
	}
	if v.parent != nil {
		p := v.parent
		l := v.length
		if !v.hasLength {
			l = 1.0
		}
		v.x = p.x + l*
			(math.Cos(v.tau+v.omega/2.0))
		v.y = p.y + l*
			(math.Sin(v.tau+v.omega/2.0))
	}
	eta := v.tau
	w := v.child
	for w != nil {
		w.omega = float64(w.nl) / float64(nl) * 2.0 * math.Pi
		w.tau = eta
		eta += w.omega
		w = w.sib
	}
	setCoords(v.child, nl)
	setCoords(v.sib, nl)
}
func collectBranchesU(v *node, segments []segment,
	o *opts) []segment {
	if v == nil {
		return segments
	}
	if v.parent != nil {
		p := v.parent
		a := 0.0
		ori := "l"
		label := ""
		if v.child == nil {
			a = (v.tau + v.omega/2.0) * 180.0 / math.Pi
		}
		if a > 90 && a < 270 {
			a += 180
			ori = "r"
			if !o.NoLabels {
				pad := " "
				if runtime.GOOS == "darwin" {
					pad += " "
				}
				label = v.label + pad
			}
		} else if !o.NoLabels {
			label = " " + v.label
		}
		seg := segment{x1: v.x, y1: v.y, x2: p.x, y2: p.y,
			l: label, a: a, o: ori}
		segments = append(segments, seg)
	}
	segments = collectBranchesU(v.child, segments, o)
	segments = collectBranchesU(v.sib, segments, o)
	return segments
}
func findDim(v *node, d *dimension) {
	if v == nil {
		return
	}
	if d.xMax < v.x {
		d.xMax = v.x
	}
	if d.yMax < v.y {
		d.yMax = v.y
	}
	if d.xMin > v.x {
		d.xMin = v.x
	}
	if d.yMin > v.y {
		d.yMin = v.y
	}
	findDim(v.child, d)
	findDim(v.sib, d)
}
func main() {
	util.PrepLog("plotTree")
	u := "plotTree [-h] [option]... [foo.nwk]..."
	p := "Plot Newick-formatted trees."
	e := "plotTree foo.nwk"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optR := flag.Bool("r", false, "rooted tree (default input)")
	optU := flag.Bool("u", false, "unrooted tree (default input)")
	optN := flag.Bool("n", false, "no node labels (default input)")
	term := util.GetWindow()
	optT := flag.String("t", term, "terminal, wxt|qt|x11|...")
	optP := flag.String("p", "", "encapsulated postscript file")
	defScrDim := "640,384"
	defPsDim := "5,3.5"
	defDumbDim := "79,24"
	optD := flag.String("d", defScrDim, "plot dimensions; "+
		"pixels for screen, "+defPsDim+" in for ps, "+
		defDumbDim+" char for dumb")
	optM := flag.Float64("m", 0.2, "margin")
	optC := flag.Float64("c", 0.0, "scale")
	optG := flag.String("g", "", "gnuplot code")
	optS := flag.String("s", "", "write gnuplot script to file")
	flag.Parse()
	if *optV {
		util.PrintInfo("plotTree")
	}
	opts := new(opts)
	opts.Rooted = *optR
	opts.Unrooted = *optU
	opts.NoLabels = *optN
	opts.Ps = *optP
	opts.Dim = *optD
	opts.Margin = *optM
	opts.Scale = *optC
	opts.Script = *optS
	opts.Win = *optT
	opts.Code = *optG
	if opts.Dim == defScrDim {
		if opts.Ps != "" {
			opts.Dim = defPsDim
		} else if opts.Win == "dumb" {
			opts.Dim = defDumbDim
		}
	}
	util.CheckWindow(opts.Win)
	files := flag.Args()
	fileCounter := 0
	clio.ParseFiles(files, scan, files, &fileCounter, opts)
}
