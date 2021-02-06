package server

import "net/http"

var helloWorld = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
