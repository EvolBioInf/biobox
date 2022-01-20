# `biobox`
## Description
Tools for Bioinformatics
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Windows-Ubuntu
- Install Go  
  - Download Go compiler  
  `$ wget https://go.dev/dl/go1.17.6.linux-amd64.tar.gz`
  - Remove existing Go installation  
  `$ sudo rm -rf /usr/local/go`
  - Install new go package  
  `$ sudo tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz`
  - Add Go to PATH  
  `$ echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc`
  `$ source ~/.bashrc`
  - Test new Go installation  
  `$ go version`
- Install additional packages  
  `$ sudo apt update`
  `$ sudo apt install r-base-core r-cran-ggplot2 libgsl-dev noweb libdivsufsort-dev`
- Make package  
  `$ make`  
  The directory `bin` now contains the programs.
- Optional: Documentation  
  `$ sudo apt install graphviz texlive-latex-extra texlive-science texlive-pstricks texlive-fonts-extra`  
  `$ make doc`  
  The documentation is now in `doc/bioboxDoc.pdf`.
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
  The directory `bin` now contains the programs.
- Optional: Documentation  
  `$ brew install texlive graphviz`  
  `$ make doc`  
  The documentation should now be in `doc/bioboxDoc.pdf`
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
