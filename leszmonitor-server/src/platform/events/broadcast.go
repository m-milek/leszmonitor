package events

import (
	"sync"

	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

const defaultBufferSize = 100

type EventBus[T any] struct {
	mu          sync.Mutex
	subscribers []chan T
	bufferSize  int
	logger      zerolog.Logger
}

func NewEventBus[T any](busName string) *EventBus[T] {
	return NewEventBusWithBuffer[T](busName, defaultBufferSize)
}

// NewEventBusWithBuffer creates a bus whose subscribers get a buffer of bufferSize messages.
func NewEventBusWithBuffer[T any](busName string, bufferSize int) *EventBus[T] {
	logger := log.New().With().Str("component", "eventBus").Str("busName", busName).Logger()
	return &EventBus[T]{
		subscribers: make([]chan T, 0),
		bufferSize:  bufferSize,
		logger:      logger,
	}
}

func (b *EventBus[T]) Subscribe() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T, b.bufferSize)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *EventBus[T]) Unsubscribe(ch <-chan T) {
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

func (b *EventBus[T]) Broadcast(message T) {
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
