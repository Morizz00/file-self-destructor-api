package storage

import "time"

type StoredFile struct {
	FileName      string        `json:"filename"`
	MIME          string        `json:"mime"`
	Data          []byte        `json:"data"`
	Password      string        `json:"password"`
	DownloadsLeft int           `json:"downloadleft"`
	Expiry        time.Duration `json:"expiry"`
}
