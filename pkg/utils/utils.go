package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"spolks/pkg/constants"
	"syscall"
)

func CreateBuffer(data []byte) []byte {
	// Create a buffer with size enough to hold a prefix and actual data
	buf := make([]byte, constants.PrefixSize+len(data))

	// State the total number of bytes (including prefix) to be transferred over
	binary.BigEndian.PutUint32(buf[:constants.PrefixSize], uint32(constants.PrefixSize+len(data)))

	// Copy data into the remaining buffer
	copy(buf[constants.PrefixSize:], data[:])

	return buf
}

func ReadData(conn net.Conn) ([]byte, error) {
	prefix := make([]byte, constants.PrefixSize)
	_, err := io.ReadFull(conn, prefix)
	if err != nil {
		return nil, err
	}

	totalDataLength := binary.BigEndian.Uint32(prefix[:])
	data := make([]byte, totalDataLength-constants.PrefixSize)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		if errors.Is(err, syscall.ECONNRESET) {
			fmt.Println("Connection closed")

		}
		return nil, err
	}
	return data, nil
}
