package rtchat

import "sync"

type Chat struct {
	subs  map[chan Message]struct{}
	mutex *sync.Mutex
}

func NewChat(chat string) *Chat {
	return &Chat{
		subs: map[chan Message]struct{}{},
	}
}

func (c *Chat) SendMessage(username string, message string) error {
	msg := Message{
		From: username,
		Data: message,
	}

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
