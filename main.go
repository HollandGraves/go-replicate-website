package main

// 																IMPORTS

import (
	"fmt"
	"io/ioutil"
)

// 																TYPES

// Page : handles information about the page will be creating
type Page struct {
	Title string
	Body  []byte
}

// 																MAIN FUNCTION

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}

// 																PERSONAL DEFINED FUNCTIONS

// Save() : creates a file with a custom name
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage() :
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
