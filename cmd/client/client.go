package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"spolks/internal/models"
	"spolks/pkg/constants"
	"spolks/pkg/utils"
	"strings"
	"time"
)

func main() {

	serverAddr, err := net.ResolveTCPAddr("tcp", ":8080")

	clientAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		fmt.Println(err)
	}
	//clientAddr := net.TCPAddr{
	//	IP:   clientAddr,
	//	Port: 1234,
	//	Zone: "",
	//}
	conn, err := net.DialTCP("tcp", clientAddr, serverAddr)
	if err != nil {
		fmt.Println("Can not connect to the server: ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("connected to ", conn.RemoteAddr())
	defer conn.Close()

	syn := models.SYN{}
	json.NewDecoder(conn).Decode(&syn)
	if syn.NeedRestore == true {
		bytes, err2 := os.ReadFile(syn.Filename)
		if err2 != nil {
			fmt.Println(err)
			return
		}
		bytes = utils.CreateBuffer(bytes[len(bytes)-syn.BytesProcessed:])
		conn.Write(bytes)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
		}
		err = sendToServer(conn, str)
		if err != nil {
			fmt.Println("Server error", err)
			break
		}
	}

}

func sendToServer(conn net.Conn, str string) error {
	str = strings.Trim(str, "\r\n")
	args := strings.Fields(str)
	msg := models.Msg{}

	switch args[0] {
	case constants.ECHO:
		msg.MsgTime = time.Now()
		msg.Cmd = constants.ECHO
		str = str[len(constants.ECHO)+1:]
		msg.Data = []byte(str)
		bytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}
		_, err = conn.Write(bytes)
		if err != nil {
			return err
		}

		data, err := utils.ReadData(conn)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case constants.TIME:
		msg.MsgTime = time.Now()
		msg.Cmd = constants.TIME
		bytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}
		_, err = conn.Write(bytes)
		if err != nil {
			return err
		}
		data, err := utils.ReadData(conn)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case constants.CLOSE:
		msg.MsgTime = time.Now()
		msg.Cmd = constants.CLOSE
		bytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}
		_, err = conn.Write(bytes)
		if err != nil {
			return err
		}
		data, err := utils.ReadData(conn)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return errors.New("Closed")
	case constants.DOWNLOAD:
		msg.MsgTime = time.Now()
		msg.Cmd = constants.DOWNLOAD
		msg.Filename = args[1]
		bytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}
		_, err = conn.Write(bytes)
		if err != nil {
			return err
		}

		data, err := utils.ReadData(conn)
		if err != nil {
			return err
		}
		bitrate := float64(len(data)) / float64(time.Since(msg.MsgTime).Milliseconds())
		fmt.Printf("bitrate: %f bpms\n", bitrate)
		file, err := os.Create("download_" + args[1])
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(data)
	case constants.UPLOAD:
		msg.MsgTime = time.Now()
		msg.Cmd = constants.UPLOAD
		msg.Filename = args[1]
		readFile, err := os.ReadFile(args[1])
		if err != nil {
			return err
		}
		msg.Data = readFile
		bytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("Unknown command,try again")
	}
	return nil
}
