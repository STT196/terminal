package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	var appsPageIndex int
	for i, page := range m.accountPages {
		if page == appsPage {
			appsPageIndex = i
			break
		}
	}

	if (m.page == shippingPage && m.state.shipping.view == shippingFormView) ||
		(m.page == paymentPage && m.state.payment.view == paymentFormView) ||
		(m.page == accountPage && m.state.account.selected == appsPageIndex && m.state.apps.editing) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			if m.page != cartPage {
				return m.CartSwitch()
			}
		case "a":
			return m.ShopSwitch()
		case "s":
			return m.AccountSwitch()
		case "k":
			return m.SkillsSwitch()
		case "o":
			return m.ContactSwitch()
		// case "f":
		// 	return m.FaqSwitch()
		case "m":
			return m.MenuSwitch()
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) HeaderView() string {

	bold := m.theme.TextAccent().Bold(true).Render
	accent := m.theme.TextAccent().Render
	base := m.theme.Base().Render
	cursor := m.theme.Base().Background(m.theme.Brand()).Render(" ")

	menu := bold("m") + base(" ☰")
	mark := bold("t") + cursor
	logo := bold("STT196")
	shop := accent("a") + base(" about")
	account := accent("p") + base(" projects")
	skills := accent("s") + base(" skills")
	contact := accent("c") + base(" contact")

	switch m.page {
	case shopPage:
		shop = accent("a about")
	case accountPage:
		account = accent("p projects")
	case skillsPage:
		skills = accent("s skills")
	case contactPage:
		contact = accent("c contact")
	}

	var tabs []string

	switch m.size {
	case small:
		tabs = []string{
			mark,
		}
	case medium:
		tabs = []string{
			menu,
			logo,
		}
	default:
		tabs = []string{
			logo,
			shop,
			account,
			skills,
			contact,
		}
	}

	var table = table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Row(tabs...).
		Width(m.widthContent).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().
				Padding(0, 1).
				AlignHorizontal(lipgloss.Center)
		}).
		Render()

	return lipgloss.Place(
		m.widthContainer,
		lipgloss.Height(table),
		lipgloss.Center,
		lipgloss.Center,
		table,
	)
}
