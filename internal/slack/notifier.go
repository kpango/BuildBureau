package slack

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/kpango/BuildBureau/internal/config"
	"github.com/slack-go/slack"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationProjectStart    NotificationType = "project_start"
	NotificationTaskComplete    NotificationType = "task_complete"
	NotificationError           NotificationType = "error"
	NotificationProjectComplete NotificationType = "project_complete"
)

// NotificationData holds data for notification templates
type NotificationData struct {
	ProjectName  string
	TaskName     string
	Agent        string
	ErrorMessage string
	Timestamp    time.Time
}

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

	// Test authentication
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	_, err := client.AuthTestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Slack: %w", err)
	}

	return &Notifier{
		client:  client,
		config:  cfg,
		enabled: true,
	}, nil
}

// Send sends a notification to Slack
func (n *Notifier) Send(ctx context.Context, notifType NotificationType, data NotificationData) error {
	if !n.enabled {
		return nil
	}

	var notifConfig config.NotificationConfig
	switch notifType {
	case NotificationProjectStart:
		notifConfig = n.config.Notifications.ProjectStart
	case NotificationTaskComplete:
		notifConfig = n.config.Notifications.TaskComplete
	case NotificationError:
		notifConfig = n.config.Notifications.Error
	case NotificationProjectComplete:
		notifConfig = n.config.Notifications.ProjectComplete
	default:
		return fmt.Errorf("unknown notification type: %s", notifType)
	}

	if !notifConfig.Enabled {
		return nil
	}

	// Parse and execute template
	data.Timestamp = time.Now()
	message, err := n.renderTemplate(notifConfig.Message, data)
	if err != nil {
		return fmt.Errorf("failed to render message template: %w", err)
	}

	// Send message with retry
	return n.sendWithRetry(ctx, message)
}

// SendProjectStart sends a project start notification
func (n *Notifier) SendProjectStart(ctx context.Context, projectName string) error {
	return n.Send(ctx, NotificationProjectStart, NotificationData{
		ProjectName: projectName,
	})
}

// SendTaskComplete sends a task completion notification
func (n *Notifier) SendTaskComplete(ctx context.Context, taskName, agent string) error {
	return n.Send(ctx, NotificationTaskComplete, NotificationData{
		TaskName: taskName,
		Agent:    agent,
	})
}

// SendError sends an error notification
func (n *Notifier) SendError(ctx context.Context, errorMsg, agent string) error {
	return n.Send(ctx, NotificationError, NotificationData{
		ErrorMessage: errorMsg,
		Agent:        agent,
	})
}

// SendProjectComplete sends a project completion notification
func (n *Notifier) SendProjectComplete(ctx context.Context, projectName string) error {
	return n.Send(ctx, NotificationProjectComplete, NotificationData{
		ProjectName: projectName,
	})
}

// renderTemplate renders a message template with data
func (n *Notifier) renderTemplate(tmplStr string, data NotificationData) (string, error) {
	tmpl, err := template.New("message").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// sendWithRetry sends a message with retry logic
func (n *Notifier) sendWithRetry(ctx context.Context, message string) error {
	var lastErr error

	for i := 0; i < n.config.RetryCount; i++ {
		_, _, err := n.client.PostMessageContext(
			ctx,
			n.config.ChannelID,
			slack.MsgOptionText(message, false),
		)
		if err == nil {
			return nil
		}

		lastErr = err

		// Wait before retry
		if i < n.config.RetryCount-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second * time.Duration(i+1)):
				// Exponential backoff
			}
		}
	}

	return fmt.Errorf("failed to send Slack message after %d retries: %w", n.config.RetryCount, lastErr)
}

// IsEnabled returns whether Slack notifications are enabled
func (n *Notifier) IsEnabled() bool {
	return n.enabled
}
