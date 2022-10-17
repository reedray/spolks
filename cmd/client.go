package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"spolks/pkg/utils"
	"strings"
)

func main() {

	serverAddr := net.TCPAddr{
		IP:   nil,
		Port: 8080,
		Zone: "",
	}

	clientAddr := net.TCPAddr{
		IP:   nil,
		Port: 1234,
		Zone: "",
	}
	conn, err := net.DialTCP("tcp", &clientAddr, &serverAddr)
	if err != nil {
		fmt.Println("Can not connect to the server: ", err.Error())
		os.Exit(-1)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
		}
		err = sendToServer(conn, str)
		if err != nil {
			continue
		}
		getFromServer(conn)

	}

}

func getFromServer(conn net.Conn) {

	bytes := make([]byte, 1024)
	conn.Read(bytes)
	fmt.Println(string(bytes))
}

func sendToServer(conn net.Conn, str string) error {
	str = strings.Trim(str, "\r\n")
	args := strings.Fields(str)

	switch args[0] {
	case "ECHO":
		buffer := utils.CreateTcpBuffer([]byte(str))
		_, err := conn.Write(buffer)
		if err != nil {
			fmt.Println(err)
		}
	case "TIME":
		buffer := utils.CreateTcpBuffer([]byte(args[0]))
		_, err := conn.Write(buffer)
		if err != nil {
			fmt.Println(err)
		}
	case "CLOSE":
		_, err := conn.Write([]byte(args[0]))
		if err != nil {
			fmt.Println(err)
		}
	case "DOWNLOAD":
		_, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Println(err)
		}
	case "UPLOAD":
		all, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Println(err)
		}
		buffer := []byte(args[0] + " " + args[1] + " ")
		buffer = append(buffer, all...)
		buffer = utils.CreateTcpBuffer(buffer)
		write, err := conn.Write(buffer)
		if err != nil {
			fmt.Println(err)
		}
		if write != len(buffer) {
			fmt.Println("Failed to upload a file")
		}
	default:
		fmt.Println("Unknown command,try again")
		return error(nil)
	}
	return nil
}
