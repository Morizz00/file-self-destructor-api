package handlers

import (
	"log"
	"net/http"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/Morizz00/self-destruct-share-api/utils"
	"github.com/go-chi/chi/v5"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	storedData, err := storage.Get(id)
	if err != nil {
		log.Printf("Download error: file not found: id=%s, error=%v", id, err)
		http.Error(w, "File not found or expired", http.StatusNotFound)
		return
	}
	fileData := storedData.Data
	password := r.URL.Query().Get("password")
	if storedData.Password != "" {
		if !utils.CheckPassword(password, storedData.Password) {
			log.Printf("Download error: wrong password: id=%s", id)
			http.Error(w, "Wrong or missing password", http.StatusForbidden)
			return
		}
	}
	if storedData.DownloadsLeft <= 0 {
		log.Printf("Download error: no downloads remaining: id=%s", id)
		http.Error(w, "No downloads remaining", http.StatusGone)
		return
	}
	if storedData.DownloadsLeft == 1 {
		err := storage.Delete(id)
		if err != nil {
			log.Printf("Download error: failed to delete file: id=%s, error=%v", id, err)
			http.Error(w, "Failed to self-destruct file", http.StatusInternalServerError)
			return
		}
		log.Printf("File self-destructed after download: id=%s", id)
	} else {
		storedData.DownloadsLeft--
		err := storage.UpdateFilePreservingTTL(id, storedData)
		if err != nil {
			log.Printf("Download error: failed to update download count: id=%s, error=%v", id, err)
			http.Error(w, "Failed to update download count", http.StatusInternalServerError)
			return
		}
		log.Printf("File downloaded: id=%s, downloads left=%d", id, storedData.DownloadsLeft)
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+storedData.FileName)
	w.Header().Set("Content-Type", storedData.MIME)
	w.Write(fileData)
}
