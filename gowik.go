package main

import (
	"fmt"
	"github.com/gokyle/goconfig"
	"github.com/gokyle/webshell"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
)

var baseNameRegexp = regexp.MustCompile("^([\\w_-]+)\\.\\w+$")
var (
	wiki_dir = "pages"
	wiki_ext = "md"
)

type WikiPage struct {
	FileName string
	RPath    string
	Title    string
	Content  template.HTML
	Raw      string
}

func init() {
	webshell.SERVER_ADDR = ""
	webshell.SERVER_PORT = "8080"
	conf, err := config.ParseFile("wiki.conf")
	if err != nil {
		fmt.Println("[!] config: ", err.Error())
	} else if _, ok := conf["source"]; ok {
		if item := conf["source"]["pages"]; item != "" {
			wiki_dir = item
		}
		if item := conf["source"]["extension"]; item != "" {
			wiki_ext = item
		}
	}
}

func main() {
	webshell.AddRoute("/", wiki)
	webshell.Serve(false, nil)
}

func wiki(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println("teh path: ", path)
	if path == "/" {
		path = "/index"
	}
	page := wiki_dir + path + "." + wiki_ext
	var wp *WikiPage
	if wp = loadWikiPage(page); wp == nil {
		http.NotFound(w, r)
		return
	}
	wp.RPath = r.URL.Path
	var err error
	var out []byte
	if r.Method == "GET" {
		if r.URL.RawQuery == "" {
			showWikiPage(wp, w, r)
			return
		} else if r.URL.RawQuery == "edit" {
			out, err = webshell.ServeTemplate("templates/edit.html", wp)
		}
	} else {
		fmt.Println("[-] updaterfy")
		err = r.ParseForm()
		updateWikiPage(wp, w, r)
		return
	}

	if err != nil {
		webshell.Error500(err.Error(), "text/plain", w, r)
	} else {
		w.Write(out)
	}
}

func showWikiPage(wp *WikiPage, w http.ResponseWriter, r *http.Request) {
	out, err := webshell.ServeTemplate("templates/index.html", wp)
	if err != nil {
		webshell.Error500(err.Error(), "text/plain", w, r)
	} else {
		w.Write(out)
	}
}

func loadWikiPage(filename string) *WikiPage {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("[!] error: ", err.Error())
		return nil
	}

	wp := new(WikiPage)
	wp.Raw = string(in)
	wp.Content = template.HTML(string(blackfriday.MarkdownCommon(in)))
	wp.Title = baseNameRegexp.ReplaceAllString(filepath.Base(filename),
		"$1")
	wp.FileName = filename
	return wp
}

func updateWikiPage(wp *WikiPage, w http.ResponseWriter, r *http.Request) {
	body := r.Form.Get("new_body")
	if body == "" {
		http.NotFound(w, r)
		return
	}

	fmt.Println("the page is ", wp.FileName)
	err := ioutil.WriteFile(wp.FileName, []byte(body), 0600)
	if err != nil {
		webshell.Error500(err.Error(), "text/plain", w, r)
	} else {
		wp = loadWikiPage(wp.FileName)
		showWikiPage(wp, w, r)
	}
}
