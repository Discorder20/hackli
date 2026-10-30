package pkg

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var editorTabs = []string{"Styles", "Layout", "State", "Commands"}

var (
	accentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	activeTab   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).Bold(true).Padding(0, 1)
	inactiveTab = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Padding(0, 1)
	panelStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).Padding(0, 1)
)

type editorModel struct {
	folderPath string
	project    *Project
	tab        int
	width      int
	height     int
	counter    int
	enabled    bool
	requests   int
	completed  int
}

type exampleResultMsg struct{}

func OpenEditor(folderPath string, project *Project) error {
	_, err := tea.NewProgram(editorModel{folderPath: folderPath, project: project}).Run()
	if err != nil {
		return fmt.Errorf("run editor: %w", err)
	}
	return nil
}

func (m editorModel) Init() tea.Cmd {
	return nil
}

func (m editorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case exampleResultMsg:
		m.completed++
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.tab = (m.tab + 1) % len(editorTabs)
		case "shift+tab":
			m.tab = (m.tab + len(editorTabs) - 1) % len(editorTabs)
		default:
			return m.updateTab(msg.String())
		}
	}
	return m, nil
}

func (m editorModel) updateTab(key string) (tea.Model, tea.Cmd) {
	switch m.tab {
	case 2:
		switch key {
		case "up":
			m.counter++
		case "down":
			m.counter--
		case "space":
			m.enabled = !m.enabled
		}
	case 3:
		if key == "enter" {
			m.requests++
			return m, exampleCommand
		}
	}
	return m, nil
}

// A command does work outside Update and returns a message back to the program.
func exampleCommand() tea.Msg {
	return exampleResultMsg{}
}

func (m editorModel) View() tea.View {
	view := tea.NewView(m.render())
	view.AltScreen = true
	return view
}

func (m editorModel) render() string {
	if (m.width > 0 && m.width < 50) || (m.height > 0 && m.height < 16) {
		return "Please resize the terminal to at least 50x16.\nTab: next tab | q: quit"
	}

	name := "Unnamed project"
	if m.project != nil && m.project.Name != "" {
		name = m.project.Name
	}
	path := m.folderPath
	if m.width > 0 {
		name = ansi.Truncate(name, m.width-len("hackli playground  "), "...")
		path = ansi.Truncate(path, m.width, "...")
	}
	header := accentStyle.Render("hackli playground") + "  " + name
	body := m.renderTab()
	return lipgloss.JoinVertical(lipgloss.Left,
		header, mutedStyle.Render(path), "", renderNavbar(m.tab), "", body, "",
		mutedStyle.Render("tab / shift+tab: switch  |  q: quit"),
	)
}

func renderNavbar(selected int) string {
	tabs := make([]string, len(editorTabs))
	for i, name := range editorTabs {
		if i == selected {
			tabs[i] = activeTab.Render(name)
		} else {
			tabs[i] = inactiveTab.Render(name)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m editorModel) renderTab() string {
	switch m.tab {
	case 0:
		return renderStylesTab()
	case 1:
		return renderLayoutTab()
	case 2:
		return m.renderStateTab()
	case 3:
		return m.renderCommandsTab()
	default:
		return ""
	}
}

func renderStylesTab() string {
	bold := lipgloss.NewStyle().Bold(true).Render("bold")
	color := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render("colored")
	badge := lipgloss.NewStyle().Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).Padding(0, 1).Render(" padded ")
	return panelStyle.Render(strings.Join([]string{
		accentStyle.Render("Lip Gloss: styles"),
		"Styles transform strings before rendering.",
		bold + "  " + color + "  " + badge,
		"This panel has a rounded border and padding.",
		"See renderStylesTab and panelStyle in code.",
	}, "\n"))
}

func renderLayoutTab() string {
	box := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	columns := lipgloss.JoinHorizontal(lipgloss.Top,
		box.Render("Left\ncolumn"), box.Render("Right\ncolumn"))
	return panelStyle.Render(strings.Join([]string{
		accentStyle.Render("Lip Gloss: layout"),
		"JoinHorizontal places blocks side by side:",
		columns,
		"JoinVertical stacks the header, tabs and body.",
	}, "\n"))
}

func (m editorModel) renderStateTab() string {
	status := "off"
	if m.enabled {
		status = "on"
	}
	return panelStyle.Render(strings.Join([]string{
		accentStyle.Render("Bubble Tea: model and update"),
		"Keys change the model; View renders it.",
		fmt.Sprintf("Counter: %d    Toggle: %s", m.counter, status),
		"Up/down: change counter  |  space: toggle",
		"Switch tabs and return to preserve state.",
	}, "\n"))
}

func (m editorModel) renderCommandsTab() string {
	return panelStyle.Render(strings.Join([]string{
		accentStyle.Render("Bubble Tea: commands and messages"),
		"Enter returns a command from Update.",
		"The command sends a message back to Update.",
		fmt.Sprintf("Commands sent: %d    Messages received: %d", m.requests, m.completed),
		"Try pressing enter several times.",
	}, "\n"))
}
