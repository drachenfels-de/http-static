package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Listen string
	Root   string
	Debug  bool
}

type MyResponseWriter struct {
	http.ResponseWriter
	Config *Config
}

func (m *MyResponseWriter) WriteHeader(code int) {
	log.Println(code, http.StatusText(code))
	m.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	if m.Config.Debug {
		if err := m.ResponseWriter.Header().Write(os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	m.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(cfg *Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.String())
		// wrap the respose writer
		wr := &MyResponseWriter{Config: cfg, ResponseWriter: w}
		next.ServeHTTP(wr, r)
	})
}

func main() {
	cfg := Config{}

	flag.StringVar(&cfg.Listen, "listen", ":8088", "Listen address `[host]:<port>`")
	flag.StringVar(&cfg.Root, "root", ".", "Path to webserver root directory.`")
	flag.BoolVar(&cfg.Debug, "debug", false, "Dump response headers.`")
	flag.Parse()

	// Simple static webserver:
	fmt.Printf("%#v\n", cfg)
	fileServer := http.FileServer(http.Dir(cfg.Root))
	log.Fatal(http.ListenAndServe(cfg.Listen, loggingMiddleware(&cfg, fileServer)))
}
