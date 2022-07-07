set terminal qt persist size 640,384
set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb 'white' behind
unset xtics
unset ytics
unset border
set label " " l rotate by 0 at 0.017,0.300 front
set label " 75" l rotate by 0 at 0.307,0.376 front
set label " One" l rotate by -21 at 0.494,0.303 front
set label " Two" l rotate by 51 at 0.497,0.608 front
set label " 69" l rotate by 0 at -0.169,0.372 front
set label "Three " r rotate by 303 at -0.439,0.793 front
set label "Four " r rotate by 375 at -0.459,0.296 front
set label "Five " r rotate by 447 at -0.040,-0.699 front
set label "0.1" c rotate by 0 at 0.447,1.017 front
set title "newick_1"
plot "-" t "" w l lc "black"
0.017 0.300
0.000 0.000

0.307 0.376
0.017 0.300

0.494 0.303
0.307 0.376

0.497 0.608
0.307 0.376

-0.169 0.372
0.017 0.300

-0.439 0.793
-0.169 0.372

-0.459 0.296
-0.169 0.372

-0.040 -0.699
0.000 0.000

0.497 0.942
0.397 0.942

0.447 1.017
0.447 1.017

0.593 0

0 1.641

0 -0.848

-0.555 0
