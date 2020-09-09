package main

import (
	"fmt"
	"os"

	_ "github.com/erikstmartin/erikbotdev/modules/keylight" // TODO: Remove this after we have cobra cmd
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	// initHueCmd()
}

var rootCmd = &cobra.Command{
	Use:   "erikbotdev",
	Short: "Twitch Bot",
	Long:  `Twitch bot for ErikDotDev`,
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
