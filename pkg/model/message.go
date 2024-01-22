package model

type Message struct {
	FileId   int64
	FileName string
}

type FileData struct {
	FileName string `json:"fileName"`
	Quality  int    `json:"quality"`
}

type FileManagementMessage struct {
	FileId int64      `json:"fileId"`
	Files  []FileData `json:"files"`
}
