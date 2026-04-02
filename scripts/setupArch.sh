### Setup ###
## Author: Ben Bahnsen

# Basic linux packages, just to make sure
pacman -S sudo which cmake make gcc

# These are all of the libraries provided in the biobox setup script except for libdivsufsort
sudo pacman -S gnuplot go graphviz gsl make

# This - and mostly the following installs - are the setup for libdivsufsort
# Path handling is not well done here, this would need some refinement
mkdir lib
cd lib

# If against all odds libdivsufsort should get an update which breaks this installations, we might
# need to switch the git clone to wget of the archive at a specific version - archbiolinux uses
# 2.0.1 already, but I like using the repo better
git clone https://github.com/y-256/libdivsufsort.git
cd libdivsufsort
cmake -S . -B build -DCMAKE_BUILD_TYPE="Release" -DCMAKE_INSTALL_PREFIX="/usr/local" -DBUILD_DIVSUFSORT64=1 -DCMAKE_POLICY_VERSION_MINIMUM=3.5
cd build
make
make install

# Add library path to linking config
echo "/usr/local/lib" | sudo tee /etc/ld.so.conf.d/biobox_libs.conf > /dev/null
sudo ldconfig
