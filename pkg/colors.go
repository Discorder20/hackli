package pkg

import "charm.land/lipgloss/v2"

const PrimaryColor = "#47e985"
const DisabledColor = "#535353"
const NeutralColor = "#FAFAFA"
const SelectedColor = "#4fcdff"

var TabPrimary = lipgloss.NewStyle().
	Foreground(lipgloss.Color(NeutralColor))

var TabHighlighted = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Black).
	Background(lipgloss.Color(PrimaryColor))

var Dark = lipgloss.NewStyle().
	Foreground(lipgloss.Color(DisabledColor))

var Input = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(PrimaryColor)).
	Width(50)

var InputDisabled = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(DisabledColor)).
	Width(50)

var RegularInput = lipgloss.NewStyle().
	Foreground(lipgloss.Color(NeutralColor)).
	Width(24)

var SelectedParam = lipgloss.NewStyle().
	Foreground(lipgloss.Color(SelectedColor)).
	Width(24).
	Bold(true)
