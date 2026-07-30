package notification

import (
	"testing"
	"time"
)

func TestHubPublishesOnlyToMatchingRecipient(t *testing.T) {
	hub := NewHub()
	userSub := hub.Subscribe("user", 7)
	defer userSub.Close()
	adminSub := hub.Subscribe("admin", 7)
	defer adminSub.Close()

	hub.Publish("user", 7)

	select {
	case <-userSub.C:
	case <-time.After(time.Second):
		t.Fatal("expected matching user subscriber to receive signal")
	}

	select {
	case <-adminSub.C:
		t.Fatal("did not expect admin subscriber to receive user signal")
	default:
	}
}

func TestHubPublishDoesNotBlockOnFullSubscriber(t *testing.T) {
	hub := NewHub()
	sub := hub.Subscribe("user", 7)
	defer sub.Close()

	hub.Publish("user", 7)
	done := make(chan struct{})
	go func() {
		hub.Publish("user", 7)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected publish to skip full subscriber instead of blocking")
	}
}
