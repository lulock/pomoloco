package main

import (
	"fmt"
	"time"
	"strings"
	"os"
	"bufio"
	//"math"
	"strconv"
)

type cliCommand struct {
	name string
	description string
	handler func(...string) error
}

func cleanInput(text string) []string {
	text_trimmed := strings.TrimSpace(text)
	text_lowercase := strings.ToLower(text_trimmed)
	result := strings.Fields(text_lowercase)

	return result // removes whitespace, slice of all words in lowercase
}

func commandExit(args ...string) error {
	fmt.Println("closing this pomoloco session - goodbye!")
	os.Exit(0)
	return nil
}

func getDuration(s string) (int, error) {
	// currently, this function only expects a string indicating minutes
	// TODO: take a duration like 1h30m and parse using time module

	if len(s) > 0 {	
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return d, nil	
	}
	return 25, nil
}

func commandWork(args ...string) error {
	duration, err := getDuration(args[0])
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Pomo Go Go Go! %d minutes of focus.", duration))
	
	timeTicker := time.NewTicker(1 * time.Second)

	// we want 60 blocks only
	durInSeconds := duration * 60 
	secondsPerBlock := durInSeconds / 60
	
	defer timeTicker.Stop()
	block := strings.Repeat("█", 60)
	squashed := strings.Repeat("-", 0)
	minsLeft := durInSeconds / 60
	secondsLeft:= durInSeconds % 60
	fmt.Printf("\r\033[K%02d:%02d * %v", minsLeft, secondsLeft, block)

	for i := durInSeconds; i >= 0; {
		select {
		case <- timeTicker.C:
			block = strings.Repeat("█", i/secondsPerBlock)
			squashed = strings.Repeat("-", (durInSeconds/secondsPerBlock) - (i/secondsPerBlock))
			minsLeft = i / 60
			secondsLeft = i % 60

			fmt.Printf("\r\033[K%02d:%02d * %v%v", minsLeft, secondsLeft, block, squashed)
			i--

		}
	
	}
	fmt.Println()
	return nil
}

func commandBreak(args ...string) error {
	duration, err := getDuration(args[0])
	if err != err {
		return err
	}
	fmt.Println(fmt.Sprintf("Go loco! %d minute break.", duration))
	
	timeTicker := time.NewTicker(1 * time.Second)
	
	defer timeTicker.Stop()
	block := strings.Repeat("█", 60)
	durInSeconds := duration * 60
	secondsPerBlock := durInSeconds / 60

	squashed := strings.Repeat("-", 0)
	minsLeft := duration
	secondsLeft:= durInSeconds % 60
	fmt.Printf("\r\033[K%02d:%02d * %v", minsLeft, secondsLeft, block)

	for i := durInSeconds; i >= 0; {
		select {
		case <- timeTicker.C:
			block = strings.Repeat("█", i/secondsPerBlock)
			squashed = strings.Repeat("-", (durInSeconds/secondsPerBlock) - (i/secondsPerBlock))
			minsLeft = i / 60
			secondsLeft = i % 60

			fmt.Printf("\r\033[K%02d:%02d * %v%v", minsLeft, secondsLeft, block, squashed)
			i--

		}
	
	}
	fmt.Println()
	return nil
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	validCommands := map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit pomoloco",
			handler: commandExit,
		},
		"pomo": {
			name: "pomo",
			description: "start pomo work session",
			handler: commandWork,
		},
		"loco": {
			name: "loco",
			description: "start loco break session",
			handler: commandBreak,
		},
	}

	validCommands["help"] = cliCommand{
		name: "help",
		description: "Displays a help message",
		handler: func(args ...string) error {
			fmt.Println()
			fmt.Println("***************************")
			fmt.Println("  Welcome to pomoloco! 🍅")
			fmt.Println("***************************")
			fmt.Println()
			fmt.Println("Usage:")

			for _, v := range (validCommands) {
				fmt.Println(fmt.Sprintf(" * %v: %v", v.name, v.description))
			}
			fmt.Println()
			return nil
		},
	}
	for i := 0; ;i++ {
		fmt.Print("pomoloco > ")
		if scanner.Scan() {
			userinput := scanner.Text()
			cleanUserInput := cleanInput(userinput)
			if len(cleanUserInput) > 0 {
				command := cleanUserInput[0]
				args := ""

				if len(cleanUserInput) > 1 {
					args = cleanUserInput[1]
				}

				cmd, ok := validCommands[command]
				if !ok {
					fmt.Println("unknown command")
				} else {
					errs := make(chan error, 1)
					defer close(errs)

					go func() {
						errs <- cmd.handler(args)
					}()
					// later:

					if err := <-errs; err != nil {
						fmt.Println(err)// handle error
					}

				//	err := cmd.handler(args)
				//	if err != nil {
				//		fmt.Println(err)
				//	}
				}

			}
		}

	}
}
