package main

// A simple example that shows how to render a progress bar in a "pure"
// fashion. In this example we bump the progress by 25% every second,
// maintaining the progress state on our top level model using the progress bar
// model's ViewAs method only for rendering.
//
// The signature for ViewAs is:
//
//     func (m Model) ViewAs(percent float64) string
//
// So it takes a float between 0 and 1, and renders the progress bar
// accordingly. When using the progress bar in this "pure" fashion and there's
// no need to call an Update method.
//
// The progress bar is also able to animate itself, however. For details see
// the progress-animated example.

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func main() {
	
	sessionType := "default"
	text := "testing"

	if len(os.Args) > 1 { 
		sessionType = os.Args[1]
	}
	
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"), progress.WithoutPercentage())
	prog.SetPercent(1.0)
	durStr := "30s"
	
	if sessionType == "pomo"{
		durStr = "25m"
		text = "Go go go! 25 minutes of focus."
	}
	if sessionType == "loco"{
		durStr = "5m"
		text =  "Go loco! 5 minute break."
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		fmt.Println("Damn...", err)
		os.Exit(1)
	}
	if _, err = tea.NewProgram(model{message: text, duration: dur, countdown: dur, percent: 1.0, progress: prog}).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	message string
	duration time.Duration
	countdown time.Duration
	percent  float64
	progress progress.Model
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:

		m.percent -= float64(1.0/m.duration.Seconds())
		m.countdown -= 1*time.Second
		if m.percent < 0.0 {
			m.percent = 0.0
			m.countdown = 0 * time.Second
			return m, tea.Quit
		}
		return m, tickCmd()

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.message + "\n\n" +
		pad + m.countdown.String() + pad +  "*" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
