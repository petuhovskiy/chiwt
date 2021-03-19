package myflv

import (
	"github.com/go-chi/chi/v5"
	"github.com/gwuhaolin/livego/av"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	handler av.Handler
}

type stream struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type streams struct {
	Publishers []stream `json:"publishers"`
	Players    []stream `json:"players"`
}

func NewServer(h av.Handler) *Server {
	return &Server{
		handler: h,
	}
}

func (server *Server) LiveFLV(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("http flv handleConn panic: ", r)
		}
	}()

	url := r.URL.String()
	//u := r.URL.Path
	//if pos := strings.LastIndex(u, "."); pos < 0 || u[pos:] != ".flv" {
	//	http.Error(w, "invalid path", http.StatusBadRequest)
	//	return
	//}
	//path := strings.TrimSuffix(strings.TrimLeft(u, "/"), ".flv")
	//paths := strings.SplitN(path, "/", 2)
	//log.Debug("url:", u, "path:", path, "paths:", paths)
	//
	//if len(paths) != 2 {
	//	http.Error(w, "invalid path", http.StatusBadRequest)
	//	return
	//}
	//
	//msgs := server.getStreams(w, r)
	//if msgs == nil || len(msgs.Publishers) == 0 {
	//	http.Error(w, "invalid path", http.StatusNotFound)
	//	return
	//} else {
	//	include := false
	//	for _, item := range msgs.Publishers {
	//		if item.Key == path {
	//			include = true
	//			break
	//		}
	//	}
	//	if include == false {
	//		http.Error(w, "invalid path", http.StatusNotFound)
	//		return
	//	}
	//}

	app := "live"
	title := chi.URLParam(r, "name")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	writer := httpflv.NewFLVWriter(app, title, url, w)

	server.handler.HandleWriter(writer)
	writer.Wait()
}
