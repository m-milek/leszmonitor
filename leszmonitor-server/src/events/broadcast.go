package events

import (
	"sync"

	"github.com/m-milek/leszmonitor/log"
	"github.com/rs/zerolog"
)

type eventBus[T any] struct {
	mu          sync.Mutex
	subscribers []chan T
	logger      zerolog.Logger
}

func newEventBus[T any](busName string) *eventBus[T] {
	logger := log.New().With().Str("component", "eventBus").Str("busName", busName).Logger()
	return &eventBus[T]{
		subscribers: make([]chan T, 0),
		logger:      logger,
	}
}

func (b *eventBus[T]) Subscribe() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T, 100) // Buffered channel to avoid blocking
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *eventBus[T]) Unsubscribe(ch <-chan T) {
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

func (b *eventBus[T]) Broadcast(message T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.logger.Trace().Msgf("Broadcasting message: %v", message)

	for _, subscriber := range b.subscribers {
		select {
		case subscriber <- message:
		default:
			// If the channel is full, we skip sending the message to avoid blocking
			continue
		}
	}
}
