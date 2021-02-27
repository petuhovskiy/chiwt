package main

import (
	"github.com/gwuhaolin/livego/configure"
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

func startHTTPFlv(stream *rtmp.RtmpStream) {
	httpflvAddr := ":7001"

	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		log.Fatal(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HTTP-FLV server panic: ", r)
			}
		}()
		log.Info("HTTP-FLV listen On ", httpflvAddr)
		hdlServer.Serve(flvListen)
	}()
}

func startHls() *hls.Server {
	hlsAddr := ":7002"
	hlsListen, err := net.Listen("tcp", hlsAddr)
	if err != nil {
		log.Fatal(err)
	}

	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HLS server panic: ", r)
			}
		}()
		log.Info("HLS listen On ", hlsAddr)
		hlsServer.Serve(hlsListen)
	}()
	return hlsServer
}

func startRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
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

func createWebServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//http://127.0.0.1:7002/live/movie.m3u8
		//http://127.0.0.1:7001/live/movie.flv
		content := `
<video src="http://127.0.0.1:7001/live/movie.flv">
    <!-- Fallback here -->
    No video :(
</video>

`

		w.Header().Add("Content-type", "text/html")
		w.Write([]byte(content))
	})

	return mux
}

func main() {
	stream := rtmp.NewRtmpStream()

	msg, err := configure.RoomKeys.GetKey("movie")
	if err != nil {
		log.WithError(err).Fatal("failed to get key")
	}

	log.WithField("key", msg).Info("key for stream")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Website server panic: ", r)
			}
		}()
		websiteAddr := ":8080"
		log.Info("Website listen On ", websiteAddr)
		websiteServer := createWebServer()

		httpServer := &http.Server{
			Handler: websiteServer,
			Addr: websiteAddr,
		}

		httpListen, err := net.Listen("tcp", websiteAddr)
		if err != nil {
			log.Fatal(err)
		}

		httpServer.Serve(httpListen)
	}()

	startHTTPFlv(stream)
	hlsServer := startHls()
	startRtmp(stream, hlsServer)
}