package styles

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var QuoteStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#626262")).
	Padding(1).
	PaddingLeft(2).
	PaddingRight(2).
	MarginLeft(2).
	MarginRight(2).
	Width(maxWidth/2).Render
