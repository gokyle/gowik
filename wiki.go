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
	"sort"
	"strings"
)

var titleRegex *regexp.Regexp
var pageRegex *regexp.Regexp

var Wiki struct {
	WikiDir      string
	Extension    string
	Stylesheet   string
	PageTemplate string
	EditTemplate string
}

var NoPage = `
    <h1>%s doesn't exist!</h1>
    <p>You can <a href="?mode=edit">create it</a>.</p>
`

var pageTemplates = []string{"head.html", "navbar.html", "body.html", "footer.html"}
var editTemplates = []string{"head.html", "navbar.html", "edit.html", "footer.html"}

func initWiki(wikiCfg map[string]string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic("could not get working directory: " + err.Error())
	}
	Wiki.WikiDir = cwd
	Wiki.Extension = ".md"
	Wiki.Stylesheet = "http://twitter.github.com/bootstrap/assets/css/bootstrap.css"

	for key, val := range wikiCfg {
		switch key {
		case "wikipath":
			pageDir, err := filepath.Abs(val)
			if err != nil {
				fmt.Println("[!] could not find ", val)
				os.Exit(1)
			}
			Wiki.WikiDir = pageDir
		case "extension":
			Wiki.Extension = val
		case "stylesheet":
			Wiki.Stylesheet = val
		}
	}

	initDefaultPaths()
	for _, tFile := range pageTemplates {
		pageFile := filepath.Join(Wiki.WikiDir, "templates", tFile)
		page, err := ioutil.ReadFile(pageFile)
		if err != nil {
			fmt.Println("[!] couldn't read ", pageFile)
			os.Exit(1)
		}
		Wiki.PageTemplate += string(page)
	}

	for _, tFile := range editTemplates {
		pageFile := filepath.Join(Wiki.WikiDir, "templates", tFile)
		page, err := ioutil.ReadFile(pageFile)
		if err != nil {
			fmt.Println("[!] couldn't read ", pageFile)
			os.Exit(1)
		}
		Wiki.EditTemplate += string(page)
	}

	pageRegex = regexp.MustCompile(fmt.Sprintf("^%s/pages/(.+)$", Wiki.WikiDir))
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
	AuthRequired  bool
	Authenticated bool
	Title         string
	Filename      string
	Body          template.HTML
	Content       string
	Special       bool
	Error         error
	ShowError     bool
	Message       string
	ShowMessage   bool
}

func LoadPage(r *http.Request) (page *Page) {
	page = new(Page)
	page.AuthRequired = Security.Enabled
	page.Authenticated = authenticated(r)
	page.RequestToFile(r.URL.Path)
	page.RenderMarkdown()
	if page.Error != nil && os.IsNotExist(page.Error) {
		page.Body = template.HTML(fmt.Sprintf(NoPage, page.Title))
	}
	return
}

func LoadPageFile(path string, r *http.Request) (page *Page) {
	page = new(Page)
	if r != nil {
		page.AuthRequired = Security.Enabled
		page.Authenticated = authenticated(r)
	}
	page.RequestToFile(path)
	page.RenderMarkdown()
	if page.Error != nil && os.IsNotExist(page.Error) {
		page.Body = template.HTML(fmt.Sprintf(NoPage, page.Title))
	}
	return page
}

// RequestToFile translates an incoming request to a markdown filename.
func (page *Page) RequestToFile(requestPath string) {
	if requestPath == "/" {
		requestPath = "/index" + Wiki.Extension
	} else if filepath.Ext(requestPath) == "" {
		requestPath = requestPath + ".md"
	} else {
		// convert extension -> md
	}
	page.Filename = filepath.Join(Wiki.WikiDir, "pages", requestPath)
}

func (page *Page) RenderMarkdown() {
	page.Title = titleRegex.ReplaceAllString(
		pageRegex.ReplaceAllString(page.Filename, "$1"), "$1")
	mdContent, err := ioutil.ReadFile(page.Filename)
	if err != nil {
		page.Error = err
		return
	}
	page.Content = strings.TrimSpace(string(mdContent))
	page.Body = template.HTML(string(blackfriday.MarkdownCommon(mdContent)))
}

func ServePage(t *template.Template, page *Page, w http.ResponseWriter) {
	out, err := webshell.BuildTemplate(t, page)
	if err != nil {
		panic("template error: " + err.Error())
	}
	w.Write(out)
}

func UpdatePage(page *Page, newBody string) (err error) {
	_, err = os.Stat(filepath.Dir(page.Filename))
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(page.Filename), 0755)
	}

	if err != nil {
		page.Error = err
		page.ShowError = true
		return
	}
	err = ioutil.WriteFile(page.Filename, []byte(newBody), 0666)
	if err != nil {
		return
	}
	return nil
}

func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("GET", "/", nil)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func Template(tplSource string, w http.ResponseWriter, r *http.Request) *template.Template {
	t := template.New("page")
	t, err := t.Parse(tplSource)
	if err != nil {
		webshell.Error500(err.Error(), "text/plain", w, r)
		return nil
	}
	return t
}

func PageList() []string {
	pageDir := filepath.Join(Wiki.WikiDir, "pages")
	pages := make([]string, 0)
	filepath.Walk(pageDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && titleRegex.MatchString(path) {
			title := titleRegex.ReplaceAllString(path, "$1")
			pages = append(pages, pageRegex.ReplaceAllString(title, "/$1"))
		}
		return err
	})
	sort.Strings(pages)
	return pages
}
