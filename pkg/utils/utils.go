package utils

import (
	"encoding/binary"
	"spolks/pkg/constants"
)

func CreateTcpBuffer(data []byte) []byte {
	// Create a buffer with size enough to hold a prefix and actual data
	buf := make([]byte, constants.PrefixSize+len(data))

	// State the total number of bytes (including prefix) to be transferred over
	binary.BigEndian.PutUint32(buf[:constants.PrefixSize], uint32(constants.PrefixSize+len(data)))

	// Copy data into the remaining buffer
	copy(buf[constants.PrefixSize:], data[:])

	return buf
}
