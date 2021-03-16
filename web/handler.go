package web

import (
	"embed"
	"fmt"
	"github.com/petuhovskiy/chiwt/conf"
	"net/http"
	"net/url"
)

const resourcesPrefix = "resources/"

//go:embed resources
var resources embed.FS

type Handler struct {
	cfg    *conf.App
	render *Render
	httpFS http.Handler
}

func NewHandler(cfg *conf.App, render *Render) *Handler {
	return &Handler{
		cfg:    cfg,
		render: render,
		httpFS: http.FileServer(http.FS(resources)),
	}
}

func (h *Handler) streamSource(name string) string {
	//http://127.0.0.1:7002/live/movie.m3u8
	//http://127.0.0.1:7001/live/movie.flv

	const host = "http://127.0.0.1:7001"
	return fmt.Sprintf("%s/live/%s.flv", host, name)
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	h.render.Main(w)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	h.render.SignIn(w)
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	p := resourcesPrefix + r.URL.Path
	rp := resourcesPrefix + r.URL.RawPath

	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = p
	r2.URL.RawPath = rp
	h.httpFS.ServeHTTP(w, r2)
}

func (h *Handler) DoSignIn(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	fmt.Fprint(w, username)
}
