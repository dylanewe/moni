package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dylanewe/moni/internal/config"
	"github.com/dylanewe/moni/internal/db"
	"github.com/dylanewe/moni/internal/service"
	"github.com/dylanewe/moni/internal/store"
	"github.com/dylanewe/moni/internal/tui"
	"github.com/dylanewe/moni/internal/util"
)

const (
	width       = 96
	columnWidth = 30
)

type mode string

const (
	modeFilePicker mode = "filepicker"
	modeCategorize mode = "categorize"
	modeSaving     mode = "saving"
	modeLoading    mode = "loading"
	modeDefault    mode = ""
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
	fileStatements   []string
	fileCursor       int
	mode             mode
	extractedTx      *db.ExtractStatementMsg
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
		mode:    modeLoading,
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
		m.mode = modeDefault
		return m, nil

	case db.ExtractStatementMsg:
		m.extractedTx = &msg
		if m.extractedTx != nil {
			if m.extractedTx.Err != nil {
				m.stateStatus = tui.StatusBarStateRed
				m.stateDescription = shortenErr(m.extractedTx.Err, 50)
				m.mode = modeDefault
			} else {
				transactions := m.extractedTx.Transactions
				var uncategorizedTx []store.Transaction
				for _, tx := range transactions {
					if !slices.Contains(m.cfg.Categories, tx.CategoryName) {
						uncategorizedTx = append(uncategorizedTx, tx)
					}
				}

				if len(uncategorizedTx) > 0 {
					m.stateDescription = "Categorize these transactions"
					m.stateStatus = tui.StatusBarStateBlue
					m.mode = modeCategorize
				} else {
					m.stateDescription = "Confirm to add these transactions?"
					m.stateStatus = tui.StatusBarStateBlue
					m.mode = modeSaving
				}

			}
		}
		m.loading = false
		return m, nil

	case tea.KeyMsg:

		switch msg.String() {

		case tea.KeyEnter.String():
			if m.mode == modeFilePicker {
				m.stateDescription = "Parsing statement..."
				m.stateStatus = tui.StatusBarStateYellow
				filepath := "../statements/" + m.fileStatements[m.fileCursor]
				m.mode = modeLoading
				m.loading = true
				return m, db.ExtractStatement(m.service.LLMParser, m.cfg.Categories, filepath)
			}
			if m.cursor == 1 {
				statements, err := util.ReadFilesFromFolder("../statements/", []string{".pdf"})
				if err != nil {
					m.stateDescription = "Error reading statements folder"
					m.stateStatus = tui.StatusBarStateRed
					return m, nil
				}
				if len(statements) == 0 {
					m.stateDescription = "No statements found"
					m.stateStatus = tui.StatusBarStateYellow
					return m, nil
				}
				m.fileStatements = statements
				m.stateDescription = "Pick a financial statement to add"
				m.stateStatus = tui.StatusBarStateBlue
				m.mode = modeFilePicker
			}

		case "ctrl+c", "q":
			if m.mode == modeFilePicker {
				m.mode = modeDefault
				return m, nil
			}
			return m, tea.Quit

		case "up", "k":
			if m.mode == modeFilePicker {
				if m.fileCursor > 0 {
					m.fileCursor--
				}
			} else {
				if m.cursor > 0 {
					if m.commands[m.cursor-1].disabled {
						m.cursor--
					}
					m.cursor--
				}
			}

		case "down", "j":
			if m.mode == modeFilePicker {
				if m.fileCursor < len(m.fileStatements)-1 {
					m.fileCursor++
				}
			} else {
				if m.cursor < len(m.commands)-1 {
					if m.commands[m.cursor+1].disabled {
						m.cursor++
					}
					m.cursor++
				}
			}
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd
}

func (m model) View() string {
	doc := &strings.Builder{}

	tui.RenderTitleRow(width, doc, tui.TitleRowProps{Title: "Moni: Your Financial Planner"})
	doc.WriteString("\n\n")

	var stateDescription string
	if !m.loading {
		stateDescription = m.stateDescription
		renderLists(doc, m)
	} else {
		stateDescription = m.spinner.View()
	}

	doc.WriteString("[q] Quit		[Enter] Select")
	doc.WriteString("\n\n")

	tui.RenderStatusBar(doc, tui.NewStatusBarProps(&tui.StatusBarProps{
		Description: stateDescription,
		User:        "NONE",
		StatusState: m.stateStatus,
		Width:       width,
	}))

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

	leftList := tui.RenderListCommands(doc, &tui.ListProps{
		Items:    items,
		Selected: m.cursor,
	})

	var rightList string
	if m.mode == modeFilePicker {
		var fileList []string
		for i, f := range m.fileStatements {
			if i == m.fileCursor {
				fileList = append(fileList, fmt.Sprintf("> %s", f))
			} else {
				fileList = append(fileList, fmt.Sprintf("  %s", f))
			}
		}
		rightList = tui.RenderListDisplay("Statements", fileList)
	} else {
		rightList = tui.RenderListDisplay(m.secondListHeader, m.secondListValues)
	}

	lists := lipgloss.JoinHorizontal(lipgloss.Top, leftList, rightList)

	doc.WriteString(lists)
	doc.WriteString("\n\n")
}

func shortenErr(err error, length int) string {
	if len(err.Error()) < length {
		return err.Error()
	}

	return err.Error()[:length] + "..."
}
