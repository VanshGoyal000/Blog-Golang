package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


type Page struct {
	Title string
	Body  []byte //For dynamic increasing in byte
}

// func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         m := validPath.FindStringSubmatch(r.URL.Path)
//         if m == nil {
//             http.NotFound(w, r)
//             return
//         }
//         fn(w, r, m[2])
//     }
// }

// here we refer page as P and taking it as a pointer
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, error := os.ReadFile(filename)
	if error != nil {
		return nil, error
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {

	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("URL Path:", r.URL.Path)
	m := validPath.FindStringSubmatch(r.URL.Path)
    log.Println("Matched Substrings:", m)

	title, err := getTitle(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	errr := p.save()
	if errr != nil {
		http.Error(w, errr.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTitle( r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return "", errors.New("invaild")
	}
	return m[2], nil
}

func main() {
	// p1 := &Page{Title: "TestPage" , Body:  []byte("This is a test")}
	// p1.save()
	// p2 , _ :=loadPage("TestPage")
	// fmt.Println(string(p2.Body) )

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Println("Starting server on port :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
