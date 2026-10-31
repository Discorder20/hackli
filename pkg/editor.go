package pkg

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var tabNames = [tabsCount]string{"Overview", "Styles"}

const (
	OverviewTab = iota
	StylesTab
)

const tabsCount = 2

const primaryColor = "#0ebb08"
const disabledColor = "#535353"

var regularTab = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA"))

var highlightedTab = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#000000")).
	Background(lipgloss.Color(primaryColor))

var darkerText = lipgloss.NewStyle().
	Foreground(lipgloss.Color(disabledColor))

var border = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(primaryColor)).
	Width(50)

var disabledBorder = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(disabledColor)).
	Width(50)

type model struct {
	folderPath   string
	project      *Project
	currentTab   int
	input        string
	inputEnabled bool
}

func (m *model) renderInput() string {
	if !m.inputEnabled {
		return disabledBorder.Render("")
	} else {
		return border.Render("")
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q":
			fmt.Print("\033[2J\033[H")
			return m, tea.Quit
		case "tab":
			m.currentTab = (m.currentTab + 1) % tabsCount
		}
	}
	return m, nil
}

func (m *model) View() tea.View {
	// render navbar with lipgloss
	var navbar string
	for i := 0; i < tabsCount; i++ {
		if i == m.currentTab {
			navbar += highlightedTab.Render(fmt.Sprintf(" %s ", tabNames[i]))
		} else {
			navbar += regularTab.Render(fmt.Sprintf(" %s ", tabNames[i]))
		}
	}

	navbar += "\n\n" + m.renderInput()
	navbar += darkerText.Render("\n\n\ntab to switch tabs | q to quit")

	return tea.NewView(navbar)
}

func OpenEditor(folderPath string, project *Project) error {
	fmt.Print("\033[2J\033[H")
	m := &model{
		folderPath: folderPath,
		project:    project,
		currentTab: OverviewTab,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return nil
}
