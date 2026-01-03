package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dylanewe/moni/internal/config"
	"github.com/dylanewe/moni/internal/db"
	service "github.com/dylanewe/moni/internal/services"
	"github.com/dylanewe/moni/internal/store"
	"github.com/dylanewe/moni/internal/tui"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const (
	width       = 96
	columnWidth = 30
)

type command struct {
	disabled bool
	name     string
}

type model struct {
	stateDescription string
	stateStatus      tui.StatusBarState
	commands         []command
	cursor           int
	secondListHeader string
	secondListValues []string
}

func initModel(cfg config.Config, service service.Service, store store.Store) model {
	return model{
		stateDescription: "Initializing...",
		commands: []command{
			{name: "View Dashboard"},
			{name: "Add Statement"},
			{name: "Add Category"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				if m.commands[m.cursor-1].disabled {
					m.cursor--
				}
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				if m.commands[m.cursor+1].disabled {
					m.cursor++
				}
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	doc := &strings.Builder{}

	tui.RenderTitleRow(width, doc, tui.TitleRowProps{Title: "Moni"})
	doc.WriteString("\n\n")

	doc.WriteString(fmt.Sprintf("Cursor: %d", m.cursor))
	doc.WriteString("\n\n")

	renderLists(doc, m)

	tui.RenderStatusBar(doc, tui.NewStatusBarProps(&tui.StatusBarProps{
		Description: m.stateDescription,
		User:        "NONE",
		StatusState: tui.StatusBarStateBlue,
		Width:       width,
	}))

	doc.WriteString("Press q to quit")
	doc.WriteString("\n\n")

	return doc.String()
}

func renderLists(doc *strings.Builder, m model) {
	var items []tui.Item
	for _, c := range m.commands {
		items = append(items, tui.Item{
			Value:    c.name,
			Disabled: c.disabled,
		})
	}

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		tui.RenderListCommands(doc, &tui.ListProps{
			Items:    items,
			Selected: m.cursor,
		}),
		tui.RenderListDisplay(m.secondListHeader, m.secondListValues),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists))
	doc.WriteString("\n\n")
}

func main() {
	cfg, err := config.GetConfig("../config.toml")
	if err != nil {
		log.Fatalf("get config error: %v", err)
	}

	db, err := db.New(cfg.DB.Address)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer db.Close()

	store := store.NewStore(db)

	llmClient := openai.NewClient(option.WithAPIKey(cfg.LLM.APIKey))
	service := service.NewService(&llmClient)

	p := tea.NewProgram(initModel(cfg, service, store))
	if _, err := p.Run(); err != nil {
		log.Fatalf("TUI run error: %v", err)
	}
}
