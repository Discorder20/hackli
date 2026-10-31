package pkg

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	OverviewTab = iota
	StylesTab
)
const tabsCount = 2

var tabNames = [tabsCount]string{"Overview", "Styles"}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

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

func selectedOverviewField(m *model) *string {
	switch m.cursorOverview {
	case 0:
		return &m.project.Name
	case 1:
		return &m.project.Url
	default:
		panic("invalid overview cursor")
	}
}

type model struct {
	folderPath     string
	project        *Project
	currentTab     int
	input          string
	inputEnabled   bool
	cursorOverview int
}

func (m *model) optionCount() int {
	switch m.currentTab {
	case OverviewTab:
		return 2
	default:
		return 0
	}
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
		nameParam := textParam("Name", m.project.Name, m.cursorOverview == 0)
		urlParam := textParam("Url", m.project.Url, m.cursorOverview == 1)
		return OverviewCardStyle(color).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				nameParam,
				urlParam,
			),
		)
	case StylesTab:
		return StylesCardStyle(color).Render("")
	}
	return ""
}

func (m *model) HandleInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		switch key {
		case "enter":
			*selectedOverviewField(m) = m.input
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
		switch msg.String() {
		case "tab":
			m.currentTab = (m.currentTab + 1) % tabsCount
		case "e":
			if m.currentTab == OverviewTab {
				m.input = *selectedOverviewField(m)
				m.inputEnabled = true
			}
		case "down":
			if m.currentTab == OverviewTab {
				m.cursorOverview = (m.cursorOverview + 1) % m.optionCount()
			}
		case "up":
			if m.currentTab == OverviewTab {
				m.cursorOverview = (m.cursorOverview - 1 + m.optionCount()) % m.optionCount()
			}
		}
	}
	return m, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok && (key.String() == "ctrl+c" || key.String() == "q" || key.String() == "Q") {
		return m, tea.Quit
	}
	if m.inputEnabled {
		return m.HandleInput(msg)
	} else {
		return m.HandleMove(msg)
	}
}

func (m *model) View() tea.View {
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
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("run editor: %w", err)
	}
	return nil
}
