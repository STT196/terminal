package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminaldotshop/terminal/go/pkg/tui/theme"
)

type aboutState struct {
	selected int
}

const aboutASCIIArt = `
   ▓▓▓▓▓ ▓▓▓▓▓▓ ▓▓▓▓▓▓                  
  ▓        ▓▓     ▓▓                    
  ▓        ▓▓     ▓▓                    
  ▓▓▓▓▓    ▓▓     ▓▓                    
      ▓    ▓▓     ▓▓                    
      ▓    ▓▓     ▓▓                    
  ▓▓▓▓     ▓▓     ▓▓                    
                                                  
`

const aboutHeroASCII = ``

func (m model) AboutSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(aboutPage)
	m.theme = theme.BasicTheme(m.renderer, nil)
	m.state.footer.commands = []footerCommand{}
	return m, nil
}

func (m model) AboutUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) AboutView() string {
	accent := m.theme.TextAccent()
	brand := m.theme.TextBrand()
	body := m.theme.TextBody()

	prompt := func(command string, text string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,

			accent.Render(command),
			body.Render(" "+text),
		)
	}

	totalWidth := m.widthContent
	if totalWidth <= 0 {
		totalWidth = 80
	}

	leftWidth := int(float64(totalWidth) * 0.58)
	if leftWidth <= 0 {
		leftWidth = totalWidth
	}
	if leftWidth < 48 {
		leftWidth = totalWidth
	}
	rightWidth := totalWidth - leftWidth
	if rightWidth < 0 {
		rightWidth = 0
	}

	logo := brand.Render(aboutASCIIArt)
	rightArt := accent.Render(aboutHeroASCII)

	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		"",
		prompt("welcome", ""),
		body.Render("Welcome to my terminal portfolio. (Version 1.0.0)"),
		body.Render("----"),
		
		
		"",
		prompt("about", ""),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			body.Render("Hi, my name is "),
			brand.Render("Thisara Tharinda"),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			body.Render("I'm a "),
			accent.Render("Devops Engineer"),
		),
	)

	rightColumn := lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Right).Render(rightArt)

	content := lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.NewStyle().Width(leftWidth).Render(leftColumn), rightColumn)

	return lipgloss.Place(
		m.widthContent,
		m.heightContent,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m model) UpdateSelectedTheme() model {
	m.theme = theme.BasicTheme(m.renderer, nil)
	return m
}
