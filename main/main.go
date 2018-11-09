package main

// 																1. IMPORTS

import (
	"log"
	"net/http"
)

// 																3. MAIN FUNCTION

func main() {
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/view/", MakeHandler(ViewHandler))
	http.HandleFunc("/edit/", MakeHandler(EditHandler))
	http.HandleFunc("/save/", MakeHandler(SaveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
