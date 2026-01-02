package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type accountState struct {
	selected       int
	focused        bool
	menuViewport   viewport.Model
	detailViewport viewport.Model
	viewportsReady bool
}

func (m model) updateAccountViewports() model {
	headerHeight := lipgloss.Height(m.HeaderView())
	footerHeight := lipgloss.Height(m.FooterView())
	verticalMarginHeight := headerHeight + footerHeight

	availableHeight := m.heightContainer - verticalMarginHeight

	var menuWidth, detailWidth int

	if m.size < large {
		// For small screens, full width for both
		menuWidth = m.widthContent
		detailWidth = m.widthContent
	} else {
		// For large screens, split the available width
		menuWidth = 20
		spacer := 2 // for "  " separator
		detailWidth = m.widthContent - menuWidth - spacer
		if detailWidth < 60 {
			detailWidth = 60
		}
	}

	if !m.state.account.viewportsReady {
		// Initialize viewports for the first time
		m.state.account.menuViewport = viewport.New(menuWidth, availableHeight)
		m.state.account.menuViewport.KeyMap = viewport.KeyMap{}
		m.state.account.detailViewport = viewport.New(detailWidth, availableHeight)
		m.state.account.detailViewport.KeyMap = modifiedKeyMap

		m.state.account.viewportsReady = true
	} else {
		// Update existing viewports
		m.state.account.menuViewport.Width = menuWidth
		m.state.account.menuViewport.Height = availableHeight

		m.state.account.detailViewport.Width = detailWidth
		m.state.account.detailViewport.Height = availableHeight
	}

	return m
}

func (m model) AccountSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(accountPage)
	m.state.account.selected = 0
	m.state.account.focused = false

	m.state.footer.commands = []footerCommand{
		{key: "↑/↓", value: "navigate"},
		{key: "PgUp/PgDn", value: "scroll"},
	}

	m = m.updateAccountViewports()
	m.state.account.menuViewport.GotoTop()
	m.state.account.detailViewport.GotoTop()
	return m, nil
}

func (m model) AccountUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Update viewport dimensions if window size changed
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		m = m.updateAccountViewports()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			if m.state.account.selected < len(m.projects)-1 {
				m.state.account.selected++
				// Scroll menu to keep selected item visible
				m.state.account.menuViewport.LineDown(2)
			}
		case "shift+tab", "up", "k":
			if m.state.account.selected > 0 {
				m.state.account.selected--
				// Scroll menu to keep selected item visible
				m.state.account.menuViewport.LineUp(2)
			}
		case "pgup":
			m.state.account.detailViewport.LineUp(5)
		case "pgdown":
			m.state.account.detailViewport.LineDown(5)
		}
	}

	// Update viewports with new content
	if m.state.account.viewportsReady {
		menuContent := m.getProjectsMenuContent()
		m.state.account.menuViewport.SetContent(menuContent)

		detailContent := m.getProjectsContent(m.state.account.detailViewport.Width - 4)
		m.state.account.detailViewport.SetContent(detailContent)
	}

	// Only update detail viewport for scroll; menu stays static
	m.state.account.detailViewport, cmd = m.state.account.detailViewport.Update(msg)
	cmds = append(cmds, cmd)

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m model) AccountView() string {
	if !m.state.account.viewportsReady {
		m = m.updateAccountViewports()
	}

	// Always show projects menu on the left
	menuContent := m.getProjectsMenuContent()
	m.state.account.menuViewport.SetContent(menuContent)

	// Show selected project details on the right
	detailContent := m.getProjectsContent(m.state.account.detailViewport.Width - 4)
	m.state.account.detailViewport.SetContent(detailContent)

	// Combine viewport views
	if m.size < large {
		// For small screens, stack the viewports vertically
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.state.account.menuViewport.View(),
			m.state.account.detailViewport.View(),
		)
	} else {
		// For large screens, place viewports side by side
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.state.account.menuViewport.View(),
			"  ",
			m.state.account.detailViewport.View(),
		)
	}
}
