package main

import (
	"fmt"
	"log"
	"net/http"
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

/*
	First, this function extracts the page title from r.URL.Path, the path component of the request URL. 
	The Path is re-sliced with [len("/view/"):] to drop the leading "/view/" component of the request path. 
	This is because the path will invariably begin with "/view/", which is not part of the page's title.
	The function then loads the page data, formats the page with a string of simple HTML, and writes it to w, the http.ResponseWriter. 
*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", viewHandler) // Move to http://localhost:8080/view/test (where test is the name of the file we want to print)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
