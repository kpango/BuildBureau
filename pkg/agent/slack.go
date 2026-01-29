package agent

import (
	"log"
)

type SlackService struct {
	Enabled  bool
	Token    string
	Channels []string
	NotifyOn []string // "task_assigned", "task_completed", "error"
}

func (s *SlackService) Notify(event, message string) {
	if !s.Enabled || s.Token == "" {
		return
	}
	// Check if we should notify on this event
	shouldNotify := false
	for _, e := range s.NotifyOn {
		if e == event {
			shouldNotify = true
			break
		}
	}
	if !shouldNotify {
		return
	}

	// Mock sending
	log.Printf("[SLACK NOTIFICATION] Channel(s): %v | Event: %s | Message: %s", s.Channels, event, message)
}
