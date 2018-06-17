package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"zvelo.io/ttlru"
)

var cache ttlru.Cache
var srv http.Server

type resp struct {
	success bool
	value   interface{}
}

func api(w http.ResponseWriter, r *http.Request) {
	ok := false
	ret := *new(interface{})
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if r.Method == "PUT" {
		cache = ttlru.New(128, ttlru.WithTTL(20*time.Second))
	} else if r.Method == "POST" {
		ok = cache.Set(key, value)
	} else if r.Method == "GET" {
		ret, ok = cache.Get(key)
	}
	if err := json.NewEncoder(w).Encode(resp{success: ok, value: ret}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	cache = ttlru.New(128, ttlru.WithTTL(20*time.Second))
	http.HandleFunc("/api", api)
	if err := http2.ConfigureServer(&srv, &http2.Server{}); err != nil {
		panic(err)
	}
	// openssl genrsa -out server.key 2048
	// openssl req -new -x509 -key server.key -out server.pem -days 3650
	// testing: h2load -n10000 -c100 -m10 "https://localhost/api?key=foo"
	if err := srv.ListenAndServeTLS("server.pem", "server.key"); err != nil {
		panic(err)
	}
}
