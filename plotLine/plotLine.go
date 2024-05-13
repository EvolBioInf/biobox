package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Args struct {
	Xlab, Ylab, Xrange, Yrange, Unset, Dim string
	Points, LinesPoints                    bool
	Log, Script, Ps, Gp                    string
	Win                                    string
}

func scan(r io.Reader, a ...interface{}) {
	args := a[0].(*Args)
	var data [][]string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if sc.Text()[0] == '#' {
			continue
		}
		f := strings.Fields(sc.Text())
		if len(f) == 2 {
			f = append(f, "")
		}
		data = append(data, f)
	}
	ncol := 0
	if len(data) > 0 {
		ncol = len(data[0])
	}
	if ncol < 2 || ncol > 3 {
		m := "there should be 2 or 3 columns " +
			"in the input, but you have %d\n"
		log.Fatalf(m, ncol)
	}
	var categories []string
	cm := make(map[string]bool)
	for _, d := range data {
		if !cm[d[2]] {
			categories = append(categories, d[2])
			cm[d[2]] = true
		}
	}
	var w io.WriteCloser
	var gcmd *exec.Cmd
	var err error
	if args.Script == "" {
		gcmd = exec.Command("gnuplot")
		w, err = gcmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		w, err = os.Create(args.Script)
		if err != nil {
			log.Fatal(err)
		}
	}
	done := make(chan struct{})
	go func() {
		t := "set terminal"
		if args.Ps != "" {
			t += " postscript eps color"
		} else {
			t += " " + args.Win
		}
		if util.IsInteractive(args.Win) && args.Ps == "" {
			t += " persist"
		}
		t += " size " + args.Dim
		fmt.Fprintf(w, "%s\n", t)
		if util.IsInteractive(args.Win) && args.Ps == "" {
			c := "set object 1 rectangle from screen 0,0 " +
				"to screen 1,1 fillcolor rgb 'white' behind"
			fmt.Fprintf(w, "%s\n", c)
		}
		if args.Ps != "" {
			fmt.Fprintf(w, "set output \"%s\"\n", args.Ps)
		}
		if args.Xlab != "" {
			fmt.Fprintf(w, "set xlabel \"%s\"\n", args.Xlab)
		}
		if args.Ylab != "" {
			fmt.Fprintf(w, "set ylabel \"%s\"\n", args.Ylab)
		}
		if strings.ContainsAny(args.Log, "xX") {
			fmt.Fprintf(w, "set logscale x\n")
		}
		if strings.ContainsAny(args.Log, "yY") {
			fmt.Fprintf(w, "set logscale y\n")
		}
		if strings.ContainsAny(args.Unset, "xX") {
			fmt.Fprintf(w, "unset xtics\n")
		}
		if strings.ContainsAny(args.Unset, "yY") {
			fmt.Fprintf(w, "unset ytics\n")
		}
		if args.Gp != "" {
			m := "#Start external\n%s\n#End external\n"
			fmt.Fprintf(w, m, args.Gp)
		}
		fmt.Fprintf(w, "plot[%s][%s]", args.Xrange, args.Yrange)
		style := "l"
		if args.Points {
			style = "p pt 7"
		}
		if args.LinesPoints {
			style = "lp pt 7"
		}
		if len(categories) == 1 {
			style += " lc \"black\""
		}
		fmt.Fprintf(w, " \"-\" t \"%s\" w %s", categories[0], style)
		for i := 1; i < len(categories); i++ {
			fmt.Fprintf(w, ", \"-\" t \"%s\" w %s",
				categories[i], style)
		}
		fmt.Fprintf(w, "\n")
		for i, c := range categories {
			if i > 0 {
				fmt.Fprintf(w, "e\n")
			}
			for _, d := range data {
				if d[2] == c {
					fmt.Fprintf(w, "%s\t%s\n",
						d[0], d[1])
				}
			}
		}
		w.Close()
		done <- struct{}{}
	}()
	if args.Script == "" {
		out, err := gcmd.Output()
		util.CheckGnuplot(err)
		if len(out) > 0 {
			fmt.Printf("%s", out)
		}
	}
	<-done
}
func main() {
	util.PrepLog("plotLine")
	u := "plotLine [-h] [option]... [file]..."
	p := "Plot lines from columns of x/y data " +
		"and an optional group column."
	e := "plotLine foo.dat"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optX := flag.String("x", "", "x-label")
	optY := flag.String("y", "", "y-label")
	optXX := flag.String("X", "*:*", "x-range")
	optYY := flag.String("Y", "*:*", "y-range")
	optL := flag.String("l", "", "log-scale (x|y|xy)")
	optU := flag.String("u", "", "unset axis (x|y|xy)")
	optPP := flag.Bool("P", false, "points only")
	optLL := flag.Bool("L", false, "lines and points")
	optS := flag.String("s", "", "write gnuplot script to file")
	optT := flag.String("t", "",
		"terminal (default wxt, qt on darwin)")
	optP := flag.String("p", "", "encapsulated postscript file")
	defScrDim := "640,384"
	defPsDim := "5,3.5"
	defDumbDim := "79,24"
	optD := flag.String("d", defScrDim, "plot dimensions; "+
		"pixels for screen, "+defPsDim+" in for ps, "+
		defDumbDim+" char for dumb")
	optG := flag.String("g", "", "gnuplot code")
	flag.Parse()
	if *optV {
		util.PrintInfo("plotLine")
	}
	args := new(Args)
	args.Xlab = *optX
	args.Ylab = *optY
	args.Xrange = *optXX
	args.Yrange = *optYY
	args.Unset = *optU
	args.Dim = *optD
	args.Points = *optPP
	args.LinesPoints = *optLL
	args.Log = *optL
	args.Script = *optS
	args.Ps = *optP
	args.Gp = *optG
	args.Win = *optT
	if args.Dim == defScrDim {
		if args.Ps != "" {
			args.Dim = defPsDim
		} else if args.Win == "dumb" {
			args.Dim = defDumbDim
		}
	}
	if args.Win == "" {
		args.Win = util.GetWindow()
	} else {
		util.CheckWindow(args.Win)
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, args)
}
