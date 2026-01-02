package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminaldotshop/terminal/go/pkg/tui/theme"
)

type shopState struct {
	selected int
}

func (m model) ShopSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(shopPage)
	m.theme = theme.BasicTheme(m.renderer, nil)
	return m, nil
}

func (m model) ShopUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) ShopView() string {
	return lipgloss.Place(
		m.widthContent,
		m.heightContent,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.theme.TextAccent().Render("Hi! I'm Thisara Tharinda"),
			m.theme.Base().Render("DevOps Engineer"),
		),
	)
}

func (m model) UpdateSelectedTheme() model {
	m.theme = theme.BasicTheme(m.renderer, nil)
	return m
}
