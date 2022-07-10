package main

import (
	"fmt"
	"net/http"
)

func main()  {
	e := new(Engine)
	http.ListenAndServe("localhost:8000", e)
}

type Engine struct {

}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/lp":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

