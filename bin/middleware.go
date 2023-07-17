package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func middleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Printf("HTTP request sent to %s from %s", r.URL.Path, r.RemoteAddr)

		w.Header().Set("Content-type", "Application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		// call registered handler
		n(w, r, ps)
	}
}
