package main

import (
	"flag"
	"github.com/codegangsta/negroni"
	"github.com/phyber/negroni-gzip/gzip"
	"net/http"
)

type Config struct {
	Listen string
	Root   string
}

func main() {

	cfg := Config{}

	flag.StringVar(&cfg.Listen, "listen", ":8088", "Listen address `[host]:<port>`")
	flag.StringVar(&cfg.Root, "root", ".", "Path to webserver root directory.`")
	flag.Parse()

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.Use(&negroni.Static{Dir: http.Dir(cfg.Root)})
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Run(cfg.Listen)
}
