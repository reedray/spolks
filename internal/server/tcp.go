package server

import (
	"encoding/binary"
	"encoding/json"
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
	queue    chan int
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
			fmt.Println(err)
			break
		}

		select {
		case t.queue <- 1:
			fmt.Println("Started new thread")
			go t.handleConnection(conn)
		default:
			fmt.Println("Pool is full!", conn.RemoteAddr(), "can`t be added to the pool,disconnecting")
			conn.Close()
		}

	}
	return nil
}

func (t *TCPServer) handleConnection(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), " connected")
	defer fmt.Println(conn.RemoteAddr(), "disconnected")
	defer conn.Close()
	//

	defer func() {
		fmt.Println("Thread removed from pool")
		<-t.queue
	}()
	defer func() {
		<-t.queue
	}()
	if val, ok := t.sessions[conn.RemoteAddr().String()]; ok {
		switch val.Cmd {
		case constants.DOWNLOAD:
			fileInfo, _ := os.Stat(val.Filename)
			toSend := int(fileInfo.Size()) - val.BytesProcessed
			bytes, err := os.ReadFile(val.Filename)
			if err != nil {
				return
			}
			bytes = bytes[toSend:]
			conn.Write(bytes)
		case constants.UPLOAD:
			syn := models.SYN{
				Filename:       val.Filename,
				BytesProcessed: val.BytesProcessed,
				NeedRestore:    true,
			}
			bytes, _ := json.Marshal(syn)
			conn.Write(bytes)

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
			file, err := os.OpenFile("upload_"+val.Filename, os.O_APPEND|os.O_WRONLY, 666)
			defer file.Close()
			if err != nil {
				fmt.Println(err)
			}
			_, err = file.Write(data)
			if err != nil {
				fmt.Println(err)
			}
		}
		delete(t.sessions, conn.RemoteAddr().String())
		fmt.Println("Restored session for: ", conn.RemoteAddr())
	} else {
		bytes, err := json.Marshal(models.SYN{})
		if err != nil {
			return
		}
		conn.Write(bytes)
	}
	for {
		msg := models.Msg{}
		err := json.NewDecoder(conn).Decode(&msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		//fmt.Println(time.Since(msg.MsgTime))
		//bitrate := float64(len(msg.Data)) / float64(time.Since(msg.MsgTime).Milliseconds())
		//fmt.Printf("bitrate: %f Mbps", bitrate)
		err = t.sendResponse(conn, msg)
		if err != nil {
			fmt.Println(err)
			buf := utils.CreateBuffer([]byte(err.Error()))
			conn.Write(buf)
			break
		}
	}

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

func (t *TCPServer) sendResponse(conn net.Conn, msg models.Msg) error {
	switch msg.Cmd {
	case constants.ECHO:
		buf := utils.CreateBuffer(msg.Data)
		conn.Write(buf)
	case constants.TIME:
		buf := utils.CreateBuffer([]byte(time.Now().String()))
		conn.Write(buf)
	case constants.CLOSE:
		buf := utils.CreateBuffer([]byte("Connection closed"))
		conn.Write(buf)
		conn.Close()
	case constants.DOWNLOAD:
		bytes, err := os.ReadFile(msg.Filename)
		if err != nil {
			return err
		}
		i := simulateErr()
		if i < 4 {
			msg.Data = msg.Data[:len(msg.Data)/2]
			msg.BytesProcessed = len(msg.Data) / 2
			t.sessions[conn.RemoteAddr().String()] = msg
			buffer := utils.CreateBuffer(bytes)
			conn.Write(buffer)
			return errors.New("Network error")
		}
		buffer := utils.CreateBuffer(bytes)
		conn.Write(buffer)
	case constants.UPLOAD:
		bitrate := float64(len(msg.Data)) / float64(time.Since(msg.MsgTime).Milliseconds())
		fmt.Printf("bitrate: %f Bpms\n", bitrate)
		file, err := os.OpenFile("upload_"+msg.Filename, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		i := simulateErr()
		if i < 9 {
			file.Write(msg.Data[:len(msg.Data)/2])
			msg.BytesProcessed = len(msg.Data) / 2
			t.sessions[conn.RemoteAddr().String()] = msg
			return errors.New("Network error")
		}
		file.Write(msg.Data)
	}
	return nil
}
