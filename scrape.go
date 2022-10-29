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
	if len(os.Args) != 2 {
		fmt.Println("This program takes exactly one argument, the ID of the server you want to scrape!")
		os.Exit(1)
	}
	guildId, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Server ID must be an int!")
		os.Exit(1)
	}
	joutName := fmt.Sprintf("./%d-levels.json", guildId)
	soutName := fmt.Sprintf("./%d-levels.sql", guildId)
	jout, err := os.OpenFile(joutName, os.O_CREATE|os.O_WRONLY, 0644)
	report(err)
	sout, err := os.OpenFile(soutName, os.O_CREATE|os.O_WRONLY, 0644)
	report(err)
	page := 0
	users := make([]Player, 0)
	keepGoing := true
	go func() {
		<-c
		fmt.Println("\nExiting...")
		keepGoing = false
	}()
	for keepGoing {
		thisPage := new(InputData)
		resp, err := http.Get(fmt.Sprintf("https://mee6.xyz/api/plugins/levels/leaderboard/%d?page=%d", guildId, page))
		if err != nil {
			fmt.Printf("Error fetching page %d: %s", page, err.Error())
			break
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &thisPage)
		if err != nil {
			fmt.Printf("Error: %s with response: %s", err.Error(), body)
			break
		}
		for _, player := range thisPage.Players {
			id, err := strconv.ParseUint(player.ID, 10, 64)
			if err != nil {
				fmt.Printf("Error converting %s to int: %s", player.ID, err.Error())
				break
			}
			sout.WriteString(fmt.Sprintf("INSERT INTO levels (id, xp, guild) VALUES (%d, %d, %d);\n", id, player.Xp, guildId))
			users = append(users, player)
			fmt.Printf("\r Current user level: %d (%d total users)", player.Level, len(users))
		}
		if thisPage.Players[len(thisPage.Players)-1].Level < 5 {
			break
		}
		page = page + 1
		time.Sleep(1 * time.Second)

	}
	fmt.Printf("Have %d users, writing to disk..\n", len(users))
	usersJson, err := json.MarshalIndent(users, "", "\t")
	report(err)
	_, err = jout.Write(usersJson)
	report(err)
	err = jout.Sync()
	report(err)
	err = sout.Sync()
	report(err)
	err = jout.Close()
	report(err)
	err = sout.Close()
	report(err)
	fmt.Printf("Done! Data written to %s, SQL queries written to %s\n", joutName, soutName)
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
	ID    string `json:"id"`
	Level uint64 `json:"level"`
	Xp    uint64 `json:"xp"`
}
