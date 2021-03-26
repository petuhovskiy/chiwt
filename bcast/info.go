package bcast

type SetupInfo struct {
	Server    string
	StreamKey string
}

type WatchRequest struct {
	Name    string
	Quality string
}

type WatchInfo struct {
	Quality   string
	StreamURL string

	// "240", "360", "480", "720", "1080", ""
	AvailableQuality []string
}
