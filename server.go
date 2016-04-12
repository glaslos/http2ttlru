package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/zvelo/ttlru"
	"golang.org/x/net/http2"
)

var l = ttlru.New(128, 20*time.Second)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	fmt.Fprintf(w, "Hi tester %q\n", html.EscapeString(r.URL.Path))
}

func api(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	ret := false
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if r.Method == "PUT" {
		l = ttlru.New(128, 2*time.Second)
	} else if r.Method == "POST" {
		ret = l.Set(key, value)
	} else if r.Method == "GET" {
		inter, _ := l.Get(key)
		value, _ = inter.(string)
	}
	fmt.Fprintf(w, "%q, %q", value, strconv.FormatBool(ret))
}

func main() {
	var srv http.Server
	http.HandleFunc("/", index)
	http.HandleFunc("/api", api)
	http2.ConfigureServer(&srv, &http2.Server{})
	// openssl genrsa -out server.key 2048
	// openssl req -new -x509 -key server.key -out server.pem -days 3650
	log.Fatal(srv.ListenAndServeTLS("server.pem", "server.key"))
}
