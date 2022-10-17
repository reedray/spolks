package server

import (
	"net"
	"spolks/internal/models"
	"strings"
)

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
		}
	case "udp":
		return &UDPServer{
			addr:     addr,
			sessions: make(map[net.Addr]models.Msg),
		}
	}
	return nil
}
