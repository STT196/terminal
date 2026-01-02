package tui

func (m model) LogoView() string {
	return m.theme.TextAccent().Bold(true).Render("STT196") + m.CursorView()
}
