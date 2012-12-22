## gowik

This is a partially-functional pseudo-wiki backed by markdown. It is
designed as a replacement for [gitit](http://www.gitit.net).

What's missing?

* proper page linking and highlighting links that don't exist. I'm reading
[blackfriday](https://github.com/russross/blackfriday) to create a
wiki-markdown library, but this is going to take some work (and time).
* installer support. this will allow the wiki to create a set of default
templates and directory layout if one is missing.

So what can you actually do? Edit and create new pages. This program
is designed for a single user, so you can specify a wiki user and
password for authentication. 

You can link to pages in the wiki with
```
        [page title](/page)
```

