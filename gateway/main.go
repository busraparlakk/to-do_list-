package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	apiURL, err := url.Parse("http://localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	webURL, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	apiProxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(apiURL)
		},
	}

	webProxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(webURL)
		},
	}

	mux := http.NewServeMux()

	mux.Handle("/api/", apiProxy)
	mux.Handle("/", webProxy)

	log.Println("Gateway http://localhost:8080 adresinde başlatılıyor")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
