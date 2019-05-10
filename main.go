package main

import "net/http"

func main() {
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":4997", http.DefaultServeMux)
}
