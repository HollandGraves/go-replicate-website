package main

// 																1. IMPORTS

import (
	"html/template"
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

// renderTemplate : creates a new template, parses the template definitions,
// and applies the template to the data object (i.e. struct, interface, etc e.g. *Page)
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

// viewHandler : extracts the page title from the path without /view/ from the path,
// prints the Title and Boby of the file into some HTML, and displays that at the path
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		// http.StatusFound = 302
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// editHandler : loads the page, and if it doesn't exist creates an empty Page struct,
// and displays an HTML form
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
