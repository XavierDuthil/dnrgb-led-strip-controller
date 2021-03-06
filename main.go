package main

import (
	"fmt"
	"log"
	"net"
)

const (
	// Server setting
	port string = ":7000"

	// LED strip settings
	ledCount   = 62
	brightness = 128
)

func main() {
	// Setup server
	udpAddress, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)

	var server Server
	server.conn, err = net.ListenUDP("udp", udpAddress)
	checkError(err)
	defer func() {
		fmt.Print("closing")
		_ = server.conn.Close()
	}()
	log.Printf("Listening via UDP on %s", udpAddress)

	// Setup LED strip
	strip := &Strip{}
	checkError(strip.setup())
	checkError(strip.Init())
	defer strip.Fini()

	// Handle requests
	for {
		msg := server.handleMessage()
		strip.update(msg)
		log.Printf("Strip updated with %q", msg)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
