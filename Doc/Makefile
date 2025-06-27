all:
	xelatex fichain.tex && \
		bibtex fichain && \
		xelatex fichain.tex && \
		xelatex fichain.tex

clean:
	rm -f *.aux *.bbl *.blg *.log *.out *.pdf *.run.xml *.bcf *.toc *.lof *.lot *.lol *.synctex.gz *.gz *.fls *.fdb_latexmk *.cut *.spl *.blg *.bbl *.dvi *.ps *.bak *.backup *.sav *.swp *.synctex.gz

.PHONY: clean
