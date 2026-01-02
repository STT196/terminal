package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Project struct {
	Name         string
	Year         string
	Technologies string
	Description  string
}

var projects = []Project{
	{
		Name:         "Spring Boot & React Full-Stack Deployment with Docker",
		Year:         "2025",
		Technologies: "Spring Boot, React, Docker, Docker Compose, Nginx, GitHub Actions, Linux",
		Description: `• Containerized both Spring Boot backend and React frontend using multi-stage Docker builds to create optimized, production-ready images.
• Deployed the full-stack system using Docker Compose with isolated containers and an internal Docker network for secure service communication.
• Configured Nginx as a reverse proxy for frontend hosting and API routing to ensure stable, scalable backend access.
• Implemented CI/CD workflow with GitHub Actions to automate builds, containerization, testing, and deployment to a Linux VPS environment.
• Enhanced application performance and deployment reliability by reducing image size, standardizing environments, and automating repetitive tasks.`,
	},
	{
		Name:         "Flask Web Application DevOps Pipeline",
		Year:         "2025",
		Technologies: "Flask, GitHub Actions, Nginx, Gunicorn, Linux VPS, Bash scripting",
		Description: `• Deployed Flask web application with automated CI/CD pipeline using GitHub Actions for continuous integration and deployment.
• Configured Nginx reverse proxy and Gunicorn WSGI server for production hosting with load balancing and SSL termination.
• Implemented monitoring and logging solutions for application performance tracking and issue resolution.
• Created Bash scripts for automated deployment processes reducing manual intervention by 90%.`,
	},
	{
		Name:         "Laravel Application Infrastructure & CI/CD",
		Year:         "2024",
		Technologies: "Laravel, GitHub Actions, MySQL, Nginx, PHP-fpm, Linux administration",
		Description: `• Designed and implemented complete DevOps pipeline for Laravel web application with automated testing and deployment.
• Configured PHP-fpm and Nginx for optimal PHP application performance and resource management.
• Managed MySQL database operations including backups, performance tuning, and replication setup.
• Implemented infrastructure monitoring and alerting systems for proactive issue detection and resolution.`,
	},
	{
		Name:         "Windows Server Infrastructure Automation",
		Year:         "2024",
		Technologies: "Windows Server 2019, PowerShell, Active Directory, Network services",
		Description: `• Automated infrastructure setup using PowerShell scripts for user management and system configuration.
• Configured network services including DHCP, NAT, and RAS for comprehensive infrastructure management.
• Implemented automated user provisioning system creating 1,000+ Active Directory accounts with proper permissions.
• Documented infrastructure procedures and created runbooks for operational continuity and knowledge transfer.`,
	},
}

func (m model) ProjectsUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			if m.state.account.selected < len(projects)-1 {
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

func (m model) getProjectsContent() string {
	if len(projects) == 0 {
		return m.theme.Base().Render("No projects available")
	}

	project := projects[m.state.account.selected]
	accent := m.theme.TextAccent().Render
	base := m.theme.Base().Render
	bold := m.theme.TextHighlight().Bold(true).Render

	content := strings.Builder{}
	content.WriteString(accent(project.Name) + "\n")
	content.WriteString(base(project.Year) + "\n\n")
	content.WriteString(bold("Technologies:\n"))
	content.WriteString(base(project.Technologies) + "\n\n")
	content.WriteString(bold("Details:\n"))
	content.WriteString(base(project.Description))

	return m.theme.Base().Padding(1, 2).Render(content.String())
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
		menuWidth := 35
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
	for i, p := range projects {
		var item string
		if i == m.state.account.selected {
			item = highlightedMenuItem.Render(p.Name + " (" + p.Year + ")")
		} else {
			item = menuItem.Render(p.Name + " (" + p.Year + ")")
		}
		content.WriteString(item + "\n")
	}

	return m.theme.Base().Padding(0, 1).Render(content.String())
}
