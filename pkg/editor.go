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

const primaryColor = "#47e985"
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

var regularParam = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Width(24)

var highlightedParam = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Width(24).
	Bold(true)

func param(prop, value string, selected bool) string {
	var style lipgloss.Style
	if selected {
		style = highlightedParam
	} else {
		style = regularParam
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		style.Align(lipgloss.Left).Render("  "+prop+":"),
		style.Align(lipgloss.Right).Render(value+"  "),
	)
}

func overviewCard(color string) lipgloss.Style {
	var overViewCard = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(50).
		Height(25)
	return overViewCard
}

func stylesCard(color string) lipgloss.Style {
	var stylesCard = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(50).
		Height(25)
	return stylesCard
}

type model struct {
	folderPath     string
	project        *Project
	currentTab     int
	input          string
	inputEnabled   bool
	cursorOverview int
}

func (m *model) renderInput() string {
	if !m.inputEnabled {
		return disabledBorder.Render("")
	} else {
		return border.Render(m.input)
	}
}

func (m *model) renderCard() string {
	var color string
	if m.inputEnabled {
		color = disabledColor
	} else {
		color = primaryColor
	}
	switch m.currentTab {
	case OverviewTab:
		nameParam := param("Name", m.project.Name, m.cursorOverview == 0)
		return overviewCard(color).Render(nameParam)
	case StylesTab:
		return stylesCard(color).Render("")
	}
	return ""
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "q" || key == "Q" {
			fmt.Print("\033[2J\033[H")
			return m, tea.Quit
		}
		if m.inputEnabled {
			switch key {
			case "enter":
				m.project.Name = m.input
				m.inputEnabled = false
			case "esc":
				m.inputEnabled = false
			case "backspace":
				runes := []rune(m.input)
				if len(runes) > 0 {
					m.input = string(runes[:len(runes)-1])
				}
			case "space":
				m.input += " "
			default:
				if len([]rune(key)) == 1 {
					m.input += key
				}
			}
			break
		}
		switch msg.String() {
		case "tab":
			m.currentTab = (m.currentTab + 1) % tabsCount
		case "e":
			if m.currentTab == OverviewTab {
				m.input = m.project.Name
				m.inputEnabled = true
			}
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

	navbar += "\n\n" + m.renderCard()
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
