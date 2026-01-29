package slack

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/kpango/BuildBureau/pkg/types"
	"github.com/slack-go/slack"
)

// Notifier handles Slack notifications with real API integration.
type Notifier struct {
	config  *types.SlackConfig
	client  *slack.Client
	enabled bool
}

// NewNotifier creates a new Slack notifier with real API client.
func NewNotifier(config *types.SlackConfig, token string) (*Notifier, error) {
	if config == nil || !config.Enabled {
		return &Notifier{enabled: false}, nil
	}

	if token == "" {
		return nil, fmt.Errorf("slack token is required when Slack is enabled")
	}

	client := slack.New(token)

	// Test the connection
	_, err := client.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Slack: %w", err)
	}

	return &Notifier{
		config:  config,
		enabled: true,
		client:  client,
	}, nil
}

// Notify sends a notification to Slack.
func (n *Notifier) Notify(ctx context.Context, notificationType, message string) error {
	if !n.enabled {
		// Notifications disabled, skip silently
		return nil
	}

	// Check if this notification type should be sent
	shouldNotify := slices.Contains(n.config.NotifyOn, notificationType)

	if !shouldNotify {
		return nil
	}

	// Send to all configured channels
	var lastErr error
	for _, channel := range n.config.Channels {
		_, _, err := n.client.PostMessage(
			channel,
			slack.MsgOptionText(message, false),
			slack.MsgOptionAsUser(true),
		)
		if err != nil {
			lastErr = fmt.Errorf("failed to send to %s: %w", channel, err)
			fmt.Printf("Warning: %v\n", lastErr)
		}
	}

	return lastErr
}

// NotifyTaskAssigned sends a task assigned notification.
func (n *Notifier) NotifyTaskAssigned(ctx context.Context, taskID, assignedTo string) error {
	message := fmt.Sprintf("✅ Task `%s` assigned to *%s* at %s",
		taskID, assignedTo, time.Now().Format(time.RFC3339))
	return n.Notify(ctx, "task_assigned", message)
}

// NotifyTaskCompleted sends a task completed notification.
func (n *Notifier) NotifyTaskCompleted(ctx context.Context, taskID string, status string) error {
	message := fmt.Sprintf("🎉 Task `%s` completed with status: *%s* at %s",
		taskID, status, time.Now().Format(time.RFC3339))
	return n.Notify(ctx, "task_completed", message)
}

// NotifyError sends an error notification.
func (n *Notifier) NotifyError(ctx context.Context, taskID string, err error) error {
	message := fmt.Sprintf("❌ Error in task `%s`: %v at %s",
		taskID, err, time.Now().Format(time.RFC3339))
	return n.Notify(ctx, "error", message)
}
