package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

// ROUTE ALL SIGNAL
func WikiServe(w http.ResponseWriter, r *http.Request) {
	mode := r.FormValue("mode")
	if r.Method == "POST" {
		WikiPost(w, r)
	} else if authorised(false, r) {
		if mode != "" {
			switch mode {
			case "edit":
				WikiEdit(w, r)
			case "delete":
				WikiDelete(w, r)
			case "list":
				WikiList(w, r)
			case "search":
				WikiSearch(w, r)
			}
		} else {
			WikiView(w, r)
		}
	} else {
		page := LoadPage(r)
		t := WikiNotAuthorised(page, w, r)
		ServePage(t, page, w)
	}
}

// WikiPost is the top level responder for all POST requests. There are
// three kinds that can occur: login, logout, and page updates.
func WikiPost(w http.ResponseWriter, r *http.Request) {
	mode := r.FormValue("mode")
	switch mode {
	case "login":
		cookie := Login(w, r)
		if cookie != nil {
			http.SetCookie(w, cookie)
			RedirectToIndex(w, r)
		}
		WikiPage("/", w, r)
	case "logout":
		Logout(r)
		WikiPage("/", w, r)
	default:
		if r.FormValue("new_body") == "" || !authorised(true, r) {
			RedirectToIndex(w, r)
			return
		}
		page := LoadPage(r)
		err := UpdatePage(page, r.FormValue("new_body"))
		t := Template(Wiki.PageTemplate, w, r)
		if err != nil {
			page.Error = err
			page.ShowError = true
		} else {
			page = LoadPage(r)
		}
		ServePage(t, page, w)
	}
}

// WikiEdit is responsible for dealing with displaying the editor.
func WikiEdit(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	page := LoadPage(r)
	if Security.Enabled && !page.Authenticated {
		t = WikiNotAuthorised(page, w, r)
	} else {
		t = Template(Wiki.EditTemplate, w, r)
	}

	if t == nil {
		RedirectToIndex(w, r)
	}
	ServePage(t, page, w)
}

// Show the requested wiki page.
func WikiView(w http.ResponseWriter, r *http.Request) {
	WikiPage(r.URL.Path, w, r)
}

// Backend handler for loading a page and templating it. Essentially the
// process is some other function calls this with a path; it loads the
// page from the path (i.e. generating a *Page), builds the template,
// and hands it off to be rendered.
func WikiPage(path string, w http.ResponseWriter, r *http.Request) {
	page := LoadPageFile(path, r)
	t := Template(Wiki.PageTemplate, w, r)
	if t == nil {
		return
	}
	ServePage(t, page, w)
}

// WikiDelete handles the deletion of pages; it verifies the user has proper
// credentials to do so.
func WikiDelete(w http.ResponseWriter, r *http.Request) {
	page := LoadPage(r)
	var t *template.Template
	if Security.Enabled && !page.Authenticated {
		t = WikiNotAuthorised(page, w, r)
		ServePage(t, page, w)
	} else {
		err := os.Remove(page.Filename)
		page = LoadPageFile("/", r)
		if err != nil {
			page.Error = err
			page.ShowError = true
		}
		t := Template(Wiki.PageTemplate, w, r)
		if t == nil {
			return
		}
		ServePage(t, page, w)
	}
}

// WikiNotAuthorised sets a page's body to the unauthorised message and templates it.
func WikiNotAuthorised(page *Page, w http.ResponseWriter, r *http.Request) *template.Template {
	t := Template(Wiki.PageTemplate, w, r)
	page.Body = template.HTML(NotAuthorised)
	page.Content = NotAuthorised
	return t
}

// WikiList builds the listing of all pages in the wiki.
func WikiList(w http.ResponseWriter, r *http.Request) {
	pages := PageList()
	var body string
	for _, pageString := range pages {
		body += "    <li><a href=\"" + pageString + "\">" + pageString[1:] + "</a></li>\n"
	}
	body = fmt.Sprintf(`<h1>Page Listing</h1>
  <ul>
%s
</ul>`, body)
	page := LoadPageFile("/", r)
	page.Body = template.HTML(body)
	page.Special = true
	t := Template(Wiki.PageTemplate, w, r)
	if t == nil {
		return
	}
	ServePage(t, page, w)
}

// WikiSearch searches the wiki for the given keywords.
func WikiSearch(w http.ResponseWriter, r *http.Request) {
	results := SearchPages(r.FormValue("terms"))
	var body string
	for _, res := range results {
		body += "    <li><a href=\"/%s\">%s</a> (%d matches)</li>"
		body = fmt.Sprintf(body, res.Page, res.Page, res.Hits)
	}
	body = fmt.Sprintf(`<h1>Search Results</h1>
  <p>There were <strong>%d</strong> matches:</p>
  <ul>
%s
</ul>`, len(results), body)
	page := LoadPageFile("/", r)
	page.Body = template.HTML(body)
	page.Special = true
	t := Template(Wiki.PageTemplate, w, r)
	if t == nil {
		return
	}
	ServePage(t, page, w)
}
