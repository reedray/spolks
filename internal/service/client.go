package service

import "net"

type client struct {
	conn     net.Conn
	commands chan<- commandID
}
