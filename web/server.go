package web

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartHTTP(name string, addr string, h http.Handler) {
	log.WithField("name", name).WithField("addr", addr).Info("starting http server")

	err := http.ListenAndServe(addr, h)
	if err != http.ErrServerClosed {
		log.WithError(err).Panic("server finished")
	}
}
