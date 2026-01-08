package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/Morizz00/self-destruct-share-api/utils"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Upload error: failed to get file from form: %v", err)
		http.Error(w, "Upload fail", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if err := utils.ValidateFileSize(fileHeader.Size); err != nil {
		log.Printf("Upload error: %v (size: %d)", err, fileHeader.Size)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	slug := r.FormValue("slug")
	if slug != "" {
		matched, _ := regexp.MatchString("^[a-z0-9-]+$", slug)
		if !matched {
			http.Error(w, "Invalid slug format", http.StatusBadRequest)
			return
		}

		_, err := storage.Get(slug)
		if err == nil {
			http.Error(w, "this custom link is already taken, try another one", http.StatusBadRequest)
			return
		}
	}

	downloads := 1
	if parsed, err := strconv.Atoi(r.FormValue("downloads")); err == nil && parsed > 0 {
		downloads = parsed
	}
	if err := utils.ValidateDownloads(downloads); err != nil {
		log.Printf("Upload error: invalid downloads: %d", downloads)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expiryMinutes := 5
	if parsed, err := strconv.Atoi(r.FormValue("expiry")); err == nil && parsed > 0 {
		expiryMinutes = parsed
	}
	if err := utils.ValidateExpiry(expiryMinutes); err != nil {
		log.Printf("Upload error: invalid expiry: %d minutes", expiryMinutes)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	expiry := time.Duration(expiryMinutes) * time.Minute

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Upload error: failed to read file: %v", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Double-check size after reading (in case Content-Length was wrong)
	if err := utils.ValidateFileSize(int64(len(fileData))); err != nil {
		log.Printf("Upload error: file size validation failed after read: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password if provided
	hashedPassword := ""
	if password != "" {
		hashedPassword, err = utils.HashPassword(password)
		if err != nil {
			log.Printf("Upload error: failed to hash password: %v", err)
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
	}

	// Sanitize filename
	sanitizedFilename := utils.SanitizeFilename(fileHeader.Filename)

	storeIt := storage.StoredFile{
		FileName:      sanitizedFilename,
		MIME:          fileHeader.Header.Get("Content-Type"),
		Data:          fileData,
		Password:      hashedPassword,
		DownloadsLeft: downloads,
		Expiry:        expiry,
	}
	var id string
	if slug != "" {
		id = slug
	} else {
		id = utils.GenerateID()
	}
	err = storage.StoreFile(id, storeIt, expiry)

	if err != nil {
		log.Printf("Upload error: storage failed: %v", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	log.Printf("File uploaded successfully: id=%s, filename=%s, size=%d, downloads=%d, expiry=%v", 
		id, sanitizedFilename, len(fileData), downloads, expiry)
	fmt.Fprintf(w, "File uploaded--Download:/file/%s\n", id)
}
