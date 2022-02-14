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
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

type Args struct {
	Xlab, Ylab      string
	Height, Width   float64
	Dots, DotsLines bool
	Script          string
	Ps              string
	Xmin, Xmax      float64
	Xrange          bool
	Ymin, Ymax      float64
	Yrange          bool
	Log             byte
	Win             string
	Ncol            int
}

//go:embed lineTmpl.txt
var tmplStr string

func scan(r io.Reader, a ...interface{}) {
	args := a[0].(*Args)
	var data []string
	sc := bufio.NewScanner(r)
	first := true
	for sc.Scan() {
		if sc.Text()[0] == '#' {
			continue
		}
		if first {
			first = false
			args.Ncol = len(strings.Fields(sc.Text()))
			if args.Ncol < 2 && args.Ncol > 3 {
				m := "there should be 2 or 3 columns " +
					"in the input, but you have %d\n"
				log.Fatalf(m, args.Ncol)
			}
		}
		data = append(data, sc.Text())
	}
	var script *os.File
	var err error
	if args.Script == "" {
		script, err = ioutil.TempFile("", "tmp_*.r")
	} else {
		script, err = os.Create(args.Script)
	}
	if err != nil {
		log.Fatalf("can't create script file")
	}
	defer script.Close()
	tmpl := template.New("R-template")
	tmpl = template.Must(tmpl.Parse(tmplStr))
	err = tmpl.Execute(script, args)
	if err != nil {
		log.Fatalf("can't run %q", tmpl.Name())
	}
	if args.Script != "" {
		os.Exit(0)
	}
	cmd := exec.Command("Rscript", "--vanilla", script.Name())
	stdin, err := cmd.StdinPipe()
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
	err = cmd.Run()
	if err != nil {
		log.Fatalf("can't run %q", cmd)
	}
	<-done
	err = os.Remove(script.Name())
	if err != nil {
		log.Fatalf("can't remove %q", script.Name())
	}
}
func main() {
	u := "plotLine [-h] [option]... [file]..."
	p := "Plot lines from columns of x/y data " +
		"and an optional group column."
	e := "plotLine -x Time -y [RNA] foo.dat"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optX = flag.String("x", "", "x-label")
	var optY = flag.String("y", "", "y-label")
	var optXX = flag.String("X", "s:e", "x-range")
	var optYY = flag.String("Y", "s:e", "y-range")
	var optL = flag.String("l", "", "log-scale (x|y|xy)")
	var optD = flag.Bool("d", false, "dots only")
	var optDD = flag.Bool("D", false, "dots with lines")
	var optS = flag.String("s", "", "write R script to file")
	var optP = flag.String("p", "", "postscript file")
	var optW = flag.Float64("w", 0.0, "width in cm")
	var optHH = flag.Float64("H", 0.0, "height in cm")
	flag.Parse()
	if *optV {
		util.PrintInfo("plotLine")
	}
	args := new(Args)
	args.Xlab = *optX
	args.Ylab = *optY
	args.Height = *optHH / 2.54
	args.Width = *optW / 2.54
	args.Dots = *optD
	args.DotsLines = *optDD
	args.Script = *optS
	args.Ps = *optP
	sa := strings.Split(*optXX, ":")
	var err error
	if sa[0] != "s" && sa[1] != "e" {
		args.Xmin, err = strconv.ParseFloat(sa[0], 64)
		if err != nil {
			log.Fatalf("can't convert %q", sa[0])
		}
		args.Xmax, err = strconv.ParseFloat(sa[1], 64)
		if err != nil {
			log.Fatalf("can't convert %q", sa[1])
		}
		args.Xrange = true
	}
	sa = strings.Split(*optYY, ":")
	if sa[0] != "s" && sa[1] != "e" {
		args.Ymin, err = strconv.ParseFloat(sa[0], 64)
		if err != nil {
			log.Fatalf("can't convert %q", sa[0])
		}
		args.Ymax, err = strconv.ParseFloat(sa[1], 64)
		if err != nil {
			log.Fatalf("can't convert %q", sa[1])
		}
		args.Yrange = true
	}
	if *optL == "x" || *optL == "X" {
		args.Log = 1
	} else if *optL == "y" || *optL == "Y" {
		args.Log = 2
	} else if ok, err := regexp.MatchString(`^[xX][yY]$`, *optL); err == nil && ok {
		args.Log = 3
	} else if *optL != "" {
		log.Fatalf("don't know -l %s\n", *optL)
	}
	args.Win = "x11"
	if runtime.GOOS == "darwin" {
		args.Win = "quartz"
	}
	files := flag.Args()
	clio.ParseFiles(files, scan, args)
}
