package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

type PageData struct {
	Title   string
	Content template.HTML
}

func main() {
	os.MkdirAll("docs", 0755)

	tmpl := template.Must(template.ParseFiles("template.html"))
	files, _ := filepath.Glob("posts/*.md")

	// Create index.html header
	indexFile, _ := os.Create("docs/index.html")
	defer indexFile.Close()
	indexFile.WriteString("<!DOCTYPE html><html lang='en'><head><title>Betterment</title><link rel='stylesheet' href='style.css'></head><body><h1>Betterment podcast/blog</h1><h2>Latest posts</h2><ul>")

	for _, file := range files {
		md, _ := os.ReadFile(file)

		var buf bytes.Buffer
		if err := goldmark.Convert(md, &buf); err != nil {
			log.Fatal(err)
		}

		base := filepath.Base(file)
		title := strings.TrimSuffix(base, filepath.Ext(base))

		// Generate post HTML
		out, _ := os.Create("docs/" + title + ".html")
		data := PageData{Title: title, Content: template.HTML(buf.String())}
		tmpl.Execute(out, data)
		out.Close()

		// Add link to index
		indexFile.WriteString("<li><a href='" + title + ".html'>" + title + "</a></li>")
	}

	indexFile.WriteString("</ul></body></html>")
	log.Println("Blog successfully generated in /docs")
}
