package main

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

type Metadata struct {
	Title     string
	Date      string
	Author    string
	WordCount int
	Path      string
}

func main() {
	inputDir := "posts/raw"
	outputDir := "posts"
	layoutPostTemplate := "html/layout-post.html"
	layoutHomeTemplate := "html/layout-home.html"
	layoutPost, err := os.ReadFile(layoutPostTemplate)
	if err != nil {
		panic(err)
	}
	postTmpl, err := template.New("layout").Parse(string(layoutPost))
	if err != nil {
		panic(err)
	}
	layoutHome, err := os.ReadFile(layoutHomeTemplate)
	if err != nil {
		panic(err)
	}
	homeTmpl, err := template.New("home").Parse(string(layoutHome))
	if err != nil {
		panic(err)
	}
	var posts []Metadata
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			content, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			mdContent, metadata, err := parseMarkdown(content)
			if err != nil {
				panic(err)
			}
			metadataHTML := fmt.Sprintf("<h1>%s</h1>", metadata.Title)
			metadataHTML += fmt.Sprintf("<div>By <b>%s</b> on <b>%s</b> | %d words, %s to read</div><hr>", metadata.Author, metadata.Date, metadata.WordCount, calculateReadTime(metadata.WordCount))
			htmlContent := []byte(metadataHTML + string(mdContent))
			outputFile := filepath.Join(outputDir, filepath.Base(path[:len(path)-len(filepath.Ext(path))])+".html")
			outputFileDir := filepath.Dir(outputFile)
			if err := os.MkdirAll(outputFileDir, 0755); err != nil {
				panic(err)
			}

			output, err := os.Create(outputFile)
			if err != nil {
				panic(err)
			}
			defer output.Close()
			if err := postTmpl.Execute(output, struct {
				Title   string
				Content template.HTML
				Year    int
			}{
				Title:   metadata.Title,
				Content: template.HTML(htmlContent),
				Year:    time.Now().Year(),
			}); err != nil {
				panic(err)
			}
			metadata.Path = "posts/" + filepath.Base(path[:len(path)-len(filepath.Ext(path))]) + ".html"
			posts = append(posts, metadata)
		}
		return nil
	})
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})
	homeFile, err := os.Create(filepath.Join("./", "index.html"))
	if err != nil {
		panic(err)
	}
	defer homeFile.Close()
	if err := homeTmpl.Execute(homeFile, struct {
		Posts []Metadata
		Year  int
	}{
		Posts: posts,
		Year:  time.Now().Year(),
	}); err != nil {
		panic(err)
	}
}

func calculateReadTime(count int) string {
	readTime := float64(count) / float64(250)
	if readTime < 1 {
		readTime = 1
	}
	if readTime > 1 {
		return fmt.Sprintf("%d minutes", int(math.Ceil(readTime)))
	}
	return fmt.Sprintf("%d minute", int(math.Round(readTime)))
}

func parseMarkdown(content []byte) ([]byte, Metadata, error) {
	var mdContent []byte
	var metadata Metadata

	parts := strings.SplitN(string(content), "\n---\n", 2)
	if len(parts) < 2 {
		return nil, metadata, errors.New("invalid Markdown format")
	}
	err := yaml.Unmarshal([]byte(parts[0]), &metadata)
	if err != nil {
		return nil, metadata, err
	}
	mdContent = blackfriday.Run([]byte(strings.TrimSpace(parts[1])))
	metadata.WordCount = len(strings.Fields(string(mdContent)))
	return mdContent, metadata, nil
}
