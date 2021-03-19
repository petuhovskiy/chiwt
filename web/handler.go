package web

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gwuhaolin/livego/configure"
	"github.com/petuhovskiy/chiwt/conf"
	"github.com/petuhovskiy/chiwt/rtchat"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net/http"
	"net/url"
)

type Handler struct {
	cfg    *conf.App
	render *Render
	httpFS http.Handler
	auth   *Auth
	chat   *rtchat.Server
}

func NewHandler(cfg *conf.App, render *Render, auth *Auth, chat *rtchat.Server) *Handler {
	return &Handler{
		cfg:    cfg,
		render: render,
		httpFS: http.FileServer(http.FS(resources)),
		auth:   auth,
		chat:   chat,
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

	message := r.FormValue("message")
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

	ctx := h.renderCtx(r)
	ctx.Stream = CurrentStream{
		Name:      username,
		VideoLink: h.streamSource(username),
		IsLive:    true,
	}

	if ctx.Auth.LoggedIn && ctx.Auth.Username == username {
		streamKey, err := configure.RoomKeys.GetKey(username)
		if err != nil {
			log.WithError(err).Fatal("failed to get stream key")
		}

		ctx.SetupInfo = &SetupInfo{
			Server:    "rtmp://localhost:1935/live",
			StreamKey: streamKey,
		}
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

func (h *Handler) streamSource(name string) string {
	//http://127.0.0.1:7002/live/movie.m3u8
	//http://127.0.0.1:7001/live/movie.flv

	//const host = "http://127.0.0.1:7001"
	//return fmt.Sprintf("%s/live/%s.flv", host, name)

	return "/live/" + name + ".flv"
}

func (h *Handler) renderCtx(r *http.Request) RenderContext {
	return RenderContext{
		Auth: GetAuth(r),
	}
}
