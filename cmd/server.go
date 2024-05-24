package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mdhuy17/project_netcentric_g5/internal/handlers"
	"github.com/mdhuy17/project_netcentric_g5/utils"
	"log"
	"net"
	"strings"
)

// Server represents the file server
type Server struct{}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	fileName, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Error reading file name:", err)
		return
	}
	fileName = strings.TrimSpace(fileName)

	// Read the file content
	fileContent, err := s.GetPokemon(fileName)
	if err != nil {
		log.Println("Error reading file:", err)
		conn.Write([]byte("Error reading file\n"))
		return
	}

	// Send the file content to the client
	_, err = conn.Write(fileContent)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}

// ReadFile reads the content of the specified file
func (s *Server) GetPokemon(name string) ([]byte, error) {
	var pokedex handlers.Pokedex
	data, err := pokedex.GetMonsterByID(utils.PokeMap[name])
	if err != nil {
		return []byte{}, err
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return []byte{}, err
	}
	return jsonData, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on port 8080...")

	server := &Server{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go server.HandleConnection(conn)
	}
}
