package service

type commandID int

const (
	ECHO commandID = iota
	TIME
	CLOSE
	UPLOAD
	DOWNLOAD
)

func (cid commandID) String() string {
	return [...]string{"ECHO", "TIME", "CLOSE", "UPLOAD", "DOWNLOAD"}[cid]
}
