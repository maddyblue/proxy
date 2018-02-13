// fileserve
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	base = flag.String("base", ".", "base directory")
	addr = flag.String("addr", ":1123", "HTTP listen address")
)

func main() {
	flag.Parse()

	sendErr := func(w http.ResponseWriter, r *http.Request, err string, code int) {
		log.Printf("ERR: %s %s: %s", r.Method, r.URL.Path, err)
		http.Error(w, err, code)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		localfile := filepath.Join(*base, r.URL.Path)
		log.Printf("% 7s %s (%s)", r.Method, r.URL.Path, localfile)
		switch r.Method {
		case "PUT":
			if err := os.MkdirAll(filepath.Dir(localfile), 0755); err != nil {
				sendErr(w, r, err.Error(), 500)
				return
			}
			f, err := os.Create(localfile)
			if err != nil {
				sendErr(w, r, err.Error(), 500)
				return
			}
			defer f.Close()
			if _, err := io.Copy(f, r.Body); err != nil {
				sendErr(w, r, err.Error(), 500)
				return
			}
			w.WriteHeader(201)
		case "GET", "HEAD":
			http.ServeFile(w, r, localfile)
		case "DELETE":
			if err := os.Remove(localfile); err != nil {
				sendErr(w, r, err.Error(), 500)
				return
			}
			w.WriteHeader(204)
		default:
			sendErr(w, r, "unsupported method", 400)
		}
	})

	log.Printf("listening on %s", *addr)
	log.Printf("basedir: %s", *base)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
