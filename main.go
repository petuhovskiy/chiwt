package main

import (
	"github.com/gwuhaolin/livego/configure"
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	"github.com/petuhovskiy/chiwt/bcast"
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

	msg, err := configure.RoomKeys.GetKey(cfg.DefaultChannel)
	if err != nil {
		log.WithError(err).Fatal("failed to get key")
	}
	log.WithField("key", msg).Info("key for stream")

	webHandler := web.NewHandler(cfg)
	webRouter := web.NewRouter(webHandler)
	go web.StartHTTP("web", cfg.WebAddr, webRouter)

	flvServer := httpflv.NewServer(stream)
	go bcast.StartFlvServer(flvServer, cfg.FlvAddr)

	hlsServer := hls.NewServer()
	go bcast.StartHlsServer(hlsServer, cfg.HlsAddr)

	bcast.StartRtmp(stream, hlsServer)
}
