package main

import (
	"fmt"
	"github.com/mdhuy17/project_netcentric_g5/battleServer/usermanager"
	"net"
	"strconv"
	"strings"
)

func main() {
	startTCPServer()
}

func startTCPServer() {
	// Listen on TCP port
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Failed to listen on port 8000:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on :8000")

	userManager := usermanager.NewUserManager()

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, userManager)
	}
}

func handleConnection(conn net.Conn, userManager *usermanager.UserManager) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		// Read data from the connection
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		// Parse the received data
		data := strings.TrimSpace(string(buf[:n]))
		parts := strings.Split(data, ": ")

		switch len(parts) {
		case 1:
			// New client connection
			username := parts[0]
			user := userManager.AddUser(username, conn)
			fmt.Printf("New client connected: %s\n", user.Username)
		case 3:
			// Received Pokemon information
			username := parts[0]
			pokemonNumber, _ := strconv.Atoi(parts[1][len("Pokemon "):])
			pokemonName := parts[2]

			err := userManager.UpdatePokemonData(username, pokemonName, pokemonNumber)
			if err != nil {
				fmt.Printf("Error updating Pokemon data for %s: %v\n", username, err)
				continue
			}

			userManager.UpdatePokemons(username, pokemonName, pokemonNumber)
			fmt.Printf("%s added Pokemon %d: %s\n", username, pokemonNumber, pokemonName)

			if userManager.AllPokemonsProvided() && len(userManager.Users) == 2 {
				userManager.StartBattle()
			}

		default:
			fmt.Printf("Received unknown message: %s\n", data)
		}
	}
}
