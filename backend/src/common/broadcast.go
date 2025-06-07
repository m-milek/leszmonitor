package common

import "sync"

type Broadcaster[T any] struct {
	mu          sync.Mutex
	subscribers []chan T
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		subscribers: make([]chan T, 0),
	}
}

func (b *Broadcaster[T]) Subscribe() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T, 100) // Buffered channel to avoid blocking
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broadcaster[T]) Unsubscribe(ch <-chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, subscriber := range b.subscribers {
		if subscriber == ch {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			close(subscriber) // Close the channel to signal no more messages
			return
		}
	}
}

func (b *Broadcaster[T]) Broadcast(message T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, subscriber := range b.subscribers {
		select {
		case subscriber <- message:
		default:
			// If the channel is full, we skip sending the message to avoid blocking
			continue
		}
	}
}
