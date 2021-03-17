package web

type RenderContext struct {
	Auth      AuthContext
	Stream    CurrentStream
	SetupInfo *SetupInfo
}

type CurrentStream struct {
	Name      string
	VideoLink string
	IsLive    bool
}

type SetupInfo struct {
	Server    string
	StreamKey string
}
