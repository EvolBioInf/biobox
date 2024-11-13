set terminal x11 persist size 640,384
set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb 'white' behind
set xlabel "x"
set ylabel "y"
plot[*:*][*:*] "-" t "g1" w l, "-" t "g2" w l
1	1
2	2
4	4
e
1	2
2	4
4	8
