package skim_monsters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Move represents a Pokémon move.
type ListMapObject struct {
	Name        string `json:"name"`
	ResourceURI string `json:"resource_uri"`
}

type Data struct {
	Description []ListMapObject `json:"descriptions"`
	Type        []ListMapObject `json:"types"`
	Abilities   []ListMapObject `json:"abilities"`
	/*
		"attack":49,"defense":49,"speed":45,"sp_atk":65,"sp_def":65,"hp":45,"weight":"69","height":"7","national_id":1,"name":"Bulbasaur","male_female_ratio":"87.5/12.5","abilities":[{"name":"chlorophyll","resource_uri":"/api/v1/ability/34/"
	*/
	Attack          int    `json:"attack"`
	Defense         int    `json:"defense"`
	Speed           int    `json:"speed"`
	SpAtk           int    `json:"sp_atk"`
	SpDef           int    `json:"sp_def"`
	HP              int    `json:"hp"`
	Weight          string `json:"weight"`
	Height          string `json:"height"`
	NationalID      int    `json:"national_id"`
	MaleFemaleRatio string `json:"male_female_ratio"`
	CatchRate       int    `json:"catch_rate"`
	ID              string `json:"_id"`
	Name            string `json:"name"`
}

// InputData represents the structure of the input text file.
type InputData struct {
	Docs []Data `json:"docs"`
	Seq  int    `json:"seq"`
}

func (i *InputData) crawl() {
	// URL to fetch the data from
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("https://pokedex.org/assets/skim-monsters-%d.txt", i)

		// Create a new HTTP request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %s", err)
		}

		// Set the headers
		req.Header.Set("Referer", "https://pokedex.org/js/worker.js")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %s", err)
		}
		defer resp.Body.Close()

		// Read the response body
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %s", err)
		}

		// Split the content into individual JSON strings
		parts := strings.Split(string(content), "\n")

		// Initialize a slice to hold all moves

		// Iterate over each part and unmarshal the JSON into InputData
		for _, part := range parts {
			if strings.TrimSpace(part) == "" {
				continue
			}

			var inputData InputData
			err := json.Unmarshal([]byte(part), &inputData)
			if err != nil {
				log.Printf("Failed to unmarshal part: %s\nError: %s", part, err)
				continue
			}

			// Save each move to a separate JSON file
			for _, move := range inputData.Docs {
				name, err := strconv.Atoi(move.ID)
				if err != nil {
					return
				}
				filename := fmt.Sprintf("./skim_monsters/data/%s.json", strconv.Itoa(name))
				moveJSON, err := json.MarshalIndent(move, "", "  ")
				if err != nil {
					log.Printf("Failed to marshal move to JSON: %s\nError: %s", move.ID, err)
					continue
				}

				err = ioutil.WriteFile(filename, moveJSON, 0644)
				if err != nil {
					log.Printf("Failed to write move to file: %s\nError: %s", filename, err)
					continue
				}
			}
		}

		fmt.Println("Moves have been saved to individual JSON files.")
	}

}
