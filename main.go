package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Series struct {
	Id	int	`json:"id"`
	Type	string	`json:"type"`
	Seasons	int	`json:"seasons"`
	Author	string	`json:"author"`
	Title	string	`json:"title"`
}

type SoManySeries []Series


func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
