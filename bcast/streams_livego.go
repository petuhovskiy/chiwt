package bcast

import (
	"github.com/gwuhaolin/livego/configure"
	"github.com/petuhovskiy/chiwt/conf"
	log "github.com/sirupsen/logrus"
)

type LivegoStreams struct {
	cfg *conf.App
}

func NewLivegoStreams(cfg *conf.App) *LivegoStreams {
	return &LivegoStreams{cfg: cfg}
}

func (s *LivegoStreams) SetupInfo(username string) *SetupInfo {
	streamKey, err := configure.RoomKeys.GetKey(username)
	if err != nil {
		log.WithError(err).Fatal("failed to get stream key")
	}

	return &SetupInfo{
		Server:    "rtmp://localhost:1935/live", // TODO:
		StreamKey: streamKey,
	}
}

func (s *LivegoStreams) WatchInfo(req WatchRequest) *WatchInfo {
	//http://127.0.0.1:7002/live/movie.m3u8
	//http://127.0.0.1:7001/live/movie.flv

	//const host = "http://127.0.0.1:7001"
	//return fmt.Sprintf("%s/live/%s.flv", host, name)

	//http://localhost:8080/live/dimaskovas.flv

	return &WatchInfo{
		Quality:          "",
		StreamURL:        "/live/" + req.Name + ".flv",
		AvailableQuality: []string{""},
	}
}
