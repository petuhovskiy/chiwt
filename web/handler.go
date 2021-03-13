package web

import (
	"fmt"
	"github.com/petuhovskiy/chiwt/conf"
	"net/http"
)

type Handler struct {
	cfg *conf.App
}

func NewHandler(cfg *conf.App) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) streamSource(name string) string {
	//http://127.0.0.1:7002/live/movie.m3u8
	//http://127.0.0.1:7001/live/movie.flv

	const host = "http://127.0.0.1:7001"
	return fmt.Sprintf("%s/live/%s.flv", host, name)
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	content := `
<video src="%s">
    <!-- Fallback here -->
    No video :(
</video>

`

	w.Header().Add("Content-type", "text/html")
	fmt.Fprintf(w, content, h.streamSource(h.cfg.DefaultChannel))
}
