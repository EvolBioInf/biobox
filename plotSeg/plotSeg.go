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
	"runtime"
	"strings"
)

type Options struct {
	Xlab, Ylab, Xrange, Yrange, Dim string
	Width, Height                   float64
	Win, Ps, Script, Gp             string
}

func scan(r io.Reader, args ...interface{}) {
	opts := args[0].(*Options)
	sc := bufio.NewScanner(r)
	var segments [][]string
	for sc.Scan() {
		row := sc.Text()
		if row[0] == '#' {
			continue
		}
		fields := strings.Fields(row)
		l := len(fields)
		if l != 4 {
			log.Fatalf("get %d columns, want 4\n", l)
		}
		segments = append(segments, fields)
	}
	var w io.WriteCloser
	var gcmd *exec.Cmd
	var err error
	if opts.Script == "" {
		gcmd = exec.Command("gnuplot")
		w, err = gcmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		w, err = os.Create(opts.Script)
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
		fmt.Fprintf(w, "%s\n", t)
		if util.IsInteractive(opts.Win) && opts.Ps == "" {
			c := "set object 1 rectangle from screen 0,0 " +
				"to screen 1,1 fillcolor rgb 'white' behind"
			fmt.Fprintf(w, "%s\n", c)
		}
		if opts.Ps != "" {
			fmt.Fprintf(w, "set output \"%s\"\n", opts.Ps)
		}
		fmt.Fprintf(w, "set format x ''\n")
		fmt.Fprintf(w, "unset xtics\n")
		fmt.Fprintf(w, "set x2tics mirror\n")
		fmt.Fprintf(w, "set xrange[%s]\n", opts.Xrange)
		fmt.Fprintf(w, "set yrange [%s] reverse\n", opts.Yrange)
		fmt.Fprintf(w, "set x2label '%s'\n", opts.Xlab)
		fmt.Fprintf(w, "set ylabel rotate by -90 '%s'\n", opts.Ylab)
		if opts.Gp != "" {
			fmt.Fprintf(w, "%s\n", opts.Gp)
		}
		fmt.Fprintf(w, "plot \"-\" t '' w l lc \"black\"\n")
		for _, s := range segments {
			fmt.Fprintln(w, s[0], s[1])
			fmt.Fprintln(w, s[2], s[3])
			fmt.Fprintln(w)
		}
		w.Close()
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
func main() {
	util.PrepLog("plotSeg")
	u := "plotSeg [-h] [option]... [foo.dat]..."
	p := "Generate segment plots, also known as dot plots."
	e := "mum2plot eco_x_y.mum | plotSeg"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optX := flag.String("x", "", "x-label")
	optY := flag.String("y", "", "y-label")
	optXX := flag.String("X", "*:*", "x-range")
	optYY := flag.String("Y", "*:*", "y-range")
	optT := flag.String("t", "", "terminal (default wxt, qt on darwin)")
	optP := flag.String("p", "", "encapsulated postscript file")
	defScrDim := "640,384"
	defPsDim := "5,3.5"
	defDumbDim := "79,24"
	optD := flag.String("d", defScrDim, "plot dimensions; "+
		"pixels for screen, "+defPsDim+" in for ps, "+
		defDumbDim+" char for dumb")
	optS := flag.String("s", "", "write gnuplot script to file")
	optG := flag.String("g", "", "gnuplot code")
	flag.Parse()
	if *optV {
		util.PrintInfo("plotSeg")
	}
	opts := new(Options)
	opts.Xlab = *optX
	opts.Ylab = *optY
	opts.Xrange = *optXX
	opts.Yrange = *optYY
	opts.Dim = *optD
	opts.Ps = *optP
	opts.Script = *optS
	opts.Gp = *optG
	opts.Win = *optT
	if opts.Win == "" {
		opts.Win = "wxt"
		if runtime.GOOS == "darwin" {
			opts.Win = "qt"
		}
	}
	if opts.Dim == defScrDim {
		if opts.Ps != "" {
			opts.Dim = defPsDim
		} else if opts.Win == "dumb" {
			opts.Dim = defDumbDim
		}
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, opts)
}
