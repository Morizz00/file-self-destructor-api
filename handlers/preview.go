package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/Morizz00/self-destruct-share-api/utils"
	"github.com/go-chi/chi/v5"
)

type PreviewRequest struct {
	FileName      string `json:"filename"`
	FileSize      int    `json:"filesize"`
	MIME          string `json:"mime"`
	DownloadsLeft int    `json:"downloadleft"`
	HasPassword   bool   `json:"haspassword"`
	FileData      string `json:"filedata,omitempty"`
}

func Preview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	storedData, err := storage.Get(id)
	if err != nil {
		http.Error(w, "File not found or expired", http.StatusNotFound)
		return
	}
	password := r.URL.Query().Get("password")
	if storedData.Password != "" {
		if !utils.CheckPassword(password, storedData.Password) {
			http.Error(w, "Wrong or missing password", http.StatusForbidden)
			return
		}
	}
	if storedData.DownloadsLeft <= 0 {
		http.Error(w, "No downloads remaining", http.StatusGone)
		return
	}

	response := PreviewRequest{
		FileName:      storedData.FileName,
		FileSize:      len(storedData.Data),
		MIME:          storedData.MIME,
		DownloadsLeft: storedData.DownloadsLeft,
		HasPassword:   storedData.Password != "",
	}

	if len(storedData.Data) < 5*1024*1024 {
		w.Header().Set("Content-Type", storedData.MIME)
		w.Header().Set("X-File-Name", storedData.FileName)
		w.Header().Set("X-File-Size", strconv.Itoa(len(storedData.Data)))
		w.Header().Set("X-Downloads-Left", strconv.Itoa(storedData.DownloadsLeft))
		w.Write(storedData.Data)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
