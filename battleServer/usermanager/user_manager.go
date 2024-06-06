package usermanager

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Username string
	Pokemons []string
	Conn     net.Conn
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
