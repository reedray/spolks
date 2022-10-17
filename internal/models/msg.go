package models

type Msg struct {
	Cmd      string
	Filename *string
	Data     *[]byte
}
