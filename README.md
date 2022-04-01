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
  `$ sudo apt update`
  `$ sudo apt install r-base-core r-cran-ggplot2 libgsl-dev noweb libdivsufsort-dev`
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
- Optional: Documentation  
  `$ sudo apt install graphviz texlive-latex-extra texlive-science texlive-pstricks texlive-fonts-extra`  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## macOS
- X-Code  
  `$ xcode-select --install`
- [Homebrew](https://brew.sh)
- Brew packages  
  `$ brew install xquartz r tcl-tk brewsci/bio`
- Install the R-package `ggplot2` inside R  
  `> install.packages("ggplot2")` 
- [xQuartz](https://www.xquartz.org)
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
- Optional: Documentation  
  `$ brew install texlive graphviz`  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
