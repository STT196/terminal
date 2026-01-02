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
	breadcrumbsHeight := lipgloss.Height(m.BreadcrumbsView())
	footerHeight := lipgloss.Height(m.FooterView())
	verticalMarginHeight := headerHeight + footerHeight + breadcrumbsHeight

	availableHeight := m.heightContainer - verticalMarginHeight

	// Calculate menu width based on projects
	menuWidth := 10 // "projects" has 8 chars
	if menuWidth > 0 {
		menuWidth += 4
	}

	// For small screens, make the menu full width
	if m.size < large {
		menuWidth = m.widthContent
	}

	detailWidth := m.widthContent - menuWidth
	if m.size < large {
		detailWidth = m.widthContent
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
			if m.state.account.selected < len(projects)-1 {
				m.state.account.selected++
			}
		case "shift+tab", "up", "k":
			if m.state.account.selected > 0 {
				m.state.account.selected--
			}
		}
	}

	// Update viewports with new content
	if m.state.account.viewportsReady {
		menuContent := m.getProjectsMenuContent()
		m.state.account.menuViewport.SetContent(menuContent)

		detailContent := m.getProjectsContent()
		m.state.account.detailViewport.SetContent(detailContent)
	}

	// Update the detailViewport with the message
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
	detailContent := m.getProjectsContent()
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

// Helper function to scroll the detail viewport to show the selected item in focused account pages
// func (m model) scrollToAccountDetailItem(model model, accountPage page) model {
// 	// If orders page is in detail view, we don't need to scroll to a specific item
// 	if accountPage == ordersPage && model.state.orders.viewing {
// 		return model
// 	}

// 	var itemHeight int
// 	var itemCount int
// 	var selectedIndex int

// 	// Different item heights and counts based on the page
// 	switch accountPage {
// 	case subscriptionsPage:
// 		itemHeight = 5 // Estimated height of a subscription item with padding
// 		itemCount = len(model.subscriptions)
// 		selectedIndex = model.state.subscriptions.selected
// 	case tokensPage:
// 		itemHeight = 7                    // Estimated height of a token item with padding
// 		itemCount = len(model.tokens) + 1 // +1 for "add token" button
// 		selectedIndex = model.state.tokens.selected
// 	case appsPage:
// 		itemHeight = 8                  // Estimated height of an app item with padding
// 		itemCount = len(model.apps) + 1 // +1 for "create app" button
// 		selectedIndex = model.state.apps.selected
// 	case ordersPage:
// 		itemHeight = 4 // Reduced height for order item with just date (instead of all products)
// 		itemCount = len(model.orders)
// 		selectedIndex = model.state.orders.selected
// 	default:
// 		return model // No scrolling for other pages
// 	}

// 	if itemCount == 0 {
// 		return model // No items to scroll to
// 	}

// 	// Calculate approximate position of selected item
// 	targetY := (selectedIndex * itemHeight) + 2

// 	// Calculate offset to position item in the visible area
// 	viewportHeight := model.state.account.detailViewport.Height
// 	currentOffset := model.state.account.detailViewport.YOffset

// 	// If item is above viewport, scroll up to show it
// 	if targetY < currentOffset {
// 		model.state.account.detailViewport.SetYOffset(targetY - 2)
// 	}

// 	// If item is below viewport, scroll down to show it
// 	if targetY+itemHeight > currentOffset+viewportHeight {
// 		model.state.account.detailViewport.SetYOffset(targetY - viewportHeight + itemHeight)
// 	}

// 	return model
// }
