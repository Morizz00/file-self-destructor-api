package handlers

import (
	"net/http"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/go-chi/chi/v5"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := storage.GetAndDelete(id)
	if err != nil {
		http.Error(w, "file expired or not foundszies", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+id)
	w.Write(data)
}
