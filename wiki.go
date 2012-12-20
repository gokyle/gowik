package main

import (
        "fmt"
        "github.com/gokyle/webshell"
	"github.com/russross/blackfriday"
        "io/ioutil"
        "net/http"
	"os"
	"path/filepath"
	"regexp"
        "html/template"
)

var titleRegex *regexp.Regexp

var Wiki struct {
	PageDir   string
	Extension string
}

func initWiki(wikiCfg map[string]string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic("could not get working directory: " + err.Error())
	}
	Wiki.PageDir = filepath.Join(cwd, "pages")
	Wiki.Extension = ".md"
	for key, val := range wikiCfg {
		switch key {
		case "pages":
			pageDir, err := filepath.Abs(val)
			if err != nil {
				fmt.Println("[!] could not find ", val)
				os.Exit(1)
			}
			Wiki.PageDir = pageDir
		case "extension":
			Wiki.Extension = val
		}
	}

	var extRegex string
	for i := 0; i < len(Wiki.Extension); i++ {
                char := string(Wiki.Extension[i])
		if char == "." {
			extRegex += "\\."
		} else {
			extRegex += char
		}
	}
	titleRegex = regexp.MustCompile("(.+)" + extRegex)
}

type Page struct {
	Authenticated bool
	Title         string
	Filename      string
	Body          template.HTML
	Content       string
	Error         error
	Message       string
	ShowMessage   bool
}

func ServeWikiPage(w http.ResponseWriter, r *http.Request) {
        page := new(Page)
        page.RequestToFile(r.URL.Path)
        page.RenderMarkdown()
        fmt.Printf("[-] wiki -> %+v\n", Wiki)
        fmt.Printf("[-] page -> %+v\n", page)
        body, err := webshell.BuildTemplateFile("templates/index.html", page)
        if err != nil {
                webshell.Error500(err.Error(), "text/plain", w, r)
        } else {
                w.Write(body)
        }
}

// RequestToFile translates an incoming request to a markdown filename.
func (page *Page) RequestToFile(requestPath string) {
	if requestPath == "/" {
		requestPath = "/index" + Wiki.Extension
	}
	page.Filename = filepath.Join(Wiki.PageDir, requestPath)
}

func (page *Page) RenderMarkdown() {
	mdContent, err := ioutil.ReadFile(page.Filename)
	if err != nil {
		page.Error = err
		return
	}
	page.Content = string(mdContent)
	page.Body = template.HTML(string(blackfriday.MarkdownCommon(mdContent)))
	page.Title = titleRegex.ReplaceAllString(filepath.Base(page.Filename),
		"$1")
}
