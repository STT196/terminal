package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SplashState struct {
	data  bool
	delay bool
}

type splashDoneMsg struct{}

// For the simplified site, show logo briefly then go to shop.
func (m model) SplashInit() tea.Cmd {
	return tea.Batch(
		m.CursorInit(),
		func() tea.Msg { return tea.DisableMouse() },
		tea.Tick(8000*time.Millisecond, func(time.Time) tea.Msg { return splashDoneMsg{} }),
	)
}

func (m model) SplashUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg.(type) {
	case splashDoneMsg:
		return m.ShopSwitch()
	}
	return m, nil
}

func (m model) SplashView() string {
	var msg string
	if m.error != nil {
		msg = m.error.message
	} else {
		msg = ""
	}

	var hint string
	if m.error != nil {
		hint = lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.theme.TextAccent().Bold(true).Render("esc"),
			" ",
			"quit",
		)
	} else {
		hint = ""
	}

	if m.error == nil {
		return lipgloss.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			lipgloss.Center,
			m.LogoView(),
		)
	}

	return lipgloss.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			"",
			"",
			"",
			m.LogoView(),
			"",
			"",
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.theme.TextError().Render(msg),
			),
			hint,
		),
	)
}
