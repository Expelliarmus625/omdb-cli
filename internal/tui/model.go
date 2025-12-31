package tui

import (
	"fmt"

	"github.com/Expelliamus625/omdb-cli/internal/api"
	"github.com/Expelliamus625/omdb-cli/internal/config"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	config      *config.Config
	movieClient api.MovieService
	textInput   textinput.Model
	infoBox     textinput.Model
	table       table.Model
	height      int
	width       int
	active      int
}

func NewModel(client api.MovieService) *model {
	// movie.GetMovie("avatar")
	columns := []table.Column{
		{Title: "Key", Width: 1},
		{Title: "Value", Width: 3},
	}

	rows := []table.Row{}
	table := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(7), table.WithFocused(false))
	textInput := textinput.New()
	textInput.Focus()
	textInput.Width = 20
	textInput.Prompt = "> "
	textInput.Placeholder = "Search for a movie"

	infoBox := textinput.New()
	infoBox.Prompt = " "

	return &model{
		movieClient: client,
		table:       table,
		textInput:   textInput,
		infoBox:     infoBox,
		active:      0,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tableCmd tea.Cmd
	var textInputCmd tea.Cmd
	var infoBoxCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update model width
		m.width = msg.Width - 2
		m.height = msg.Height - 2

		// update table Width and Heigh
		m.table.SetWidth(int(float64(m.width) * 0.7))
		m.table.SetHeight(m.height - 4)

		// Set Column Width. Currently setting width of all columns equally
		for i := range m.table.Columns() {
			newColWidth := int(float64(m.table.Width() / len(m.table.Columns())))
			m.table.Columns()[i].Width = newColWidth
			fmt.Printf("%s: %d", m.table.Columns()[i].Title, newColWidth)
		}

		// Set textinput width (30% of available width, leaving 70% for table
		m.textInput.Width = int(float64(m.width)*0.3) - 10
		m.infoBox.Width = int(float64(m.width)*0.3) - 10

	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "tab":
			// Toggle active window
			if m.active == 0 {
				m.active = 1
				m.textInput.Blur()
				m.table.Focus()
			} else {
				m.active = 0
				m.textInput.Focus()
				m.table.Blur()
			}

		case "enter":
			switch m.active {
			case 0:
				input := m.textInput.Value()
				movie, err := m.movieClient.GetMovie(input)
				if err != nil {
					// fmt.Printf("%s", err.Error())
					m.infoBox.SetValue(err.Error())
					return m, nil
				}
				// TODO: Implement getPlotRows
				// getPlotRows(movie.Plot, m.table.Columns()[1].Width)
				m.table.SetRows([]table.Row{
					{"Title", movie.Title},
					{"Year", movie.Year},
					{"Rated", movie.Rated},
					{"Director", movie.Director},
					{"Runtime", movie.Runtime},
					{"Released", movie.Released},
					{"Genre", movie.Genre},
					{"Country", movie.Country},
					{"Writer", movie.Writer},
					{"Actors", movie.Actors},
					{"Plot", movie.Plot},
					{"Country", movie.Country},
					{"Language", movie.Language},
					{"IMDB Rating", movie.ImdbRating},
					{"Awards", movie.Awards},
					{"DVD", movie.DVD},
				})
				m.textInput.Reset()
				m.infoBox.Reset()
				return m, nil
			case 1:
				return m, nil
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, tableCmd = m.table.Update(msg)
	m.textInput, textInputCmd = m.textInput.Update(msg)
	m.infoBox, infoBoxCmd = m.infoBox.Update(msg)
	return m, tea.Batch(tableCmd, textInputCmd, infoBoxCmd)
}

func (m *model) View() string {
	// Table style
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)

	m.table.SetStyles(s)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left, baseStyle.Render(m.textInput.View()), m.infoBox.View()),
		baseStyle.Render(m.table.View())) + "\nPress ctrl + c to quit"
	// return baseStyle.Render(m.table.View()) + "\n" + m.textInput.View()
}

// func getPlotRows(plot string, width int) []table.Row {
// 	plotRows := []table.Row{
// 		{"Plot", plot[:width]},
// 	}
//
// 	for _, str := range plot[width+1:] {
// 		row := table.Row{
// 			"", plot[width:]
// 		}
// 	}
// }
