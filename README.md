## gowik

This is a partially-functional pseudo-wiki backed by markdown. It is
designed as a replacement for [gitit](http://www.gitit.net).

### What's missing?

* proper page linking and highlighting links that don't exist. I'm reading
[blackfriday](https://github.com/russross/blackfriday) to create a
wiki-markdown library, but this is going to take some work (and time).

So what can you actually do? Edit and create new pages. This program
is designed for a single user, so you can specify a wiki user and
password for authentication. 

You can link to pages in the wiki with
```
        [page title](/page)
```

### Installing

```
$ go get github.com/gokyle/gowik
$ mkdir ~/wiki
$ cd ~/wiki
$ gowik
```

Source control isn't handled by `gowik` (*yet*), but I use git to backup
and sync my wiki.

## How did this come about?
I couldn't sleep one night and decided to do something about my annoyances
with gitit. It took 30 minutes for the first POC app, which had no stylesheets,
required manually changing the url to edit a page, and only supported editing
and viewing existing pages. From the original index page:

> ## why would you even what is this
> 
> * couldn't sleep
> * using gitit, but h37 ghc and cabal
> * couldn't sleep
> * want to play with webshell moar
> * couldn't sleep
> * why not?
> * couldn't sleep
> * are you not making a wiki?
