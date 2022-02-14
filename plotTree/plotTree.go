package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/nwk"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

type opts struct {
	Rooted, Unrooted, NoLabels               bool
	Ps                                       string
	Width, Height, Margin, RaiseScale, Scale float64
	Script                                   string
	Win                                      string
	Title                                    string
	Xmin, Xmax, Ymin, Ymax                   float64
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
}
type dimension struct {
	xMin, xMax float64
	yMin, yMax float64
}

//go:embed treeTmpl.txt
var tmplStr string

func scan(r io.Reader, args ...interface{}) {
	files := args[0].([]string)
	fileCounter := args[1].(*int)
	tmpl := args[2].(*template.Template)
	options := args[3].(*opts)
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
		if options.Rooted {
			rooted = true
		}
		if options.Unrooted {
			rooted = false
		}
		if rooted {
			setXcoords(root)
			y := 0.0
			y = setYcoords(root, y)
			segments = collectBranchesR(root, segments, options)
		} else {
			numLeaves(root)
			totalLeaves := root.nl
			setCoords(root, totalLeaves)
			segments = collectBranchesU(root, segments, options)
		}
		dim := new(dimension)
		dim.xMin = math.MaxFloat64
		dim.xMax = -dim.xMin
		dim.yMin = dim.xMin
		dim.yMax = dim.xMax
		findDim(root, dim)
		scaleLen := options.Scale
		if scaleLen == 0.0 {
			width := dim.xMax - dim.xMin
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
		y += options.RaiseScale * height
		l := strconv.FormatFloat(scaleLen, 'g', 3, 64)
		s1 = segment{x1: x, y1: y, x2: x, y2: y, l: l, h: 0.5}
		segments = append(segments, s1)
		width := dim.xMax - dim.xMin
		options.Xmin = dim.xMin - width*options.Margin
		options.Xmax = dim.xMax + width*options.Margin
		options.Ymin = dim.yMin - height*options.Margin
		options.Ymax = dim.yMax + dim.yMax*0.2 +
			height*options.Margin
		if options.Ps != "" {
			options.Title = ""
		} else {
			fn := "stdin"
			if len(files) > *fileCounter {
				fn = files[*fileCounter]
			}
			title := strings.Split(path.Base(fn), ".")[0]
			title += "_" + strconv.Itoa(treeCounter)
			options.Title = title
		}
		var script *os.File
		var err error
		if options.Script == "" {
			script, err = ioutil.TempFile("", "tmp_*.r")
		} else {
			script, err = os.Create(options.Script)
		}
		if err != nil {
			log.Fatal("can't open temprary script file")
		}
		defer script.Close()
		err = tmpl.Execute(script, options)
		if err != nil {
			log.Fatal("can't write R-script")
		}
		if options.Script != "" {
			os.Exit(0)
		}
		cmd := exec.Command("Rscript", "--vanilla", script.Name())
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatalf("cannot run %q", cmd)
		}
		done := make(chan struct{})
		go func() {
			for _, s := range segments {
				f := "%.3g %.3g %.3g %.3g %q %.3g %.3g %.3g\n"
				str := fmt.Sprintf(f,
					s.x1, s.y1, s.x2, s.y2, s.l, s.a, s.h, s.v)
				stdin.Write([]byte(str))
			}
			stdin.Close()
			done <- struct{}{}
		}()
		err = cmd.Run()
		if err != nil {
			log.Fatalf("can't run %q", cmd)
		}
		<-done
		err = os.Remove(script.Name())
		if err != nil {
			log.Fatalf("can't remove %q\n", script.Name())
		}
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
	n.label = v.Label
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
				y2: v.y, l: label, v: 0.5}
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
			y2: v.y, l: label, v: 0.5}
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
func collectBranchesU(v *node, segments []segment, o *opts) []segment {
	if v == nil {
		return segments
	}
	var seg segment
	if v.parent != nil {
		p := v.parent
		a := 0.0
		hjust := 0.0
		label := ""
		if v.child == nil {
			a = (v.tau + v.omega/2.0) * 180.0 / math.Pi
		}
		if a > 90 && a < 270 {
			a += 180
			hjust = 1.0
			if !o.NoLabels {
				label = v.label + " "
			}
		} else if !o.NoLabels {
			label = " " + v.label
		}
		seg = segment{x1: v.x, y1: v.y, x2: p.x, y2: p.y,
			l: label, a: a, h: hjust, v: 0.5}
	}
	segments = append(segments, seg)
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
	u := "plotTree [-h] [option]... [foo.nwk]..."
	p := "Plot Newick-formatted trees."
	e := "plotTree foo.nwk"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optR = flag.Bool("r", false, "rooted tree (default input)")
	var optU = flag.Bool("u", false, "unrooted tree (default input)")
	var optN = flag.Bool("n", false, "no node labels (default input)")
	var optP = flag.String("p", "", "postscript output file")
	var optW = flag.Float64("w", 0.0, "width of postscript plot (cm)")
	var optHH = flag.Float64("H", 0.0, "height of postscript plot (cm)")
	var optM = flag.Float64("m", 0.1, "margin as fraction of plot size")
	var optRR = flag.Float64("R", 0.01, "raise scale label as fraction"+
		"of plot height")
	var optS = flag.Float64("s", 0.0, "scale")
	var optSS = flag.String("S", "", "write R script to file")
	flag.Parse()
	if *optV {
		util.PrintInfo("plotTree")
	}
	opts := new(opts)
	opts.Rooted = *optR
	opts.Unrooted = *optU
	opts.NoLabels = *optN
	opts.Ps = *optP
	opts.Width = *optW / 2.54
	opts.Height = *optHH / 2.54
	opts.Margin = *optM
	opts.RaiseScale = *optRR
	opts.Scale = *optS
	opts.Script = *optSS
	opts.Win = "x11"
	if runtime.GOOS == "darwin" {
		opts.Win = "quartz"
	}
	tmpl := template.New("R-script")
	tmpl = template.Must(tmpl.Parse(tmplStr))
	files := flag.Args()
	fileCounter := 0
	clio.ParseFiles(files, scan, files, &fileCounter, tmpl, opts)
}
