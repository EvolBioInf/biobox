set terminal quartz persist
set logscale x
set logscale y
plot[0.1:10][0.2:100] "-" t "g1" w lp, "-" t "g2" w lp
1	1
2	2
4	4
e
1	2
2	4
4	8
