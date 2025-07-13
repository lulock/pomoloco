/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// marinaraCmd represents the marinara command
var marinaraCmd = &cobra.Command{
	Use:   "marinara",
	Short: "Turn your pomodori into marinara with a reflective note-taking session.",
	Long: `Simmer your thoughts after focused work with this reflection / note-taking feature.
Jot down anything that was on your mind. What went well? What didn't? What did you learn?
Pomoloco records any notes imported as markdown (.md) files and saves them in the dir specified in your config.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("marinara called")
	},
}

func init() {
	rootCmd.AddCommand(marinaraCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// marinaraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// marinaraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
