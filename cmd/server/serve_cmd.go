package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erikstmartin/erikbotdev/http"
	"github.com/spf13/cobra"
)

func init() {
	// runCmd.Flags().BoolVarP(
	// 	&forceStreamingOn,
	// 	"streaming-on",
	// 	"s",
	// 	false,
	// 	"Whether to force the bot to consider the stream on. Only valid if you don't have the 'OBS' module running",
	// )
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run the erikbotdev server",
	Long:  `Use this command to start up the chatbot server.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		// TODO: register routes
		portStr := fmt.Sprintf(":%s", port)
		log.Printf("Serving on port %s", portStr)
		go http.Start(portStr, "./web")

		// err := bot.InitDatabase(bot.DatabasePath(), 0600)
		// if err != nil {
		// 	if err.Error() == "timeout" {
		// 		log.Fatal("Timeout opening database. Check to ensure another process does not have the database file open")
		// 	}
		// 	log.Fatal("Failed to initialize database: ", err)
		// }

		// sig := make(chan os.Signal, 1)
		// signal.Notify(sig, os.Interrupt)
		// go func() {
		// 	<-sig

		// 	bot.ExecuteTrigger("bot::Shutdown", bot.Params{
		// 		Command: "shutdown",
		// 	})

		// 	if bot.IsModuleEnabled("OBS") {
		// 		obs.Disconnect()
		// 	}
		// 	os.Exit(0)
		// }()

		// // TODO: Handle scenario where startup trigger contains a twitch action
		// bot.ExecuteTrigger("bot::Startup", bot.Params{
		// 	Command: "startup",
		// })

		// if err := twitch.Run(); err != nil {
		// 	panic(err)
		// }
	},
}
