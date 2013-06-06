package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"regexp"
	"text/template"
)

func assertMap(i interface{}) (map[interface{}]interface{}, error) {
	v, ok := i.(map[interface{}]interface{})
	if !ok {
		return v, errors.New("Not a map!")
	}
	return v, nil
}

func assertSlice(i interface{}) ([]interface{}, error) {
	v, ok := i.([]interface{})
	if !ok {
		return v, errors.New("Not a slice!")
	}
	return v, nil
}

func texEscape(text []byte) []byte {
	escapeYAML := regexp.MustCompile(`&|%|\$|#|_|{|}|\\`) // also replace slashes with double slashes for the backslash escape below
	escapeYAMLtilde := regexp.MustCompile(`~`)
	escapeYAMLcircum := regexp.MustCompile(`\^`)
	escapeYAMLbackslash := regexp.MustCompile(`\\\\`) // conflicts with above, so look for \\
	text = escapeYAML.ReplaceAll(text, []byte(`\${0}`))
	text = escapeYAMLtilde.ReplaceAll(text, []byte(`\textasciitilde{}`))
	text = escapeYAMLcircum.ReplaceAll(text, []byte(`\textasciicircum{}`))
	text = escapeYAMLbackslash.ReplaceAll(text, []byte(`\textbackslash{}`))
	return text
}

func generatePDF(yaml []byte) {
	yaml = texEscape(yaml)
	var all map[string]interface{}
	err := goyaml.Unmarshal(yaml, &all)
	if err != nil {
		panic(err)
	}

	tmpl, err := ioutil.ReadFile("tex/resume.tex.tmpl")
	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBuffer(make([]byte, 0, len(tmpl)))
	funcMap := template.FuncMap{
		"assertMap":   assertMap,
		"assertSlice": assertSlice,
	}
	resume := template.Must(template.New("resume").Funcs(funcMap).Delims("[[", "]]").Parse(string(tmpl)))
	resume.Execute(buffer, all)

	links := regexp.MustCompile(`\[([^]]*)\]\(([^)]*)\)`)
	final := links.ReplaceAll(buffer.Bytes(), []byte(`\href{$2}{$1}`)) // returns copy

	outfile, err := os.Create("output/resume.tex")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	outfile.Write(final)
}

func main() {
	yaml, err := ioutil.ReadFile("resume.yaml")
	if err != nil {
		panic(err)
	}
	generatePDF(yaml)
}
