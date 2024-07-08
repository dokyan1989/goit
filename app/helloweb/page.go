package helloweb

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dokyan1989/goit/misc/httpparser"
	"github.com/dokyan1989/goit/misc/projectpath"
)

var root = projectpath.Root()

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := filepath.Join(root, "app/helloweb/data", p.Title+".txt")
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := filepath.Join(root, "app/helloweb/data", title+".txt")
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "page_view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "page_edit", p)
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

var templates = template.Must(template.ParseFS(mustSubFS("tmpl"), "edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var title string

		if err := httpparser.RouteParams(r, "title", &title); err != nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, title)
	}
}
