package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"buildbureau/internal/agents"
	"buildbureau/internal/protocol"
	"buildbureau/pkg/a2a"
)

type State int

const (
	StateInput State = iota
	StateRunning
	StateDone
	StateError
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))
)

type LogMsg a2a.Message
type ResultMsg protocol.ProjectSummary
type ErrorMsg error

type Model struct {
	state      State
	textInput  textinput.Model
	viewport   viewport.Model
	system     *agents.System
	sub        <-chan a2a.Message
	logs       []string
	err        error
	result     protocol.ProjectSummary
	projectReq string
}

func NewModel(sys *agents.System, sub <-chan a2a.Message) Model {
	ti := textinput.New()
	ti.Placeholder = "Describe your project (e.g., 'Build a Todo App in Go')"
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 50

	vp := viewport.New(80, 20)
	vp.SetContent("Waiting for input...")

	return Model{
		state:     StateInput,
		textInput: ti,
		viewport:  vp,
		system:    sys,
		sub:       sub,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == StateInput {
				m.projectReq = m.textInput.Value()
				m.state = StateRunning
				m.logs = append(m.logs, "Starting project: "+m.projectReq)
				m.viewport.SetContent(strings.Join(m.logs, "\n"))

				// Start the project in a goroutine and wait for logs/result
				cmds = append(cmds, m.runProjectCmd(), m.waitForLogCmd())
			}
		}

	case LogMsg:
		if m.state == StateRunning {
			line := fmt.Sprintf("[%s] %s -> %s: %v", msg.Timestamp.Format("15:04:05"), msg.From, msg.Type, msg.Payload)
			m.logs = append(m.logs, line)
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			m.viewport.GotoBottom()
			// Continue waiting for logs
			cmds = append(cmds, m.waitForLogCmd())
		}

	case ResultMsg:
		m.state = StateDone
		m.result = protocol.ProjectSummary(msg)
		m.logs = append(m.logs, "\n\n=== PROJECT COMPLETED ===\nSuccess: "+fmt.Sprintf("%v", m.result.Success))
		for f, c := range m.result.AllArtifacts {
			m.logs = append(m.logs, fmt.Sprintf("File: %s\n%s\n---", f, c))
		}
		m.viewport.SetContent(strings.Join(m.logs, "\n"))
		m.viewport.GotoBottom()

	case ErrorMsg:
		m.state = StateError
		m.err = msg
		m.logs = append(m.logs, fmt.Sprintf("\nERROR: %v", m.err))
		m.viewport.SetContent(strings.Join(m.logs, "\n"))
	}

	// Update components
	switch m.state {
	case StateInput:
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	case StateRunning, StateDone, StateError:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case StateInput:
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			titleStyle.Render(" BuildBureau "),
			"Please enter your project requirement:",
			m.textInput.View(),
		)
	case StateRunning:
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			titleStyle.Render(" BuildBureau - Running "),
			m.viewport.View(),
			infoStyle.Render("Processing... (Ctrl+C to quit)"),
		)
	case StateDone:
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			titleStyle.Render(" BuildBureau - Done "),
			m.viewport.View(),
			infoStyle.Render("Project Complete! (Ctrl+C to quit)"),
		)
	case StateError:
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			titleStyle.Render(" BuildBureau - Error "),
			m.viewport.View(),
			infoStyle.Render("An error occurred. (Ctrl+C to quit)"),
		)
	}
	return ""
}

// runProjectCmd triggers the system logic
func (m Model) runProjectCmd() tea.Cmd {
	return func() tea.Msg {
		req := protocol.RequirementSpec{
			ProjectName: "UserProject", // Could ask for this too
			Details:     m.projectReq,
		}
		res, err := m.system.RunProject(context.Background(), req)
		if err != nil {
			return ErrorMsg(err)
		}
		return ResultMsg(res)
	}
}

// waitForLogCmd listens for one message from the bus
func (m Model) waitForLogCmd() tea.Cmd {
	return func() tea.Msg {
		msg := <-m.sub
		return LogMsg(msg)
	}
}
