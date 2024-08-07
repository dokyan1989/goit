package helloweb

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

var templates = template.Must(template.ParseFS(mustSubFS("tmpl"), "*.html"))

type Page struct {
	Tmpl        string
	Data        any
	RedirectURL string
	Err         error
}

func page(tmpl string, data any) *Page {
	return &Page{Tmpl: tmpl, Data: data}
}

func pageError(err error) *Page {
	return &Page{Err: err}
}

// func pageRedirect(url string) *Page {
// 	return &Page{RedirectURL: url}
// }

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p == nil {
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	if p.Err != nil {
		code := http.StatusInternalServerError
		http.Error(w, p.Err.Error(), code)
		return
	}

	if p.RedirectURL != "" {
		http.Redirect(w, r, p.RedirectURL, http.StatusFound)
		return
	}

	if p.Tmpl == "" {
		code := http.StatusNotImplemented
		http.Error(w, http.StatusText(code), code)
		return
	}

	name := fmt.Sprintf("%s.html", p.Tmpl)
	err := templates.ExecuteTemplate(w, name, p.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PageHandler func(w http.ResponseWriter, r *http.Request) *Page

func makePageHandler(h PageHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())

		// inject request_id to the logger
		l := zerolog.Ctx(r.Context())
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Logger().With().Str("request_id", reqID)
		})

		p := h(w, r)
		p.ServeHTTP(w, r)
	}
}
