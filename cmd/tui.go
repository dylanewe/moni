package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dylanewe/moni/internal/config"
	"github.com/dylanewe/moni/internal/db"
	"github.com/dylanewe/moni/internal/service"
	"github.com/dylanewe/moni/internal/store"
	"github.com/dylanewe/moni/internal/tui"
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
	db               *db.DBConnectionMsg
	store            *store.Store
	service          *service.Service
	cfg              *config.Config
	loading          bool
	spinner          spinner.Model
}

func initModel(cfg *config.Config, service *service.Service) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner:          s,
		loading:          true,
		stateDescription: "Initializing...",
		stateStatus:      tui.StatusBarStateBlue,
		commands: []command{
			{name: "View Dashboard"},
			{name: "Add Statement"},
			{name: "Add Category"},
		},
		cfg:     cfg,
		service: service,
	}
}

func (m model) Init() tea.Cmd {
	addr := m.cfg.DB.Address
	return tea.Batch(
		db.Init(addr),
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case db.DBConnectionMsg:
		m.stateDescription = ""
		m.db = &msg
		if m.db != nil {
			if m.db.Err != nil {
				m.stateStatus = tui.StatusBarStateRed
				m.stateDescription = shortenErr(m.db.Err, 100)

			} else {
				m.stateStatus = tui.StatusBarStateGreen
				m.stateDescription = "Connected to database"
			}
		}
		m.loading = false
		return m, nil

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

	var stateDescription string
	if !m.loading {
		stateDescription = m.stateDescription
		renderLists(doc, m)
	} else {
		stateDescription = m.spinner.View()
	}

	tui.RenderStatusBar(doc, tui.NewStatusBarProps(&tui.StatusBarProps{
		Description: stateDescription,
		User:        "NONE",
		StatusState: m.stateStatus,
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

func shortenErr(err error, length int) string {
	if len(err.Error()) < length {
		return err.Error()
	}

	return err.Error()[:length] + "..."
}
