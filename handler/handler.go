package handler

// 																1. IMPORTS

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// 																2. TYPES

// Page : handles information about the page will be creating
type Page struct {
	Title string
	Body  []byte
}

// 																3. PERSONALLY DEFINED FUNCTIONS

// Save creates a file with a custom name
func (p *Page) Save() error {
	os.Mkdir("data", 0700)
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// LoadPage loads page from directory
func LoadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// template caching. templates will panic if unable to ParseFiles()
// templates is global so it only has to cach once
//
// replaced template.ParseFiles("tmpl/edit.html", "tmpl/view.html") with
// template.ParseGlob("tmpl/*.html") so as to grab every .html file in /tmpl/
var templates = template.Must(template.ParseGlob("tmpl/*.html"))

// RenderTemplate creates a new template, parses the template definitions,
// and applies the template to the data object (i.e. struct, interface, etc e.g. *Page)
func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		// http.StatusInternalServerError = 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validPath : will make sure that the path that a user can type in is restricted,
// so if they try to go to a different path the program will panic and exit
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// MakeHandler function that wraps the handler functions to keep scope of variables local
// and to error check before being
func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// RootHandler redirects the page to the root directory first
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// http.StatusFound = 302
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

// ViewHandler extracts the page title from the path without /view/ from the path,
// prints the Title and Boby of the file into some HTML, and displays that at the path
func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := LoadPage(title)
	if err != nil {
		fmt.Println("Error:", err)
		// http.StatusNotFound = 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// if I want to potentiall redirect to a different page e.g. custom 404 page
		// http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	RenderTemplate(w, "view", p)
}

// EditHandler loads the page, and if it doesn't exist creates an empty Page struct,
// and displays an HTML form
func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := LoadPage(title)
	if err != nil {
		fmt.Println("Error:", err)
		// p = &Page{Title: title}
	}
	RenderTemplate(w, "edit", p)
}

// SaveHandler saves the editing page and redirects to the view path
func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		fmt.Println("Error:", err)
		// http.StatusInternalServerError = 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// http.StatusFound = 302
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
