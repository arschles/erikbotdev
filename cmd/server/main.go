package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/erikstmartin/erikbotdev/bot"
	_ "github.com/erikstmartin/erikbotdev/modules/bot"
)

func main() {
	file, err := os.Open(findConfigFile())
	if err != nil {
		fmt.Println(err)
		return
	}
	err = bot.LoadConfig(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = bot.Init(); err != nil {
		fmt.Println(err)
		return
	}

	err = bot.InitDatabase(bot.DatabasePath(), 0600)
	if err != nil {
		if err.Error() == "timeout" {
			log.Fatal("Timeout opening database. Check to ensure another process does not have the database file open")
		}
		log.Fatal("Failed to initialize database: ", err)
	}

	execute()
}

func findConfigFile() string {
	configFileName := os.Getenv("ERIKBOTDEV_CONFIG_FILE_NAME")
	if configFileName == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Couldn't find home directory")
		}
		return filepath.Join(home, "config.json")
	}
	log.Printf("Using config %s", configFileName)
	return configFileName
}
