package cli

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/plutov/ultrafocus/hosts"
)

type sessionState uint

const (
	menuView sessionState = iota
	blacklistView
	timeView
)

type model struct {
	textarea              textarea.Model
	timer                 timer.Model
	fatalErr              error
	status                hosts.FocusStatus
	domains               []string
	commandsListSelection int
	state                 sessionState
}

func NewModel() model {
	domains, status, err := hosts.ExtractDomainsFromHostsFile()

	if len(domains) == 0 {
		domains = hosts.DefaultDomains
	}

	return model{
		textarea: GetTextareModel(),
		timer:    GetTimerModel(time.Second * 5),
		domains:  domains,
		state:    menuView,
		status:   status,
		fatalErr: err,
	}
}

func (m model) Init() tea.Cmd {
	if m.fatalErr != nil {
		return tea.Quit
	}

	return nil
}

func (m *model) getCommandsList() []command {
	if m.status == hosts.FocusStatusOn {
		return []command{commandFocusOff, commandConfigureBlacklist, startTime, stopTime}
	}

	return []command{commandFocusOn, commandConfigureBlacklist, startTime, stopTime}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		commands := m.getCommandsList()
		switch msg.String() {

		case "up", "k":
			if m.state == menuView && m.commandsListSelection > 0 {
				m.commandsListSelection--
			}

		case "down", "j":
			if m.state == menuView && m.commandsListSelection < len(commands)-1 {
				m.commandsListSelection++
			}

		case "enter", " ":
			if m.state == menuView {
				m = commands[m.commandsListSelection].Run(m)
				if m.fatalErr != nil {
					return m, tea.Quit
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.state == blacklistView {
				domains := strings.Split(m.textarea.Value(), "\n")
				domains = hosts.CleanDomainsList(domains)

				if err := hosts.WriteDomainsToHostsFile(domains, m.status); err != nil {
					m.fatalErr = err
					return m, tea.Quit
				}

				m.commandsListSelection = 0
				m.domains = domains
				m.state = menuView
				m.textarea.Blur()
			}
		}
	}

	return m, tea.Batch(cmds...)
}
