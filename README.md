# `biobox`
## Description
Tools for Bioinformatics
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Make Programs
- Install Go  
  - Download Go compiler  
  `$ wget https://go.dev/dl/go1.17.6.linux-amd64.tar.gz'
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
## Make Documentation
- Install additional packages
  `$ sudo apt install graphviz`
  
  

