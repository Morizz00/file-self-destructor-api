package storage

type StoredFile struct {
	FileName      string `json:"filename"`
	MIME          string `json:"mime"`
	Data          []byte `json:"date"`
	Password      string `json:"password"`
	DownloadsLeft int    `json:"downloads left"`
}
