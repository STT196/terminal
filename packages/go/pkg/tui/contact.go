package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) ContactSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(contactPage)
	m.state.footer.commands = []footerCommand{
	}
	return m, nil
}

func (m model) ContactUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) ContactView(totalWidth int) string {
	base := m.theme.Base().Width(m.widthContent).Render
	accent := m.theme.TextAccent().Render
	bold := m.theme.TextAccent().Bold(true).Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		bold("Contact Me"),
		"",
		base("Feel free to reach out:"),
		"",
		base("GitHub: ")+accent("https://github.com/STT196")+m.CursorView(),
		"",
		base("Email: ")+accent("thisaratharindaciscoitn@gmail.com"),
	)
}
