package main

import (
	"log"
	"net/http"

	"github.com/go-udemy-course-exercises/exercise-2/handler"
)

func main() {
	http.HandleFunc("/", handler.RootHandler)
	http.HandleFunc("/view/", handler.MakeHandler(handler.ViewHandler))
	http.HandleFunc("/edit/", handler.MakeHandler(handler.EditHandler))
	http.HandleFunc("/save/", handler.MakeHandler(handler.SaveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
