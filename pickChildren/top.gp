set terminal postscript eps monochrome size 3cm,2cm
set output "top.ps"
unset xtics
unset ytics
unset border
set label " 7" l rotate by 0 at 0,1.5 front
set label " 6" l rotate by 0 at 1,0.5 front
set label " 2" l rotate by 0 at 3,0 front
set label " 1" l rotate by 0 at 3,1 front
set label " 5" l rotate by 0 at 2,2.5 front
set label " 3" l rotate by 0 at 3,2 front
set label " 4" l rotate by 0 at 3,3 front
plot "-" t "" w l lc "black" lw 3
0 1.5
0 1.5

0 1.5
0 0.5

1 0.5
0 0.5

1 0.5
1 0

3 0
1 0

1 0.5
1 1

3 1
1 1

0 1.5
0 2.5

2 2.5
0 2.5

2 2.5
2 2

3 2
2 2

2 2.5
2 3

3 3
2 3