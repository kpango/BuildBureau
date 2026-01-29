package a2a

import (
	"context"
	"testing"
	"time"
)

func TestBus(t *testing.T) {
	bus := NewBus()

	// Subscribe
	ch1 := bus.Subscribe("agent1")
	ch2 := bus.SubscribeGlobal()

	msg := Message{
		ID: "1",
		From: "test",
		To: "agent1",
		Type: "PING",
	}

	// Send
	go func() {
		bus.Send(context.Background(), msg)
	}()

	// Verify agent1 received
	select {
	case m := <-ch1:
		if m.ID != "1" {
			t.Errorf("Agent1 got wrong message ID")
		}
	case <-time.After(1 * time.Second):
		t.Error("Agent1 timeout")
	}

	// Verify Global received
	select {
	case m := <-ch2:
		if m.ID != "1" {
			t.Errorf("Global got wrong message ID")
		}
	case <-time.After(1 * time.Second):
		t.Error("Global timeout")
	}
}
