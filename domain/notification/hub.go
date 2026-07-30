package notification

import (
	"fmt"
	"sync"
)

type Subscription struct {
	C     <-chan struct{}
	close func()
	once  sync.Once
}

func (s *Subscription) Close() {
	if s == nil {
		return
	}
	s.once.Do(s.close)
}

type Hub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan struct{}]struct{}
}

func NewHub() *Hub {
	return &Hub{subscribers: make(map[string]map[chan struct{}]struct{})}
}

func (h *Hub) Subscribe(recipientType string, recipientID uint) *Subscription {
	key := recipientKey(recipientType, recipientID)
	ch := make(chan struct{}, 1)
	h.mu.Lock()
	if h.subscribers[key] == nil {
		h.subscribers[key] = make(map[chan struct{}]struct{})
	}
	h.subscribers[key][ch] = struct{}{}
	h.mu.Unlock()

	return &Subscription{
		C: ch,
		close: func() {
			h.mu.Lock()
			defer h.mu.Unlock()
			delete(h.subscribers[key], ch)
			if len(h.subscribers[key]) == 0 {
				delete(h.subscribers, key)
			}
			close(ch)
		},
	}
}

func (h *Hub) Publish(recipientType string, recipientID uint) {
	key := recipientKey(recipientType, recipientID)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subscribers[key] {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

var defaultHub = NewHub()

func Subscribe(recipientType string, recipientID uint) *Subscription {
	return defaultHub.Subscribe(recipientType, recipientID)
}

func Publish(recipientType string, recipientID uint) {
	defaultHub.Publish(recipientType, recipientID)
}

func recipientKey(recipientType string, recipientID uint) string {
	return fmt.Sprintf("%s:%d", recipientType, recipientID)
}
