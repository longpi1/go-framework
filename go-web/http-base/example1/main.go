package main

import (
	"fmt"
	"net/http"
)

func main()  {
	http.HandleFunc("/", firstHandler)
	http.HandleFunc("/lp",secondHandler)
	http.ListenAndServe("localhost:8000", nil)
}

func firstHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func secondHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}
