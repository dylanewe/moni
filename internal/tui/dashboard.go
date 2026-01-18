package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dylanewe/moni/internal/store"
)

type DashboardModel struct {
	store         *store.Store
	keys          KeyMap
	width, height int

	// state
	currentDate   time.Time
	totalIncome   float64
	totalExpense  float64
	monthlyIncome float64
	monthlyExpense float64
	err           error
}

func NewDashboardModel(store *store.Store, keys KeyMap) *DashboardModel {
	now := time.Now()
	return &DashboardModel{
		store:       store,
		keys:        keys,
		currentDate: now,
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	return m.fetchData()
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case dataFetchedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.totalIncome = msg.totalIncome
		m.totalExpense = msg.totalExpense
		m.monthlyIncome = msg.monthlyIncome
		m.monthlyExpense = msg.monthlyExpense
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Prev):
			m.currentDate = m.currentDate.AddDate(0, -1, 0)
			return m, m.fetchData()
		case key.Matches(msg, m.keys.Next):
			now := time.Now()
			if m.currentDate.Year() < now.Year() || (m.currentDate.Year() == now.Year() && m.currentDate.Month() < now.Month()) {
				m.currentDate = m.currentDate.AddDate(0, 1, 0)
				return m, m.fetchData()
			}
		}
	}
	return m, nil
}

func (m *DashboardModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	totalView := m.renderTotalView()
	monthlyView := m.renderMonthlyView()

	return lipgloss.JoinHorizontal(lipgloss.Top, totalView, monthlyView)
}

func (m *DashboardModel) renderTotalView() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(m.width/2 - 2).
		Height(m.height - 2).
		Padding(1)

	income := fmt.Sprintf("Total Income: %.2f", m.totalIncome)
	expense := fmt.Sprintf("Total Expense: %.2f", m.totalExpense)

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, income, expense))
}

func (m *DashboardModel) renderMonthlyView() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(m.width/2 - 2).
		Height(m.height - 2).
		Padding(1)

	monthStr := m.currentDate.Format("January 2006")
	income := fmt.Sprintf("Income for %s: %.2f", monthStr, m.monthlyIncome)
	expense := fmt.Sprintf("Expense for %s: %.2f", monthStr, m.monthlyExpense)

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, income, expense))
}

type dataFetchedMsg struct {
	totalIncome   float64
	totalExpense  float64
	monthlyIncome float64
	monthlyExpense float64
	err           error
}

func (m *DashboardModel) fetchData() tea.Cmd {
	return func() tea.Msg {
		totalIncome, totalExpense, err := m.store.Dashboard.GetTotalIncomeAndExpense()
		if err != nil {
			return dataFetchedMsg{err: err}
		}

		monthlyIncome, monthlyExpense, err := m.store.Dashboard.GetMonthlyIncomeAndExpense(m.currentDate.Year(), int(m.currentDate.Month()))
		if err != nil {
			return dataFetchedMsg{err: err}
		}

		return dataFetchedMsg{
			totalIncome:   totalIncome,
			totalExpense:  totalExpense,
			monthlyIncome: monthlyIncome,
			monthlyExpense: monthlyExpense,
		}
	}
}

func (m *DashboardModel) updateData(msg dataFetchedMsg) {
	if msg.err != nil {
		m.err = msg.err
		return
	}
	m.totalIncome = msg.totalIncome
	m.totalExpense = msg.totalExpense
	m.monthlyIncome = msg.monthlyIncome
	m.monthlyExpense = msg.monthlyExpense
}
