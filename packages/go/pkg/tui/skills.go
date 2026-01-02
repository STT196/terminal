package tui

import (
	"embed"
	"encoding/json"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed skills.json
var skillsData embed.FS

func LoadSkills() []string {
	data, err := skillsData.ReadFile("skills.json")
	if err != nil {
		log.Fatalf("Failed to read embedded file: %s", err)
	}
	var skills []string
	if err := json.Unmarshal(data, &skills); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
	return skills
}

func (m model) SkillsSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(skillsPage)
	m.state.footer.commands = []footerCommand{
		{key: "c", value: "cart"},
	}
	return m, nil
}

func (m model) SkillsView(totalWidth int) string {
	title := m.theme.TextAccent().Bold(true).Render("Skills & Technologies")

	// Display each skill on a separate line with bullet point
	var skillLines []string
	for _, skill := range m.skills {
		skillLines = append(skillLines, m.theme.Base().Render("• "+skill))
	}

	return lipgloss.Place(
		m.widthContent,
		m.heightContent,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			strings.Join(skillLines, "\n"),
		),
	)
}
