package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

type server struct {
	CmdCh          chan commandID
	ErrCh          chan error
	FailedSessions map[*Msg]bool
}

func NewServer() *server {
	return &server{
		CmdCh:    make(chan commandID),
		ErrCh:    make(chan error),
		Sessions: make(map[*Msg]bool),
	}
}

//func (s *server) start() {
//	for cmd := range s.cmd {
//		switch cmd.id {
//		case ECHO:
//
//		case TIME:
//
//		case CLOSE:
//
//		case DOWNLOAD:
//
//		case UPLOAD:
//
//		}
//	}
//}

func (s *server) Start(network, addr string) (<-chan commandID, <-chan error) {
	listener, err := net.Listen(network, addr)
	if err != nil {
		fmt.Println("Failed to start server:", err.Error())
		os.Exit(-1)
	}
	defer listener.Close()
	fmt.Println("Server started on ", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			return nil, nil
		}
		go s.HandleRequests(conn)

	}
	return nil, nil
}
func (s *server) HandleRequests(conn net.Conn) {
	fmt.Println(conn.RemoteAddr().String(), " connected")
	//1 принять сообщение от клиента
	//2 проверить есть ли сессия для него в проваленных
	//3 если есть то выполнить заново msg
	//4

	m := Msg{}
	err := json.NewDecoder(conn).Decode(&m)
	if err != nil {
		fmt.Println("Failed to parse data", err.Error())
		s.ErrCh <- err
	}
	switch m.CmdID {
	case ECHO:
		s.echo(conn, m.Data)
	case TIME:
		s.time(conn)
	case CLOSE:
		s.close(conn)
	case DOWNLOAD:
		s.download(conn, m.FileName)
	case UPLOAD:
		s.upload(conn, m.FileName)
	}
}

func (s *server) Shutdown() {

}

func (s *server) echo(conn net.Conn, data []byte) {
	Msg{
		FileName: "",
		CmdID:    0,
		Data:     nil,
	}
}

func (s *server) time(conn net.Conn) {

}

func (s *server) close(conn net.Conn) {

}

func (s *server) download(conn net.Conn, name string) {

}

func (s *server) upload(conn net.Conn, name string) {

}

func (s *server) handleTCP(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			s.ErrCh <- errors.New("type assertion failed")
		}
		go s.HandleRequests(tcpConn)
	}

}

func (s *server) handleUDP(listener net.Listener) {

}
