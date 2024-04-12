package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var story Story

const startingArc = "intro"

type Story map[string]StoryPart

type StoryPart struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {
	b, err := os.ReadFile("./internal/gopher.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &story)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", loadStoryPart)
	fmt.Println("Listening on port 8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func loadStoryPart(w http.ResponseWriter, r *http.Request) {
	storyArc := strings.TrimPrefix(r.URL.Path, "/")
	if storyArc == "" {
		storyArc = startingArc
	}

	part, ok := story[storyArc]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	serveTemplate(w, part)
}

func serveTemplate(w http.ResponseWriter, storyPart StoryPart) {
	tmpl, err := template.ParseFiles("web/template/part.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tmpl == nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, storyPart); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
