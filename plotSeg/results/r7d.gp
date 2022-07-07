set terminal qt persist size 640,384
set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb 'white' behind
set format x ''
unset xtics
set x2tics mirror
set xrange[100:500]
set yrange [100:500] reverse
set x2label ''
set ylabel rotate by -90 ''
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

