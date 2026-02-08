set terminal x11 persist size 640,384
set object 1 rectangle from screen 0,0 to screen 1,1 fillcolor rgb 'white' behind
unset xtics
unset ytics
unset border
set label " " l rotate by 0 at 0.01725,0.2995 front
set label " 75" l rotate by 0 at 0.3074,0.3757 front
set label " One" l rotate by -21 at 0.4938,0.303 front
set label " Two" l rotate by 51 at 0.4974,0.6078 front
set label " 69" l rotate by 0 at -0.1691,0.3721 front
set label "Three " r rotate by 303 at -0.4392,0.7929 front
set label "Four " r rotate by 375 at -0.4593,0.296 front
set label "Five " r rotate by 447 at -0.04024,-0.6988 front
set label "0.1" c rotate by 0 at 0.4474,1.017 front
set title "newick1_1"
plot "-" t "" w l lc "black"
0.01725 0.2995
0 0

0.3074 0.3757
0.01725 0.2995

0.4938 0.303
0.3074 0.3757

0.4974 0.6078
0.3074 0.3757

-0.1691 0.3721
0.01725 0.2995

-0.4392 0.7929
-0.1691 0.3721

-0.4593 0.296
-0.1691 0.3721

-0.04024 -0.6988
0 0

0.4974 0.942
0.3974 0.942

0.4474 1.017
0.4474 1.017

0.6888 0

0 1.79

0 -0.9972

-0.6506 0
