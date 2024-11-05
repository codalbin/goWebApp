package main

import (
	// "fmt"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
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
template.Must is a convenience wrapper that panics when passed a non-nil error value, and otherwise returns the *Template unaltered.
A panic is appropriate here; if the templates can't be loaded the only sensible thing to do is exit the program
*/
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
First, this function extracts the page title from r.URL.Path, the path component of the request URL.
The Path is re-sliced with [len("/view/"):] to drop the leading "/view/" component of the request path.
This is because the path will invariably begin with "/view/", which is not part of the page's title.
The function then loads the page data, formats the page with a string of simple HTML, and writes it to w, the http.ResponseWriter.
*/
// func viewHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/view/"):]
// 	p, _ := loadPage(title)
// 	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
// }
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// Function to load the page (or, if it doesn't exist, create an empty Page struct), and displays an HTML form
// Works but not a good way to print html like this => better to use html file

//	func editHandler(w http.ResponseWriter, r *http.Request) {
//		title := r.URL.Path[len("/edit/"):]
//		p, err := loadPage(title)
//		if err != nil {
//			p = &Page{Title: title}
//		}
//		fmt.Fprintf(w, "<h1>Editing %s</h1>"+
//			"<form action=\"/save/%s\" method=\"POST\">"+
//			"<textarea name=\"body\">%s</textarea><br>"+
//			"<input type=\"submit\" value=\"Save\">"+
//			"</form>",
//			p.Title, p.Title, p.Body)
//	}
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

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

/*
regexp.MustCompile will parse and compile the regular expression, and return a regexp.Regexp.
MustCompile is distinct from Compile in that it will panic if the expression compilation fails, while Compile returns an error as a second parameter
*/
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", makeHandler(viewHandler)) // Move to http://localhost:8080/view/test (where test is the name of the file we want to print)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
