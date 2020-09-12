package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/erikstmartin/erikbotdev/bot"
	"github.com/erikstmartin/erikbotdev/http"
	"github.com/erikstmartin/erikbotdev/modules/twitch"
	twitchcl "github.com/gempir/go-twitch-irc/v2"
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
		bot.Status.Streaming = true
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		twitchMainChannel := os.Getenv("TWITCH_MAIN_CHANNEL")
		twitchOauthToken := os.Getenv("TWITCH_OAUTH_TOKEN")

		if twitchMainChannel == "" {
			log.Fatalf("No TWITCH_MAIN_CHANNEL set")
		}
		if twitchOauthToken == "" {
			log.Fatalf("No TWITCH_OAUTH_TOKEN set")
		}

		twitchClient := twitchcl.NewClient(twitchMainChannel, twitchOauthToken)
		bot.StartCountersListTimer(
			twitchClient,
			[]string{"PEBKAC", "credleak"},
			30*time.Second,
		)
		go func() {
			if err := twitch.Run(twitchClient); err != nil {
				log.Fatalf("Error running twitch chat listener (%s)", err)
			}
		}()

		go func() {
			portStr := fmt.Sprintf(":%s", port)
			log.Printf("Serving on port %s", portStr)
			http.Start(portStr, "./web")
		}()

		select {}

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

	},
}
