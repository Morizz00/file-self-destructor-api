package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/Morizz00/self-destruct-share-api/utils"
)

type StoredFile struct {
	FileName      string `json:"filename"`
	MIME          string `json:"mime"`
	Data          []byte `json:"date"`
	Password      string `json:"password"`
	DownloadsLeft int    `json:"downloads_left"`
}

func Upload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Upload fail", http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	downloads := 1
	if parsed, err := strconv.Atoi(r.FormValue("downloads")); err == nil && parsed > 0 {
		downloads = parsed
	}

	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read ts shi", http.StatusInternalServerError)
		return
	}

	storeIt := storage.StoredFile{
		FileName:      fileHeader.Filename,
		MIME:          fileHeader.Header.Get("Content-Type"),
		Data:          fileData,
		Password:      password,
		DownloadsLeft: downloads,
	}
	id := utils.GenerateID()
	err = storage.StoreFile(id, storeIt)

	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File uploaded--Download:/file/%s\n", id)
}
