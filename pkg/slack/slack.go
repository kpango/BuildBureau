package slack

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

type Notifier struct {
	Client    *slack.Client
	ChannelID string
	Enabled   bool
}

func NewNotifier(token string, channelID string, enabled bool) *Notifier {
	var client *slack.Client
	if enabled && token != "" {
		client = slack.New(token)
	}

	return &Notifier{
		Client:    client,
		ChannelID: channelID,
		Enabled:   enabled,
	}
}

func (n *Notifier) Notify(msg string) {
	if !n.Enabled || n.Client == nil || n.ChannelID == "" {
		return
	}

	// Run in goroutine to not block main thread
	go func() {
		_, _, err := n.Client.PostMessage(
			n.ChannelID,
			slack.MsgOptionText(msg, false),
		)
		if err != nil {
			log.Printf("Failed to send Slack notification: %v", err)
		}
	}()
}

func (n *Notifier) NotifyError(errStr string) {
    if !n.Enabled {
        return
    }
    msg := fmt.Sprintf(":x: *Error Occurred*\n%s", errStr)
    n.Notify(msg)
}

func (n *Notifier) NotifySuccess(taskName string) {
    if !n.Enabled {
        return
    }
    msg := fmt.Sprintf(":white_check_mark: *Task Completed*\n%s", taskName)
    n.Notify(msg)
}
