package main

import (
	"fmt"
	"log"
	"net/http"
)

func mains() {
	// registers the handler function to handle all requests (the "/" pattern matches all paths)
	http.HandleFunc("/", handler)
	// http.ListenAndServe("localhost:8000", nil) starts the server on localhost port 8000
	// log.Fatal() will terminate the program and log an error if the server fails to start
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
// r *http.Request: Contains all the request information
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)		// %q: Formats the string with double quotes
}