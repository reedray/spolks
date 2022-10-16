package service

type Msg struct {
	FileName string    `json:"file_name"`
	CmdID    commandID `json:"cmd_id"`
	Data     []byte    `json:"data,omitempty"`
}

func NewMsg() *Msg {
	return &Msg{}

}
