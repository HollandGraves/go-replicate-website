package main

// 																1. IMPORTS

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// 																2. TYPES

// Page : handles information about the page will be creating
type Page struct {
	Title string
	Body  []byte
}

// 																3. MAIN FUNCTION

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
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

// template caching. templates will panic if unable to ParseFiles()
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// renderTemplate : creates a new template, parses the template definitions,
// and applies the template to the data object (i.e. struct, interface, etc e.g. *Page)
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validPath : will make sure that the path that a user can type in is restricted,
// so if they try to go to a different path the program will panic and exit
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// makeHandler : function that wraps the handler functions to keep scope of variables local
// and to error check before being
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// viewHandler : extracts the page title from the path without /view/ from the path,
// prints the Title and Boby of the file into some HTML, and displays that at the path
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
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
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// saveHandler : saves the editing page and redirects to the view path
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
