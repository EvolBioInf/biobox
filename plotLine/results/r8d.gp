set terminal eps color size 5,5
set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb 'white' behind
set output "test.ps"
plot[*:*][*:*] "-" t "g1" w l, "-" t "g2" w l
1	1
2	2
4	4
e
1	2
2	4
4	8
