package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	if len(os.Args) != 3 {
		fmt.Println("This program takes exactly two arguments, the ID of the server you want to scrape, and the level you want to scrape until!")
		os.Exit(1)
	}
	guildId, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Server ID must be an int!")
		os.Exit(1)
	}
	levelToScrapeUntil, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		fmt.Println("Server ID must be an int!")
		os.Exit(1)
	}
	joutName := fmt.Sprintf("./%d-levels.json", guildId)
	jout, err := os.OpenFile(joutName, os.O_CREATE|os.O_WRONLY, 0644)
	report(err)
	page := 0
	users := make([]IntPlayer, 0)
	keepGoing := true
	hadError := false
	go func() {
		<-c
		fmt.Println("\nExiting...")
		keepGoing = false
	}()
	for keepGoing {
		thisPage := new(InputData)
		resp, err := http.Get(fmt.Sprintf("https://mee6.xyz/api/plugins/levels/leaderboard/%d?limit=1000&page=%d", guildId, page))
		if err != nil {
			if resp.StatusCode == 429 {
				dur, err := strconv.Atoi(resp.Header.Get("Retry-After"))
				if err != nil {
					fmt.Printf("\nError: %s\n", err.Error())
					hadError = true
					break
				}
				time.Sleep(time.Duration(dur+1) * time.Second)
				continue
			}
			fmt.Printf("\nError fetching page %d: %s\n", page, err.Error())
			hadError = true
			break
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("\nError: %s with response: %s\n", err.Error(), string(body[:]))
			hadError = true
			break
		}
		err = json.Unmarshal(body, &thisPage)
		if err != nil {
			fmt.Printf("\nError: %s with response: %s\n", err.Error(), string(body[:]))
			hadError = true
			break
		}
		if len(thisPage.Players) == 0 {
			break
		}
		for _, player := range thisPage.Players {
			id, err := strconv.ParseUint(player.ID, 10, 64)
			intPlayer := IntPlayer{ID: id, Level: player.Level, Xp: player.Xp, Discriminator: player.Discriminator, Username: player.Username, Avatar: player.Avatar}
			if err != nil {
				fmt.Printf("\nError converting %s to int: %s\n", player.ID, err.Error())
				hadError = true
				break
			}
			users = append(users, intPlayer)
		}
		finalLevel := thisPage.Players[len(thisPage.Players)-1].Level
		if finalLevel < levelToScrapeUntil {
			break
		}
		fmt.Printf("\r Current user level: %d (%d total users)", finalLevel, len(users))

		page = page + 1
		time.Sleep(1 * time.Second)

	}
	fmt.Printf("\nHave %d users, writing to disk..\n", len(users))
	usersJson, err := json.MarshalIndent(users, "", "\t")
	report(err)
	_, err = jout.Write(usersJson)
	report(err)
	err = jout.Sync()
	report(err)
	err = jout.Close()
	report(err)
	fmt.Printf("Done! Data written to %s\n", joutName)
	if hadError {
		os.Exit(1)
	}
}

func report(e error) {
	if e != nil {
		fmt.Printf("There was an error: %s\n", e.Error())
		os.Exit(1)
	}
}

type InputData struct {
	Players []Player `json:"players"`
}

type Player struct {
	ID            string `json:"id"`
	Level         uint64 `json:"level"`
	Xp            uint64 `json:"xp"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
        Avatar        string `json:"avatar,omitempty"`
}

type IntPlayer struct {
	ID            uint64 `json:"id"`
	Level         uint64 `json:"level"`
	Xp            uint64 `json:"xp"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
        Avatar        string `json:"avatar,omitempty"`
}
