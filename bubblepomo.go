package main

// Building upon Bubbleatea's simple rendering of a progrerss bar in a "pure" fashion. 
// This is a pomodoro app that generates visual countdowns for pomo "focus" sessions and loco "breaks"

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
	
	//prog := progress.New(progress.WithScaledGradient("#ff9933", "#6600cc"), progress.WithoutPercentage())
	prog := progress.New(progress.WithScaledGradient("#99ff99", "#ff99ff"), progress.WithoutPercentage())
	//prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"), progress.WithoutPercentage())
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
		m.progress.Width = msg.Width - padding*3 - 6
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
	mins := m.countdown / time.Minute
	sec := (m.countdown % time.Minute) / time.Second // remaining duration after subtracting full minutes / seconds gives remaining seconds
	time := fmt.Sprintf("%02d:%02d", mins, sec)
	return "\n" +
		pad + m.message + "\n\n" +
		pad + time + pad +  "*" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
