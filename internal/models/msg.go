package models

import "time"

type Msg struct {
	Cmd            string    `json:"cmd"`
	Filename       string    `json:"filename,omitempty"`
	Data           []byte    `json:"data,omitempty"`
	MsgTime        time.Time `json:"msg_time"`
	BytesProcessed int       `json:"bytes_read,omitempty"`
}
type SYN struct {
	Filename       string `json:"filename"`
	BytesProcessed int    `json:"bytes_processed"`
	NeedRestore    bool   `json:"need_restore"`
}
