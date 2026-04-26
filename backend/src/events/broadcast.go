package events

import (
	"sync"

	"github.com/m-milek/leszmonitor/log"
)

type broadcaster[T any] struct {
	mu          sync.Mutex
	subscribers []chan T
}

func newBroadcaster[T any]() *broadcaster[T] {
	return &broadcaster[T]{
		subscribers: make([]chan T, 0),
	}
}

func (b *broadcaster[T]) Subscribe() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T, 100) // Buffered channel to avoid blocking
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *broadcaster[T]) Unsubscribe(ch <-chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, subscriber := range b.subscribers {
		if subscriber == ch {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			close(subscriber)
			return
		}
	}
}

func (b *broadcaster[T]) Broadcast(message T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	log.Main.Trace().Msgf("Broadcasting message: %v", message)

	for _, subscriber := range b.subscribers {
		select {
		case subscriber <- message:
		default:
			// If the channel is full, we skip sending the message to avoid blocking
			continue
		}
	}
}
