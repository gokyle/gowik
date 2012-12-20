package main

import (
	"flag"
	config "github.com/gokyle/goconfig"
	"github.com/gokyle/webshell"
	"log"
)

var cfg config.ConfigMap
var app *webshell.WebApp

func init() {
	confFile := flag.String("f", "wiki.conf", "config file")
	flag.Parse()

	var err error
	cfg, err = config.ParseFile(*confFile)
	if err != nil {
		panic("could not open config file")
		// defaultConfig()
	}
	initSecurity(cfg["security"])
        initWiki(cfg["wiki"])
	initServer(cfg["server"])
}

func initServer(serverCfg map[string]string) {
	var (
		address = "127.0.0.1"
		port    = "8080"
	)

	for key, val := range serverCfg {
		switch key {
		case "port":
			port = val
		case "address":
			address = val
		}
	}
	if Security.TLS.Enabled {
		app = webshell.NewTLSApp("gowik", address, port,
			Security.TLS.Key, Security.TLS.Cert)
	} else {
		app = webshell.NewApp("gowik", address, port)
	}
}

func main() {
        app.AddRoute("/", ServeWikiPage)
	log.Fatal(app.Serve())
}
