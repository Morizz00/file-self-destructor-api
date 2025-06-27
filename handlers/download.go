package handlers

import (
	"net/http"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/go-chi/chi/v5"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	storedData, err := storage.GetAndDelete(id)
	if err != nil {
		http.Error(w, "file expired or not foundszies", http.StatusNotFound)
		return
	}
	fileData := storedData.Data
	password := r.URL.Query().Get("password")
	if storedData.Password != "" && storedData.Password != password {
		http.Error(w, "wrong or missing password lil bud", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+storedData.FileName)
	w.Header().Set("Content-Type", storedData.MIME)
	w.Write(fileData)

}
