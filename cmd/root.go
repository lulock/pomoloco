package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)
var rootCmd = &cobra.Command{
	Use:   "session",
	Short: "A CLI pomodoro app. Go loco after each pomo.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Choose a colour theme.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	
	themeCmd.Flags().StringP(
		"theme",
		"t",
		"",
		"specify a theme for your session",
	)
	rootCmd.AddCommand(themeCmd)
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Oops. An error while executing Pomoloco '%s'\n", err)
        os.Exit(1)
    }
}


