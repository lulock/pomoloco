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

type Theme struct {
	ColourOne string
	ColourTwo string
}

var (
	WatermelonTheme = Theme {
		ColourOne: "#99ff99",
		ColourTwo: "#ff99ff",
	}
	SolarizedTheme = Theme {
		ColourOne : "#2aa198",
		ColourTwo : "#b58900",
	}
	RiverTheme = Theme {
		ColourOne : "#43cea2",
		ColourTwo : "#185a9d",
	}
	ShoreTheme = Theme {
		ColourOne : "#ffd194",
		ColourTwo : "#70e1f5",
	}
	BeachTheme = Theme {
		ColourOne : "#ffd194",
		ColourTwo : "#70e1f5",
	}
	MoonriseTheme = Theme {
		ColourOne : "#dae2f8",
		ColourTwo : "#d6a4a4",
	}
	DefaultTheme = Theme {
		ColourOne : "#ffffff",
		ColourTwo : "#ffffff",
	}
)

var Themes = map[string]Theme{
	"watermelon" : WatermelonTheme,
	"solarized" : SolarizedTheme,
	"river" : RiverTheme,
	"shore" : ShoreTheme,
	"beach" : BeachTheme,
	"moonrise" : MoonriseTheme,
	"default" : DefaultTheme,
}

func ThemeLookup(name string) Theme {
	if theme, ok := Themes[name]; ok {
		return theme
	}
	return DefaultTheme
}
