package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"time"
)

const SingSongHandler = "BottlesOfBeerOperations.SingSong"

var nextAddr string
var bottles int
var conn *rpc.Client

type BottlesOfBeerOperations struct{}
type Response struct{}

type Request struct {
	BuddyId int
	Bottles int
}

func printVerse(buddyId, bottles int) {
	fmt.Println("Buddy", buddyId, "sings:", bottles, "bottles of beer on the wall.", bottles, "bottles of beer. Take one down, pass it around,", bottles-1, "bottles of beer on the wall")
	fmt.Println()
}

func printEnding(buddyId int) {
	fmt.Println("Buddy", buddyId, "sings: No more bottles of beer on the wall, no more bottles of beer.\nWe've taken them down and passed them around; now we're drunk and passed out!")
}

func (b *BottlesOfBeerOperations) SingSong(req Request, _ *Response) (err error) {
	var buddyId int
	if bottles > 0 {
		buddyId = 1
	} else {
		buddyId = req.BuddyId + 1
	}
	if req.Bottles > 0 {
		printVerse(buddyId, req.Bottles)
		conn.Go(SingSongHandler, Request{BuddyId: buddyId, Bottles: req.Bottles - 1}, &Response{}, nil)
	} else {
		printEnding(buddyId)
	}
	return
}

func handleError(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func main() {
	port := flag.String("port", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for the next member of the round")
	flag.IntVar(&bottles, "bottles", 0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()

	// Register RPC service
	err := rpc.Register(&BottlesOfBeerOperations{})
	handleError(err, "Could not register RPC service.")

	// Create a listener to accept incoming connections
	listener, err := net.Listen("tcp", ":"+*port)
	handleError(err, "Could not create a listener.")
	defer listener.Close()

	// We sleep to allow time for all clients to connect.
	time.Sleep(time.Second * 5)

	// Attempt to connect to the next buddy.
	conn, err = rpc.Dial("tcp", nextAddr)
	handleError(err, "Could not dial next buddy.")

	// If the buddy specifies the number of bottles, start the song.
	if bottles > 0 {
		defer conn.Close()
		buddyId := 1
		printVerse(buddyId, bottles)
		conn.Go(SingSongHandler, Request{BuddyId: buddyId, Bottles: bottles - 1}, &Response{}, nil)
	}

	for {
		conn, err := listener.Accept()
		handleError(err, "Could not start the server.")
		go rpc.ServeConn(conn)
	}
}
