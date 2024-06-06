package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	startTCPClient()
}

func startTCPClient() {
	// Connect to the TCP server
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to the TCP server.")

	// Get the user's username
	username := getUsernameFromInput()

	// Send the username to the server
	_, err = conn.Write([]byte(username))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}

	// Start a goroutine to read responses from the server
	go readResponsesFromServer(conn)
	// Read user input and send it to the server
	readAndSendPokemons(conn, username)

}

func getUsernameFromInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	return strings.TrimSpace(username)
}

func readAndSendPokemons(conn net.Conn, username string) {
	reader := bufio.NewReader(os.Stdin)
	for i := 1; i < 4; i++ {
		fmt.Printf("Enter Pokemon %d: ", i)
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == "exit" {
			break
		}

		// Append the username and Pokemon number to the input text
		text = fmt.Sprintf("%s: Pokemon %d: %s", username, i, strings.TrimSpace(text))

		// Send the message to the server
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
	// Wait for the server's response before exiting
	for {
		start := time.Now()
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		response := strings.TrimSpace(string(buf[:n]))
		fmt.Println(response)

		// Check if the server has sent the final message
		if strings.Contains(response, "Opponent's Pokemon:") {
			continue
		}
		if time.Since(start) >= 3*time.Second {
			fmt.Println("Server is not responding, exiting client.")
			return
		}
	}

	fmt.Println("Exiting client.")
}
func readResponsesFromServer(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		response := strings.TrimSpace(string(buf[:n]))
		fmt.Println(response)
	}
}
