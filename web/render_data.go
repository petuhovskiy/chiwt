package web

import "github.com/petuhovskiy/chiwt/bcast"

type RenderContext struct {
	Auth      AuthContext
	Stream    CurrentStream
	SetupInfo *bcast.SetupInfo
}

type CurrentStream struct {
	Name      string
	VideoLink string
	IsLive    bool
	Info      *bcast.WatchInfo
}
