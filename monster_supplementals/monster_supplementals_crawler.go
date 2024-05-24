package main

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
type Move struct {
	/*
		{"specialAttackEV":0,"hpEV":0,"hatchSteps":5100,"defenseEV":0,"attackEV":2,"specialDefenseEV":0,"speedEV":0,"genderRatio":25,"species":"Fairy Pokémon","japaneseName":"グランブル","hepburnName":"Guranburu","eggGroups":"Field, Fairy","_id":"00210","_rev":"1-A"}

	*/
	SpecialAttackEV  int         `json:"specialAttackEV"`
	HpEV             int         `json:"hpEV"`
	HatchSteps       int         `json:"hatchSteps"`
	DefenseEV        int         `json:"defenseEV"`
	AttackEV         int         `json:"attackEV"`
	SpecialDefenseEV int         `json:"specialDefenseEV"`
	SpeedEV          int         `json:"speedEV"`
	GenderRatio      interface{} `json:"genderRatio"`
	Species          string      `json:"species"`
	JapaneseName     string      `json:"japaneseName"`
	HepburnName      string      `json:"Guranburu"`
	EggGroups        string      `json:"eggGroups"`

	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}

// InputData represents the structure of the input text file.
type InputData struct {
	Docs []Move `json:"docs"`
	Seq  int    `json:"seq"`
}

func main() {
	// URL to fetch the data from
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("https://pokedex.org/assets/monsters-supplemental-%d.txt", i)

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
				filename := fmt.Sprintf("./monster_supplementals/data/%s.json", strconv.Itoa(name))
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
