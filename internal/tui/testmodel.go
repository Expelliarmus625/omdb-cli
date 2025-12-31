package tui

import (
	"github.com/Expelliamus625/omdb-cli/internal/logger"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loaderModel struct {
	spinner spinner.Model
	width   int
	height  int
}

func NewLoaderModel() *loaderModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	return &loaderModel{
		spinner: s,
	}
}

func (lm *loaderModel) Init() tea.Cmd {
	return lm.spinner.Tick
}

func (lm *loaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		lm.width = msg.Width
		lm.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return lm, func() tea.Msg { return NavigateToMainMsg{} }
		case "ctrl+c":
			return lm, tea.Quit
		}
	}
	var cmd tea.Cmd
	lm.spinner, cmd = lm.spinner.Update(msg)
	return lm, cmd
}

func (lm *loaderModel) View() string {
	content := lm.spinner.View() + " Coming soon to an Theatre near you..."
	logger.Log.Info("Height and width of loaderModel", "height", lm.height, "width", lm.width)
	return lipgloss.Place(lm.width, lm.height, lipgloss.Center, lipgloss.Center, content)
}
