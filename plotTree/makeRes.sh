./plotTree -r -s results/r1.gp newick1.nwk
./plotTree -u -s results/r2.gp newick1.nwk
./plotTree -r -s results/r3.gp -n newick1.nwk
./plotTree -u -s results/r4.gp -n newick1.nwk
./plotTree    -s results/r5.gp -t dumb newick1.nwk
./plotTree    -s results/r6.gp newick2.nwk
./plotTree    -s results/r7.gp -t dumb newick2.nwk
for a in $(seq 7)
do
    sed 's/wxt/qt/' results/r${a}.gp > results/tmp.gp
    mv results/tmp.gp results/r${a}d.gp
done
