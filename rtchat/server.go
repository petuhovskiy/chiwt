package rtchat

import "sync"

type Server struct {
	chats map[string]*Chat
	mutex sync.Mutex
}

func NewServer() *Server {
	return &Server{chats: map[string]*Chat{}}
}

func (s *Server) SendMessage(chat string, username string, message string) error {
	ch := s.getChat(chat)
	return ch.SendMessage(username, message)
}

func (s *Server) Subscribe(chat string, f func(msg Message) error) error {
	ch := s.getChat(chat)
	return ch.Subscribe(f)
}

func (s *Server) getChat(name string) *Chat {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	chat, ok := s.chats[name]
	if ok {
		return chat
	}

	chat = NewChat(name)
	s.chats[name] = chat

	return chat
}
