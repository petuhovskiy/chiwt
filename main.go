package main

import (
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	"github.com/petuhovskiy/chiwt/bcast"
	"github.com/petuhovskiy/chiwt/bcast/myflv"
	"github.com/petuhovskiy/chiwt/conf"
	"github.com/petuhovskiy/chiwt/web"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := conf.ReadFromEnv()
	if err != nil {
		log.WithError(err).Fatal("failed to read config")
	}

	stream := rtmp.NewRtmpStream()
	flvServer := myflv.NewServer(stream)

	render, err := web.NewRender()
	if err != nil {
		log.WithError(err).Fatal("failed to create render")
	}

	auth, err := web.NewAuth()
	if err != nil {
		log.WithError(err).Fatal("failed to create auth")
	}

	webHandler := web.NewHandler(cfg, render, auth)
	webRouter := web.NewRouter(webHandler, flvServer)
	go web.StartHTTP("web", cfg.WebAddr, webRouter)

	hlsServer := hls.NewServer()
	go bcast.StartHlsServer(hlsServer, cfg.HlsAddr)

	bcast.StartRtmp(stream, hlsServer)
}
