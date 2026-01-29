package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kpango/BuildBureau/internal/agent"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

// Model represents the UI model
type Model struct {
	textarea        textarea.Model
	spinner         spinner.Model
	messages        []string
	agentStatuses   []agent.Status
	projectName     string
	projectStatus   string
	currentPhase    string
	err             error
	ready           bool
	width           int
	height          int
	lastUpdate      time.Time
}

// NewModel creates a new UI model
func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Enter your project requirements..."
	ta.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		textarea:      ta,
		spinner:       s,
		messages:      make([]string, 0),
		agentStatuses: make([]agent.Status, 0),
		projectStatus: "待機中",
		lastUpdate:    time.Now(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.spinner.Tick,
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if msg.Alt {
				// Submit the project requirements
				return m, m.submitRequirements()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
		m.ready = true

	case statusUpdateMsg:
		m.agentStatuses = msg.statuses
		m.lastUpdate = time.Now()

	case projectUpdateMsg:
		m.projectName = msg.name
		m.projectStatus = msg.status
		m.currentPhase = msg.phase

	case messageMsg:
		m.messages = append(m.messages, msg.text)
		if len(m.messages) > 10 {
			m.messages = m.messages[len(m.messages)-10:]
		}

	case errMsg:
		m.err = msg.err
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🏢 BuildBureau - Multi-Layer AI Agent System"))
	b.WriteString("\n\n")

	// Project info
	if m.projectName != "" {
		b.WriteString(fmt.Sprintf("Project: %s\n", m.projectName))
		b.WriteString(fmt.Sprintf("Status: %s\n", m.projectStatus))
		if m.currentPhase != "" {
			b.WriteString(fmt.Sprintf("Current Phase: %s\n", m.currentPhase))
		}
		b.WriteString("\n")
	}

	// Agent statuses
	if len(m.agentStatuses) > 0 {
		b.WriteString("Agent Status:\n")
		for _, status := range m.agentStatuses {
			statusIcon := "⚪"
			if status.State == "working" {
				statusIcon = m.spinner.View()
			} else if status.State == "completed" {
				statusIcon = "✅"
			} else if status.State == "error" {
				statusIcon = "❌"
			}
			b.WriteString(fmt.Sprintf("  %s %s (%s): %s\n", 
				statusIcon, status.AgentID, status.AgentType, status.Message))
		}
		b.WriteString("\n")
	}

	// Recent messages
	if len(m.messages) > 0 {
		b.WriteString("Recent Messages:\n")
		for _, msg := range m.messages {
			b.WriteString(fmt.Sprintf("  • %s\n", msg))
		}
		b.WriteString("\n")
	}

	// Error display
	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	// Input area
	b.WriteString("Requirements Input:\n")
	b.WriteString(m.textarea.View())
	b.WriteString("\n\n")
	b.WriteString(infoStyle.Render("Alt+Enter: Submit | Esc: Exit"))

	return b.String()
}

// submitRequirements submits the project requirements
func (m Model) submitRequirements() tea.Cmd {
	return func() tea.Msg {
		// This would be implemented to actually submit to the agent system
		return messageMsg{text: "Project requirements submitted"}
	}
}

// Message types for updates
type statusUpdateMsg struct {
	statuses []agent.Status
}

type projectUpdateMsg struct {
	name   string
	status string
	phase  string
}

type messageMsg struct {
	text string
}

type errMsg struct {
	err error
}

// UpdateAgentStatuses updates agent statuses in the UI
func UpdateAgentStatuses(statuses []agent.Status) tea.Cmd {
	return func() tea.Msg {
		return statusUpdateMsg{statuses: statuses}
	}
}

// UpdateProject updates project information in the UI
func UpdateProject(name, status, phase string) tea.Cmd {
	return func() tea.Msg {
		return projectUpdateMsg{name: name, status: status, phase: phase}
	}
}

// AddMessage adds a message to the UI
func AddMessage(text string) tea.Cmd {
	return func() tea.Msg {
		return messageMsg{text: text}
	}
}

// ShowError shows an error in the UI
func ShowError(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg{err: err}
	}
}
