package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/softpuff/esecret/mutate"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
}

func handleMutate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		sendError(err, w)
		return
	}

	mutated, err := mutate.Mutate(body)
	if err != nil {
		sendError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func sendError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Error %v", err)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/mutate", handleMutate)

	s := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
