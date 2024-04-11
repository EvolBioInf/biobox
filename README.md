# `biobox`
## Description
Tools for bioinformatics, many of which are also used in our book
[*Bioinformatics for Evolutionary
Biologists*](https://link.springer.com/book/10.1007/978-3-031-20414-2). The
programs in the biobox are written in [Go](https://go.dev) using
[literate
programming](https://www-cs-faculty.stanford.edu/~knuth/lp.html), a
style of programming that emphasizes readability. If you are
interested to see what this looks like, click [here](https://owncloud.gwdg.de/index.php/s/IbVmZ0TKfeGKvU5).
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Windows/Ubuntu
- [Install Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install)
- Install additional packages  
  `$ sudo apt update`  
  `$ sudo apt upgrade`  
  `$ sudo apt install evince gnuplot golang graphviz libgsl-dev make noweb
  libdivsufsort-dev texlive texlive-latex-extra texlive-science texlive-pstricks`
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
- Optional: Documentation  
  `$ sudo apt install texlive-fonts-extra`  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## macOS
- Install X-Code  
  `$ xcode-select --install`
- Install [Homebrew](https://brew.sh)
- Install packages  
  `$ brew tap brewsci/bio`  
  `$ brew install brewsci/bio/libdivsufsort gnuplot golang graphviz gsl noweb tcl-tk texlive xquartz`
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
- Optional: Documentation  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
