package cli

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

var appNameStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF5F87")).
	Bold(true)

var errorAlertStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#FF5F87")).
	Padding(0, 1)

var statusStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#009933")).
	Padding(0, 1)

var errorInfoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250")).
	Padding(0, 1)

var listEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
var listItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

var faint = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Faint(true)

func (m model) View() string {
	if m.fatalErr != nil {
		return errorAlertStyle.Render("ERROR") + errorInfoStyle.Render(m.fatalErr.Error()) + "\n"
	}

	s := appNameStyle.Render("ultrafocus") + faint.Render(" - Reclaim your time.") + "\n\n"
	s += statusStyle.Render("STATUS") + errorInfoStyle.Render(string(m.status)) + "\n\n"

	if m.state == blacklistView {
		s += "Edit/add domains:\n\n" + m.textarea.View() + "\n\n"
		s += "press Esc to save.\n"
	}

	if m.state == menuView {
		commands := m.getCommandsList()

		l := list.New().Enumerator(func(items list.Items, i int) string {
			if i == m.commandsListSelection {
				return "→"
			}
			return " "
		}).
			EnumeratorStyle(listEnumeratorStyle).
			ItemStyle(listItemStyle)
		for _, c := range commands {
			l.Item(c.Name + faint.Render(" - "+c.Desc))
		}
		s += l.String() + "\n\n"
	}

	if m.state == timeView {
		s += "Timer: " + m.timer.View() + "\n\n"
		m.timer.Start()
	}

	s += "press q to quit.\n"

	return s
}
