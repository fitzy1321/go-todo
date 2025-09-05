package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type EntryFormModel struct {
	text textinput.Model
}

/* tea.Model Interface: Init, Update, View*/

func (m EntryFormModel) Init() tea.Cmd { return nil }

func (m EntryFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmdq
	var cmd tea.Cmd
	m.text, cmd = m.text.Update(msg)
	return m, cmd
}

func (m EntryFormModel) View() string {
	return fmt.Sprintf("Enter Title for new Todo Item\n%s", m.text.View())
}

/* Public Functions */

func NewEntryForm() EntryFormModel {
	todoForm := EntryFormModel{}
	todoForm.text = textinput.New()
	return todoForm
}
func (m *EntryFormModel) Blur() {
	m.text.Blur()
}

func (m *EntryFormModel) Focus() {
	m.text.Focus()
}

func (m *EntryFormModel) Value() string {
	return m.text.Value()
}
