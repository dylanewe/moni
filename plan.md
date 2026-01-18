# Plan for TUI Dashboard Implementation

This document outlines the plan to implement a dashboard view in the Moni TUI.

## 1. Data Access Layer (`internal/store`)

- [ ] Create a new file `internal/store/dashboard.go` to house the data access logic for the dashboard.

- [ ] In `internal/store/dashboard.go`, define a `DashboardStore` and add the following methods:
  - `GetTotalIncomeAndExpense() (float64, float64, error)`: This will query the `transactions` table to get the sum of all income and expenses.
  - `GetMonthlyIncomeAndExpense(year int, month int) (float64, float64, error)`: This will query the `transactions` table to get the sum of income and expenses for a given month and year.

- [ ] Update `internal/store/store.go` to include the `DashboardStore` in the main `Store` struct. This will make it accessible to the TUI.

## 2. TUI Components (`internal/tui`)

- [ ] Create a new file `internal/tui/dashboard.go` for the main dashboard component. This component will manage the state of the dashboard, including the currently selected month and year.

- [ ] The `Dashboard` model will implement the `bubbletea.Model` interface (`Init`, `Update`, `View`).

- [ ] The `Update` method will handle key presses for navigating between months (`h`, `l`, left arrow, right arrow) and will fetch the new data from the `DashboardStore`.

- [ ] The `View` method will use `lipgloss` to create a two-pane layout:
  - The left pane will display the total income and expense.
  - The right pane will display the monthly income and expense for the selected month.

- [ ] The `Dashboard` model will prevent navigation to future months.

## 3. Main Application (`cmd/`)

- [ ] Modify `cmd/tui.go` to manage the different views (the existing transaction list and the new dashboard).

- [ ] Add a new state to the main `model` to track the current view (e.g., `listView`, `dashboardView`).

- [ ] Add a keybinding (e.g., `d`) to switch from the list view to the dashboard view. A corresponding keybinding (e.g., `l` for list) should be available in the dashboard to switch back.

- [ ] The main model's `Update` function will delegate to the active view's `Update` function. The `View` function will similarly delegate to the active view's `View` function.

## 4. File Breakdown

### New Files:
- `internal/store/dashboard.go`: For database queries related to the dashboard.
- `internal/tui/dashboard.go`: The main component for the dashboard view.

### Modified Files:
- `internal/store/store.go`: To integrate the new `DashboardStore`.
- `cmd/tui.go`: To manage the new dashboard view and switch between views.
- `internal/tui/constants.go`: To add new keybindings.
