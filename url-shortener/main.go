package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type routes struct {
	Path string
	URL  string
}

func parseYaml(file string) ([]routes, error) {
	ymlFile, _ := filepath.Abs(file + ".yml")
	ymlData, err := ioutil.ReadFile(ymlFile)

	if err != nil {
		return nil, err
	}

	var urlmap []routes

	err = yaml.Unmarshal(ymlData, &urlmap)
	if err != nil {
		return nil, err
	}

	return urlmap, nil
}

func mapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := paths[r.URL.Path]

		if url != "" {
			http.Redirect(w, r, url, 301)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func buildMap(routes []routes) map[string]string {
	paths := make(map[string]string)

	for _, y := range routes {
		paths[y.Path] = y.URL
	}

	return paths
}

func yamlHandler(yml string, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedYaml)
	return mapHandler(pathMap, fallback), nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mate)
	return mux
}

func mate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Whoa mate!")
}

func main() {
	yaml := flag.String("yml", "ymlData", "filename of csv")

	flag.Parse()
	fallbackHandler := defaultMux()
	mapHandler, _ := yamlHandler(*yaml, fallbackHandler)

	http.HandleFunc("/", mapHandler)
	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", mapHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
