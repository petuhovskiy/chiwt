package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/petuhovskiy/chiwt/bcast"
	"github.com/petuhovskiy/chiwt/conf"
	"github.com/petuhovskiy/chiwt/rtchat"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net/http"
	"net/url"
)

type Streams interface {
	SetupInfo(username string) *bcast.SetupInfo
	WatchInfo(req bcast.WatchRequest) *bcast.WatchInfo
}

type Handler struct {
	cfg     *conf.App
	render  *Render
	httpFS  http.Handler
	auth    *Auth
	chat    *rtchat.Server
	streams Streams
}

func NewHandler(cfg *conf.App, render *Render, auth *Auth, chat *rtchat.Server, streams Streams) *Handler {
	return &Handler{
		cfg:     cfg,
		render:  render,
		httpFS:  http.FileServer(http.FS(resources)),
		auth:    auth,
		chat:    chat,
		streams: streams,
	}
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authCtx := h.auth.FromRequest(r)

		ctx := r.Context()
		ctx = context.WithValue(ctx, authKey, authCtx)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (h *Handler) DoSignIn(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	// TODO: verify username

	token, err := h.auth.IssueToken(username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  jwtCookie,
		Value: token,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) DoLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   jwtCookie,
		MaxAge: -1,
		Value:  "",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) ChatSubscribe(conn *websocket.Conn) {
	chat := chi.URLParam(conn.Request(), "chat")

	log.WithField("chat", chat).Info("chat subscription started")

	err := h.chat.Subscribe(chat, func(msg rtchat.Message) error {
		err := websocket.JSON.Send(conn, msg)
		return err
	})

	log.WithField("chat", chat).WithError(err).Info("chat subscription finished")
}

func (h *Handler) ChatSend(w http.ResponseWriter, r *http.Request) {
	chat := chi.URLParam(r, "chat")

	auth := GetAuth(r)
	if !auth.LoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req rtchat.SendRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	message := req.Message
	h.chat.SendMessage(chat, auth.Username, message)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	auth := GetAuth(r)
	if auth.LoggedIn {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	h.render.SignIn(w)
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	ctx := h.renderCtx(r)

	h.render.Main(w, ctx)
}

func (h *Handler) WatchStream(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	quality := r.URL.Query().Get("quality")

	info := h.streams.WatchInfo(bcast.WatchRequest{
		Name:    username,
		Quality: quality,
	})

	ctx := h.renderCtx(r)
	ctx.Stream = CurrentStream{
		Name:      username,
		VideoLink: info.StreamURL,
		Info:      info,
		IsLive:    true,
	}

	if ctx.Auth.LoggedIn && ctx.Auth.Username == username {
		ctx.SetupInfo = h.streams.SetupInfo(username)
	}

	h.render.Stream(w, ctx)
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

func (h *Handler) renderCtx(r *http.Request) RenderContext {
	return RenderContext{
		Auth: GetAuth(r),
	}
}
