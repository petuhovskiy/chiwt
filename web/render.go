package web

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type Render struct {
	templates *template.Template
}

func NewRender() (*Render, error) {
	templates, err := template.New("stream").ParseFS(templates, "templates/*")
	if err != nil {
		return nil, err
	}

	return &Render{
		templates: templates,
	}, nil
}

func (r *Render) Main(w http.ResponseWriter, data RenderContext) {
	w.Header().Add("Content-type", "text/html")

	err := r.templates.ExecuteTemplate(w, "main.tmpl", data)
	if err != nil {
		log.WithError(err).Error("failed to execute template")
		fmt.Fprint(w, err)
	}
}

func (r *Render) SignIn(w http.ResponseWriter) {
	w.Header().Add("Content-type", "text/html")

	b, err := resources.ReadFile("resources/signin.html")
	if err != nil {
		log.WithError(err).Error("failed to read file")
		return
	}

	w.Write(b)
}

func (r *Render) Stream(w http.ResponseWriter, data RenderContext) {
	w.Header().Add("Content-type", "text/html")

	err := r.templates.ExecuteTemplate(w, "stream.tmpl", data)
	if err != nil {
		log.WithError(err).Error("failed to execute template")
		fmt.Fprint(w, err)
	}
}
