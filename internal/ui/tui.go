package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kpango/BuildBureau/pkg/types"
)

// Model represents the TUI model
type Model struct {
	viewport    viewport.Model
	textarea    textarea.Model
	events      []types.AgentEvent
	ready       bool
	width       int
	height      int
	inputMode   bool
	taskChannel chan types.Task
}

// TickMsg represents a periodic tick for updating the UI
type TickMsg time.Time

// EventMsg wraps an agent event for the UI
type EventMsg types.AgentEvent

// NewModel creates a new TUI model
func NewModel(taskChannel chan types.Task) Model {
	ta := textarea.New()
	ta.Placeholder = "Enter client request here..."
	ta.Focus()
	ta.CharLimit = 1000
	ta.SetWidth(80)
	ta.SetHeight(3)

	return Model{
		textarea:    ta,
		events:      make([]types.AgentEvent, 0),
		inputMode:   true,
		taskChannel: taskChannel,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		tickCmd(),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.inputMode && m.textarea.Value() != "" {
				// Submit task
				task := types.Task{
					ID:          generateTaskID(),
					Title:       "Client Request",
					Description: m.textarea.Value(),
					CreatedAt:   time.Now(),
					Status:      types.StatusPending,
					CreatedBy:   types.RoleClient,
				}
				m.taskChannel <- task
				m.textarea.Reset()
				m.inputMode = false
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-10)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 10
		}

		m.textarea.SetWidth(msg.Width - 4)

	case EventMsg:
		m.events = append(m.events, types.AgentEvent(msg))
		m.viewport.SetContent(m.renderEvents())
		m.viewport.GotoBottom()

	case TickMsg:
		return m, tickCmd()
	}

	if m.inputMode {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Width(m.width)

	sb.WriteString(headerStyle.Render("BuildBureau - Multi-Agent System"))
	sb.WriteString("\n\n")

	// Event log viewport
	sb.WriteString(m.viewport.View())
	sb.WriteString("\n\n")

	// Input area
	if m.inputMode {
		inputStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1)

		sb.WriteString(inputStyle.Render(m.textarea.View()))
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Faint(true).Render("Press Enter to submit, Esc to quit"))
	} else {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

		sb.WriteString(statusStyle.Render("🔄 Agents working... Press Esc to quit"))
	}

	return sb.String()
}

// renderEvents renders the event log
func (m Model) renderEvents() string {
	var sb strings.Builder

	for _, event := range m.events {
		sb.WriteString(m.formatEvent(event))
		sb.WriteString("\n")
	}

	return sb.String()
}

// formatEvent formats a single event for display
func (m Model) formatEvent(event types.AgentEvent) string {
	timestamp := event.Timestamp.Format("15:04:05")

	// Color coding based on agent role
	roleColor := m.getRoleColor(event.Agent)
	eventColor := m.getEventColor(event.Type)

	roleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(roleColor)).
		Bold(true)

	eventStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(eventColor))

	timeStyle := lipgloss.NewStyle().
		Faint(true)

	return fmt.Sprintf("%s [%s] %s: %s",
		timeStyle.Render(timestamp),
		roleStyle.Render(string(event.Agent)),
		eventStyle.Render(string(event.Type)),
		event.Message,
	)
}

// getRoleColor returns the color for an agent role
func (m Model) getRoleColor(role types.AgentRole) string {
	switch role {
	case types.RoleCEO:
		return "205" // Pink/magenta
	case types.RoleManager:
		return "214" // Orange
	case types.RoleLead:
		return "45" // Cyan
	case types.RoleEmployee:
		return "42" // Green
	case types.RoleSecretary:
		return "183" // Purple
	default:
		return "15" // White
	}
}

// getEventColor returns the color for an event type
func (m Model) getEventColor(eventType types.EventType) string {
	switch eventType {
	case types.EventTaskAssigned:
		return "226" // Yellow
	case types.EventTaskCompleted:
		return "46" // Bright green
	case types.EventTaskStarted:
		return "51" // Bright cyan
	case types.EventError:
		return "196" // Red
	case types.EventMessage:
		return "250" // Light gray
	default:
		return "15" // White
	}
}

// tickCmd returns a command that sends a tick message
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// generateTaskID generates a simple task ID
func generateTaskID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}

// AddEvent adds an event to the UI (thread-safe way to add events)
func AddEvent(event types.AgentEvent) tea.Cmd {
	return func() tea.Msg {
		return EventMsg(event)
	}
}
