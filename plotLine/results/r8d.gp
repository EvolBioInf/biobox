set terminal eps color size 340,340
set output "test.ps"
plot[*:*][*:*] "-" t "g1" w l, "-" t "g2" w l
1	1
2	2
4	4
e
1	2
2	4
4	8
