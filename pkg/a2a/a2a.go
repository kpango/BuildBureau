package a2a

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Message represents a generic message passed between agents.
type Message struct {
	ID        string      `json:"id"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

// Bus is the communication channel for A2A interaction.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string]chan Message // map[agentID]channel
	globalSub   []chan Message          // for UI/Logging to listen to all
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string]chan Message),
		globalSub:   make([]chan Message, 0),
	}
}

// Subscribe registers an agent to receive messages.
func (b *Bus) Subscribe(agentID string) <-chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Message, 100)
	b.subscribers[agentID] = ch
	return ch
}

// SubscribeGlobal registers a listener for all messages (e.g. TUI, Logger).
func (b *Bus) SubscribeGlobal() <-chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Message, 100)
	b.globalSub = append(b.globalSub, ch)
	return ch
}

// Send sends a message to a specific agent.
func (b *Bus) Send(ctx context.Context, msg Message) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Notify global subscribers
	for _, sub := range b.globalSub {
		select {
		case sub <- msg:
		default:
			// Non-blocking drop if full
		}
	}

	targetCh, ok := b.subscribers[msg.To]
	if !ok {
		return fmt.Errorf("agent %s not found", msg.To)
	}

	select {
	case targetCh <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("mailbox full for agent %s", msg.To)
	}
}

// Broadcast sends a message to all subscribers (optional, depending on need).
func (b *Bus) Broadcast(ctx context.Context, msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, sub := range b.globalSub {
		select {
		case sub <- msg:
		default:
		}
	}

	for _, ch := range b.subscribers {
		select {
		case ch <- msg:
		default:
		}
	}
}
