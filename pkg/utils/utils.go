package utils

import (
	"encoding/binary"
	"spolks/internal/models"
	"spolks/pkg/constants"
	"strings"
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

func CreateMsg(b []byte) models.Msg {
	s := string(b)

	msg := models.Msg{}

	if strings.HasPrefix(s, constants.ECHO) {
		msg.Cmd = constants.ECHO
		s = s[len(constants.ECHO)+1:]
		buf := make([]byte, len(s))
		buf = []byte(s)
		msg.Data = &buf
	}
	if strings.HasPrefix(s, constants.TIME) {
		msg.Cmd = constants.TIME
	}
	if strings.HasPrefix(s, constants.CLOSE) {
		msg.Cmd = constants.CLOSE
	}
	if strings.HasPrefix(s, constants.DOWNLOAD) {
		msg.Cmd = constants.DOWNLOAD
		args := strings.Fields(s)
		msg.Filename = &args[1]
	}
	if strings.HasPrefix(s, constants.UPLOAD) {
		msg.Cmd = constants.UPLOAD
		args := strings.Fields(s)
		msg.Filename = &args[1]
		s = s[len(constants.UPLOAD)+len(args[1])+2:] //+2 for spaces for cmd and filename
		bytes := make([]byte, len(s))
		bytes = []byte(s)
		msg.Data = &bytes
	}
	return msg

}
