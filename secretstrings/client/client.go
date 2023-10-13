package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"time"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func words() ([]string, error) {
	f, err := os.Open("../wordlist")
	if err != nil {
		fmt.Println("Could not open wordlist file.")
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	// Define a list of server addresses to distribute the load
	servers := []string{"127.0.0.1:8030", "127.0.0.1:8031", "127.0.0.1:8032"}

	rand.Seed(time.Now().UnixNano())

	words, err := words()
	if err != nil {
		fmt.Println("Could not read words:", err)
	}

	for _, word := range words {
		// Randomly select a server from the list.
		// There are probably better ways to load balance.
		serverAddr := servers[rand.Intn(len(servers))]

		client, err := rpc.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Printf("Failed to connect to server %s: %v\n", serverAddr, err)
			continue
		}
		defer client.Close()

		request := stubs.Request{Message: word}
		response := new(stubs.Response)
		err = client.Call(stubs.PremiumReverseHandler, request, response)
		if err != nil {
			fmt.Printf("Failed to reverse word '%s' on server %s: %v\n", word, serverAddr, err)
		} else {
			fmt.Printf("Reversed '%s' to '%s' on server %s\n", word, response.Message, serverAddr)
		}
	}
}
