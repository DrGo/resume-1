.PHONY: publish clean
all: output/resume.pdf

output/resume.pdf: output/resume.tex tex/res.cls
	TEXINPUTS=.//:$(TEXINPUTS) pdflatex -interaction=batchmode -output-directory output $<

output/resume.tex: tex/resume.tex.tmpl resume.yaml output/generate
	./output/generate

output/generate: generate.go
	go build -o output/generate

publish: output/resume.pdf
	rsync -uv output/resume.pdf ../twmb-io/static/

clean:
	rm -rf output
