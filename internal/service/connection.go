package service

import (
	"encoding/binary"
	"io"
	"net"
	"spolks/pkg/constants"
)

type connection struct {
	conn net.Conn
	id   string
}

func (c *connection) Read() ([]byte, error) {
	prefix := make([]byte, constants.PrefixSize)

	// Read the prefix, which contains the length of data expected
	_, err := io.ReadFull(c.conn, prefix)
	if err != nil {
		return nil, err
	}

	totalDataLength := binary.BigEndian.Uint32(prefix[:])

	// Buffer to store the actual data
	data := make([]byte, totalDataLength-constants.PrefixSize)

	// Read actual data without prefix
	_, err = io.ReadFull(c.conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
