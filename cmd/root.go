/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	//"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lulock/pomoloco/internal/styles"

	"net/http"
	"io"
	"encoding/json"
	"github.com/gen2brain/beeep"
)

type DailyQuote []struct {
	Quote string `json:"q"`
	Author string `json:"a"`
	H string `json:"h"`
}

// Building upon Bubbleatea's simple rendering of a progrerss bar in a "pure" fashion. 
// This is a pomodoro app that generates visual countdowns for pomo "focus" sessions and loco "breaks"
const (
	padding  = 2
	maxWidth = 80
)

type tickMsg time.Time

type model struct {
	theme styles.Theme
	randomQuote DailyQuote
	pomo bool
	message string
	pomoDuration time.Duration
	locoDuration time.Duration
	start time.Time
	timeLeft time.Duration
	progressBar progress.Model
	percent float64
}

// type model struct {
// 	randomQuote DailyQuote
// 	pomo bool
//
// 	pomoMessage string
// 	locoMessage string
// 	initPomoCountdown time.Duration
// 	initLocoCountdown time.Duration
// 	pomoCountdown time.Duration
// 	locoCountdown time.Duration
// 	pomoProgress progress.Model
// 	locoProgress progress.Model
// 	duration time.Duration
// 	percent  float64
// }
//
func newModel(pomoDur, locoDur time.Duration, theme styles.Theme) model {
	pomoText := "Go go go! Time to focus."

	prog := progress.New(progress.WithScaledGradient(theme.ColourOne, theme.ColourTwo), progress.WithoutPercentage())

	prog.SetPercent(1.0)

	m := model{
		theme: theme,
	//	randomQuote: quote, 
		pomo: true, 
		message: pomoText, 
		pomoDuration: pomoDur, 
		locoDuration: locoDur,
		timeLeft: pomoDur,
		progressBar: prog,
		percent: 1.0,
		start: time.Now(),
	}

	return m

}

func (m *model) nextSession() {

	width := m.progressBar.Width
	if m.pomo {
		m.pomo = false
		m.timeLeft, _ = time.ParseDuration("0s")
		m.percent = 0.0
		m.message = "Go loco! Time for a break."
		m.progressBar = progress.New(progress.WithScaledGradient(m.theme.ColourTwo, m.theme.ColourOne), progress.WithoutPercentage())
		m.progressBar.SetPercent(0.0)
	} else {
		m.pomo = true
		m.timeLeft = m.pomoDuration
		m.percent = 1.0
		m.message = "Go go go! Time to focus."
		m.progressBar = progress.New(progress.WithScaledGradient(m.theme.ColourOne, m.theme.ColourTwo), progress.WithoutPercentage())
		m.progressBar.SetPercent(1.0)
	}
	m.progressBar.Width = width
	m.start = time.Now()
} 

func (m *model) notify() {
	session := ""
	if m.pomo {
		session = "Pomodoro session"
	} else {
		session = "Break"
	}

	beeep.AppName = "Pomoloco"
	
	err := beeep.Notify("Times up!", fmt.Sprintf("%s is over.", session), "./internal/imgs/catmato.png")

	if err != nil {
		// ignore ... no notification required
	}
}

// because m implements Init from the tea.Model interface ... it's a tea.Model 
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		getQuote(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DailyQuote:
		m.randomQuote = msg
		return m, nil
	case tea.KeyMsg:
		 // Cool, what was the actual key pressed?
        	switch msg.String() {
		// These keys should exit the program.
        	case "ctrl+c", "q", "esc":
            		return m, tea.Quit
		case "n", "enter":
			m.nextSession()
			return m, nil
		case "r":
			
			return m, getQuote()
		default:
			return m, nil
		}

	case tea.WindowSizeMsg:
		
		m.progressBar.Width = msg.Width - padding * 3 - 6

		if m.progressBar.Width > maxWidth {
			m.progressBar.Width = maxWidth
		}

		return m, nil

	case tickMsg:
		
		if m.pomo {
			m.timeLeft -= 1 * time.Second
			m.percent -= float64(1.0/m.pomoDuration.Seconds())
			if m.percent < 0.0 {
				m.notify()
				m.nextSession()
			}
		} else {
			m.timeLeft += 1 * time.Second
			m.percent += float64(1.0/m.locoDuration.Seconds())

			if m.percent > 1.0 {
				// wait until user starts a new session
				if !strings.Contains(m.message, "Break is over.") {
					m.notify()	
				} // else already notified.
				m.percent = 1.0
				//m.timeLeft, _ = time.ParseDuration("0s")
				m.message = "Break is over. Press enter to start a new session."
				
			}
		}

		return m, tickCmd()

	default:
		return m, nil
	}
}

// the main view should call the appropriate subview
func (m model) View() string {
	var message string
	var mins time.Duration
	var sec time.Duration
	var progr string
	var quote string

	if len(m.randomQuote) > 0{
		quote = styles.QuoteStyle(m.randomQuote[0].Quote + fmt.Sprintf("\n  -- %s", m.randomQuote[0].Author)) + "\n\n"
	} else {
		quote = ""
	}
	
	message = m.message
	mins = m.timeLeft / time.Minute
	sec = (m.timeLeft % time.Minute) / time.Second // remaining duration after subtracting full minutes / seconds gives remaining seconds
	progr = m.progressBar.ViewAs(m.percent)
	
	pad := strings.Repeat(" ", padding)
	time := fmt.Sprintf("%02d:%02d", mins, sec)
	return "\n" +
		quote +
		pad + message + "\n\n" +
		pad + time + pad +  "*" +
		pad + progr + "\n\n" +
		pad + styles.HelpStyle("esc to quit * enter to skip to next")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func getQuote() tea.Cmd {
	quote := DailyQuote{}
	resp, err := http.Get("https://zenquotes.io/api/random")

	if err != nil {
	// offline mode ... quote stays empty
	} else {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &quote)
		if err != nil {
			fmt.Println("could not unmarshal??")
		}	
	}
	return func() tea.Msg {
		return quote
	}
}
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomoloco",
	Short: "A CLI pomodoro timer with notetaking.",
	Long: `Start a pomodoro focus timer of 25 minutes and go loco with a 5 minute break.
	Make each pomodoro count by writing a summary or reflection. 
	Markdown files will be generated for each reflection and saved in your directory of choice.
	Example usage: 
		pomodoro --pomo=25 --loco=5 --theme=watermelon

	focus for 25 minutes followed by a break of 5 minutes.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: a lot of duplicate code here that needs refactoring!

		pomoTime, _ := cmd.Flags().GetString("pomo")
		locoTime, _ := cmd.Flags().GetString("loco")	

		conftheme := viper.GetString("theme")
		theme := styles.ThemeLookup(conftheme)

		pomoDur, err := time.ParseDuration(pomoTime + "m")
		if err != nil {
			fmt.Println("Damn...", err)
			os.Exit(1)
		}
		
		locoDur, err := time.ParseDuration(locoTime + "m")
		if err != nil {
			fmt.Println("Damn...", err)
			os.Exit(1)
		}
		
		m := newModel(pomoDur, locoDur, theme) 

		if _, err = tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Oh no!", err)
			os.Exit(1)
		}	
	},
}

var pomoduration string
var locoduration string
var theme string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Print("oops")
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pomoloco.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&pomoduration, "pomo", "p", "25", "duration of focus")
	rootCmd.PersistentFlags().StringVarP(&locoduration, "loco", "l", "5", "duration of break")
	rootCmd.PersistentFlags().StringVarP(&theme, "theme", "t", "", "watermelon")

	viper.BindPFlag("theme", rootCmd.PersistentFlags().Lookup("theme"))
	viper.AddConfigPath(".")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("yaml")   // Look for specific type
	viper.ReadInConfig()	
}
