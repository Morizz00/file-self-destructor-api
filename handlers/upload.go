package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/Morizz00/self-destruct-share-api/utils"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Upload fail", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "storage failed", http.StatusInternalServerError)
		return
	}

	id := utils.GenerateID()
	err = storage.StoreFile(id, data)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	dowloadURL := fmt.Sprintf("/file/%s", id)
	fmt.Fprintf(w, "File uploaded--Download:/file/%s\n", dowloadURL)
}
