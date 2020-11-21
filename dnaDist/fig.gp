set terminal epslatex monochrome
set output "jc.tex"
set size 5/5., 4/3.
set format xy "\\Large$%g$"
#set format y "\Large $%.0t\times 10^{%T}$"
set xlabel "\\Large$\\pi$"
set ylabel ""
#set format "\Large$%g$"
#set logscale xy
#set pointsize 2
f(x) = x
set key center top
plot [0:0.75][] "jc.dat" title "\\Large$J$" wi li lw 3,\
f(x) title "\\Large$\\pi$" wi li lw 3
