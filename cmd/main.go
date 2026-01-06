package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dylanewe/moni/internal/config"
	"github.com/dylanewe/moni/internal/service"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	cfg, err := config.GetConfig("../config.toml")
	if err != nil {
		log.Fatalf("get config error: %v", err)
	}

	llmClient := openai.NewClient(option.WithAPIKey(cfg.LLM.APIKey))
	service := service.NewService(&llmClient)

	p := tea.NewProgram(initModel(&cfg, &service))
	if _, err := p.Run(); err != nil {
		log.Fatalf("TUI run error: %v", err)
	}
}
