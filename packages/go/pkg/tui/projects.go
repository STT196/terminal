package tui

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminaldotshop/terminal/go/pkg/resource"
)

//go:embed projects.json
var fallbackProjectsData embed.FS

type Project struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Year         string `json:"year"`
	Technologies string `json:"technologies"`
	Description  string `json:"description"`
	Order        *int   `json:"order,omitempty"`
}

type projectsResponse struct {
	Data []Project `json:"data"`
}

func LoadProjects() []Project {
	// 1. Try the API first
	apiURL := resource.Resource.Api.Url + "/project"
	resp, err := http.Get(apiURL)
	if err == nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			var result projectsResponse
			if err := json.Unmarshal(body, &result); err == nil && len(result.Data) > 0 {
				return result.Data
			}
		}
	}

	// 2. Try external file (configurable via PROJECTS_FILE env var)
	projectsFile := os.Getenv("PROJECTS_FILE")
	if projectsFile == "" {
		projectsFile = "/data/projects.json"
	}
	if data, err := os.ReadFile(projectsFile); err == nil {
		var p []Project
		if err := json.Unmarshal(data, &p); err == nil && len(p) > 0 {
			return p
		}
	}

	// 3. Fall back to embedded projects.json
	return loadFallbackProjects()
}

func loadFallbackProjects() []Project {
	data, err := fallbackProjectsData.ReadFile("projects.json")
	if err != nil {
		return []Project{}
	}
	var p []Project
	if err := json.Unmarshal(data, &p); err != nil {
		return []Project{}
	}
	return p
}

func (m model) ProjectsUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			if m.state.project.selected < len(m.projects)-1 {
				m.state.project.selected++
			}
		case "shift+tab", "up", "k":
			if m.state.project.selected > 0 {
				m.state.project.selected--
			}
		}
	}
	return m, nil
}

func (m model) getProjectsContent(totalWidth int) string {
	if len(m.projects) == 0 {
		return "No projects available"
	}

	project := m.projects[m.state.project.selected]

	var content strings.Builder
	content.WriteString(m.theme.TextAccent().Render(wordWrap(project.Name, totalWidth)) + "\n")
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
		menuWidth := m.state.project.menuViewport.Width
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
		if i == m.state.project.selected {
			item = highlightedMenuItem.Render(wordWrap(p.Name, m.state.project.menuViewport.Width-2))
		} else {
			item = menuItem.Render(wordWrap(p.Name, m.state.project.menuViewport.Width-2))
		}
		content.WriteString(item + "\n")
		// Add spacing between menu items
		if i < len(m.projects)-1 {
			content.WriteString("\n")
		}
	}

	return content.String()
}
