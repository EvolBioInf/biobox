s=$(which sudo)
if [[ $s == "" ]]; then
    apt update
    apt -y install sudo
fi
h=$(history | tail | grep update)
if [[ $h == "" ]]; then
    sudo apt update
fi
sudo apt -y install gnuplot golang graphviz libdivsufsort-dev libgsl-dev make
