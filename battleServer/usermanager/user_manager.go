package usermanager

import "net"

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

func (um *UserManager) AllPokemonsProvided(username string) bool {
	user, exists := um.Users[username]
	if exists {
		for _, pokemon := range user.Pokemons {
			if pokemon == "" {
				return false
			}
		}
		return true
	}
	return false
}
