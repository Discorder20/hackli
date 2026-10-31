package pkg

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var tabNames = [tabsCount]string{"Overview", "Styles"}
var paramsOverviewPointerList []*string

const (
	OverviewTab = iota
	StylesTab
)

const tabsCount = 2

var currentTabOptions = -1

func textParam(prop, value string, selected bool) string {
	var style lipgloss.Style
	if selected {
		style = SelectedParam
	} else {
		style = RegularInput
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
		return InputDisabled.Render("")
	} else {
		return Input.Render(m.input)
	}
}

func (m *model) renderCard() string {
	var color string
	if m.inputEnabled {
		color = DisabledColor
	} else {
		color = PrimaryColor
	}
	switch m.currentTab {
	case OverviewTab:
		currentTabOptions = 2
		nameParam := textParam("Name", m.project.Name, m.cursorOverview == 0)
		solutionParam := textParam("Url", m.project.Url, m.cursorOverview == 1)
		return overviewCard(color).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				nameParam,
				solutionParam,
			),
		)
	case StylesTab:
		currentTabOptions = 0
		return stylesCard(color).Render("")
	}
	return ""
}

func (m *model) HandleInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		switch key {
		case "enter":
			*paramsOverviewPointerList[m.cursorOverview] = m.input
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
	case tea.PasteMsg:
		if m.inputEnabled {
			m.input += msg.Content
		}
	}
	return m, nil
}

func (m *model) HandleMove(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "q" || key == "Q" {
			fmt.Print("\033[2J\033[H")
			return m, tea.Quit
		}
		switch msg.String() {
		case "tab":
			m.currentTab = (m.currentTab + 1) % tabsCount
		case "e":
			if m.currentTab == OverviewTab {
				m.input = *paramsOverviewPointerList[m.cursorOverview]
				m.inputEnabled = true
			}
		case "ctrl+c":
			fmt.Print("\033[2J\033[H")
			return m, tea.Quit
		case "down":
			if m.currentTab == OverviewTab {
				m.cursorOverview = (m.cursorOverview + 1) % currentTabOptions
			}
		case "up":
			if m.currentTab == OverviewTab {
				m.cursorOverview = (m.cursorOverview - 1 + currentTabOptions) % currentTabOptions
			}
		}
	}
	return m, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.inputEnabled {
		return m.HandleInput(msg)
	} else {
		return m.HandleMove(msg)
	}
}

func (m *model) View() tea.View {
	// render navbar with lipgloss
	var navbar string
	for i := 0; i < tabsCount; i++ {
		if i == m.currentTab {
			navbar += TabHighlighted.Render(fmt.Sprintf(" %s ", tabNames[i]))
		} else {
			navbar += TabPrimary.Render(fmt.Sprintf(" %s ", tabNames[i]))
		}
	}

	navbar += "\n\n" + m.renderCard()
	navbar += "\n\n" + m.renderInput()
	navbar += Dark.Render("\n\n\ntab to switch tabs | q to quit")

	return tea.NewView(navbar)
}

func OpenEditor(folderPath string, project *Project) error {
	fmt.Print("\033[2J\033[H")
	m := &model{
		folderPath: folderPath,
		project:    project,
		currentTab: OverviewTab,
	}
	paramsOverviewPointerList = []*string{
		&m.project.Name,
		&m.project.Url,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return nil
}
