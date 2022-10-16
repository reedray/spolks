package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"spolks/pkg/utils"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
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
		sendToServer(conn, str)
		getFromServer(conn)

	}

}

func getFromServer(conn net.Conn) {

}

func sendToServer(conn net.Conn, str string) {
	str = strings.Trim(str, "\r\n")
	args := strings.Fields(str)

	switch args[0] {
	case "ECHO":
		_, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Println(err)
		}
	case "TIME":
		_, err := conn.Write([]byte(args[0]))
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
		file, err := os.Open(args[1])
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		all, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		buffer := utils.CreateTcpBuffer(all)
		_, err = conn.Write([]byte(str))
		if err != nil {
			fmt.Println(err)
		}
		write, err := conn.Write(buffer)
		if err != nil {
			fmt.Println(err)
		}
		if write != len(buffer) {
			fmt.Println("Failed to upload a file")
		}
	default:
		fmt.Println("Unknown command,try again")
	}
}
