package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TodoFormModel struct {
	title textinput.Model
}

func NewTodoInput() *TodoFormModel {
	todoForm := TodoFormModel{}
	todoForm.title = textinput.New()
	todoForm.title.Focus()
	return &todoForm
}

func (m TodoFormModel) Init() tea.Cmd { return nil }

func (m TodoFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, nil
		}
	}
	return m, nil
}

func (m TodoFormModel) View() string {
	return "Under Construction"
}
