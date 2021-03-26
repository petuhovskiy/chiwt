package bcast

import (
	"github.com/petuhovskiy/chiwt/conf"
	"sync"
)

type IngestorStreams struct {
	cfg *conf.App

	quality []string

	servers []string
	counter int
	mutex   sync.Mutex
}

func NewIngestor(cfg *conf.App) *IngestorStreams {
	return &IngestorStreams{
		cfg:     cfg,
		quality: append([]string{""}, cfg.AvailableQuality...),
		servers: cfg.IngestorWatch,
	}
}

func (s *IngestorStreams) SetupInfo(username string) *SetupInfo {
	return &SetupInfo{
		Server:    s.cfg.IngestorUpload,
		StreamKey: username,
	}
}

func (s *IngestorStreams) WatchInfo(req WatchRequest) *WatchInfo {
	link := s.selectWatchServer(req)
	if req.Quality == "" {
		link += s.cfg.OriginalPrefix
	} else {
		link += s.cfg.EncodedPrefix
	}

	link += req.Name

	if req.Quality != "" {
		link += "_" + req.Quality
	}

	link += ".flv"

	return &WatchInfo{
		Quality:          req.Quality,
		StreamURL:        link,
		AvailableQuality: s.quality,
	}
}

func (s *IngestorStreams) selectWatchServer(req WatchRequest) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// round-robin
	s.counter++
	return s.servers[s.counter%len(s.servers)]
}
