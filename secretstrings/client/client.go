package main

import (
	"bufio"
	//	"net/rpc"
	"flag"
	"net/rpc"
	"os"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"

	//	"bufio"
	//	"os"
	//	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
	"fmt"
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
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	words, err := words()
	if err != nil {
		fmt.Println("Could not read words:", err)
	}

	for _, word := range words {
		request := stubs.Request{Message: word}
		response := new(stubs.Response)
		err = client.Call(stubs.PremiumReverseHandler, request, response)
		if err != nil {
			fmt.Printf("Failed to reverse word '%s': %v\n", word, err)
		} else {
			fmt.Printf("Reversed '%s' to '%s'\n", word, response.Message)
		}
	}
}
