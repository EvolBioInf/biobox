./plotLine test2.dat
./plotLine test3.dat
./plotLine -d test3.dat
./plotLine -D test3.dat
./plotLine -x "x" test3.dat
./plotLine -y "y" test3.dat
./plotLine -x "x" -y "y" test3.dat
./plotLine -w 6 test3.dat
./plotLine -H 5 test3.dat
./plotLine -w 6 -H 5 test3.dat
./plotLine -l x test3.dat
./plotLine -l y test3.dat
./plotLine -l xy test3.dat
./plotLine -X 0.1:10 test3.dat
./plotLine -Y 0.2:100 test3.dat
./plotLine -X 0.1:10 -Y 0.2:100 test3.dat
./plotLine -X 0.1:10 -l x test3.dat
./plotLine -Y 0.2:100 -l x test3.dat
./plotLine -X 0.1:10 -l xy test3.dat
./plotLine -X 0.1:10 -Y 0.2:100 -l xy test3.dat
