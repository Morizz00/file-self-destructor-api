package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Morizz00/self-destruct-share-api/storage"
	"github.com/go-chi/chi/v5"
)

type MetaResponse struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	URL           string `json:"url"`
	Type          string `json:"type"`
	SiteName      string `json:"site_name"`
	FileSize      string `json:"file_size"`
	FileType      string `json:"file_type"`
	DownloadsLeft int    `json:"downloads_left"`
	ExpiresAt     string `json:"expires_at"`
}

func GetMeta(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	storedData, err := storage.Get(id)
	if err != nil {
		defaultMeta := MetaResponse{
			Title:       "File Not Found - FileOrcha",
			Description: "This file has expired or doesn't exist.",
			Image:       "https://via.placeholder.com/1200x630/ef4444/ffffff?text=File+Not+Found",
			URL:         fmt.Sprintf("%s/download.html?id=%s", getBaseURL(r), id),
			Type:        "website",
			SiteName:    "FileOrcha",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defaultMeta)
		return
	}

	if storedData.DownloadsLeft <= 0 {
		defaultMeta := MetaResponse{
			Title:       "File Expired - FileOrcha",
			Description: "This file has no downloads remaining.",
			Image:       "https://via.placeholder.com/1200x630/f59e0b/ffffff?text=File+Expired",
			URL:         fmt.Sprintf("%s/download.html?id=%s", getBaseURL(r), id),
			Type:        "website",
			SiteName:    "FileOrcha",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defaultMeta)
		return
	}

	fileName := storedData.FileName
	fileSize := formatFileSize(len(storedData.Data))
	fileType := getFileTypeDisplay(storedData.MIME, fileName)

	title := fmt.Sprintf("%s - FileOrcha", fileName)
	if len(title) > 60 {
		title = fileName[:57] + "... - FileOrcha"
	}

	description := fmt.Sprintf("Download %s (%s) • %d downloads left • Self-destructing file",
		fileName, fileSize, storedData.DownloadsLeft)
	if len(description) > 160 {
		description = fmt.Sprintf("Download %s • %d downloads left • Self-destructing file",
			fileName, storedData.DownloadsLeft)
	}

	previewImage := generatePreviewImage(fileName, fileType)

	expiryTime := time.Now().Add(storedData.Expiry)
	expiresAt := expiryTime.Format("2006-01-02 15:04:05 UTC")

	meta := MetaResponse{
		Title:         title,
		Description:   description,
		Image:         previewImage,
		URL:           fmt.Sprintf("%s/download.html?id=%s", getBaseURL(r), id),
		Type:          "website",
		SiteName:      "FileOrcha",
		FileSize:      fileSize,
		FileType:      fileType,
		DownloadsLeft: storedData.DownloadsLeft,
		ExpiresAt:     expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}

func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func formatFileSize(bytes int) string {
	if bytes == 0 {
		return "0 Bytes"
	}

	const k = 1024
	sizes := []string{"Bytes", "KB", "MB", "GB"}
	i := 0
	size := float64(bytes)

	for size >= k && i < len(sizes)-1 {
		size /= k
		i++
	}

	return fmt.Sprintf("%.1f %s", size, sizes[i])
}

func getFileTypeDisplay(mimeType, fileName string) string {

	ext := strings.ToLower(fileName[strings.LastIndex(fileName, ".")+1:])

	extMap := map[string]string{
		"pdf": "PDF Document",
		"doc": "Word Document", "docx": "Word Document",
		"xls": "Excel Spreadsheet", "xlsx": "Excel Spreadsheet",
		"ppt": "PowerPoint Presentation", "pptx": "PowerPoint Presentation",
		"txt": "Text File", "md": "Markdown File",
		"jpg": "Image", "jpeg": "Image", "png": "Image", "gif": "Image", "webp": "Image",
		"mp4": "Video", "avi": "Video", "mov": "Video", "wmv": "Video",
		"mp3": "Audio", "wav": "Audio", "flac": "Audio",
		"zip": "Archive", "rar": "Archive", "7z": "Archive",
		"js": "JavaScript", "ts": "TypeScript", "py": "Python", "go": "Go",
		"html": "HTML", "css": "CSS", "json": "JSON",
	}

	if displayName, exists := extMap[ext]; exists {
		return displayName
	}

	if mimeType != "" {
		parts := strings.Split(mimeType, "/")
		if len(parts) > 1 {
			return strings.Title(parts[0]) + " " + strings.Title(parts[1])
		}
	}

	return "File"
}

func generatePreviewImage(fileName, fileType string) string {
	title := fileName
	if len(title) > 30 {
		title = title[:27] + "..."
	}
	imageURL := fmt.Sprintf("https://via.placeholder.com/1200x630/3b82f6/ffffff?text=%s",
		strings.ReplaceAll(title+" - "+fileType, " ", "+"))

	return imageURL
}
