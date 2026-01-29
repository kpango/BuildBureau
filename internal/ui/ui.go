package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kpango/BuildBureau/internal/config"
)

// Model represents the UI model
type Model struct {
	config       *config.Config
	viewport     viewport.Model
	textarea     textarea.Model
	messages     []Message
	ready        bool
	width        int
	height       int
	inputMode    bool
	onSubmit     func(string) error
}

// Message represents a message in the conversation
type Message struct {
	Timestamp time.Time
	Role      string
	Agent     string
	Content   string
	Level     int // Hierarchy level for indentation
}

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			PaddingLeft(2)

	ceoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	deptHeadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true)

	managerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#45B7D1")).
			Bold(true)

	workerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#96CEB4"))

	secretaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDA15E")).
			Italic(true)

	timestampStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)
)

// NewModel creates a new UI model
func NewModel(cfg *config.Config, onSubmit func(string) error) Model {
	ta := textarea.New()
	ta.Placeholder = "Enter client request..."
	ta.Focus()
	ta.CharLimit = 2000
	ta.SetWidth(80)
	ta.SetHeight(3)

	vp := viewport.New(80, 20)

	return Model{
		config:    cfg,
		textarea:  ta,
		viewport:  vp,
		messages:  []Message{},
		inputMode: true,
		onSubmit:  onSubmit,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textarea.Blink
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
			if m.inputMode && !m.textarea.Focused() {
				m.textarea.Focus()
			} else if m.inputMode {
				// Submit the input
				input := m.textarea.Value()
				if input != "" && m.onSubmit != nil {
					m.AddMessage(Message{
						Timestamp: time.Now(),
						Role:      "User",
						Agent:     "You",
						Content:   input,
						Level:     0,
					})
					m.textarea.Reset()
					m.inputMode = false
					go m.onSubmit(input)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 3
		footerHeight := 5
		if m.inputMode {
			footerHeight = 7
		}

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight - footerHeight

		m.textarea.SetWidth(msg.Width - 4)

		if !m.ready {
			m.ready = true
		}

	case AddMessageMsg:
		m.AddMessage(Message(msg))
		m.viewport.GotoBottom()

	case ProcessingCompleteMsg:
		m.inputMode = true
		m.textarea.Focus()
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	// Update textarea if in input mode
	if m.inputMode {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "\nInitializing..."
	}

	// Header
	header := titleStyle.Render("🏢 BuildBureau - Multi-Agent AI System")

	// Build message view
	content := m.renderMessages()
	m.viewport.SetContent(content)

	// Footer with input
	var footer string
	if m.inputMode {
		footer = fmt.Sprintf("\n%s\n%s\n%s",
			promptStyle.Render("Enter your request:"),
			m.textarea.View(),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("Press Enter to submit • Ctrl+C to quit"),
		)
	} else {
		footer = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("Processing... Please wait.")
	}

	return fmt.Sprintf("%s\n\n%s\n%s", header, m.viewport.View(), footer)
}

// renderMessages renders all messages with proper styling
func (m Model) renderMessages() string {
	if len(m.messages) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			Render("Waiting for client request...")
	}

	var lines []string
	for _, msg := range m.messages {
		lines = append(lines, m.renderMessage(msg))
	}

	return strings.Join(lines, "\n\n")
}

// renderMessage renders a single message with appropriate styling
func (m Model) renderMessage(msg Message) string {
	// Indentation based on hierarchy level
	indent := strings.Repeat("  ", msg.Level)

	// Choose style based on role
	var roleStyle lipgloss.Style
	switch msg.Role {
	case "CEO", "CEOSecretary":
		if strings.Contains(msg.Role, "Secretary") {
			roleStyle = secretaryStyle
		} else {
			roleStyle = ceoStyle
		}
	case "DeptHead", "DeptHeadSecretary":
		if strings.Contains(msg.Role, "Secretary") {
			roleStyle = secretaryStyle
		} else {
			roleStyle = deptHeadStyle
		}
	case "Manager", "ManagerSecretary":
		if strings.Contains(msg.Role, "Secretary") {
			roleStyle = secretaryStyle
		} else {
			roleStyle = managerStyle
		}
	case "Worker":
		roleStyle = workerStyle
	default:
		roleStyle = lipgloss.NewStyle().Bold(true)
	}

	timestamp := timestampStyle.Render(msg.Timestamp.Format("15:04:05"))
	agentName := roleStyle.Render(msg.Agent)
	
	// Wrap content if too long
	content := msg.Content
	if len(content) > 100 {
		content = content[:100] + "..."
	}

	return fmt.Sprintf("%s%s [%s]: %s", indent, timestamp, agentName, content)
}

// AddMessage adds a message to the UI
func (m *Model) AddMessage(msg Message) {
	m.messages = append(m.messages, msg)
	
	// Limit history
	maxLines := m.config.System.UI.MaxHistoryLines
	if len(m.messages) > maxLines {
		m.messages = m.messages[len(m.messages)-maxLines:]
	}
}

// AddMessageMsg is a message type for adding messages to the UI
type AddMessageMsg Message

// ProcessingCompleteMsg signals that processing is complete
type ProcessingCompleteMsg struct{}

// AddAgentMessage adds an agent message to the UI from outside
func AddAgentMessage(role, agent, content string, level int) tea.Cmd {
	return func() tea.Msg {
		return AddMessageMsg{
			Timestamp: time.Now(),
			Role:      role,
			Agent:     agent,
			Content:   content,
			Level:     level,
		}
	}
}

// SignalProcessingComplete signals that processing is complete
func SignalProcessingComplete() tea.Cmd {
	return func() tea.Msg {
		return ProcessingCompleteMsg{}
	}
}
