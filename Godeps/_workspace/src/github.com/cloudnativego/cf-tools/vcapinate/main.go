package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

func mapEnvironment() (out map[string]string) {
	env := os.Environ()
	out = make(map[string]string)
	for _, term := range env {
		s := strings.Split(term, "=")
		out[s[0]] = s[1]
	}
	return
}

func substituteEnvironment(w io.Writer, p string) (err error) {
	t, err := template.ParseFiles(p)
	if err != nil {
		return
	}

	return t.Execute(w, mapEnvironment())
}

func generateVCAPServices(w io.Writer, p string) {
	config := new(bytes.Buffer)
	err := substituteEnvironment(config, p)
	if err != nil {
		fmt.Printf("Error parsing configuration file: %s", err)
		os.Exit(1)
	}

	sc, err := parseSpaceConfiguration(config.String())
	if err != nil {
		fmt.Printf("Invalid configuration file: %s", err)
		os.Exit(1)
	}

	t := template.Must(template.New("services").Funcs(fns).Parse(serviceTemplate))
	t.Execute(w, sc)
}

func main() {
	path := flag.String("path", "", "path to configuration yaml file")
	flag.Parse()
	if *path == "" {
		panic("vcapinate -path=vcap.yml")
	}

	generateVCAPServices(os.Stdout, *path)
}
