package main

import (
  "errors"
  "fmt"
  "os"
	"io/ioutil"
	"launchpad.net/goyaml"
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

func main() {
	yaml, err := ioutil.ReadFile("resume.yaml")
	if err != nil {
		panic(err)
	}

  funcMap := template.FuncMap{
    "assertMap": assertMap,
    "assertSlice": assertSlice,
  }

  var u map[string]interface{}

	err = goyaml.Unmarshal(yaml, &u)
	if err != nil {
		panic(err)
	}
  fmt.Println(u)

  const tmpl = `
  {{ $contact := assertMap .contact}}
  {{ $work := assertSlice .work }}
  {{ $education := assertSlice .education }}
  {{ $interests := .interests }}
  {{ $contact.name }}
  {{ range $work }}
    {{ .employer }}
    {{ $fluff := assertSlice .fluff }}
    {{ range $fluff }}
      {{ . }}
    {{ end }}
  {{ end }}
  {{ $interests }}
  `
  t := template.Must(template.New("t").Funcs(funcMap).Parse(tmpl))
  t.Execute(os.Stdout, u)
}
