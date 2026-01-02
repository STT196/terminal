package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	terminal "github.com/terminaldotshop/terminal-sdk-go"
)

type footerState struct {
	commands []footerCommand
}

type footerCommand struct {
	key   string
	value string
}

// wordWrap breaks a string into multiple lines to fit within maxWidth
func wordWrap(text string, maxWidth int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	lines := []string{}
	currentLine := words[0]

	for _, word := range words[1:] {
		// Check if adding this word would exceed the width
		testLine := currentLine + " " + word
		if lipgloss.Width(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			// Line would be too long, start a new line
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	// Add the last line
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

// ToggleRegion switches between regions and creates a new client with the updated region header
func (m model) ToggleRegion() (model, tea.Cmd) {
	// Toggle between "na", "eu", and "global"
	var newRegion terminal.Region

	if m.region == nil || *m.region == terminal.RegionNa {
		newRegion = terminal.RegionEu
	} else if *m.region == terminal.RegionEu {
		newRegion = terminal.RegionGlobal
	} else {
		newRegion = terminal.RegionNa
	}

	// Update the model's region
	m.region = &newRegion

	// Create new client with updated region
	m.client = m.CreateSDKClient()

	// Return command to reload data
	cmd := func() tea.Msg {
		_, err := m.client.Cart.Clear(m.context)
		if err != nil {
			return err
		}

		response, err := m.client.View.Init(m.context)
		if err != nil {
			return err
		}
		return response.Data
	}

	return m, cmd
}

func (m model) FooterView() string {
	bold := m.theme.TextAccent().Bold(true).Render
	base := m.theme.Base().Render

	table := m.theme.Base().
		Width(m.widthContent).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.theme.Border()).
		PaddingBottom(1).
		Align(lipgloss.Center)

	if m.size == small {
		return table.Render(bold("m") + base(" menu"))
	}

	// Add commands
	commands := []string{}
	for _, cmd := range m.state.footer.commands {
		commands = append(commands, bold(" "+cmd.key+" ")+base(cmd.value+"  "))
	}

	lines := []string{}
	lines = append(lines, commands...)
	lines = append(lines, base("  "))
	lines = append(lines, base("powered by terminal.shop"))

	var content = "STT196"
	footer := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		content,
		table.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				lines...,
			),
		))

	// Add the region selector and the rest of the commands
	return lipgloss.Place(
		m.widthContainer,
		lipgloss.Height(footer),
		lipgloss.Center,
		lipgloss.Center,
		footer,
	)
}
