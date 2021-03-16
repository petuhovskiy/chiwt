package web

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const templatesPrefix = "templates/"

//go:embed templates
var templates embed.FS

type Render struct {
}

func (r *Render) Main(w http.ResponseWriter) {
	b, err := resources.ReadFile("resources/index.html")
	if err != nil {
		log.WithError(err).Error("failed to read file")
		return
	}

	w.Write(b)
}

func (r *Render) SignIn(w http.ResponseWriter) {
	b, err := resources.ReadFile("resources/signin.html")
	if err != nil {
		log.WithError(err).Error("failed to read file")
		return
	}

	w.Write(b)
}