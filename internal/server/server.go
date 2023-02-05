package server

import (
	"net"
	"spolks/internal/models"
	"strings"
)

const maxCons = 2

type Server interface {
	Start() error
	Shutdown()
}

func NewServer(protocol, addr string) Server {
	switch strings.ToLower(protocol) {
	case "tcp":
		return &TCPServer{
			addr:     addr,
			sessions: make(map[string]models.Msg),
			queue:    make(chan int, maxCons),
		}
	case "udp":
		return &UDPServer{
			addr:     addr,
			sessions: make(map[net.Addr]models.Msg),
		}
	}
	return nil
}
