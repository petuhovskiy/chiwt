package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.HandleFunc("/js/*", h.Static)
	r.HandleFunc("/css/*", h.Static)

	r.Get("/", h.MainPage)
	r.Get("/signin", h.SignIn)
	r.Post("/signin", h.DoSignIn)

	return r
}
