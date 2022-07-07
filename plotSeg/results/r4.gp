set terminal wxt persist size 640,384
set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb 'white' behind
set format x ''
unset xtics
set x2tics mirror
set xrange[*:*]
set yrange [*:*] reverse
set x2label 'x'
set ylabel rotate by -90 'y'
plot "-" t '' w l lc "black"
1 1
57 57

65 65
229 229

214 226
235 247

226 244
309 327

303 320
392 409

