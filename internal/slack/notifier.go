package slack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kpango/BuildBureau/internal/config"
	"github.com/slack-go/slack"
)

// Notifier handles Slack notifications
type Notifier struct {
	client *slack.Client
	config *config.Config
	logger *log.Logger
}

// EventType represents different types of events that can be notified
type EventType string

const (
	EventTaskAssigned     EventType = "task_assigned"
	EventTaskCompleted    EventType = "task_completed"
	EventProjectStarted   EventType = "project_started"
	EventProjectCompleted EventType = "project_completed"
	EventErrorOccurred    EventType = "error_occurred"
	EventMilestoneReached EventType = "milestone_reached"
)

// NotificationPayload contains information about an event to notify
type NotificationPayload struct {
	EventType EventType
	Role      string
	AgentName string
	Message   string
	Details   map[string]string
	Timestamp time.Time
}

// NewNotifier creates a new Slack notifier
func NewNotifier(cfg *config.Config, logger *log.Logger) (*Notifier, error) {
	if !cfg.Slack.Enabled {
		return &Notifier{
			config: cfg,
			logger: logger,
		}, nil
	}

	if cfg.Slack.Token == "" {
		return nil, fmt.Errorf("slack token is required when slack is enabled")
	}

	client := slack.New(cfg.Slack.Token)

	// Test the connection
	_, err := client.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Slack: %w", err)
	}

	logger.Println("Successfully connected to Slack")

	return &Notifier{
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// Notify sends a notification to Slack if configured to do so
func (n *Notifier) Notify(ctx context.Context, payload NotificationPayload) error {
	if !n.config.Slack.Enabled || n.client == nil {
		// Slack is disabled, just log
		n.logger.Printf("[Notification] %s - %s: %s\n", payload.Role, payload.EventType, payload.Message)
		return nil
	}

	// Check if we should notify for this event and role
	if !n.config.ShouldNotify(string(payload.EventType), payload.Role) {
		return nil
	}

	// Get the appropriate channel
	channel := n.config.GetChannelForEvent(string(payload.EventType))
	if channel == "" {
		return fmt.Errorf("no channel configured for event type: %s", payload.EventType)
	}

	// Format the message
	message := n.formatMessage(payload)

	// Send the message
	options := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true),
	}

	// Use threads if configured
	if n.config.Slack.MessageFormat.UseThreads && payload.Details["thread_ts"] != "" {
		options = append(options, slack.MsgOptionTS(payload.Details["thread_ts"]))
	}

	channelID, timestamp, err := n.client.PostMessageContext(ctx, channel, options...)
	if err != nil {
		n.logger.Printf("Failed to send Slack notification: %v\n", err)
		return fmt.Errorf("failed to send slack message: %w", err)
	}

	n.logger.Printf("Sent Slack notification to %s (ts: %s)\n", channelID, timestamp)
	return nil
}

// formatMessage formats a notification message according to configuration
func (n *Notifier) formatMessage(payload NotificationPayload) string {
	var message string

	// Add prefix if configured
	if n.config.Slack.MessageFormat.Prefix != "" {
		message += n.config.Slack.MessageFormat.Prefix + " "
	}

	// Add agent name if configured
	if n.config.Slack.MessageFormat.IncludeAgentName && payload.AgentName != "" {
		message += fmt.Sprintf("[%s] ", payload.AgentName)
	}

	// Add main message
	message += payload.Message

	// Add timestamp if configured
	if n.config.Slack.MessageFormat.IncludeTimestamp {
		timestamp := payload.Timestamp
		if timestamp.IsZero() {
			timestamp = time.Now()
		}
		message += fmt.Sprintf(" (at %s)", timestamp.Format("15:04:05"))
	}

	// Add details if present
	if len(payload.Details) > 0 {
		message += "\n\nDetails:"
		for key, value := range payload.Details {
			if key != "thread_ts" { // Skip internal thread tracking
				message += fmt.Sprintf("\n• %s: %s", key, value)
			}
		}
	}

	return message
}

// NotifyTaskAssigned sends a notification when a task is assigned
func (n *Notifier) NotifyTaskAssigned(ctx context.Context, role, agentName, taskTitle, assignee string) error {
	return n.Notify(ctx, NotificationPayload{
		EventType: EventTaskAssigned,
		Role:      role,
		AgentName: agentName,
		Message:   fmt.Sprintf("Task assigned: %s → %s", taskTitle, assignee),
		Timestamp: time.Now(),
	})
}

// NotifyTaskCompleted sends a notification when a task is completed
func (n *Notifier) NotifyTaskCompleted(ctx context.Context, role, agentName, taskTitle string, details map[string]string) error {
	return n.Notify(ctx, NotificationPayload{
		EventType: EventTaskCompleted,
		Role:      role,
		AgentName: agentName,
		Message:   fmt.Sprintf("Task completed: %s", taskTitle),
		Details:   details,
		Timestamp: time.Now(),
	})
}

// NotifyProjectStarted sends a notification when a project starts
func (n *Notifier) NotifyProjectStarted(ctx context.Context, projectTitle, clientID string) error {
	return n.Notify(ctx, NotificationPayload{
		EventType: EventProjectStarted,
		Role:      "CEO",
		AgentName: "CEO",
		Message:   fmt.Sprintf("New project started: %s", projectTitle),
		Details: map[string]string{
			"Client": clientID,
		},
		Timestamp: time.Now(),
	})
}

// NotifyProjectCompleted sends a notification when a project is completed
func (n *Notifier) NotifyProjectCompleted(ctx context.Context, projectTitle string, details map[string]string) error {
	return n.Notify(ctx, NotificationPayload{
		EventType: EventProjectCompleted,
		Role:      "CEO",
		AgentName: "CEO",
		Message:   fmt.Sprintf("Project completed: %s", projectTitle),
		Details:   details,
		Timestamp: time.Now(),
	})
}

// NotifyError sends a notification when an error occurs
func (n *Notifier) NotifyError(ctx context.Context, role, agentName, errorMsg string) error {
	return n.Notify(ctx, NotificationPayload{
		EventType: EventErrorOccurred,
		Role:      role,
		AgentName: agentName,
		Message:   fmt.Sprintf("Error occurred: %s", errorMsg),
		Timestamp: time.Now(),
	})
}

// NotifyMilestone sends a notification when a milestone is reached
func (n *Notifier) NotifyMilestone(ctx context.Context, role, agentName, milestone string, details map[string]string) error {
	return n.Notify(ctx, NotificationPayload{
		EventType: EventMilestoneReached,
		Role:      role,
		AgentName: agentName,
		Message:   fmt.Sprintf("Milestone reached: %s", milestone),
		Details:   details,
		Timestamp: time.Now(),
	})
}
