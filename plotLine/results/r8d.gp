set terminal eps color size 5,5
set output "test.ps"
plot[][] "-" t "g1" w lp, "-" t "g2" w lp
1	1
2	2
4	4
e
1	2
2	4
4	8
