package main

import (
	"fmt"
	"log"
	"strconv"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Strip struct {
	ws2811.WS2811
}

func (s *Strip) setup() error {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	newStrip, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return fmt.Errorf("failed to initialize WS2811 strip: %w", err)
	}

	s.WS2811 = *newStrip
	return nil
}

func (s *Strip) update(msg []byte) {
	msgStr := string(msg)
	msgLength := len(msgStr)
	for i := 0; i+6 <= msgLength; i += 6 {
		hexInt, err := strconv.ParseInt(msgStr[i:i+6], 16, 64)
		if err != nil {
			log.Printf("failed to parse hex value %q: %s", msgStr[i:i+6], err)
		}

		ledNum := i / 6
		s.Leds(0)[ledNum] = uint32(hexInt)
	}

	checkError(s.Render())
}
