package main

import (
	"flag"
        "fmt"
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
                fmt.Println("[!] could not load %s: %s\n",
                        *confFile, err.Error())
                fmt.Println("[+] using defaults")
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
        fmt.Println("[+] wiki serving on ", app.Address())
	log.Fatal(app.Serve())
}
