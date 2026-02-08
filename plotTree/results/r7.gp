set terminal dumb size 79,24
unset xtics
unset ytics
unset border
set label " A\\_1" l rotate by 0 at 1,0 front
set label " B 2" l rotate by 0 at 1,1 front
set label "0.1" c rotate by 0 at 0.95,1.15 front
set title "newick2_1"
plot "-" t "" w l lc "black"
0 0.5
0 0

1 0
0 0

0 0.5
0 1

1 1
0 1

1 1.1
0.9 1.1

0.95 1.15
0.95 1.15

1.2 0
