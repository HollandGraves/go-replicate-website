package main

// 																1. IMPORTS

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// 																2. TYPES

// Page : handles information about the page will be creating
type Page struct {
	Title string
	Body  []byte
}

// 																3. MAIN FUNCTION

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 																4. PERSONAL DEFINED FUNCTIONS

// Save() : creates a file with a custom name
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage() : loads page from directory
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// viewHandler : extracts the page title from the path without /view/ from the path,
// prints the Title and Boby of the file into some HTML, and displays that at the path
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

// editHandler : loads the page, and if it doesn't exist creates an empty Page struct,
// and displays an HTML form
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	fmt.Fprintf(w, "<h1>Editing &s</h1>"+
		"<form action=\"/save/&s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, p.Body)
}
