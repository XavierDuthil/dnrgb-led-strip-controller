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
	ledCounts  = 62
	brightness = 128
)

func main() {
	udpAddress, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)

	// Setup server
	var server Server
	server.conn, err = net.ListenUDP("udp", udpAddress)
	checkError(err)
	defer func() {
		fmt.Print("closing")
		_ = server.conn.Close()
	}()

	// Setup LED strip
	strip := &Strip{}
	err = strip.setup()
	checkError(err)

	checkError(strip.Init())
	defer strip.Fini()

	for {
		strip.update(server.handleMessage())
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
