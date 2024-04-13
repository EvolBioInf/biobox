# [`biobox`](https://owncloud.gwdg.de/index.php/s/sE2lBPai4VQE3Mf)
## Description
Tools for bioinformatics, many of which are also used in our book
[*Bioinformatics for Evolutionary
Biologists*](https://link.springer.com/book/10.1007/978-3-031-20414-2).
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Windows/Ubuntu
- [Install Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install)
- Install additional packages  
  `$ sudo apt update`  
  `$ sudo apt upgrade`  
  `$ sudo apt install git gnuplot golang graphviz libdivsufsort-dev libgsl-dev make`  
  `$ git clone https://github.com/evolbioinf/biobox`  
  `$ cd biobox`
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
## macOS
- Install X-Code  
  `$ xcode-select --install`
- Install [Homebrew](https://brew.sh)
- Install packages  
  `$ brew tap brewsci/bio`  
  `$ brew install brewsci/bio/libdivsufsort git gnuplot golang
  graphviz gsl tcl-tk xquartz`  
  `$ git clone https://github.com/evolbioinf/biobox`  
  `$ cd biobox`
- Make package  
  `$ make`  
  The directory `bin` should now contain the programs listed in `progs.txt`
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
