package main

import (
	"fmt"
	"github.com/gokyle/webshell"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var titleRegex *regexp.Regexp

var Wiki struct {
	PageDir    string
	Extension  string
	Stylesheet string
}

func initWiki(wikiCfg map[string]string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic("could not get working directory: " + err.Error())
	}
	Wiki.PageDir = filepath.Join(cwd, "pages")
	Wiki.Extension = ".md"
	Wiki.Stylesheet = "http://twitter.github.com/bootstrap/assets/css/bootstrap.css"

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
		case "stylesheet":
			Wiki.Stylesheet = val
		}
	}
	defaultHead = fmt.Sprintf(defaultHead, Wiki.Stylesheet)

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
	update := false
	if r.Method != "GET" && r.Method != "HEAD" {
		update = true
	}
	if r.URL.Path == "/logout" {
		Logout(r)
	}
	page := new(Page)
	page.Authenticated = authorised(update, r)
	page.RequestToFile(r.URL.Path)
	page.RenderMarkdown()
	fmt.Printf("[-] wiki -> %+v\n", Wiki)
	fmt.Printf("[-] page -> %+v\n", page)
	ShowPage(page, w, r)
	//body, err := webshell.BuildTemplateFile("templates/index.html", page)
	// if err != nil {
	//         webshell.Error500(err.Error(), "text/plain", w, r)
	// } else {
	//         ShowPage(page, w, r)
	// }
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

func ShowPage(page *Page, w http.ResponseWriter, r *http.Request) {
	tplSrc := defaultHead + defaultNavBar + defaultDisplayBody + defaultFooter
	t := template.New("page")
	t, err := t.Parse(tplSrc)
	if err != nil {
		fmt.Printf("[!] template error: %s\n", err.Error())
		return
	}
	out, err := webshell.BuildTemplate(t, page)
	if err != nil {
		panic("template error: " + err.Error())
	}
	w.Write(out)
}
