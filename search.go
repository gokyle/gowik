package main

/*
   code to support searching pages for keywords
*/

import (
	"fmt"
	"regexp"
	"strings"
)

type SearchResult struct {
	Page     string
	Filename string
	Hits     int
}

type SearchTerm struct {
	Term  string
	Regex *regexp.Regexp
}

var stripQuotes = regexp.MustCompile("^\"(.+)\"$")

func SearchFile(path string, terms []SearchTerm) (sr *SearchResult) {
	title := titleRegex.ReplaceAllString(path, "$1")
	pageFile := pageRegex.ReplaceAllString(title, "/$1")
	page := LoadPageFile(pageFile, nil)
	if page == nil {
		return
	}

	sr = new(SearchResult)
	sr.Page = page.Title
	sr.Filename = page.Filename
	body := string(page.Body)

	for _, term := range terms {
		matches := term.Regex.FindAllString(body, -1)
		sr.Hits += len(matches)
		matches = term.Regex.FindAllString(page.Title, 1)
		if len(matches) > 0 {
			sr.Hits++
		}
	}

	if sr.Hits > 0 {
		return sr
	}
	return nil
}

func GetSearchTerms(searchString string) []SearchTerm {
	searchString = strings.TrimSpace(searchString) + " "
	terms := make([]SearchTerm, 0)
	var start, stop int
	var quote bool

	for i := 0; i < len(searchString); i++ {
		c := string(searchString[i])
		if c == " " && quote {
			continue
		} else if c == "\"" {
			quote = !quote
		}

		if c == " " || i == len(searchString) {
			stop = i
			term := stripQuotes.ReplaceAllString(searchString[start:stop], "$1")
			term = strings.TrimSpace(term)
			st := SearchTerm{term, regexp.MustCompile("(?i)" + term)}
			terms = append(terms, st)
			start = i
		}
	}
	return terms
}

func SearchPages(termString string) []SearchResult {
	pages := PageList()
	results := make([]SearchResult, 0)
	terms := GetSearchTerms(termString)
	for _, p := range pages {
		sr := SearchFile(p, terms)
		if sr != nil {
			results = append(results, *sr)
		}
	}
	return results
}
