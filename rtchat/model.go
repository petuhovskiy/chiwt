package rtchat

type Message struct {
	From string
	Data string
}

type SendRequest struct {
	Message string
}
