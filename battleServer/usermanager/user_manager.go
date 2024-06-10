package usermanager

import (
	"encoding/json"
	"fmt"
	"github.com/mdhuy17/project_netcentric_g5/internal/models"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type User struct {
	Username    string
	Pokemons    []string
	Conn        net.Conn
	PokemonData []*PokemonData
}

type PokemonData struct {
	Monster             *models.Monster             `json:"monster"`
	Description         []*models.Descriptions      `json:"description"`
	Evolution           *models.Evolution           `json:"evolution"`
	Types               []*models.Types             `json:"types"`
	MonsterSupplemental *models.MonsterSupplemental `json:"monster_supplemental"`
	MonsterMoves        []*models.Move              `json:"monster_moves"`
}

type UserManager struct {
	Users map[string]*User
}

func NewUserManager() *UserManager {
	return &UserManager{
		Users: make(map[string]*User),
	}
}

func (um *UserManager) AddUser(username string, conn net.Conn) *User {
	if _, exists := um.Users[username]; !exists {
		user := &User{
			Username: username,
			Pokemons: make([]string, 3),
			Conn:     conn,
		}
		um.Users[username] = user
		return user
	}
	return um.Users[username]
}

func (um *UserManager) UpdatePokemons(username, pokemon string, index int) {
	user, exists := um.Users[username]
	if exists {
		user.Pokemons[index-1] = pokemon
	}
}

func (um *UserManager) GetOpponentPokemons(username string) []string {
	for _, user := range um.Users {
		if user.Username != username {
			return user.Pokemons
		}
	}
	return nil
}

func (um *UserManager) AllPokemonsProvided() bool {
	// Check if there are exactly 2 connected players
	if len(um.Users) != 2 {
		return false
	}

	// Check if each player has provided 3 Pokemon
	for _, user := range um.Users {
		if len(user.Pokemons) != 3 {
			return false
		}
		for _, pokemon := range user.Pokemons {
			if pokemon == "" {
				return false
			}
		}
	}

	return true
}

func (um *UserManager) StartBattle() {
	// Broadcast the Pokemon information to both players
	um.broadcastPokemons()

	// Implement the battle logic here
	um.performBattle()
}

func (um *UserManager) broadcastPokemons() {
	for _, user := range um.Users {
		message := fmt.Sprintf("%s's Pokemon: %s\nOpponent's Pokemon: %s", user.Username, strings.Join(user.Pokemons, ", "), strings.Join(um.getOpponentPokemons(user), ", "))
		_, err := user.Conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", user.Username, err)
		}
		message = fmt.Sprintf("Battle Commence")
		_, err = user.Conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", user.Username, err)
		}
	}
}

func (um *UserManager) getOpponentPokemons(currentUser *User) []string {
	for _, user := range um.Users {
		if user != currentUser {
			return user.Pokemons
		}
	}
	return nil
}

func (um *UserManager) performBattle() {
	// Implement the battle logic here
	// This is where you would handle the battle between the two players
	// and determine the winner
}

func (um *UserManager) UpdatePokemonData(username, pokemonName string, pokemonIndex int) error {
	user, exists := um.Users[username]
	if !exists {
		return fmt.Errorf("user %s not found", username)
	}
	pokemonID := getPokemonIDFromName(pokemonName)

	// Construct the path to the Pokemon data file relative to the current file
	pokemonDataFilePath := filepath.Join("..", "internal", "models", "monsters", "data", fmt.Sprintf("%d.json", pokemonID))
	pokemonData, err := readPokemonJSONData(pokemonDataFilePath)
	if err != nil {
		return fmt.Errorf("error reading Pokemon data: %v", err)
	}

	// Update the user's Pokemon information
	if len(user.PokemonData) < pokemonIndex {
		user.PokemonData = append(user.PokemonData, pokemonData)
	} else {
		user.PokemonData[pokemonIndex-1] = pokemonData
	}

	return nil
}
func getPokemonIDFromName(pokemonName string) int {
	// Construct the path to the pokemonNames.json file relative to the main.go file
	jsonFilePath := filepath.Join("..", "internal", "models", "pokemonNames.json")

	// Read the pokemonNames.json file
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Printf("Error reading pokemonNames.json file: %v\n", err)
		return 0
	}

	// Unmarshal the JSON data into a slice of strings
	var pokemonNames []string
	err = json.Unmarshal(data, &pokemonNames)
	if err != nil {
		fmt.Printf("Error unmarshaling pokemonNames.json data: %v\n", err)
		return 0
	}

	// Search for the Pokemon name in the slice and return the index (which is the ID)
	for i, name := range pokemonNames {
		if name == pokemonName {
			return i + 1 // The IDs start from 1, not 0
		}
	}

	fmt.Printf("Pokemon name '%s' not found in pokemonNames.json\n", pokemonName)
	return 0
}
func readPokemonJSONData(filePath string) (*PokemonData, error) {
	// Read the JSON data from the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var pokemonData PokemonData
	err = json.Unmarshal(data, &pokemonData)
	if err != nil {
		return nil, err
	}

	return &pokemonData, nil
}
