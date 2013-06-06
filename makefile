all: resume.pdf

resume.pdf: resume.tex res.cls
	pdflatex -interaction=batchmode $<

resume.tex: resume.tex.tmpl resume.yaml resume
	./resume

resume: build.go
	go build

clean:
	rm -f resume.pdf resume.tex resume
