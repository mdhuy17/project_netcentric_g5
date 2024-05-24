package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdhuy17/project_netcentric_g5/internal/handlers"
	"github.com/mdhuy17/project_netcentric_g5/utils"
)

func main() {

	var pokedex handlers.Pokedex

	data, err := pokedex.GetMonsterByID(utils.PokeMap["Ivysaur"])
	if err != nil {
		return
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print JSON data
	fmt.Println(string(jsonData))

}
