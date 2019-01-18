package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type PlayersResults struct {
	Data struct {
		Players struct {
			Results []struct {
				Player struct {
					PlayerID string `json:"playerId"`
					Handle   string `json:"handle"`
				} `json:"player"`
			} `json:"results"`
		} `json:"players"`
	} `json:"data"`
}

type PlayerData struct {
	Data struct {
		Player struct {
			ID       string `json:"id"`
			Metadata []struct {
				Key   string      `json:"key"`
				Name  string      `json:"name"`
				Value interface{} `json:"value"`
			} `json:"metadata"`
			Stats []struct {
				Metadata struct {
					Key        string `json:"key"`
					Name       string `json:"name"`
					IsReversed bool   `json:"isReversed"`
				} `json:"metadata"`
				Value float64 `json:"value"`
			} `json:"stats"`
			Segments []struct {
				Metadata []struct {
					Key   string `json:"key"`
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"metadata"`
				Stats []struct {
					Metadata struct {
						Key        string `json:"key"`
						Name       string `json:"name"`
						IsReversed bool   `json:"isReversed"`
					} `json:"metadata"`
					Value float64 `json:"value"`
				} `json:"stats"`
			} `json:"segments"`
		} `json:"player"`
	} `json:"data"`
}

var (
	Token          = "Bot NTM1MTc4MTA2NDA3NzQ3NjA0.DyEbaQ.3eORw3CVDCqn4NbR66VF6wDpwxA"
	BotName        = "<@535178106407747604>"
	stopBot        = make(chan bool)
	vcsession      *discordgo.VoiceConnection
	HelloWorld     = "!helloworld"
	GetId          = "!getId"
	GetSquadWeekly = "!getSquadWeekly"
	GetStats       = "!getStats"
	Search         = "!search"
)

func main() {
	discord, err := discordgo.New()
	discord.Token = Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	<-stopBot
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, HelloWorld)):
		sendMessage(s, c, "Hello world！")
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, GetId)):
		messageContentList := strings.Fields(m.Content)
		if len(messageContentList) > 2 {
			sendMessage(s, c, string(getId(messageContentList[2])))
		} else {
			sendMessage(s, c, "Input your username! ex.'@fnlc getId {YOUR USERNAME}")
		}
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, Search)):
		messageContentList := strings.Fields(m.Content)
		if len(messageContentList) > 2 {
			sendMessage(s, c, string(getSearch(messageContentList[2])))
		} else {
			sendMessage(s, c, "Input your username! ex.'@fnlc getId {YOUR USERNAME}")
		}
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, GetSquadWeekly)):
		messageContentList := strings.Fields(m.Content)
		if len(messageContentList) > 2 {
			head := "SquadWeekly\n"
			kd := "K/D: " + string(getSquadWeekly(messageContentList[2]))
			message := head + kd
			sendMessage(s, c, message)
		} else {
			sendMessage(s, c, "Input your userid! ex.'@fnlc getId {YOUR USERID}")
		}
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, GetStats)):
		messageContentList := strings.Fields(m.Content)
		if len(messageContentList) > 2 {
			head := messageContentList[2] + "'s Stats:doughnut:\n"
			stats := string(getStats(getId(messageContentList[2])))
			message := head + stats
			sendMessage(s, c, message)
		} else {
			sendMessage(s, c, "Input your username! ex.'@fnlc getId {YOUR USERNAME}")
		}
	}
}

func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func getId(userName string) string {
	url := "https://api.scoutsdk.com/graph"
	var jsonStr = []byte(`{"query":"query id($title: String, $platform: String, $identifier: String) { players(title: $title, platform: $platform, identifier: $identifier) {results {player { playerId handle } } } }","variables":{ "title": "fortnite", "identifier": "` + userName + `", "platform": "epic"}}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Host", "api.scoutsdk.com")
	req.Header.Set("Accept", "application/com.scoutsdk.graph+json; version=1.1.0; charset=utf8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "313")
	req.Header.Set("Scout-App", "1fbb8b74-2a24-4855-9a64-60fb4acfdd78")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	var r PlayersResults
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
	}

	return string(r.Data.Players.Results[0].Player.PlayerID)
}

func getSearch(userName string) string {
	url := "https://api.scoutsdk.com/graph"
	var jsonStr = []byte(`{"query":"query id($title: String, $platform: String, $identifier: String) { players(title: $title, platform: $platform, identifier: $identifier) {results {player { playerId handle } } } }","variables":{ "title": "fortnite", "identifier": "` + userName + `", "platform": "epic"}}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Host", "api.scoutsdk.com")
	req.Header.Set("Accept", "application/com.scoutsdk.graph+json; version=1.1.0; charset=utf8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "313")
	req.Header.Set("Scout-App", "1fbb8b74-2a24-4855-9a64-60fb4acfdd78")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	var r PlayersResults
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
	}

	return string(body)
}

func getSquadWeekly(userId string) string {
	url := "https://api.scoutsdk.com/graph"
	var jsonStr = []byte(`{"query":"query player($title: String, $id: String, $segment: String) { player(title: $title, id: $id, segment: $segment)  { id metadata { key name value displayValue } stats { metadata { key name isReversed } value displayValue } segments { metadata { key name value displayValue } stats { metadata { key name isReversed } value displayValue } } }}","variables":{ "title": "fortnite", "segment": "p2.br.m0.weekly", "id": "` + userId + `"}}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Host", "api.scoutsdk.com")
	req.Header.Set("Accept", "application/com.scoutsdk.graph+json; version=1.1.0; charset=utf8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "313")
	req.Header.Set("Scout-App", "1fbb8b74-2a24-4855-9a64-60fb4acfdd78")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	var d PlayerData
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &d)
	if err != nil {
	}

	s := fmt.Sprintf("%f", d.Data.Player.Stats[8].Value)
	return s
}

func getStats(userId string) string {
	s := ""
	url := "https://api.scoutsdk.com/graph"
	var jsonStr = []byte(`{"query":"query player($title: String, $id: String, $segment: String) { player(title: $title, id: $id, segment: $segment)  { id metadata { key name value displayValue } stats { metadata { key name isReversed } value displayValue } segments { metadata { key name value displayValue } stats { metadata { key name isReversed } value displayValue } } }}","variables":{ "title": "fortnite", "segment": "*", "id": "` + userId + `"}}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Host", "api.scoutsdk.com")
	req.Header.Set("Accept", "application/com.scoutsdk.graph+json; version=1.1.0; charset=utf8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "313")
	req.Header.Set("Scout-App", "1fbb8b74-2a24-4855-9a64-60fb4acfdd78")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	var d PlayerData
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &d)
	if err != nil {
	}

	if len(d.Data.Player.Segments) > 5 {
		duoAlltime := "duoAlltime: " + fmt.Sprintf("kills %5.0f, matchPlayed %5.0f, wins %5.0f, K/D %f\n", d.Data.Player.Segments[0].Stats[0].Value, d.Data.Player.Segments[0].Stats[2].Value, d.Data.Player.Segments[0].Stats[3].Value, d.Data.Player.Segments[0].Stats[8].Value)
		squadAlltime := "squAlltime: " + fmt.Sprintf("kills %5.0f, matchPlayed %5.0f, wins %5.0f, K/D %f\n", d.Data.Player.Segments[1].Stats[0].Value, d.Data.Player.Segments[1].Stats[2].Value, d.Data.Player.Segments[1].Stats[3].Value, d.Data.Player.Segments[1].Stats[8].Value)
		soloAlltime := "solAlltime: " + fmt.Sprintf("kills %5.0f, matchPlayed %5.0f, wins %5.0f, K/D %f\n", d.Data.Player.Segments[2].Stats[0].Value, d.Data.Player.Segments[2].Stats[2].Value, d.Data.Player.Segments[2].Stats[3].Value, d.Data.Player.Segments[2].Stats[8].Value)
		duoWeekly := "duoWeekly: " + fmt.Sprintf("kills %5.0f, matchPlayed %5.0f, wins %5.0f, K/D %f\n", d.Data.Player.Segments[3].Stats[0].Value, d.Data.Player.Segments[3].Stats[2].Value, d.Data.Player.Segments[3].Stats[3].Value, d.Data.Player.Segments[3].Stats[8].Value)
		squadWeekly := "squWeekly: " + fmt.Sprintf("kills %5.0f, matchPlayed %5.0f, wins %5.0f, K/D %f\n", d.Data.Player.Segments[4].Stats[0].Value, d.Data.Player.Segments[4].Stats[2].Value, d.Data.Player.Segments[4].Stats[3].Value, d.Data.Player.Segments[4].Stats[8].Value)
		soloWeekly := "solWeekly: " + fmt.Sprintf("kills %5.0f, matchPlayed %5.0f, wins %5.0f, K/D %f\n", d.Data.Player.Segments[5].Stats[0].Value, d.Data.Player.Segments[5].Stats[2].Value, d.Data.Player.Segments[5].Stats[3].Value, d.Data.Player.Segments[5].Stats[8].Value)
		s = "```\n" + soloAlltime + duoAlltime + squadAlltime + "\n" + soloWeekly + duoWeekly + squadWeekly + "```\n"
	} else {
		s = "```\n" + "Error" + "```\n"
	}
	return s
}
