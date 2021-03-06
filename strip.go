package main

import (
	"fmt"
	"log"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Strip struct {
	ws2811.WS2811
	ledCount      int
	ledBrightness int
}

const (
	// Supported protocols
	// c.f. https://github.com/Aircoookie/WLED/wiki/UDP-Realtime-Control
	DNRGBProtocol = "DNRGB"

	// Common indexes for all protocols
	protocolByteIndex = 0
	timeoutByteIndex  = 1

	// Indexes specific to DNRGB protocol
	DNRGBStartLedHighByteIndex = 2
	DNRGBStartLedLowByteIndex  = 3
	DNRGBLedValuesStartIndex   = 4
)

func (s *Strip) setup() error {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = s.ledBrightness
	opt.Channels[0].LedCount = s.ledCount

	newStrip, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return fmt.Errorf("failed to initialize WS2811 strip: %w", err)
	}

	s.WS2811 = *newStrip
	return nil
}

func determineProtocol(b byte) string {
	switch b {
	case '\x04':
		return DNRGBProtocol
	default:
		log.Printf("Unsupported protocol %q", b)
		return ""
	}
}

func combineTwoBytes(high, low uint16) uint16 {
	return (high << 8) | low
}
func combineThreeBytes(high uint32, middle uint32, low uint32) uint32 {
	return (high << 16) | (middle << 8) | low
}

func (s *Strip) update(msg []byte) {
	protocolByte := msg[protocolByteIndex]
	// TODO: implement flush after timeout
	//timeoutSecondsByte := msg[timeoutByteIndex]

	protocol := determineProtocol(protocolByte)

	switch protocol {
	case DNRGBProtocol:
		s.updateDNRGB(msg)
	default:
		log.Printf("Unsupported protocol %q", protocol)
		return
	}

	checkError(s.Render())
}

func (s *Strip) updateDNRGB(msg []byte) {
	ledIndexHigh := msg[DNRGBStartLedHighByteIndex]
	ledIndexLow := msg[DNRGBStartLedLowByteIndex]
	ledStartIndex := combineTwoBytes(uint16(ledIndexHigh), uint16(ledIndexLow))

	for i := uint16(DNRGBLedValuesStartIndex); i+3 <= uint16(len(msg)); i += 3 {
		ledIndex := ledStartIndex + (i-DNRGBLedValuesStartIndex)/3
		if ledIndex >= uint16(s.ledCount) {
			log.Printf("Tried to assign LED #%d which is out of bounds", ledIndex)
			break
		}

		s.Leds(0)[ledIndex] = combineThreeBytes(uint32(msg[i]), uint32(msg[i+1]), uint32(msg[i+2]))
	}
}
