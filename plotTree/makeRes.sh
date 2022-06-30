./plotTree -r -s results/r1.gp newick.nwk
./plotTree -u -s results/r2.gp newick.nwk
./plotTree -r -s results/r3.gp -n newick.nwk
./plotTree -u -s results/r4.gp -n newick.nwk
for a in $(seq 4)
do
    sed 's/wxt/qt/' results/r${a}.gp > results/tmp.gp
    mv results/tmp.gp results/r${a}d.gp
done
