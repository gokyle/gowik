package main

/*
   code to support searching pages for keywords
 */

import ("regexp")

type SearchResult struct {
        Page    string
        Filename string
        Hits    int
}

type SearchTerm struct {
        Term string
        Regex *regexp.Regexp
}
