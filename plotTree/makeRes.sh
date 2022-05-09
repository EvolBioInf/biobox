./plotTree -r -S r1.r newick.nwk
./plotTree -u -S r2.r newick.nwk
./plotTree -r -S r3.r -n newick.nwk
./plotTree -u -S r4.r -n newick.nwk
for a in $(seq 4)
do
    sed 's/x11/quartz/' results/r${a}.r > results/tmp.r
    mv results/tmp.r results/r${a}d.r
done
