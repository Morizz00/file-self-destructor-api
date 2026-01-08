package utils

import (
	"path/filepath"
	"strings"
)

const (
	MaxFileSize      = 50 * 1024 * 1024
	MaxDownloads     = 10
	MaxExpiryMinutes = 10080
)

func SanitizeFilename(filename string) string {
	// Remove any path components
	filename = filepath.Base(filename)

	// Remove null bytes and control characters
	filename = strings.ReplaceAll(filename, "\x00", "")
	filename = strings.Map(func(r rune) rune {
		if r < 32 && r != 9 && r != 10 && r != 13 {
			return -1
		}
		return r
	}, filename)

	// Limit length
	if len(filename) > 255 {
		filename = filename[:255]
	}

	// Remove leading/trailing spaces and dots
	filename = strings.TrimSpace(filename)
	filename = strings.Trim(filename, ".")

	// If empty after sanitization, use default
	if filename == "" {
		filename = "file"
	}

	return filename
}

// ValidateFileSize checks if file size is within limits
func ValidateFileSize(size int64) error {
	if size > MaxFileSize {
		return ErrFileTooLarge
	}

	if size <= 0 {
		return ErrInvalidFileSize
	}

	return nil
}

// ValidateDownloads checks if download count is within limits
func ValidateDownloads(downloads int) error {
	if downloads < 1 {
		return ErrInvalidDownloads
	}
	if downloads > MaxDownloads {
		return ErrDownloadsExceeded
	}
	return nil
}

// ValidateExpiry checks if expiry time is within limits
func ValidateExpiry(expiryMinutes int) error {
	if expiryMinutes < 1 {
		return ErrInvalidExpiry
	}
	if expiryMinutes > MaxExpiryMinutes {
		return ErrExpiryExceeded
	}
	return nil
}
