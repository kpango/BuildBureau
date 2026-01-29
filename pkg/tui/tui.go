package tui

import (
	"context"
	"fmt"
	"strings"

	"buildbureau/pkg/protocol"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	presidentClient protocol.AgentServiceClient
	textInput       textinput.Model
	viewport        viewport.Model
	messages        []string
	err             error
}

func InitialModel(presidentClient protocol.AgentServiceClient) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter instructions for the President..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	vp := viewport.New(100, 20)
	vp.SetContent("Welcome to BuildBureau.\n")

	return Model{
		presidentClient: presidentClient,
		textInput:       ti,
		viewport:        vp,
		messages:        []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input != "" {
				m.messages = append(m.messages, fmt.Sprintf("> %s", input))
				m.textInput.Reset()

				// Send to President
				return m, func() tea.Msg {
					resp, err := m.presidentClient.AssignTask(context.Background(), &protocol.Task{
						ID:          fmt.Sprintf("task-%d", len(m.messages)),
						Description: input,
						AssignedBy:  "Client",
					})
					if err != nil {
						return errMsg(err)
					}
					return responseMsg(resp)
				}
			}
		}

	case errMsg:
		m.messages = append(m.messages, fmt.Sprintf("Error: %v", msg))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))

	case responseMsg:
		m.messages = append(m.messages, fmt.Sprintf("System: Task %s status: %s", msg.TaskID, msg.Status))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n(esc to quit)",
		m.viewport.View(),
		m.textInput.View(),
	)
}

type errMsg error

type responseMsg *protocol.TaskResponse

func Start(presidentClient protocol.AgentServiceClient) error {
	p := tea.NewProgram(InitialModel(presidentClient))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
