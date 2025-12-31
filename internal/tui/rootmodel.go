package tui

import (
	"github.com/Expelliamus625/omdb-cli/internal/api"
	tea "github.com/charmbracelet/bubbletea"
)

// Rootmodel encapsulates the actual model to be displayed on the screen. This will allow switching out the displayed model during runtime
type rootModel struct {
	TUIModel
	client api.MovieService
	width  int
	height int
}

func RootScreen(client api.MovieService) rootModel {
	model := NewModel(client)

	return rootModel{
		TUIModel: model,
		client:   client,
	}
}

func (rm rootModel) Init() tea.Cmd {
	return rm.TUIModel.Init()
}

func (rm rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Capture window size changes
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		rm.width = msg.Width
		rm.height = msg.Height
	}

	newModel, cmd := rm.TUIModel.Update(msg)

	// Check if navigational messages are persent, else return updated model
	switch msg.(type) {
	case NavigateToMainMsg:
		rm.TUIModel = NewModel(rm.client)
		// Send the cached window size to the new model
		if rm.width > 0 && rm.height > 0 {
			rm.TUIModel.Update(tea.WindowSizeMsg{Width: rm.width, Height: rm.height})
		}
		return rm, rm.TUIModel.Init()
	case NavigateToSecondMsg:
		rm.TUIModel = NewLoaderModel()
		if rm.width > 0 && rm.height > 0 {
			rm.TUIModel.Update(tea.WindowSizeMsg{Width: rm.width, Height: rm.height})
		}

		return rm, rm.TUIModel.Init()
	}

	rm.TUIModel = newModel.(TUIModel)
	return rm, cmd
}

func (rm rootModel) View() string {
	return rm.TUIModel.View()
}
