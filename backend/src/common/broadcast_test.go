package common

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBroadcaster(t *testing.T) {
	t.Run("Creates New Broadcaster", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()
		assert.NotNil(t, broadcaster)
		assert.NotNil(t, broadcaster.subscribers)
		assert.Len(t, broadcaster.subscribers, 0)
	})

	t.Run("Creates Broadcaster with Different Types", func(t *testing.T) {
		stringBroadcaster := NewBroadcaster[string]()
		intBroadcaster := NewBroadcaster[int]()
		structBroadcaster := NewBroadcaster[struct{ Value int }]()

		assert.NotNil(t, stringBroadcaster)
		assert.NotNil(t, intBroadcaster)
		assert.NotNil(t, structBroadcaster)
	})
}

func TestBroadcaster_Subscribe(t *testing.T) {
	t.Run("Single Subscription", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()

		ch := broadcaster.Subscribe()
		assert.NotNil(t, ch)
		assert.Len(t, broadcaster.subscribers, 1)

		// Verify channel buffer size
		select {
		case <-ch:
			t.Fatal("Channel should be empty")
		default:
			// Expected - channel is empty
		}
	})

	t.Run("Multiple Subscriptions", func(t *testing.T) {
		broadcaster := NewBroadcaster[int]()

		ch1 := broadcaster.Subscribe()
		ch2 := broadcaster.Subscribe()
		ch3 := broadcaster.Subscribe()

		assert.NotNil(t, ch1)
		assert.NotNil(t, ch2)
		assert.NotNil(t, ch3)
		assert.Len(t, broadcaster.subscribers, 3)

		// Verify all channels are different
		assert.NotEqual(t, ch1, ch2)
		assert.NotEqual(t, ch2, ch3)
		assert.NotEqual(t, ch1, ch3)
	})

	t.Run("Concurrent Subscriptions", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()
		numSubscribers := 10
		var wg sync.WaitGroup
		channels := make([]<-chan string, numSubscribers)

		wg.Add(numSubscribers)
		for i := 0; i < numSubscribers; i++ {
			go func(index int) {
				defer wg.Done()
				channels[index] = broadcaster.Subscribe()
			}(i)
		}

		wg.Wait()
		assert.Len(t, broadcaster.subscribers, numSubscribers)

		// Verify all channels are different
		for i := 0; i < numSubscribers; i++ {
			assert.NotNil(t, channels[i])
			for j := i + 1; j < numSubscribers; j++ {
				assert.NotEqual(t, channels[i], channels[j])
			}
		}
	})
}

func TestBroadcaster_Unsubscribe(t *testing.T) {
	t.Run("Unsubscribe Existing Channel", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()

		ch1 := broadcaster.Subscribe()
		ch2 := broadcaster.Subscribe()
		assert.Len(t, broadcaster.subscribers, 2)

		broadcaster.Unsubscribe(ch1)
		assert.Len(t, broadcaster.subscribers, 1)

		// Verify the correct channel was removed
		assert.Equal(t, (<-chan string)(broadcaster.subscribers[0]), ch2)

		// Verify channel is closed
		select {
		case _, ok := <-ch1:
			assert.False(t, ok, "Channel should be closed")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Channel should be closed and readable")
		}
	})

	t.Run("Unsubscribe Non-Existing Channel", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()

		ch1 := broadcaster.Subscribe()
		ch2 := make(chan string) // Not subscribed to broadcaster

		assert.Len(t, broadcaster.subscribers, 1)
		broadcaster.Unsubscribe(ch2) // Should not panic or affect existing subscribers
		assert.Len(t, broadcaster.subscribers, 1)
		assert.Equal(t, (<-chan string)(broadcaster.subscribers[0]), ch1)
	})

	t.Run("Unsubscribe All Channels", func(t *testing.T) {
		broadcaster := NewBroadcaster[int]()

		ch1 := broadcaster.Subscribe()
		ch2 := broadcaster.Subscribe()
		ch3 := broadcaster.Subscribe()
		assert.Len(t, broadcaster.subscribers, 3)

		broadcaster.Unsubscribe(ch2)
		assert.Len(t, broadcaster.subscribers, 2)

		broadcaster.Unsubscribe(ch1)
		assert.Len(t, broadcaster.subscribers, 1)

		broadcaster.Unsubscribe(ch3)
		assert.Len(t, broadcaster.subscribers, 0)
	})

	t.Run("Concurrent Unsubscribe", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()
		numSubscribers := 10
		channels := make([]<-chan string, numSubscribers)

		// Subscribe all channels
		for i := 0; i < numSubscribers; i++ {
			channels[i] = broadcaster.Subscribe()
		}
		assert.Len(t, broadcaster.subscribers, numSubscribers)

		// Unsubscribe concurrently
		var wg sync.WaitGroup
		wg.Add(numSubscribers)
		for i := 0; i < numSubscribers; i++ {
			go func(ch <-chan string) {
				defer wg.Done()
				broadcaster.Unsubscribe(ch)
			}(channels[i])
		}

		wg.Wait()
		assert.Len(t, broadcaster.subscribers, 0)
	})
}

func TestBroadcaster_Broadcast(t *testing.T) {
	t.Run("Broadcast to Single Subscriber", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()
		ch := broadcaster.Subscribe()

		message := "test message"
		broadcaster.Broadcast(message)

		select {
		case received := <-ch:
			assert.Equal(t, message, received)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Message should have been received")
		}
	})

	t.Run("Broadcast to Multiple Subscribers", func(t *testing.T) {
		broadcaster := NewBroadcaster[int]()
		numSubscribers := 5
		channels := make([]<-chan int, numSubscribers)

		for i := 0; i < numSubscribers; i++ {
			channels[i] = broadcaster.Subscribe()
		}

		message := 42
		broadcaster.Broadcast(message)

		for i, ch := range channels {
			select {
			case received := <-ch:
				assert.Equal(t, message, received, "Subscriber %d should receive message", i)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Subscriber %d should have received message", i)
			}
		}
	})

	t.Run("Broadcast to No Subscribers", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()

		// Should not panic
		broadcaster.Broadcast("test message")
		assert.Len(t, broadcaster.subscribers, 0)
	})

	t.Run("Broadcast Multiple Messages", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()
		ch := broadcaster.Subscribe()

		messages := []string{"msg1", "msg2", "msg3"}

		for _, msg := range messages {
			broadcaster.Broadcast(msg)
		}

		for i, expectedMsg := range messages {
			select {
			case received := <-ch:
				assert.Equal(t, expectedMsg, received, "Message %d should match", i)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Message %d should have been received", i)
			}
		}
	})

	t.Run("Broadcast with Full Channel Buffer", func(t *testing.T) {
		broadcaster := NewBroadcaster[int]()
		ch := broadcaster.Subscribe()

		// Fill the channel buffer (100 messages)
		for i := 0; i < 100; i++ {
			broadcaster.Broadcast(i)
		}

		// This broadcast should be skipped due to full buffer
		broadcaster.Broadcast(999)

		// Read messages from channel
		receivedMessages := make([]int, 0, 100)
		for i := 0; i < 100; i++ {
			select {
			case msg := <-ch:
				receivedMessages = append(receivedMessages, msg)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Should have received message %d", i)
			}
		}

		// Verify we received the first 100 messages
		for i, msg := range receivedMessages {
			assert.Equal(t, i, msg)
		}

		// Verify no more messages (the 999 should have been skipped)
		select {
		case unexpected := <-ch:
			t.Fatalf("Should not have received additional message: %v", unexpected)
		case <-time.After(50 * time.Millisecond):
			// Expected - no more messages
		}
	})

	t.Run("Concurrent Broadcast", func(t *testing.T) {
		broadcaster := NewBroadcaster[int]()
		ch := broadcaster.Subscribe()

		numMessages := 10
		var wg sync.WaitGroup

		wg.Add(numMessages)
		for i := 0; i < numMessages; i++ {
			go func(msg int) {
				defer wg.Done()
				broadcaster.Broadcast(msg)
			}(i)
		}

		wg.Wait()

		// Collect all received messages
		receivedMessages := make([]int, 0, numMessages)
		for i := 0; i < numMessages; i++ {
			select {
			case msg := <-ch:
				receivedMessages = append(receivedMessages, msg)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Should have received %d messages, got %d", numMessages, len(receivedMessages))
			}
		}

		assert.Len(t, receivedMessages, numMessages)
	})
}

func TestBroadcaster_Integration(t *testing.T) {
	t.Run("Subscribe, Broadcast, Unsubscribe Flow", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()

		// Subscribe multiple channels
		ch1 := broadcaster.Subscribe()
		ch2 := broadcaster.Subscribe()
		ch3 := broadcaster.Subscribe()

		// Broadcast message to all
		broadcaster.Broadcast("message1")

		// Verify all receive the message
		for i, ch := range []<-chan string{ch1, ch2, ch3} {
			select {
			case msg := <-ch:
				assert.Equal(t, "message1", msg, "Channel %d should receive message", i+1)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Channel %d should have received message", i+1)
			}
		}

		// Unsubscribe one channel
		broadcaster.Unsubscribe(ch2)

		// Broadcast another message
		broadcaster.Broadcast("message2")

		// Verify only remaining channels receive the message
		select {
		case msg := <-ch1:
			assert.Equal(t, "message2", msg)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Channel 1 should have received message")
		}

		select {
		case msg := <-ch3:
			assert.Equal(t, "message2", msg)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Channel 3 should have received message")
		}

		// Verify ch2 is closed
		select {
		case _, ok := <-ch2:
			assert.False(t, ok, "Channel 2 should be closed")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Channel 2 should be closed and readable")
		}
	})

	t.Run("Complex Type Broadcasting", func(t *testing.T) {
		type Message struct {
			ID      int
			Content string
			Data    []byte
		}

		broadcaster := NewBroadcaster[Message]()
		ch := broadcaster.Subscribe()

		message := Message{
			ID:      123,
			Content: "test content",
			Data:    []byte("test data"),
		}

		broadcaster.Broadcast(message)

		select {
		case received := <-ch:
			assert.Equal(t, message.ID, received.ID)
			assert.Equal(t, message.Content, received.Content)
			assert.Equal(t, message.Data, received.Data)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should have received message")
		}
	})
}

// Helper function to setup a test broadcaster with subscribers
func setupBroadcasterWithSubscribers[T any](broadcaster *Broadcaster[T], count int) []<-chan T {
	channels := make([]<-chan T, count)
	for i := 0; i < count; i++ {
		channels[i] = broadcaster.Subscribe()
	}
	return channels
}

func TestBroadcaster_EdgeCases(t *testing.T) {
	t.Run("Unsubscribe Same Channel Twice", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()
		ch := broadcaster.Subscribe()

		broadcaster.Unsubscribe(ch)
		assert.Len(t, broadcaster.subscribers, 0)

		// Should not panic
		broadcaster.Unsubscribe(ch)
		assert.Len(t, broadcaster.subscribers, 0)
	})

	t.Run("Broadcast After Unsubscribe", func(t *testing.T) {
		broadcaster := NewBroadcaster[string]()
		ch := broadcaster.Subscribe()

		broadcaster.Unsubscribe(ch)

		// Should not panic
		broadcaster.Broadcast("test")
		assert.Len(t, broadcaster.subscribers, 0)
	})
}
