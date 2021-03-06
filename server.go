package main

import "net"

type Server struct {
	conn     *net.UDPConn
	messages chan string
}

func (server *Server) handleMessage() []byte {
	var buf [512]byte

	n, _, err := server.conn.ReadFromUDP(buf[0:])
	if err != nil {
		return nil
	}

	return buf[0:n]
}
