package main

import (
	"fmt"
	"os"
)

// Define the structure of a page
type Page struct {
	Title string
	Body  []byte
}

// Method called by the pointer of a page, no parameter, return value of type error
// Method to save page body to a text file
// If everything goes well it will return "nil"
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600) // 0600  indicates that the file should be created with read-write permissions for the current user only
}

// Function to construct the file name from the title parameter
// Reads the file's contents into a new variable body
// And returns a pointer to a Page literal constructed with the proper title and body values
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}


