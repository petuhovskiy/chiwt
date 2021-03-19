package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/petuhovskiy/chiwt/bcast/myflv"
	"golang.org/x/net/websocket"
	"net/http"
)

func NewRouter(h *Handler, flv *myflv.Server) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(h.AuthMiddleware)

	// static content
	r.HandleFunc("/js/*", h.Static)
	r.HandleFunc("/css/*", h.Static)

	// main page
	r.Get("/", h.MainPage)

	// watch stream
	r.Get("/w/{username}", h.WatchStream)

	// auth block
	r.Get("/signin", h.SignIn)
	r.Post("/signin", h.DoSignIn)
	r.Get("/logout", h.DoLogout) // TODO: POST?

	// video streams
	r.Get("/live/{name}.flv", flv.LiveFLV)

	// chat
	r.Handle("/chat/{chat}/subscribe", websocket.Handler(h.ChatSubscribe))
	r.Post("/chat/{chat}/send", h.ChatSend)

	return r
}
