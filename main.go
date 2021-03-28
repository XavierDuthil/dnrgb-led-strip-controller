package main

import (
	"log"
	"net"
)

const (
	// Server setting
	port = ":7000"

	// LED strip settings
	ledCount      = 62
	ledBrightness = 128
)

func main() {
	// Setup server
	udpAddress, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddress)
	checkError(err)
	defer func() {
		log.Println("Closing UDP listener")
		_ = conn.Close()
	}()
	log.Printf("Listening via UDP on %s", udpAddress)

	// Setup LED strip
	strip := &Strip{
		ledCount:      ledCount,
		ledBrightness: ledBrightness,
	}
	checkError(strip.setup())
	checkError(strip.Init())
	defer strip.Fini()

	// Maximum size of a message (in bytes) is prefixSize + nbLeds * nbColors
	//   Where prefixSize = 4 (protocol + ledIndex)
	//   And   nbColors   = 3 (RGB)
	var buf [4 + ledCount * 3]byte

	// Handle requests
	for {
		n, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Printf("Failed to read from UDP: %s", err)
			continue
		}
		strip.updateDNRGB(buf[0:n])
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
