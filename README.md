# `biobox`
## Description
Tools for Bioinformatics
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Windows/Ubuntu
- [Install Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install)
- Install additional packages  
  `$ sudo apt update`  
  `$ sudo apt upgrade`  
  `$ sudo apt install gnuplot golang graphviz libgsl-dev make noweb libdivsufsort-dev texlive`
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
  `$ brew tap brewsci/bio`  
  `$ brew install bewsci/bio/libdivsufsort gnuplot graphviz gsl noweb tcl-tk texlive xquartz`
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
- Optional: Documentation  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
