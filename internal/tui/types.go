package tui

import tea "github.com/charmbracelet/bubbletea"

type TUIModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

// NavigationMessages
type (
	NavigateToMainMsg   struct{}
	NavigateToSecondMsg struct{}
)
