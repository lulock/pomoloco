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
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lulock/pomoloco/internal/styles"

	"net/http"
	"io"
	"encoding/json"
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
	randomQuote DailyQuote
	pomo bool
	pomoMessage string
	locoMessage string
	pomoCountdown time.Duration
	locoCountdown time.Duration
	pomoProgress progress.Model
	locoProgress progress.Model
	duration time.Duration
	percent  float64
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		
		m.pomoProgress.Width = msg.Width - padding * 3 - 6
		m.locoProgress.Width = msg.Width - padding * 3 - 6
		if m.pomoProgress.Width > maxWidth || m.locoProgress.Width > maxWidth {
			m.pomoProgress.Width, m.locoProgress.Width = maxWidth, maxWidth
		}
		return m, nil

	case tickMsg:
		if m.pomo {
			m.pomoCountdown -= 1 * time.Second
		} else {
			m.locoCountdown -= 1 * time.Second
		}
		m.percent -= float64(1.0/m.duration.Seconds())
		if m.percent < 0.0 {
			if m.pomo {
				m.pomo = false
				m.percent = 1.0
				m.duration = m.locoCountdown
				m.pomoCountdown = 0 * time.Second
			} else {
				m.percent = 0.0
				m.locoCountdown = 0 * time.Second
				return m, tea.Quit
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
	
	if m.pomo {
		message = m.pomoMessage
		mins = m.pomoCountdown / time.Minute
		sec = (m.pomoCountdown % time.Minute) / time.Second // remaining duration after subtracting full minutes / seconds gives remaining seconds
		progr = m.pomoProgress.ViewAs(m.percent)
	} else {
		message = m.locoMessage
		mins = m.locoCountdown / time.Minute
		sec = (m.locoCountdown % time.Minute) / time.Second // remaining duration after subtracting full minutes / seconds gives remaining seconds
		progr = m.locoProgress.ViewAs(m.percent)
	}
	
	pad := strings.Repeat(" ", padding)
	time := fmt.Sprintf("%02d:%02d", mins, sec)
	return "\n" +
		quote +
		pad + message + "\n\n" +
		pad + time + pad +  "*" +
		pad + progr + "\n\n" +
		pad + styles.HelpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomoloco",
	Short: "A CLI pomodoro timer with notetaking.",
	Long: `Start a pomodoro focus timer of 25 minutes and go loco with a 5 minute break.
	Make each pomodoro count by writing a summary or reflection. 
	Markdown files will be generated for each reflection and saved in your directory of choice.
	Example usage: 
		pomodoro pomo=25 loco=5 marinara=3

	focus for 25 minutes followed by a break of 5 minutes 3 times before prompted to write a reflection.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: a lot of duplicate code here that needs refactoring!
		pomoTime, _ := cmd.Flags().GetString("pomo")
		locoTime, _ := cmd.Flags().GetString("loco")	
	
		conftheme := viper.GetString("theme")
		theme := styles.ThemeLookup(conftheme)
		
		dat := DailyQuote{}
		resp, err := http.Get("https://zenquotes.io/api/random")

		if err != nil {
			// offline mode ... dat just stays empty for now
		} else {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &dat)
			if err != nil {
				fmt.Println("could not unmarshal??")
			}	
		}
	
		pomoProg := progress.New(progress.WithScaledGradient(theme.ColourOne, theme.ColourTwo), progress.WithoutPercentage())
		locoProg := progress.New(progress.WithScaledGradient(theme.ColourTwo, theme.ColourOne), progress.WithoutPercentage())
		pomoProg.SetPercent(1.0)
		locoProg.SetPercent(1.0)
		pomoText := fmt.Sprintf("Go go go! %s minutes of focus.", pomoTime)
		locoText := fmt.Sprintf("Go loco! %s-minute break.", locoTime)
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
		if _, err = tea.NewProgram(model{ randomQuote: dat, pomo: true, pomoMessage: pomoText, locoMessage: locoText, duration: pomoDur, pomoCountdown: pomoDur, locoCountdown: locoDur, percent: 1.0, pomoProgress: pomoProg, locoProgress: locoProg}).Run(); err != nil {
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
