package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kpango/BuildBureau/internal/agent"
)

const (
	// Default UI dimensions.
	defaultWidth          = 80
	defaultHeight         = 20
	defaultTextareaHeight = 3
	defaultCharLimit      = 1000
)

var (
	//nolint:gochecknoglobals // TUI styles are package-level configuration
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginLeft(2)

	//nolint:gochecknoglobals // TUI styles are package-level configuration
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)

	//nolint:gochecknoglobals // TUI styles are package-level configuration
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			MarginLeft(2)

	//nolint:gochecknoglobals // TUI styles are package-level configuration
	outputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1).
			MarginLeft(2).
			MarginTop(1)
)

type Model struct {
	textarea   textarea.Model
	err        error
	org        *agent.Organization
	output     string
	viewport   viewport.Model
	width      int
	height     int
	ready      bool
	processing bool
}

func NewModel(org *agent.Organization) Model {
	ta := textarea.New()
	ta.Placeholder = "Enter your task or instruction..."
	ta.Focus()
	ta.CharLimit = defaultCharLimit
	ta.SetWidth(defaultWidth)
	ta.SetHeight(defaultTextareaHeight)

	vp := viewport.New(defaultWidth, defaultHeight)
	vp.SetContent("Welcome to BuildBureau!\n\nEnter your task and press Ctrl+S to submit.\nPress Ctrl+C or Esc to quit.")

	return Model{
		org:      org,
		textarea: ta,
		viewport: vp,
		output:   vp.View(),
		ready:    true,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

type taskResultMsg struct {
	err    error
	result string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		//nolint:exhaustive // Key handling intentionally only covers specific cases
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlS:
			if !m.processing && m.textarea.Value() != "" {
				m.processing = true
				instruction := m.textarea.Value()
				m.textarea.Reset()

				// Process task asynchronously
				return m, func() tea.Msg {
					ctx := context.Background()
					response, err := m.org.ProcessClientTask(ctx, instruction)
					if err != nil {
						return taskResultMsg{err: err}
					}
					return taskResultMsg{result: response.Result}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update viewport and textarea sizes
		headerHeight := 3
		footerHeight := 8
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - headerHeight - footerHeight
		m.textarea.SetWidth(msg.Width - 6)

	case taskResultMsg:
		m.processing = false
		if msg.err != nil {
			m.output = fmt.Sprintf("Error: %v\n\n%s", msg.err, m.output)
		} else {
			m.output = fmt.Sprintf("=== Task Result ===\n%s\n\n%s", msg.result, m.output)
		}
		m.viewport.SetContent(m.output)
		m.viewport.GotoTop()
	}

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(taCmd, vpCmd)
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🏢 BuildBureau - Multi-Agent Development System"))
	b.WriteString("\n\n")

	// Output viewport
	b.WriteString(outputStyle.Render(m.viewport.View()))
	b.WriteString("\n\n")

	// Input area
	b.WriteString(inputStyle.Render(m.textarea.View()))
	b.WriteString("\n")

	// Help text
	status := ""
	if m.processing {
		status = " [Processing...]"
	}
	b.WriteString(helpStyle.Render(fmt.Sprintf("Ctrl+S: Submit | Ctrl+C/Esc: Quit%s", status)))

	return b.String()
}
