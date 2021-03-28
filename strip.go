package main

import (
	"fmt"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Strip struct {
	ws2811.WS2811
	ledCount      uint32
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
	DNRGBLedValuesStartIndex   = uint32(4)
)

func (s *Strip) setup() error {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = s.ledBrightness
	opt.Channels[0].LedCount = int(s.ledCount)

	newStrip, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return fmt.Errorf("failed to initialize WS2811 strip: %w", err)
	}

	s.WS2811 = *newStrip
	return nil
}

func combineTwoBytes(high, low uint32) uint32 {
	return (high << 8) | low
}
func combineThreeBytes(high uint32, middle uint32, low uint32) uint32 {
	return (high << 16) | (middle << 8) | low
}

func (s *Strip) updateDNRGB(msg []byte) {
	msgLength := uint32(len(msg))
	ledIndexHigh := uint32(msg[DNRGBStartLedHighByteIndex])
	ledIndexLow := uint32(msg[DNRGBStartLedLowByteIndex])
	ledIndex := combineTwoBytes(ledIndexHigh, ledIndexLow)

	for i := DNRGBLedValuesStartIndex; i+3 <= msgLength; i += 3 {
		s.Leds(0)[ledIndex] = combineThreeBytes(uint32(msg[i]), uint32(msg[i+1]), uint32(msg[i+2]))
		ledIndex++
	}
	_ = s.Render()
}
