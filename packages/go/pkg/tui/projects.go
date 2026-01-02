package tui

import (
	"embed"
	"encoding/json"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed projects.json
var projectsData embed.FS

type Project struct {
	Name         string `json:"name"`
	Year         string `json:"year"`
	Technologies string `json:"technologies"`
	Description  string `json:"description"`
}

func LoadProjects() []Project {
	data, err := projectsData.ReadFile("projects.json")
	if err != nil {
		log.Fatalf("Failed to read embedded file: %s", err)
	}
	var p []Project
	if err := json.Unmarshal(data, &p); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
	return p
}

func (m model) ProjectsUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			if m.state.account.selected < len(m.projects)-1 {
				m.state.account.selected++
			}
		case "shift+tab", "up", "k":
			if m.state.account.selected > 0 {
				m.state.account.selected--
			}
		}
	}
	return m, nil
}

func (m model) getProjectsContent(totalWidth int) string {
	if len(m.projects) == 0 {
		return "No projects available"
	}

	project := m.projects[m.state.account.selected]

	var content strings.Builder
	content.WriteString(m.theme.TextAccent().Render(wordWrap(project.Name, totalWidth)) + "\n")
	content.WriteString(m.theme.Base().Render(project.Year) + "\n\n")
	content.WriteString(m.theme.TextHighlight().Render("Technologies:") + "\n")
	content.WriteString(m.theme.Base().Render(wordWrap(project.Technologies, totalWidth)) + "\n\n")
	content.WriteString(m.theme.TextHighlight().Render("Details:") + "\n")

	// Preserve bullet lines by wrapping each line separately and add a blank line between bullets
	descLines := strings.Split(project.Description, "\n")
	for i, line := range descLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		descLines[i] = m.theme.Base().Render(wordWrap(trimmed, totalWidth))
	}
	// Join with double newline for spacing between items
	content.WriteString(strings.Join(descLines, "\n\n"))

	return content.String()
}

func (m model) getProjectsMenuContent() string {
	var menuItem lipgloss.Style
	var highlightedMenuItem lipgloss.Style

	if m.size < large {
		menuItem = m.theme.Base().
			Width(m.widthContent - 1).
			Align(lipgloss.Center)
		highlightedMenuItem = m.theme.Base().
			Width(m.widthContent - 1).
			Align(lipgloss.Center).
			Background(m.theme.Highlight()).
			Foreground(m.theme.Accent())
	} else {
		// Use the actual menu viewport width
		menuWidth := m.state.account.menuViewport.Width
		if menuWidth == 0 {
			menuWidth = 20 // fallback
		}
		menuItem = m.theme.Base().
			Width(menuWidth).
			Padding(0, 1)
		highlightedMenuItem = m.theme.Base().
			Width(menuWidth).
			Padding(0, 1).
			Background(m.theme.Highlight()).
			Foreground(m.theme.Accent())
	}

	var content strings.Builder
	for i, p := range m.projects {
		var item string
		if i == m.state.account.selected {
			item = highlightedMenuItem.Render(wordWrap(p.Name+" ("+p.Year+")", m.state.account.menuViewport.Width-2))
		} else {
			item = menuItem.Render(wordWrap(p.Name+" ("+p.Year+")", m.state.account.menuViewport.Width-2))
		}
		content.WriteString(item + "\n")
		// Add spacing between menu items
		if i < len(m.projects)-1 {
			content.WriteString("\n")
		}
	}

	return content.String()
}
