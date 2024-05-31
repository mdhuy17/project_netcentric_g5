package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	conn     net.Conn
	username string
	pokemons []string
}

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

	clients := make(map[net.Conn]*Client)

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, clients)
	}
}

func handleConnection(conn net.Conn, clients map[net.Conn]*Client) {
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
			clients[conn] = &Client{
				conn:     conn,
				username: username,
				pokemons: make([]string, 3), // Initialize the slice with a length of 3
			}
			fmt.Printf("New client connected: %s\n", username)
		case 3:
			// Received a Pokemon
			username := parts[0]
			pokemonNumber, err := strconv.Atoi(parts[1][len("Pokemon "):])
			if err != nil {
				fmt.Printf("Invalid Pokemon number received from %s: %s\n", username, parts[1])
				return
			}
			pokemonName := parts[2]

			client, ok := clients[conn]
			if !ok {
				fmt.Printf("Received message from unknown client: %s\n", data)
				return
			}

			if pokemonNumber >= 1 && pokemonNumber <= len(client.pokemons) {
				client.pokemons[pokemonNumber-1] = pokemonName
				fmt.Printf("%s added Pokemon %d: %s\n", username, pokemonNumber, pokemonName)

				// Check if the client has provided all 3 Pokemon
				if allPokemonsProvided(client) {
					// Broadcast the updated Pokemon list to all clients
					broadcastPokemons(clients)
				}
			} else {
				fmt.Printf("Invalid Pokemon number received from %s: %d\n", username, pokemonNumber)
			}
		default:
			fmt.Printf("Received unknown message: %s\n", data)
		}
	}
}

func broadcastPokemons(clients map[net.Conn]*Client) {
	for _, client := range clients {
		message := fmt.Sprintf("%s's Pokemon: %s\nOpponent's Pokemon: %s", client.username, strings.Join(client.pokemons, ", "), getOpponentPokemons(clients, client))
		_, err := client.conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", client.username, err)
		}
	}
}

func allPokemonsProvided(client *Client) bool {
	for _, pokemon := range client.pokemons {
		if pokemon == "" {
			return false
		}
	}
	return true
}

func getOpponentPokemons(clients map[net.Conn]*Client, currentClient *Client) string {
	var opponentPokemons []string
	for _, client := range clients {
		if client != currentClient {
			opponentPokemons = client.pokemons
			break
		}
	}
	return strings.Join(opponentPokemons, ", ")
}
