set terminal qt persist size 640,384
set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb 'white' behind
plot[*:*][*:*] "-" t "g1" w lp pt 7, "-" t "g2" w lp pt 7
1	1
2	2
4	4
e
1	2
2	4
4	8
