package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"spolks/internal/models"
	"spolks/pkg/constants"
	"spolks/pkg/utils"
	"syscall"
	"time"
)

type TCPServer struct {
	addr     string
	server   net.Listener
	sessions map[string]models.Msg
}

func (t *TCPServer) Start() (err error) {
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer t.server.Close()
	fmt.Println("Server started on:", t.server.Addr())
	for {

		conn, err := t.server.Accept()
		if err != nil {
			//errors.Is(err)
			fmt.Println(err)
			break
		}
		go t.handleConnection(conn)
	}
	return nil
}

func (t *TCPServer) handleConnection(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), " connected")
	defer fmt.Println(conn.RemoteAddr(), "disconnected")
	defer conn.Close()
	for {
		prefix := make([]byte, constants.PrefixSize)
		_, err := io.ReadFull(conn, prefix)
		if err != nil {
			break
		}

		totalDataLength := binary.BigEndian.Uint32(prefix[:])
		data := make([]byte, totalDataLength-constants.PrefixSize)
		_, err = io.ReadFull(conn, data)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) {
				fmt.Println("Connection closed")
				break
			}
		}
		msg := utils.CreateMsg(data)

		if val, ok := t.sessions[conn.RemoteAddr().String()]; ok {
			if val.Cmd == msg.Cmd && *val.Filename == *msg.Filename {
				err := t.createResponse(conn, val)
				if err != nil {
					fmt.Println(err)
					conn.Write([]byte(err.Error()))
					continue
				}
				fmt.Println("RESTORED SESSION FOR:", conn.RemoteAddr())
				delete(t.sessions, conn.RemoteAddr().String())
			}
		}

		err = t.createResponse(conn, msg)
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte(err.Error()))
			conn.Close()
		}
	}

}

func (t *TCPServer) createResponse(conn net.Conn, msg models.Msg) error {
	switch msg.Cmd {
	case constants.ECHO:
		conn.Write(*msg.Data)
	case constants.TIME:
		now := time.Now().String()
		conn.Write([]byte(now))
	case constants.CLOSE:
		conn.Write([]byte("Connection closed"))
		conn.Close()
	case constants.DOWNLOAD:
		i := simulateErr()
		if i < 3 {
			t.sessions[conn.RemoteAddr().String()] = msg
			return errors.New("Network error")
		}
		bytes, err := os.ReadFile(*msg.Filename)
		if err != nil {
			conn.Write([]byte("No such file on a server"))
			return err
		}
		bytes = utils.CreateTcpBuffer(bytes)
		conn.Write(bytes)
	case constants.UPLOAD:
		i := simulateErr()
		if i < 6 {
			t.sessions[conn.RemoteAddr().String()] = msg
			return errors.New("Network error")
		}
		file, _ := os.Create(*msg.Filename)
		file.Write(*msg.Data)
		conn.Write([]byte("File created"))
	}
	return nil
}

func simulateErr() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 10
	return rand.Intn(max-min+1) + min
}

func (t *TCPServer) Shutdown() {
	fmt.Println("Shutting down server")
	err := t.server.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	os.Exit(0)

}
