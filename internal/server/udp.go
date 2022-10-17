package server

import (
	"fmt"
	"net"
	"os"
	"spolks/internal/models"
)

type UDPServer struct {
	addr     string
	server   *net.UDPConn
	sessions map[net.Addr]models.Msg
}

func (u *UDPServer) Start() (err error) {
	//TODO implement me
	panic("implement me")
}

func (u *UDPServer) Shutdown() {
	fmt.Println("Shutting down server")
	err := u.server.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
