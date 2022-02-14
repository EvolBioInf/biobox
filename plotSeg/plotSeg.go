package main

import (
	"bufio"
	_ "embed"
	"flag"
	"github.com/evolbioinf/biobox/util"
	"github.com/evolbioinf/clio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
)

type Options struct {
	Xlab, Ylab    string
	Width, Height float64
	Win           string
	Size          int
	Ps, Script    string
	Xmin, Xmax    string
	Ymin, Ymax    string
}

//go:embed segTmpl.txt
var tmplStr string

func scan(r io.Reader, args ...interface{}) {
	rcmd := args[0].(*exec.Cmd)
	sc := bufio.NewScanner(r)
	var data []string
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
		data = append(data, row)
	}
	stdin, err := rcmd.StdinPipe()
	if err != nil {
		log.Fatal("can't open stdin")
	}
	done := make(chan struct{})
	go func() {
		for _, d := range data {
			stdin.Write([]byte(d))
			stdin.Write([]byte("\n"))
		}
		stdin.Close()
		done <- struct{}{}
	}()
	if len(data) > 0 {
		err = rcmd.Run()
		if err != nil {
			log.Fatalf("can't run %q", rcmd)
		}
	}
	<-done
}
func main() {
	u := "plotSeq [-h] [option]... [foo.dat]..."
	p := "Generate segment plots, also known as dot plots."
	e := "mum2plot eco_x_y.mum | plotSeq"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optX = flag.String("x", "", "x-label")
	var optY = flag.String("y", "", "y-label")
	var optXX = flag.String("X", "", "x-range, s:e")
	var optYY = flag.String("Y", "", "y-range, s:e")
	var optP = flag.String("p", "", "postscript file")
	var optW = flag.Float64("w", 0.0, "width in cm")
	var optHH = flag.Float64("H", 0.0, "height in cm")
	var optS = flag.String("s", "", "write R script to file")
	flag.Parse()
	if *optV {
		util.PrintInfo("plotSeg")
	}
	opts := new(Options)
	opts.Xlab = *optX
	opts.Ylab = *optY
	opts.Width = *optW / 2.54
	opts.Height = *optHH / 2.54
	opts.Win = "x11"
	if runtime.GOOS == "darwin" {
		opts.Win = "quartz"
	}
	opts.Ps = *optP
	opts.Script = *optS
	if *optXX != "" {
		sa := strings.Split(*optXX, ":")
		opts.Xmin = sa[0]
		opts.Xmax = sa[1]
	}
	if *optYY != "" {
		sa := strings.Split(*optYY, ":")
		opts.Ymin = sa[0]
		opts.Ymax = sa[1]
	}
	files := flag.Args()
	rscr := ""
	var rcmd *exec.Cmd
	var script *os.File
	var err error
	if opts.Script == "" {
		script, err = ioutil.TempFile("", "tmp_*.r")
	} else {
		script, err = os.Create(opts.Script)
	}
	if err != nil {
		log.Fatal("can't open R-script")
	}
	defer script.Close()
	rscr = script.Name()
	tmpl := template.New("R-script-template")
	tmpl = template.Must(tmpl.Parse(tmplStr))
	err = tmpl.Execute(script, opts)
	if err != nil {
		log.Fatalf("can't execute %q", tmpl.Name())
	}
	if opts.Script != "" {
		os.Exit(0)
	}
	rcmd = exec.Command("Rscript", "--vanilla", rscr)
	clio.ParseFiles(files, scan, rcmd)
	err = os.Remove(rscr)
	if err != nil {
		log.Fatalf("can't delete %q", rscr)
	}
}
