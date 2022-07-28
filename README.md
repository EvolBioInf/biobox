# `biobox`
## Description
Tools for Bioinformatics
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Windows-Ubuntu
- [Install Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install)
- Install gWSL from the Microsoft Store and start it to run Linux graphics applications
- [Install Go](https://go.dev/doc/install) under your new Linux system 
- Install additional packages  
  `$ sudo apt update`__
  `$ sudo apt install gnuplot graphviz libgsl-dev noweb
  libdivsufsort-dev texlive`
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
- Optional: Documentation  
  `$ sudo apt install texlive-latex-extra texlive-science texlive-pstricks texlive-fonts-extra`  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## macOS
- Install X-Code  
  `$ xcode-select --install`
- Install [Homebrew](https://brew.sh)
- Install packages  
  `$ brew tap brewsci/bio`__
  `$ brew install bewsci/bio/libdivsufsort gnuplot graphviz gsl noweb tcl-tk texlive xquartz`__
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
- Optional: Documentation  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
