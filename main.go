package main

import (
	"fmt"
	"os"

	"github.com/Expelliamus625/omdb-cli/internal/api"
	"github.com/Expelliamus625/omdb-cli/internal/config"
	"github.com/Expelliamus625/omdb-cli/internal/logger"
	"github.com/Expelliamus625/omdb-cli/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := logger.Init("./logs/logfile.log"); err != nil {
		fmt.Printf("Unable to initialize logger %v", err)
		os.Exit(1)
	}

	logger.Log.Info("Application Starting")
	cfg, err := config.Load("./config.json")
	if err != nil {
		logger.Log.Error("Could not load config", "error", err)
		fmt.Printf("Could not load config: %v\n", err)
		os.Exit(1)
	}
	logger.Log.Info("Config Loaded Successfully")

	client := api.NewClient(cfg)
	m := tui.RootScreen(client)
	// m := tui.NewModel(client)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Log.Error("Ran into an error", "error", err)
		fmt.Printf("Ran into an error: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("Exiting. Goodbye...")
}
