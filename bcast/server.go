package bcast

import (
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	log "github.com/sirupsen/logrus"
	"net"
)

func StartFlvServer(server *httpflv.Server, addr string) {
	log.WithField("name", "flv-http").WithField("addr", addr).Info("staring server")

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.WithError(err).Panic("failed to listen port")
	}

	err = server.Serve(l)
	if err != nil {
		log.WithError(err).Panic("flv server error")
	}
}

func StartHlsServer(server *hls.Server, addr string) {
	log.WithField("name", "hls").WithField("addr", addr).Info("staring server")

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.WithError(err).Panic("failed to listen port")
	}

	err = server.Serve(l)
	if err != nil {
		log.WithError(err).Panic("hls server error")
	}
}

func StartRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
	rtmpAddr := ":1935"

	rtmpListen, err := net.Listen("tcp", rtmpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	if hlsServer == nil {
		rtmpServer = rtmp.NewRtmpServer(stream, nil)
		log.Info("HLS server disable....")
	} else {
		rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
		log.Info("HLS server enable....")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("RTMP server panic: ", r)
		}
	}()
	log.Info("RTMP Listen On ", rtmpAddr)
	rtmpServer.Serve(rtmpListen)
}
