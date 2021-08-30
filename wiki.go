package main

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

var pages = []*Page{
	{Title: "Pierwsza", Body: []byte("Opis")},
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Page methods
// func (p *Page) save() error {
// 	filename := p.Title + ".txt"

// 	return ioutil.WriteFile(filename, p.Body, 0600)
// }

func (p *Page) saveLocal() error {
	existing, err := loadLoocal(p.Title)
	if err != nil {
		pages = append(pages, p)
		return nil
	}

	existing.Title = p.Title
	existing.Body = p.Body

	return nil
}

// func loadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
// 	body, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Page{Title: title, Body: body}, nil
// }

func loadLoocal(title string) (*Page, error) {
	for _, page := range pages {
		if page.Title == title {
			return page, nil
		}
	}

	return nil, errors.New("Not found")
}

// Api methods
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadLoocal(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadLoocal(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}

	err := p.saveLocal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
