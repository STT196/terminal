package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type projectState struct {
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

	if !m.state.project.viewportsReady {
		// Initialize viewports for the first time
		m.state.project.menuViewport = viewport.New(menuWidth, availableHeight)
		m.state.project.menuViewport.KeyMap = viewport.KeyMap{}
		m.state.project.detailViewport = viewport.New(detailWidth, availableHeight)
		m.state.project.detailViewport.KeyMap = modifiedKeyMap

		m.state.project.viewportsReady = true
	} else {
		// Update existing viewports
		m.state.project.menuViewport.Width = menuWidth
		m.state.project.menuViewport.Height = availableHeight

		m.state.project.detailViewport.Width = detailWidth
		m.state.project.detailViewport.Height = availableHeight
	}

	return m
}

func (m model) ProjectSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(projectsPage)
	m.state.project.selected = 0
	m.state.project.focused = false

	m.state.footer.commands = []footerCommand{
		{key: "↑/↓", value: "navigate"},
		{key: "PgUp/PgDn", value: "scroll"},
	}

	m = m.updateAccountViewports()
	m.state.project.menuViewport.GotoTop()
	m.state.project.detailViewport.GotoTop()
	return m, nil
}

func (m model) ProjectUpdate(msg tea.Msg) (model, tea.Cmd) {
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
			if m.state.project.selected < len(m.projects)-1 {
				m.state.project.selected++
				// Scroll menu to keep selected item visible
				m.state.project.menuViewport.LineDown(2)
			}
		case "shift+tab", "up", "k":
			if m.state.project.selected > 0 {
				m.state.project.selected--
				// Scroll menu to keep selected item visible
				m.state.project.menuViewport.LineUp(2)
			}
		case "pgup":
			m.state.project.detailViewport.LineUp(5)
		case "pgdown":
			m.state.project.detailViewport.LineDown(5)
		}
	}

	// Update viewports with new content
	if m.state.project.viewportsReady {
		menuContent := m.getProjectsMenuContent()
		m.state.project.menuViewport.SetContent(menuContent)

		detailContent := m.getProjectsContent(m.state.project.detailViewport.Width - 4)
		m.state.project.detailViewport.SetContent(detailContent)
	}

	// Only update detail viewport for scroll; menu stays static
	m.state.project.detailViewport, cmd = m.state.project.detailViewport.Update(msg)
	cmds = append(cmds, cmd)

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m model) ProjectView() string {
	if !m.state.project.viewportsReady {
		m = m.updateAccountViewports()
	}

	// Always show projects menu on the left
	menuContent := m.getProjectsMenuContent()
	m.state.project.menuViewport.SetContent(menuContent)

	// Show selected project details on the right
	detailContent := m.getProjectsContent(m.state.project.detailViewport.Width - 4)
	m.state.project.detailViewport.SetContent(detailContent)

	// Combine viewport views
	if m.size < large {
		// For small screens, stack the viewports vertically
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.state.project.menuViewport.View(),
			m.state.project.detailViewport.View(),
		)
	} else {
		// For large screens, place viewports side by side
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.state.project.menuViewport.View(),
			"  ",
			m.state.project.detailViewport.View(),
		)
	}
}
