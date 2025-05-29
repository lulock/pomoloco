package main

import (
	"fmt"
	"os"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
)

const (
	padding = 2   // padding arounf bar
	maxWidth = 80 // max width of bar
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type model struct {
	choices  []string           //items on the todo list
	cursor   int                // which todo list item our cursor is pointing at
	selected map[int]struct{}   // which todo items are selected
	progress progress.Model     // track progress lol
}

func initialModel() model {
	m := model{
		// our todo list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
		progress: progress.New(progress.WithDefaultGradient()),
	}
	
	return m
}

func (m model) Init() tea.Cmd {
	return nil // no I/O right now please
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit	

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
			return m, nil

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
				cmd := m.progress.DecrPercent((1.0/float64(len(m.choices))))
				return m, cmd
			} else {
				m.selected[m.cursor] = struct{}{} // just a signal to select? to return ok
				cmd := m.progress.IncrPercent((1.0/float64(len(m.choices))))

				return m, cmd
			}
		}
		
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4 
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
	return m, nil

}

func (m model) View() string {
	// header
	s := "What should we buy at Ica?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		//render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	pad := strings.Repeat(" ", padding)
	s += fmt.Sprintf("\n %s%s \n\n %s", pad, m.progress.View(),pad)
	
	// footer
	s += "\nPress q to quit.\n"

	return s
}


func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}

}
