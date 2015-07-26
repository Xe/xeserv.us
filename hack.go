package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func init() {
	u, e := url.Parse("http://172.17.42.1:9200")
	if e != nil {
		log.Fatal("Bad destination.")
	}
	h := httputil.NewSingleHostReverseProxy(u)
	s := &http.Server{
		Addr:           ":9200",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()
}
