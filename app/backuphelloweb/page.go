package backuphelloweb

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dokyan1989/goit/misc/httpparser"
	"github.com/dokyan1989/goit/misc/projectpath"
)

var root = projectpath.Root()

var templates = template.Must(template.ParseFS(mustSubFS("tmpl"), "edit.html", "view.html"))

type Page struct {
	Title string
	Body  []byte
	Tmpl  string
}

func page(title string, body []byte, tmpl string) *Page {
	return &Page{title, body, tmpl}
}

func (p *Page) save() error {
	filename := filepath.Join(root, "app/helloweb/data", dotTXT(p.Title))
	return os.WriteFile(filename, p.Body, 0600)
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, dotHTML(p.Tmpl), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Redirect struct {
	URL  string
	Code int
}

func (re *Redirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, re.URL, http.StatusFound)
}

type NotFound struct {
	URL  string
	Code int
}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func loadBody(title string) ([]byte, error) {
	filename := filepath.Join(root, "app/helloweb/data", dotTXT(title))
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	var title string
	if err := httpparser.RouteParams(r, "title", &title); err != nil {
		return &NotFound{}, nil
	}

	b, err := loadBody(title)
	if err != nil {
		return &Redirect{URL: "/edit/" + title, Code: http.StatusFound}, nil
	}

	return page(title, b, "view"), nil
}

func editHandler(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	var title string
	if err := httpparser.RouteParams(r, "title", &title); err != nil {
		return &NotFound{}, nil
	}

	b, err := loadBody(title)
	if err != nil {
		b = nil
	}

	return page(title, b, "edit"), nil
}

func saveHandler(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	var title string
	if err := httpparser.RouteParams(r, "title", &title); err != nil {
		return &NotFound{}, nil
	}

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		return nil, err
	}

	return &Redirect{URL: "/view/" + title, Code: http.StatusFound}, nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request) (http.Handler, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h, err := fn(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if h != nil {
			h.ServeHTTP(w, r)
		}
	}
}

func dotHTML(name string) string {
	return fmt.Sprintf("%s.html", name)
}

func dotTXT(name string) string {
	return fmt.Sprintf("%s.txt", name)
}
