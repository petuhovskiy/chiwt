package rtchat

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

type Chat struct {
	name  string
	subs  map[chan Message]struct{}
	mutex sync.Mutex
}

func NewChat(name string) *Chat {
	return &Chat{
		name: name,
		subs: map[chan Message]struct{}{},
	}
}

func (c *Chat) SendMessage(username string, message string) error {
	msg := Message{
		From: username,
		Data: message,
	}

	log.WithField("name", c.name).WithField("message", msg).Info("new chat message")

	c.broadcast(msg)
	return nil
}

func (c *Chat) Subscribe(f func(msg Message) error) error {
	ch := make(chan Message)

	c.sub(ch)
	defer c.unsub(ch)

	for msg := range ch {
		err := f(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chat) sub(ch chan Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.subs[ch] = struct{}{}
}

func (c *Chat) unsub(ch chan Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.subs, ch)
}

func (c *Chat) broadcast(msg Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for ch := range c.subs {
		ch <- msg
	}
}
