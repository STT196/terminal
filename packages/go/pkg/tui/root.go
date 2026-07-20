package tui

import (
	"context"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminaldotshop/terminal/go/pkg/tui/theme"
)

type page = int
type size = int

const (
	menuPage page = iota
	splashPage
	aboutPage
	projectsPage
	contactPage

	skillsPage
)

const (
	undersized size = iota
	small
	medium
	large
)

type model struct {
	ready           bool
	command         []string
	switched        bool
	page            page
	state           state
	context         context.Context
	renderer        *lipgloss.Renderer
	theme           theme.Theme
	fingerprint     string
	anonymous       bool
	viewportWidth   int
	viewportHeight  int
	widthContainer  int
	heightContainer int
	widthContent    int
	heightContent   int
	size            size
	accessToken     string

	projects []Project
	skills   []string
	error    *VisibleError
}

type VisibleError struct {
	message string
}

type state struct {
	splash SplashState
	cursor cursorState

	about   aboutState
	project projectState
	footer  footerState

	menu menuState
}

type children struct {
}

func NewModel(
	renderer *lipgloss.Renderer,
	fingerprint string,
	anonymous bool,
	clientIP *string,
	command []string,
) (tea.Model, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "client_ip", clientIP)

	result := model{
		command:     command,
		context:     ctx,
		page:        splashPage,
		renderer:    renderer,
		fingerprint: fingerprint,
		anonymous:   anonymous,
		theme:       theme.BasicTheme(renderer, nil),
		projects:    LoadProjects(),
		skills:      LoadSkills(),
		// Initialize with reasonable defaults for large screens
		viewportWidth:   100,
		viewportHeight:  40,
		widthContainer:  80,
		heightContainer: 30,
		widthContent:    78,
		heightContent:   25,
		size:            large,
	}
	return result, nil
}

func (m model) Init() tea.Cmd {
	return m.SplashInit()
}

func (m model) SwitchPage(page page) model {
	m.page = page
	m.switched = true
	return m
}

func (m model) InitialDataLoaded() (model, tea.Cmd) {
	if len(m.command) == 0 {
		return m.AboutSwitch()
	}

	// Search projects by command
	command := strings.ToLower(m.command[0])

	for index, project := range m.projects {
		if strings.ToLower(project.Name) == command {
			m.state.project.selected = index
			return m.SwitchPage(projectsPage), nil
		}
	}

	return m.AboutSwitch()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case VisibleError:
		m.error = &msg
	case error:
		m.error = &VisibleError{
			message: msg.Error(),
		}
	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		m.viewportHeight = msg.Height

		switch {
		case m.viewportWidth < 20 || m.viewportHeight < 10:
			m.size = undersized
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 60:
			m.size = small
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 90:
			m.size = medium
			m.widthContainer = 70
			m.heightContainer = int(math.Min(float64(msg.Height), 35))
		default:
			m.size = large
			m.widthContainer = int(math.Min(float64(m.viewportWidth-10), 100))
			m.heightContainer = int(math.Min(float64(msg.Height), 40))
		}

		m.widthContent = m.widthContainer - 2
		m.heightContent = m.heightContainer - lipgloss.Height(m.HeaderView()) - lipgloss.Height(m.FooterView()) - 2
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.error != nil {
				if m.page == splashPage {
					return m, tea.Quit
				}
				m.error = nil
				return m, nil
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	case CursorTickMsg:
		m, cmd := m.CursorUpdate(msg)
		return m, cmd

	}

	var cmd tea.Cmd
	switch m.page {
	case menuPage:
		m, cmd = m.MenuUpdate(msg)
	case splashPage:
		m, cmd = m.SplashUpdate(msg)
	case projectsPage:
		m, cmd = m.ProjectUpdate(msg)

	case contactPage:
		m, cmd = m.ContactUpdate(msg)
	case aboutPage:
		m, cmd = m.AboutUpdate(msg)

	}

	var headerCmd tea.Cmd
	m, headerCmd = m.HeaderUpdate(msg)
	cmds = append(cmds, headerCmd)

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if m.switched {
		m.switched = false
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.size == undersized {
		return m.ResizeView()
	}

	switch m.page {
	case splashPage:
		return m.SplashView()
	case menuPage:
		return m.MenuView()
	default:
		header := m.HeaderView()
		footer := m.FooterView()

		// Get content based on current page
		content := m.getContent()

		height := m.heightContainer
		height -= lipgloss.Height(header)
		height -= lipgloss.Height(footer)

		body := m.theme.Base().Width(m.widthContainer).Height(height).Render(content)
		// bodyHeight := lipgloss.Height(body)
		// if bodyHeight < height {
		// 	body += lipgloss.NewStyle().Height(height - bodyHeight).Render(" ")
		// }

		items := []string{}
		items = append(items, header)
		items = append(items, body)
		items = append(items, footer)

		child := lipgloss.JoinVertical(
			lipgloss.Left,
			items...,
		)

		return m.renderer.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			lipgloss.Center,
			m.theme.Base().
				MaxWidth(m.widthContainer).
				MaxHeight(m.heightContainer).
				Render(child),
		)
	}
}

func (m model) getContent() string {
	page := "unknown"
	switch m.page {
	case aboutPage:
		page = m.AboutView()

	case projectsPage:
		page = m.ProjectView()
	case skillsPage:
		page = m.SkillsView(m.widthContent)
	case contactPage:
		page = m.ContactView(m.widthContent)
	}
	return page
}

var modifiedKeyMap = viewport.KeyMap{
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "½ page down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
}
