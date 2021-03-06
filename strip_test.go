package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrip_updateDNRGB(t *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name     string
		args     args
		wantLeds []uint32
	}{
		{
			name: "DNRGB: index 0",
			args: args{
				msg: []byte("\x04\xff\x00\x00\x11\x11\x11\x22\x22\x22\x33\x33\x33"),
			},
			wantLeds: []uint32{
				0x111111,
				0x222222,
				0x333333,
				0x000000,
			},
		}, {
			name: "DNRGB: index 1",
			args: args{
				msg: []byte("\x04\xff\x00\x01\x11\x11\x11\x22\x22\x22\x33\x33\x33"),
			},
			wantLeds: []uint32{
				0x000000,
				0x111111,
				0x222222,
				0x333333,
			},
		}, {
			name: "Unsupported protocol",
			args: args{
				msg: []byte("\x05\xff\x00\x00\x11\x11\x11\x22\x22\x22\x33\x33\x33"),
			},
			wantLeds: []uint32{
				0x000000,
				0x000000,
				0x000000,
				0x000000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Strip{
				ledCount:      4,
				ledBrightness: 255,
			}
			checkError(s.setup())
			checkError(s.Init())

			s.update(tt.args.msg)
			assert.Equal(t, tt.wantLeds, s.Leds(0))
		})
	}
}
