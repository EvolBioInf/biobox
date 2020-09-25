all: biobox.tex
	latex biobox
	dvipdf biobox
biobox.tex: cchar.tex
	bash ../scripts/org2nw biobox.org | noweave -n -x > cchar.tex
clean:
	rm -f cchar.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc 
