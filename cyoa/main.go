package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

type chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []option `json:"options"`
}

type option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type storyHandler struct {
	Chapters map[string]chapter
	Template *template.Template
}

func (s storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	if path == "" {
		http.Redirect(w, r, "/intro", 301)
		return
	}

	chapter := s.Chapters[path]

	s.Template.Execute(w, chapter)
}

func jsonHandler() map[string]chapter {
	jsonFile, _ := os.Open("stories.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var chapters map[string]chapter

	json.Unmarshal(byteValue, &chapters)

	return chapters
}

func main() {
	tmpl, _ := template.ParseFiles("sample.html")

	chapters := jsonHandler()

	handler := storyHandler{
		Chapters: chapters,
		Template: tmpl,
	}
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", handler)
}
