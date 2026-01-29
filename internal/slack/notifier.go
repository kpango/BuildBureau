package slack

import (
	"fmt"
	"log"

	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/pkg/types"
	"github.com/slack-go/slack"
)

// Notifier handles Slack notifications
type Notifier struct {
	client  *slack.Client
	config  config.SlackConfig
	enabled bool
}

// NewNotifier creates a new Slack notifier
func NewNotifier(cfg config.SlackConfig) (*Notifier, error) {
	if !cfg.Enabled {
		return &Notifier{
			enabled: false,
		}, nil
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("Slack token is required when Slack is enabled")
	}

	client := slack.New(cfg.Token)

	return &Notifier{
		client:  client,
		config:  cfg,
		enabled: true,
	}, nil
}

// IsEnabled returns whether Slack notifications are enabled
func (n *Notifier) IsEnabled() bool {
	return n.enabled
}

// NotifyEvent sends a notification for an agent event
func (n *Notifier) NotifyEvent(event types.AgentEvent) error {
	if !n.enabled {
		return nil
	}

	shouldNotify := n.shouldNotifyForEvent(event)
	if !shouldNotify {
		return nil
	}

	channel := n.getChannelForEvent(event)
	message := n.formatEventMessage(event)

	return n.sendMessage(channel, message)
}

// shouldNotifyForEvent determines if we should send a notification for this event
func (n *Notifier) shouldNotifyForEvent(event types.AgentEvent) bool {
	switch event.Type {
	case types.EventTaskAssigned:
		for _, role := range n.config.Notifications.NotifyOnTaskAssigned {
			if role == string(event.Agent) {
				return true
			}
		}
	case types.EventTaskCompleted:
		for _, role := range n.config.Notifications.NotifyOnTaskCompleted {
			if role == string(event.Agent) {
				return true
			}
		}
	case types.EventError:
		return n.config.Notifications.NotifyOnError
	}
	return false
}

// getChannelForEvent determines which channel to use for an event
func (n *Notifier) getChannelForEvent(event types.AgentEvent) string {
	switch event.Agent {
	case types.RoleCEO, types.RoleManager:
		return n.config.Channels.Management
	case types.RoleLead, types.RoleEmployee:
		return n.config.Channels.Dev
	default:
		return n.config.Channels.Updates
	}
}

// formatEventMessage creates a formatted message for an event
func (n *Notifier) formatEventMessage(event types.AgentEvent) string {
	emoji := n.getEmojiForEventType(event.Type)
	return fmt.Sprintf("%s [BuildBureau] *%s* %s: %s",
		emoji,
		event.Agent,
		event.Type,
		event.Message,
	)
}

// getEmojiForEventType returns an appropriate emoji for the event type
func (n *Notifier) getEmojiForEventType(eventType types.EventType) string {
	switch eventType {
	case types.EventTaskAssigned:
		return "📋"
	case types.EventTaskCompleted:
		return "✅"
	case types.EventTaskStarted:
		return "🚀"
	case types.EventError:
		return "❌"
	case types.EventMessage:
		return "💬"
	default:
		return "ℹ️"
	}
}

// sendMessage sends a message to a Slack channel
func (n *Notifier) sendMessage(channel, message string) error {
	if channel == "" {
		// No channel configured, skip
		return nil
	}

	_, _, err := n.client.PostMessage(
		channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		log.Printf("Failed to send Slack message: %v", err)
		return err
	}

	return nil
}

// NotifyTaskAssigned sends a notification when a task is assigned
func (n *Notifier) NotifyTaskAssigned(agent types.AgentRole, task types.Task) error {
	event := types.AgentEvent{
		Type:    types.EventTaskAssigned,
		Agent:   agent,
		Message: fmt.Sprintf("Task assigned: %s", task.Title),
		TaskID:  task.ID,
	}
	return n.NotifyEvent(event)
}

// NotifyTaskCompleted sends a notification when a task is completed
func (n *Notifier) NotifyTaskCompleted(agent types.AgentRole, task types.Task) error {
	event := types.AgentEvent{
		Type:    types.EventTaskCompleted,
		Agent:   agent,
		Message: fmt.Sprintf("Task completed: %s", task.Title),
		TaskID:  task.ID,
	}
	return n.NotifyEvent(event)
}

// NotifyError sends a notification when an error occurs
func (n *Notifier) NotifyError(agent types.AgentRole, errorMsg string) error {
	event := types.AgentEvent{
		Type:    types.EventError,
		Agent:   agent,
		Message: errorMsg,
	}
	return n.NotifyEvent(event)
}
