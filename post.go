package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func main() {
	url := "https://api.scoutsdk.com/graph"
	//AQUACAHQNjjhikHNSbppyhLI43KQ
	//var jsonStr = []byte(`{"query":"query id($title: String, $platform: String, $identifier: String) { players(title: $title, platform: $platform, identifier: $identifier) {results {player { playerId handle } } } }","variables":{ "title": "fortnite", "identifier": "` + username + `", "platform": "epic"}}`)
	//var jsonStr = []byte(`{"query":"query thing($title: String, $name: String, $parameters: [Dynamic]) { thing(title: $title, name: $name, parameters: $parameters) }","variables":{ "title": "fortnite", "parameters": ["daily"], "name": "store"}}`)
	var jsonStr = []byte(`{"query":"query player($title: String, $id: String, $segment: String) { player(title: $title, id: $id, segment: $segment)  { id metadata { key name value displayValue } stats { metadata { key name isReversed } value displayValue } segments { metadata { key name value displayValue } stats { metadata { key name isReversed } value displayValue } } }}","variables":{ "title": "fortnite", "segment": "*", "id": "AQUACAHQNjjhikHNSbppyhLI43KQ"}}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Host", "api.scoutsdk.com")
	req.Header.Set("Accept", "application/com.scoutsdk.graph+json; version=1.1.0; charset=utf8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "313")
	req.Header.Set("Scout-App", "1fbb8b74-2a24-4855-9a64-60fb4acfdd78")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var d PlayerData
	err = json.Unmarshal(body, &d)
	if err != nil {
		panic(err)
	}
	squadAlltime := fmt.Sprintf("kills %.0f\nmatchPlayed %f\nwins %f\nK/D %f\n", d.Data.Player.Segments[1].Stats[0].Value, d.Data.Player.Segments[1].Stats[2].Value, d.Data.Player.Segments[1].Stats[3].Value, d.Data.Player.Segments[1].Stats[8].Value)
	fmt.Println(squadAlltime)
	fmt.Println(d)
}
