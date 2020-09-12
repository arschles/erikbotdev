package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

var builtinCommands map[string]CommandFunc = map[string]CommandFunc{
	"help":     helpCmd,
	"commands": helpCmd,
	"me":       userInfoCmd,
	"props":    givePointsCmd,
	"sounds":   soundListCmd,
	"so":       shoutoutCmd,
	"counters": listCountersCmd,
	// "rickroll": rickrollCommand,
}

func TwitchSay(cmd Params, msg string) error {
	args := map[string]string{
		"channel": cmd.Channel,
		"message": msg,
	}
	log.Printf("TwitchSay with args %v", args)
	return ExecuteAction("twitch", "Say", args, cmd)
}

func StartCountersListTimer(
	twitchClient *twitch.Client,
	counterNames []string,
	duration time.Duration,
) {
	countersChan := make(chan []CounterData)
	go CountersListTimer(counterNames, countersChan)
	go func() {
		for {
			counters := <-countersChan
			log.Printf("About to print all the counters after %s", duration)
			for _, counter := range counters {
				timeSince := time.Now().Sub(counter.TimeUpdated)
				counterStr := fmt.Sprintf("time since last %s: %s",
					counter.Name,
					timeSince,
				)
				log.Println(counterStr)
				twitchClient.Say(
					"arschles",
					counterStr,
				)
			}
		}
	}()

}

func helpCmd(cmd Params) error {
	if len(cmd.CommandArgs) > 0 {
		cname := cmd.CommandArgs[0]
		if c, ok := config.Commands[cname]; ok {
			return TwitchSay(cmd, fmt.Sprintf("%s: %s", cname, c.Description))
		}

		return nil
	}

	cmds := make([]string, 0)
	for _, c := range config.Commands {
		if c.Enabled {
			cmds = append(cmds, c.Name)
		}
	}

	return TwitchSay(cmd, strings.Join(cmds, ", "))
}

func userInfoCmd(cmd Params) error {
	u, err := GetUser(cmd.UserID)
	if err != nil {
		return err
	}

	return TwitchSay(cmd, fmt.Sprintf("%s: %d points", u.DisplayName, u.Points))
}

func givePointsCmd(cmd Params) error {
	if len(cmd.CommandArgs) != 2 {
		return nil
	}

	user, err := GetUser(cmd.UserID)
	if err != nil {
		return err
	}

	points, err := strconv.ParseUint(cmd.CommandArgs[1], 10, 64)
	if err != nil {
		return err
	}

	recipient := strings.TrimPrefix(cmd.CommandArgs[0], "@")
	twitchUser, err := GetUserByName(recipient)
	if err != nil {
		return nil
	}

	// Allow owner to give unlimited points
	if strings.ToLower(user.DisplayName) == strings.ToLower(cmd.Channel) {
		destUser, err := GetUser(twitchUser.ID)
		if err != nil {
			return err
		}
		destUser.GivePoints(points)
	} else {
		user.TransferPoints(points, twitchUser.ID)
	}

	return nil
}

func soundListCmd(cmd Params) error {
	files, err := ioutil.ReadDir(MediaPath())
	if err != nil {
		return err
	}

	sounds := make([]string, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			name := filepath.Base(f.Name())
			sounds = append(sounds, strings.TrimSuffix(name, filepath.Ext(name)))
		}
	}

	return TwitchSay(cmd, "sounds: "+strings.Join(sounds, ", "))
}

func listCountersCmd(cmd Params) error {
	countersList := ListCounters()
	counterStrings := make([]string, len(countersList))
	for i, counter := range countersList {
		counterStrings[i] = fmt.Sprintf("%s: %d", counter.Name, counter.Value)
	}
	return TwitchSay(cmd, "counters: "+strings.Join(counterStrings, "\n"))
}

// TODO; Hit Twitch API and ensure user exists
func shoutoutCmd(cmd Params) error {
	if len(cmd.CommandArgs) > 0 {
		user := cmd.CommandArgs[0]
		return TwitchSay(cmd, fmt.Sprintf("Shoutout %s! Check out their channel, shower them with follows and subs: https://twitch.tv/%s", user, user))
	}

	return fmt.Errorf("username is required")
}
